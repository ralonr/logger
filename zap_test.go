package logger

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestNewZapLogger(t *testing.T) {
	logger := NewZapLogger(InfoLevel)
	assert.NotNil(t, logger, "Logger should not be nil")
	assert.Equal(t, InfoLevel, logger.LogLevel, "LogLevel should be InfoLevel")
}

func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{DebugLevel, "debug"},
		{InfoLevel, "info"},
		{WarnLevel, "warn"},
		{ErrorLevel, "error"},
		{FatalLevel, "fatal"},
		{LogLevel(999), ""},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.level.String())
	}
}

func TestLogMethods(t *testing.T) {
	tests := []struct {
		level         LogLevel
		logMethod     func(z *ZapLogger, msg string, fields Fields)
		expectedLevel zapcore.Level
	}{
		{DebugLevel, (*ZapLogger).Debug, zap.DebugLevel},
		{InfoLevel, (*ZapLogger).Info, zap.InfoLevel},
		{WarnLevel, (*ZapLogger).Warn, zap.WarnLevel},
		{ErrorLevel, (*ZapLogger).Error, zap.ErrorLevel},
	}

	for _, test := range tests {
		t.Run(test.expectedLevel.String(), func(t *testing.T) {
			buf := new(bytes.Buffer)
			ws := zapcore.AddSync(buf)
			encoderConfig := zap.NewProductionEncoderConfig()
			core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), ws, zap.NewAtomicLevelAt(test.expectedLevel))
			logger := &ZapLogger{
				logger:   zap.New(core),
				LogLevel: test.level,
			}

			fields := Fields{"key": "value"}
			test.logMethod(logger, "test message", fields)
			assert.Contains(t, buf.String(), "test message")
			assert.Contains(t, buf.String(), `"key":"value"`)
		})
	}
}

func TestShouldLog(t *testing.T) {
	logger := NewZapLogger(InfoLevel)
	assert.False(t, logger.shouldLog(DebugLevel), "Debug level should not be logged")
	assert.True(t, logger.shouldLog(InfoLevel), "Info level should be logged")
	assert.True(t, logger.shouldLog(WarnLevel), "Warn level should be logged")
	assert.True(t, logger.shouldLog(ErrorLevel), "Error level should be logged")
	assert.True(t, logger.shouldLog(FatalLevel), "Fatal level should be logged")
}

func TestFatalLogMethod(t *testing.T) {
	buf := new(bytes.Buffer)
	ws := zapcore.AddSync(buf)
	encoderConfig := zap.NewProductionEncoderConfig()
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), ws, zap.NewAtomicLevelAt(zap.FatalLevel))
	logger := &ZapLogger{
		logger:   zap.New(core),
		LogLevel: FatalLevel,
		exitFunc: func(int) {}, // Mock the exit function
	}

	fields := Fields{"key": "value"}

	// Use a recoverable function to test Fatal
	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, buf.String(), "test message")
			assert.Contains(t, buf.String(), `"key":"value"`)
		}
	}()
	logger.Fatal("test message", fields)
}

func TestMapToZapFields(t *testing.T) {
	fields := Fields{
		"string":    "test",
		"int":       123,
		"int64":     int64(1234567890),
		"float64":   1.23,
		"bool":      true,
		"bytes":     []byte("test"),
		"interface": struct{}{},
	}

	zapFields := mapToZapFields(fields)

	assert.Equal(t, len(fields), len(zapFields), "The number of zap fields should match the number of input fields")
	for _, field := range zapFields {
		switch field.Key {
		case "string":
			assert.Equal(t, "test", field.String)
		case "int":
			assert.Equal(t, int64(123), field.Integer)
		case "int64":
			assert.Equal(t, int64(1234567890), field.Integer)
		case "float64":
			if field.Interface != nil {
				assert.Equal(t, 1.23, field.Interface.(float64))
			}
		case "bool":
			assert.Equal(t, true, field.Integer == 1)
		case "bytes":
			if field.Interface != nil {
				assert.Equal(t, []byte("test"), field.Interface.([]byte))
			}
		case "interface":
			assert.IsType(t, struct{}{}, field.Interface)
		default:
			t.Errorf("Unexpected field key: %s", field.Key)
		}
	}
}
