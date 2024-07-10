
# Logger

A structured and leveled logging package implemented using `zap`.

## Features

- Supports different log levels: `Debug`, `Info`, `Warn`, `Error`, `Fatal`.
- Structured logging with fields.
- Configurable log output formats.
- Thread-safe and high performance.
- Customizable exit function for `Fatal` log level.
- Swappable logging implementation via a common interface.

## Installation

```sh
go get -u go.uber.org/zap
go get -u github.com/stretchr/testify
```

## Logger Interface

The `Logger` interface defines the behavior of the logging methods. By using this interface, different logging implementations can be swapped without changing the caller code.

```go
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
```

## Usage

### Initializing the Logger

Create a new logger instance with the desired log level using the `Logger` interface. This allows you to easily swap out the logging implementation without changing the caller code.

```go
package main

import (
	"your/package/logger"
)

func main() {
	var log logger.Logger = logger.NewZapLogger(logger.InfoLevel)
	log.Info("Logger initialized", logger.Fields{"module": "main"})
}
```

### Logging with Different Levels

Log messages at different levels with structured fields using the `Logger` interface.

```go
log.Debug("This is a debug message", logger.Fields{"key": "value"})
log.Info("This is an info message", logger.Fields{"key": "value"})
log.Warn("This is a warn message", logger.Fields{"key": "value"})
log.Error("This is an error message", logger.Fields{"key": "value"})
log.Fatal("This is a fatal message", logger.Fields{"key": "value"}) // Exits the application
```

### Custom Exit Function for Fatal Logs

You can customize the exit behavior for fatal logs, which is useful for testing.

```go
log := logger.NewZapLogger(logger.InfoLevel)
zapLogger := log.(*logger.ZapLogger) // Type assertion to access ZapLogger specific methods
zapLogger.exitFunc = func(code int) {
    // Custom exit behavior
}
```

## Example

```go
package main

import (
	"your/package/logger"
	"time"
)

func main() {
	var log logger.Logger = logger.NewZapLogger(logger.DebugLevel)
	log.Debug("Debug message", logger.Fields{"example": "debug", "time": time.Now()})
	log.Info("Info message", logger.Fields{"example": "info", "count": 42})
	log.Warn("Warn message", logger.Fields{"example": "warn", "duration": time.Second})
	log.Error("Error message", logger.Fields{"example": "error", "status": true})
	log.Fatal("Fatal message", logger.Fields{"example": "fatal"}) // Will exit the application
}
```

## ZapLogger Implementation

The `ZapLogger` is an implementation of the `Logger` interface using `zap`.

```go
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger   *zap.Logger
	LogLevel LogLevel
	exitFunc func(int)
}

// NewZapLogger returns a new *ZapLogger
func NewZapLogger(level LogLevel) Logger {
	atomicLevel := zap.NewAtomicLevel()
	switch level {
	case DebugLevel:
		atomicLevel.SetLevel(zap.DebugLevel)
	case InfoLevel:
		atomicLevel.SetLevel(zap.InfoLevel)
	case WarnLevel:
		atomicLevel.SetLevel(zap.WarnLevel)
	case ErrorLevel:
		atomicLevel.SetLevel(zap.ErrorLevel)
	case FatalLevel:
		atomicLevel.SetLevel(zap.FatalLevel)
	default:
		atomicLevel.SetLevel(zap.InfoLevel)
	}

	loggerConfig := zap.Config{
		Level:         atomicLevel,
		Encoding:      "json",
		OutputPaths:   []string{"stdout"},
		EncoderConfig: zap.NewProductionEncoderConfig(),
	}
	loggerConfig.EncoderConfig.CallerKey = "caller"
	loggerConfig.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	loggerConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	logger, _ := loggerConfig.Build(zap.AddCallerSkip(1))

	return &ZapLogger{
		logger:   logger,
		LogLevel: level,
		exitFunc: func(int) {},
	}
}

func (z *ZapLogger) Debug(msg string, fields Fields) {
	if z.shouldLog(DebugLevel) {
		z.logger.Debug(msg, mapToZapFields(fields)...)
	}
}

func (z *ZapLogger) Info(msg string, fields Fields) {
	if z.shouldLog(InfoLevel) {
		z.logger.Info(msg, mapToZapFields(fields)...)
	}
}

func (z *ZapLogger) Warn(msg string, fields Fields) {
	if z.shouldLog(WarnLevel) {
		z.logger.Warn(msg, mapToZapFields(fields)...)
	}
}

func (z *ZapLogger) Error(msg string, fields Fields) {
	if z.shouldLog(ErrorLevel) {
		z.logger.Error(msg, mapToZapFields(fields)...)
	}
}

func (z *ZapLogger) Fatal(msg string, fields Fields) {
	if z.shouldLog(FatalLevel) {
		z.logger.Fatal(msg, mapToZapFields(fields)...)
		z.exitFunc(1)
	}
}

func (z *ZapLogger) GetLevel() LogLevel {
	return z.LogLevel
}

func (z *ZapLogger) shouldLog(level LogLevel) bool {
	return level >= z.LogLevel
}

func mapToZapFields(fields Fields) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		switch val := v.(type) {
		case string:
			zapFields = append(zapFields, zap.String(k, val))
		case int:
			zapFields = append(zapFields, zap.Int(k, val))
		case int64:
			zapFields = append(zapFields, zap.Int64(k, val))
		case float64:
			zapFields = append(zapFields, zap.Float64(k, val))
		case bool:
			zapFields = append(zapFields, zap.Bool(k, val))
		case []byte:
			zapFields = append(zapFields, zap.Binary(k, val))
		case time.Time:
			zapFields = append(zapFields, zap.Time(k, val))
		case time.Duration:
			zapFields = append(zapFields, zap.Duration(k, val))
		default:
			zapFields = append(zapFields, zap.Any(k, v))
		}
	}
	return zapFields
}
```

## Testing

Unit tests are provided to cover various logging scenarios. Use the following command to run the tests:

```sh
go test ./...
```

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss your ideas or improvements.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
