package logger

import (
	"encoding/json"
	"log"
	"os"
)

type LoggerOptions interface {
	Path() string
}

type Logger struct {
	*log.Logger
}

// Generates a new logger that writes to the given filename
func NewLogger(logPath string) (*Logger, error) {
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
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
