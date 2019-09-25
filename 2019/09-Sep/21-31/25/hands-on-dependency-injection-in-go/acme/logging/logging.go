package logging

import "fmt"

// Logger is our standard interface
type Logger interface {
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
}

// L is the global instance of the logger
var L = &LoggerStdout{}

// LoggerStdout logs to std out
type LoggerStdout struct{}

// Debug logs message at DEBUG level
func (l LoggerStdout) Debug(message string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+message, args...)
}

// Info logs message at INFO level
func (l LoggerStdout) Info(message string, args ...interface{}) {
	fmt.Printf("[INFO] "+message, args...)
}

// Warn logs message at WARN level
func (l LoggerStdout) Warn(message string, args ...interface{}) {
	fmt.Printf("[WARN] "+message, args...)
}

// Error logs message at ERROR level
func (l LoggerStdout) Error(message string, args ...interface{}) {
	fmt.Printf("[ERROR] "+message, args...)
}
