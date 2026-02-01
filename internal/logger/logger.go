package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents log level
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger provides structured logging
type Logger struct {
	verbose bool
	output  io.Writer
}

// New creates a new logger
func New(verbose bool) *Logger {
	return &Logger{
		verbose: verbose,
		output:  os.Stdout,
	}
}

// NewWithWriter creates a logger with a custom writer
func NewWithWriter(verbose bool, w io.Writer) *Logger {
	return &Logger{
		verbose: verbose,
		output:  w,
	}
}

// Debug logs a debug message (only if verbose is enabled)
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.verbose {
		l.log(LevelDebug, format, args...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LevelWarn, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	prefix := l.levelPrefix(level)

	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.output, "%s %s %s\n", timestamp, prefix, message)
}

func (l *Logger) levelPrefix(level Level) string {
	switch level {
	case LevelDebug:
		return "[DEBUG]"
	case LevelInfo:
		return "[INFO] "
	case LevelWarn:
		return "[WARN] "
	case LevelError:
		return "[ERROR]"
	default:
		return "[???]  "
	}
}

// SetOutput sets the output writer
func (l *Logger) SetOutput(w io.Writer) {
	l.output = w
}
