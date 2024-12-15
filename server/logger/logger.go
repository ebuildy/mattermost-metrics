package logger

import "fmt"

// Logger define a simple logger, implemented by pluginAPIClient.Log
type Logger interface {
	Debug(message string, keyValuePairs ...interface{})
	Info(message string, keyValuePairs ...interface{})
	Warn(message string, keyValuePairs ...interface{})
	Error(message string, keyValuePairs ...interface{})
}

// fakeLogger is used for tests
type fakeLogger struct{}

func NewFakeLogger() Logger {
	return &fakeLogger{}
}

func (f fakeLogger) Debug(message string, keyValuePairs ...interface{}) {
	fmt.Printf("[DEBUG] %s\n", message)
}

func (f fakeLogger) Info(message string, keyValuePairs ...interface{}) {
	fmt.Printf("[INFO] %s\n", message)
}

func (f fakeLogger) Warn(message string, keyValuePairs ...interface{}) {
	fmt.Printf("[WARN] %s\n", message)
}

func (f fakeLogger) Error(message string, keyValuePairs ...interface{}) {
	fmt.Printf("[ERROR] %s\n", message)
}
