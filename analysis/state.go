package analysis

import (
	"fmt"

	"github.com/harrisoncramer/go-lsp/lsp"
)

type State struct {
	Documents map[string]string
}

func NewState() State {
	return State{
		Documents: map[string]string{},
	}
}

// Creates the text for the given uri (document) in state
func (s *State) OpenDocument(uri, text string) {
	s.Documents[uri] = text
}

// Updates the text for the given uri (document) in state
func (s *State) UpdateDocument(uri, text string) {
	s.Documents[uri] = text
}

// These methods don't really do what's expected for a hover or
// a definition response, they just update the state for demo
// purposes...

func (s *State) Hover(id int, params lsp.TextDocumentPositionParams) lsp.HoverResponse {
	content := fmt.Sprintf("File: %s; Chars: %d", params.TextDocument.URI, len(s.Documents[params.TextDocument.URI]))
	msg := lsp.NewHoverResponse(id, content)
	return msg
}

func (s *State) Definition(id int, uri string, pos lsp.Position) lsp.DefinitionResponse {
	msg := lsp.NewDefinitionResponse(id, uri, pos)
	return msg
}

func (s *State) Save(uri string) {
	// Run some static analysis on the new text document content
}
