package logger

import (
	"context"
	"log"
	"os"
)

type LogLevel int

const (
	InfoLevel LogLevel = iota
	DebugLevel
	ErrorLevel
)

type Logger struct {
	*log.Logger
	level LogLevel
}

type NewLoggerParams struct {
	LogPath string
	Level   LogLevel
}

// NewLogger creates a new logger with a specified log level
func NewLogger(ctx context.Context, params NewLoggerParams) (*Logger, error) {
	f, err := os.OpenFile(params.LogPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	return &Logger{
		Logger: log.New(f, "[go-lsp]: ", log.Ldate|log.Ltime|log.Lshortfile),
		level:  params.Level,
	}, nil
}

// logf logs a formatted message at a specific log level
func (l *Logger) logf(level LogLevel, prefix string, format string, v ...interface{}) {
	if l.level >= level {
		l.Printf(prefix+" "+format, v...)
	}
}

// Info logs a formatted info level message
func (l *Logger) Info(format string, v ...interface{}) {
	l.logf(InfoLevel, "INFO:", format, v...)
}

// Debug logs a formatted debug level message
func (l *Logger) Debug(format string, v ...interface{}) {
	l.logf(DebugLevel, "DEBUG:", format, v...)
}

// Error logs a formatted error level message
func (l *Logger) Error(format string, v ...interface{}) {
	l.logf(ErrorLevel, "ERROR:", format, v...)
}
