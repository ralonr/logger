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
	Debug(msg string, fields map[string]any)
	Info(msg string, fields map[string]any)
	Warn(msg string, fields map[string]any)
	Error(msg string, fields map[string]any)
	Fatal(msg string, fields map[string]any)

	GetLevel() LogLevel
}
