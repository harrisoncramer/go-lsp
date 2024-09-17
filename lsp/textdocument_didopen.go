package lsp

type DidOpenTextDocumentNotification struct {
	Notification
	Params struct {
		TextDocumentItem textDocumentItem `json:"textDocument"`
	}
}
