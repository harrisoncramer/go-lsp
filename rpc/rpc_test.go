package rpc_test

import (
	"testing"

	"github.com/harrisoncramer/go-lsp/lsp"
	"github.com/harrisoncramer/go-lsp/rpc"
)

type EncodingExample struct {
	Method string
}

var goodMsg = "Content-Length: 54\r\n\r\n{\"jsonrpc\":\"2.0\",\"result\":{\"contents\":\"Hello World!\"}}"

func TestEncodeMessage(t *testing.T) {
	t.Run("Should encode a simple hover message", func(t *testing.T) {
		got := rpc.EncodeMessage(
			lsp.HoverResponse{
				Response: lsp.Response{
					RPC: "2.0",
				},
				Result: lsp.HoverResults{
					Contents: "Hello World!",
				},
			},
		)
		want := goodMsg
		if want != got {
			t.Fatalf("Got:\n%s\nWant:\n%s\n", got, want)
		}
	})

}

func TestDecodeMessage(t *testing.T) {
	t.Run("Should get the length of a message", func(t *testing.T) {
		_, content, err := rpc.DecodeMessage([]byte(goodMsg))
		if err != nil {
			t.Fatal(err)
		}
		want := 54
		if len(content) != want {
			t.Errorf("Got %d but wanted %d", len(content), want)
		}
	})

	t.Run("Should error for missing content length", func(t *testing.T) {
		missingContentLengthMsg := "Content-Length: \r\n\r\n{\"jsonrpc\":\"2.0\",\"result\":{\"contents\":\"Hello World!\"}}"
		_, _, err := rpc.DecodeMessage([]byte(missingContentLengthMsg))
		if err != rpc.ErrHeaderNotFound {
			t.Fatalf("Got %v but wanted %v", err, rpc.ErrHeaderNotFound)
		}
	})

}
