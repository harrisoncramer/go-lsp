package lsp

type Request struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Method string `json:"method"`
}

type Response struct {
	RPC string `json:"jsonrpc"`
	ID  *int   `json:"id,omitempty"`
}

type Notification struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
}

// Union of all possible request types
type LSPRequest interface {
	DidChangeTextDocumentNotification |
		DidSaveTextDocumentNotification |
		DidOpenTextDocumentNotification |
		HoverRequest |
		InitializeRequest |
		DefinitionRequest
}

// Union of all possible response types
type LSPResponse interface {
	HoverResponse |
		InitializeResponse |
		DefinitionResponse
}
