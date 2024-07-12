package logger

import (
	"bytes"
	"fmt"
	"testing"
)

// TestNewZap tests the NewZap function.
func TestNewZap(t *testing.T) {
	buffer := new(bytes.Buffer)
	config := Config{
		Level:    InfoLevel,
		Output:   buffer,
		ExitFunc: func(int) {},
	}
	zapLogger := NewZap(config)

	if zapLogger.Config.Level != InfoLevel {
		t.Errorf("Expected log level %v, got %v", InfoLevel, zapLogger.Config.Level)
	}

	if zapLogger.Config.ExitFunc == nil {
		t.Errorf("Expected ExitFunc to be set, but it is nil")
	}

	// Use a buffer to test log output
	zapLogger.Info("Info message", Fields{"key": "value"})
	expected := "Info message"
	if !bytes.Contains(buffer.Bytes(), []byte(expected)) {
		t.Errorf("Expected %s to contain %s", buffer.String(), expected)
	}
}

// TestZap_Debug tests the Debug method.
func TestZap_Debug(t *testing.T) {
	buffer := new(bytes.Buffer)
	config := Config{
		Level:    DebugLevel,
		Output:   buffer,
		ExitFunc: func(int) {},
	}
	zapLogger := NewZap(config)

	zapLogger.Debug("Debug message", Fields{"key": "value"})
	expected := "Debug message"
	if !bytes.Contains(buffer.Bytes(), []byte(expected)) {
		t.Errorf("Expected %s to contain %s", buffer.String(), expected)
	}
}

// TestZap_Info tests the Info method.
func TestZap_Info(t *testing.T) {
	buffer := new(bytes.Buffer)
	config := Config{
		Level:    InfoLevel,
		Output:   buffer,
		ExitFunc: func(int) {},
	}
	zapLogger := NewZap(config)

	zapLogger.Info("Info message", Fields{"key": "value"})
	expected := "Info message"
	if !bytes.Contains(buffer.Bytes(), []byte(expected)) {
		t.Errorf("Expected %s to contain %s", buffer.String(), expected)
	}
}

// TestZap_Warn tests the Warn method.
func TestZap_Warn(t *testing.T) {
	buffer := new(bytes.Buffer)
	config := Config{
		Level:    WarnLevel,
		Output:   buffer,
		ExitFunc: func(int) {},
	}
	zapLogger := NewZap(config)

	zapLogger.Warn("Warn message", Fields{"key": "value"})
	expected := "Warn message"
	if !bytes.Contains(buffer.Bytes(), []byte(expected)) {
		t.Errorf("Expected %s to contain %s", buffer.String(), expected)
	}
}

// TestZap_Error tests the Error method.
func TestZap_Error(t *testing.T) {
	buffer := new(bytes.Buffer)
	config := Config{
		Level:    ErrorLevel,
		Output:   buffer,
		ExitFunc: func(int) {},
	}
	zapLogger := NewZap(config)

	zapLogger.Error("Error message", Fields{"key": "value"})
	expected := "Error message"
	if !bytes.Contains(buffer.Bytes(), []byte(expected)) {
		t.Errorf("Expected %s to contain %s", buffer.String(), expected)
	}
}

// TestShouldLog tests the shouldLog function.
func TestShouldLog(t *testing.T) {
	tests := []struct {
		configLevel Level
		logLevel    Level
		expected    bool
	}{
		{DebugLevel, DebugLevel, true},
		{DebugLevel, InfoLevel, true},
		{DebugLevel, WarnLevel, true},
		{DebugLevel, ErrorLevel, true},
		{DebugLevel, FatalLevel, true},
		{InfoLevel, DebugLevel, false},
		{InfoLevel, InfoLevel, true},
		{InfoLevel, WarnLevel, true},
		{InfoLevel, ErrorLevel, true},
		{InfoLevel, FatalLevel, true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v-%v", test.configLevel, test.logLevel), func(t *testing.T) {
			config := Config{
				Level: test.configLevel,
			}
			zapLogger := NewZap(config)
			actual := zapLogger.shouldLog(test.logLevel)
			if actual != test.expected {
				t.Errorf("shouldLog(%v) = %v, expected %v", test.logLevel, actual, test.expected)
			}
		})
	}
}

// TestMapToZapFields tests the mapToZapFields function.
func TestMapToZapFields(t *testing.T) {
	fields := Fields{
		"string":  "value",
		"int":     42,
		"int64":   int64(64),
		"float64": 3.14,
		"bool":    true,
	}

	zapFields := mapToZapFields(fields)

	if len(zapFields) != len(fields) {
		t.Errorf("Expected %d fields, got %d", len(fields), len(zapFields))
	}
}
