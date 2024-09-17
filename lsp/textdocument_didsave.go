package lsp

type DidSaveTextDocumentNotification struct {
	Notification
	Params DidSaveTextDocumentParams
}

type DidSaveTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}
