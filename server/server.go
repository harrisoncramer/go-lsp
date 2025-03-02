package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/harrisoncramer/go-lsp/analysis"
	"github.com/harrisoncramer/go-lsp/logger"
	"github.com/harrisoncramer/go-lsp/lsp"
	"github.com/harrisoncramer/go-lsp/rpc"
)

type Server struct {
	logger *logger.Logger
	parser rpc.Rpc
}

func NewServer(ctx context.Context, logger *logger.Logger) Server {
	parser := rpc.NewParser(logger)
	return Server{
		logger: logger,
		parser: parser,
	}
}

// Starts the server and splits the RPC stream, processing each message in turn
func (s Server) Start() {
	s.logger.Println("server started")

	writer := os.Stdout
	state := analysis.NewState()
	parser := bufio.NewScanner(os.Stdin)
	parser.Split(s.parser.Split)

	for parser.Scan() {
		msg := parser.Bytes()
		method, contents, err := s.parser.DecodeMessage(msg)
		if err != nil {
			s.logger.Printf("failed to decode msg: %s", err)
			continue
		}

		s.handleMessage(writer, method, contents, state)
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
func (s Server) handleMessage(writer io.Writer, method string, contents []byte, state analysis.State) {
	s.logger.Printf("msg: %s", method)
	switch method {
	case didChangeMsg:
		if request, err := decodeMsg[lsp.DidChangeTextDocumentNotification](contents); err == nil {
			for _, change := range request.Params.ContentChanges {
				state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
			}
		}
	case didOpenMsg:
		if request, err := decodeMsg[lsp.DidOpenTextDocumentNotification](contents); err == nil {
			state.OpenDocument(
				request.Params.TextDocumentItem.URI,
				request.Params.TextDocumentItem.Text,
			)
		}
	case saveMsg:
		if request, err := decodeMsg[lsp.DidSaveTextDocumentNotification](contents); err == nil {
			state.Save(request.Params.TextDocument.URI)
		}
	case hoverMsg:
		if request, err := decodeMsg[lsp.HoverRequest](contents); err == nil {
			msg := state.Hover(request.ID, request.Params.TextDocumentPositionParams)
			sendResponse(s.parser, msg, writer, s.logger)
		}
	case definitionMsg:
		if request, err := decodeMsg[lsp.DefinitionRequest](contents); err == nil {
			msg := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
			sendResponse(s.parser, msg, writer, s.logger)
		}
	case initializeMsg:
		if request, err := decodeMsg[lsp.InitializeRequest](contents); err == nil {
			msg := lsp.NewInitializeResponse(request.ID)
			sendResponse(s.parser, msg, writer, s.logger)
			s.logger.Println("server initialized")
		}
	case initializedMsg:
		s.logger.Println("server initialized")
	}
}

// Unmarshals a message from the client into the appropriate message struct
func decodeMsg[T lsp.LSPRequest](contents []byte) (*T, error) {
	var res T
	if err := json.Unmarshal(contents, &res); err != nil {
		return nil, fmt.Errorf("failed to decode %T message: %w", res, err)
	}

	return &res, nil
}

// Marshals the response type into an RPC message and writes it to the client
func sendResponse[T lsp.LSPResponse](parser rpc.Rpc, msg T, writer io.Writer, logger *logger.Logger) {
	response := parser.EncodeMessage(msg)
	if writer == nil {
		panic("Writer not provided")
	}
	_, err := writer.Write([]byte(response))
	if err != nil {
		logger.Printf("failed to send message to client: %v", err)
	}
}
