package logger

import (
	"testing"
)

// TestInfo tests the Info function of the logger
func TestLogger(t *testing.T) {
	msg := "Test Info message"
	Info(msg)
	Error(msg)

}
