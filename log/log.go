package log

import (
	"log"
	"os"
	"path"
)

type LoggerOptions interface {
	Path() string
}

// Generates a new logger that writes to the given filename
func NewLogger() (*log.Logger, error) {
	fileName, err := makeLogPath()
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	return log.New(f, "[go-lsp]: ", log.Ldate|log.Ltime|log.Lshortfile), nil
}

func makeLogPath() (string, error) {
	return path.Join("/tmp", "go-lsp.log"), nil
}
