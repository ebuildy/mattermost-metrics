package fake

import (
	"fmt"
	"github.com/ebuildy/mattermost-plugin-minotor/server/internal/core/ports"
)

// fakeLogger is used for tests
type fakeLogger struct{}

func NewFakeLogger() ports.Logger {
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
