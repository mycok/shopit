package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// Level type represents the severity level for an entry log.
type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LevelOff
)

// String() returns a human-friendly message for each severity level.
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}


// Logger type holds the output destination that the log entries
// will be written to, the minimum severity level that log entries will be written for,
// plus a mutex for coordinating the writes.
type Logger struct {
	out      io.Writer
	minLevel Level
	mu      sync.Mutex
}

// New() returns a new Logger instance which writes log entries at or above a minimum severity
// level to a specific output destination.
func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out: out,
		minLevel: minLevel,
	}
}

func (l *Logger) LogInfo(message string, props map[string]string) {
	l.log(LevelInfo, message, props)
}

func (l *Logger) LogError(err error, props map[string]string) {
	l.log(LevelError, err.Error(), props)
}

func (l *Logger) LogFatal(err error, props map[string]string) {
	l.log(LevelFatal, err.Error(), props)
	// Terminate the program immediately.
	os.Exit(1)
}

func (l *Logger) log(level Level, message string, props map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	inputs := struct {
		Level string `json:"level"`
		Time string `json:"time"`
		Message string `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace string `json:"trace,omitempty"`

	}{
		Level: level.String(),
		Time: time.Now().UTC().Format(time.RFC3339),
		Message: message,
		Properties: props,
	}

	// Include a stack trace for entries at the ERROR and FATAL levels.
	if level >= LevelError {
		inputs.Trace = string(debug.Stack())
	}

	// Declare a line variable for holding the actual log entry text.
	var line []byte

	// Marshal the anonymous struct to JSON and store it in the line variable. If there
	// is a problem creating the JSON, set the contents of the log entry to be that
	// plain-text error message instead.
	line, err := json.Marshal(&inputs)
	if err != nil {
		line = []byte(level.String() + ": unable to marshal log data:" + err.Error())
	}

	// Lock the mutex so that no two writes to the output destination can happen concurrently.
	// If we don't do this, it's possible that the text for two or more
	// log entries will be intermingled in the output.
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}

// Write() implements io.Writer.Write(p []byte) method.
func (l *Logger) Write(message []byte) (int, error) {
	return l.log(LevelError, string(message), nil)
}