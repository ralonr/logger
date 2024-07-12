package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Zap is a logger implementation using zap.
type Zap struct {
	logger *zap.Logger
	Config Config
}

// NewZap returns a new *Zap.
func NewZap(config Config) *Zap {
	atomicLevel := zap.NewAtomicLevel()
	switch config.Level {
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

	output := zapcore.AddSync(config.Output)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.CallerKey = "caller"
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), output, atomicLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	if config.ExitFunc == nil {
		config.ExitFunc = os.Exit // default to os.Exit
	}

	return &Zap{
		logger: logger,
		Config: config,
	}
}

// Debug logs a debug message with structured fields.
func (z *Zap) Debug(msg string, fields Fields) {
	if z.shouldLog(DebugLevel) {
		z.logger.Debug(msg, mapToZapFields(fields)...)
	}
}

// Info logs an info message with structured fields.
func (z *Zap) Info(msg string, fields Fields) {
	if z.shouldLog(InfoLevel) {
		z.logger.Info(msg, mapToZapFields(fields)...)
	}
}

// Warn logs a warning message with structured fields.
func (z *Zap) Warn(msg string, fields Fields) {
	if z.shouldLog(WarnLevel) {
		z.logger.Warn(msg, mapToZapFields(fields)...)
	}
}

// Error logs an error message with structured fields.
func (z *Zap) Error(msg string, fields Fields) {
	if z.shouldLog(ErrorLevel) {
		z.logger.Error(msg, mapToZapFields(fields)...)
	}
}

// Fatal logs a fatal message with structured fields and exits the application.
func (z *Zap) Fatal(msg string, fields Fields) {
	if z.shouldLog(FatalLevel) {
		z.logger.Fatal(msg, mapToZapFields(fields)...)
		z.Config.ExitFunc(1)
	}
}

// shouldLog determines if a log entry should be logged based on the log level.
func (z *Zap) shouldLog(level Level) bool {
	return level >= z.Config.Level
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
		default:
			zapFields = append(zapFields, zap.Any(k, v))
		}
	}
	return zapFields
}
