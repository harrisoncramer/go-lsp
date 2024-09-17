package run

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/harrisoncramer/go-lsp/analysis"
	"github.com/harrisoncramer/go-lsp/logger"
	"github.com/harrisoncramer/go-lsp/lsp"
	"github.com/harrisoncramer/go-lsp/rpc"
)

// Starts the server and splits the RPC stream, processing each message in turn
func Start(logger *logger.Logger) {
	logger.Println("server started")

	writer := os.Stdout
	state := analysis.NewState()
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Split(rpc.Split)
	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("failed to decode msg: %s", err)
			continue
		}

		handleMessage(logger, writer, method, contents, state)
	}
}

// All the different types of messages the client could send
const (
	hoverMsg       = "textDocument/hover"
	didOpenMsg     = "textDocument/didOpen"
	definitionMsg  = "textDocument/definition"
	saveMsg        = "textDocument/didSave"
	didChangeMsg   = "textDocument/didChange"
	initializeMsg  = "initialize"
	initializedMsg = "initialized"
)

// Handles each message type by decoding the request and responding or updating internal state
func handleMessage(logger *logger.Logger, writer io.Writer, method string, contents []byte, state analysis.State) {
	logger.Printf("msg: %s", method)
	switch method {
	case didChangeMsg:
		if request, err := decodeMsg[lsp.DidChangeTextDocumentNotification](contents, logger); err == nil {
			for _, change := range request.Params.ContentChanges {
				state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
			}
		}
	case didOpenMsg:
		if request, err := decodeMsg[lsp.DidOpenTextDocumentNotification](contents, logger); err == nil {
			state.OpenDocument(
				request.Params.TextDocumentItem.URI,
				request.Params.TextDocumentItem.Text,
			)
		}
	case saveMsg:
		if request, err := decodeMsg[lsp.DidSaveTextDocumentNotification](contents, logger); err == nil {
			state.Save(request.Params.TextDocument.URI)
		}
	case hoverMsg:
		if request, err := decodeMsg[lsp.HoverRequest](contents, logger); err == nil {
			msg := state.Hover(request.ID, request.Params.TextDocumentPositionParams)
			sendResponse(msg, writer, logger)
		}
	case definitionMsg:
		if request, err := decodeMsg[lsp.DefinitionRequest](contents, logger); err == nil {
			msg := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
			sendResponse(msg, writer, logger)
		}
	case initializeMsg:
		if request, err := decodeMsg[lsp.InitializeRequest](contents, logger); err == nil {
			msg := lsp.NewInitializeResponse(request.ID)
			sendResponse(msg, writer, logger)
		}
	case initializedMsg:
		logger.Println("server initialized")
	}
}

// Unmarshals a message from the client into the appropriate message struct
func decodeMsg[T lsp.LSPRequest](contents []byte, logger *logger.Logger) (*T, error) {
	var res T
	if err := json.Unmarshal(contents, &res); err != nil {
		return nil, fmt.Errorf("failed to decode %T message: %w", res, err)
	}

	return &res, nil
}

// Marshals the response type into an RPC message and writes it to the client
func sendResponse[T lsp.LSPResponse](msg T, writer io.Writer, logger *logger.Logger) {
	response := rpc.EncodeMessage(msg)
	if writer == nil {
		panic("Writer not provided")
	}
	_, err := writer.Write([]byte(response))
	if err != nil {
		logger.Printf("failed to send message to client: %v", err)
	}
}
