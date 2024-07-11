package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger is a logger implementation using zap.
type ZapLogger struct {
	logger   *zap.Logger
	LogLevel LogLevel
	exitFunc func(int)
}

// NewZapLogger returns a new *ZapLogger.
func NewZapLogger(level LogLevel) *ZapLogger {
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

	// Create a zap logger with the specified level
	loggerConfig := zap.Config{
		Level:         atomicLevel,
		Encoding:      "json",
		OutputPaths:   []string{"stdout"},
		EncoderConfig: zap.NewProductionEncoderConfig(),
	}

	// Configure the logger
	loggerConfig.EncoderConfig.CallerKey = "caller"
	loggerConfig.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	loggerConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	logger, err := loggerConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		fmt.Printf("Error building logger: %v\n", err)
		return nil
	}
	return &ZapLogger{
		logger:   logger,
		LogLevel: level,
		exitFunc: func(int) {}, // default to no-op
	}
}

// GetLogLevel returns the log level as a string.
func (l LogLevel) GetLogLevel() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return ""
	}
}

// Debug logs a debug message with structured fields.
func (z *ZapLogger) Debug(msg string, fields Fields) {
	if z.shouldLog(DebugLevel) {
		z.logger.Debug(msg, mapToZapFields(fields)...)
	}
}

// Info logs an info message with structured fields.
func (z *ZapLogger) Info(msg string, fields Fields) {
	if z.shouldLog(InfoLevel) {
		z.logger.Info(msg, mapToZapFields(fields)...)
	}
}

// Warn logs a warning message with structured fields.
func (z *ZapLogger) Warn(msg string, fields Fields) {
	if z.shouldLog(WarnLevel) {
		z.logger.Warn(msg, mapToZapFields(fields)...)
	}
}

// Error logs an error message with structured fields.
func (z *ZapLogger) Error(msg string, fields Fields) {
	if z.shouldLog(ErrorLevel) {
		z.logger.Error(msg, mapToZapFields(fields)...)
	}
}

// Fatal logs a fatal message with structured fields and exits the application.
func (z *ZapLogger) Fatal(msg string, fields Fields) {
	if z.shouldLog(FatalLevel) {
		z.logger.Fatal(msg, mapToZapFields(fields)...)
		z.exitFunc(1)
	}
}

// shouldLog determines if a log entry should be logged based on the log level.
func (z *ZapLogger) shouldLog(level LogLevel) bool {
	return level >= z.LogLevel
}

// mapToZapFields converts Fields to zap.Field with type-specific handling for better performance.
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
