package reallysimplelogger

import (
	"encoding/json"
	"io"
	"runtime/debug"
	"sync"
	"time"
)

type LogLevel int8

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn 
	LevelError
	LevelFatal
)

func (l LogLevel) String() string {
	return []string{"INFO", "DEBUG", "WARN", "ERROR", "FATAL"}[l]
}


type Logger struct {
	out      io.Writer
	minLevel LogLevel
	mu       sync.Mutex
}

func New(out io.Writer, minLevel LogLevel) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

func (l *Logger) Info(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) Error(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

func (l *Logger) Fatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
}
func (l *Logger) Warn(err error, properties map[string]string) {
	l.print(LevelWarn, err.Error(), properties)
}
func (l *Logger) Debug(err error, properties map[string]string) {
	l.print(LevelDebug, err.Error(), properties)
}

func (l *Logger) print(level LogLevel, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message:" + err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}

func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
