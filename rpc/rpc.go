package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/harrisoncramer/go-lsp/lsp"
)

var ErrHeaderNotFound = errors.New("message missing header")
var headerText = "Content-Length: "
var headerSep = []byte{'\r', '\n', '\r', '\n'}

// Encodes a message struct into a string that adheres to the protocol
func EncodeMessage[T lsp.LSPResponse](msg T) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s%d%s%s", headerText, len(content), headerSep, content)
}

type BaseMessage struct {
	Method string `json:"method"`
}

// Decodes a byte slice and extracts the method and message content
func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, headerSep)
	if !found {
		return "", nil, ErrHeaderNotFound
	}

	// Content-Length: <number>
	contentLengthBytes := header[len(headerText):]
	if len(contentLengthBytes) == 0 {
		return "", nil, ErrHeaderNotFound
	}

	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", nil, err
	}

	var baseMessage BaseMessage
	if err := json.Unmarshal(content[:contentLength], &baseMessage); err != nil {
		return "", nil, err
	}

	return baseMessage.Method, content[:contentLength], nil
}

// Splits a byte slice, used to satisfy a scanner parsing the stdin data
func Split(data []byte, _ bool) (advance int, token []byte, err error) {

	header, _, contentLength, err := ParseHeader(data)
	if err != nil {
		return 0, nil, err
	}

	totalLength := len(header) + len(headerSep) + contentLength
	return totalLength, data[:totalLength], nil
}

// Parses the header for an RPC message
func ParseHeader(msg []byte) ([]byte, []byte, int, error) {

	header, content, found := bytes.Cut(msg, headerSep)
	if !found {
		return nil, nil, 0, ErrHeaderNotFound
	}

	// Content-Length: <number>
	contentLengthBytes := header[len(headerText):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return nil, nil, 0, err
	}

	if len(content) < contentLength {
		return nil, nil, 0, nil
	}

	return header, content, contentLength, nil

}
