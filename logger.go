package logger

// LogLevel represents the severity of the log message.
type LogLevel int

// Fields represents a map of key-value pairs for structured logging.
type Fields map[string]any

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Logger implements the behaviour of the logging methods
type Logger interface {
	Debug(msg string, fields Fields)
	Info(msg string, fields Fields)
	Warn(msg string, fields Fields)
	Error(msg string, fields Fields)
	Fatal(msg string, fields Fields)

	GetLevel() LogLevel
}
