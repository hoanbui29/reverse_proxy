package logger

import (
	"encoding/json"
	"io"
	"runtime/debug"
	"sync"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type logMessage struct {
	Level      string
	Time       string
	Message    string
	Properties map[string]string
	Trace      string
}

func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
		mu:       &sync.Mutex{},
	}
}

type Logger struct {
	out      io.Writer
	minLevel Level
	mu       *sync.Mutex
}

func (l *Logger) Debug(message string, properties map[string]string) {
	l.print(LevelDebug, message, properties)
}

func (l *Logger) Info(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) Warn(message string, properties map[string]string) {
	l.print(LevelWarn, message, properties)
}

func (l *Logger) Error(message string, properties map[string]string) {
	l.print(LevelError, message, properties)
}

func (l *Logger) Fatal(message string, properties map[string]string) {
	l.print(LevelFatal, message, properties)
}

func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	logMessage := logMessage{
		Level:      level.String(),
		Time:       time.Now().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	if level >= LevelError {
		logMessage.Trace = string(debug.Stack())
	}

	lines := make([]byte, 0)

	lines, err := json.Marshal(logMessage)

	if err != nil {
		lines = []byte(LevelError.String() + ": Unable to marshal JSON for logging" + err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(lines, '\n'))
}
