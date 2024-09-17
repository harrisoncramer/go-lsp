package main

import (
	"github.com/harrisoncramer/go-lsp/logger"
	"github.com/harrisoncramer/go-lsp/run"
)

func main() {

	logger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	run.Start(logger)
}
