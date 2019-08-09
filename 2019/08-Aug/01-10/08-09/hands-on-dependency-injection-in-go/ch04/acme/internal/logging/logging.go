package logging

import "fmt"

var L = &LoggerStdout{}

type LoggerStdout struct{}

func (l LoggerStdout) Debug(message string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+message, args...)
}

func (l LoggerStdout) Info(message string, args ...interface{}) {
	fmt.Printf("[INFO] "+message, args...)
}

func (l LoggerStdout) Warn(message string, args ...interface{}) {
	fmt.Printf("[WARN] "+message, args...)
}

func (l LoggerStdout) Error(message string, args ...interface{}) {
	fmt.Printf("[ERROR] "+message, args...)
}
