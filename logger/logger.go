package logger

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

type LoggerOptions interface {
	Path() string
}

type Logger struct {
	*log.Logger
}

// Generates a new logger that writes to the given filename
func NewLogger() (*Logger, error) {
	fileName, err := makeLogPath()
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	return &Logger{
		log.New(f, "[go-lsp]: ", log.Ldate|log.Ltime|log.Lshortfile),
	}, nil
}

func (l *Logger) PrintJSON(v any) {
	prettyJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		l.Printf("could not pretty-print: %#v", v)
		return
	}

	l.Println(string(prettyJSON))
}

func makeLogPath() (string, error) {
	return path.Join("/tmp", "go-lsp.log"), nil
}
