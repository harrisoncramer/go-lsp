package lsp

type DefinitionRequest struct {
	Request
	Params DefinitionParams `json:"params"`
}

type DefinitionParams struct {
	TextDocumentPositionParams
}

type DefinitionResponse struct {
	Response
	Result Location `json:"result"`
}

type DefinitionResults struct {
	Contents string `json:"contents"`
}

func NewDefinitionResponse(id int, uri string, pos Position) DefinitionResponse {
	return DefinitionResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: Location{
			TextDocumentIdentifier: TextDocumentIdentifier{
				URI: uri,
			},
			Range: Range{
				Start: Position{
					Line:      pos.Line - 1, // Stubbed behavior...
					Character: pos.Character - 1,
				},
				End: Position{
					Line:      pos.Line - 1,
					Character: pos.Character - 1,
				},
			},
		},
	}
}
