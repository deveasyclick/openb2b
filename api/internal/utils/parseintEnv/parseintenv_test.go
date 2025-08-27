package parseintenv

import (
	"os"
	"testing"

	"github.com/deveasyclick/openb2b/pkg/interfaces"
)

// mock logger
type mockLogger struct {
	warnCalled bool
	lastMsg    string
}

func (m *mockLogger) Warn(msg string, keysAndValues ...interface{}) {
	m.warnCalled = true
	m.lastMsg = msg
}

// implement other methods so it satisfies interfaces.Logger if needed
func (m *mockLogger) Info(string, ...interface{})  {}
func (m *mockLogger) Debug(string, ...interface{}) {}
func (m *mockLogger) Error(string, ...interface{}) {}
func (m *mockLogger) Fatal(string, ...interface{}) {}
func (m *mockLogger) WithValues(keysAndValues ...interface{}) interfaces.Logger {
	return m
}

func TestParseIntEnv_DefaultWhenUnset(t *testing.T) {
	logger := &mockLogger{}
	key := "TEST_INT_ENV"
	os.Unsetenv(key)

	got := ParseIntEnv(key, 10, logger)
	want := 10

	if got != want {
		t.Errorf("expected %d, got %d", want, got)
	}
	if logger.warnCalled {
		t.Errorf("logger.Warn should NOT be called when unset")
	}
}

func TestParseIntEnv_ValidInt(t *testing.T) {
	logger := &mockLogger{}
	key := "TEST_INT_ENV"
	os.Setenv(key, "42")
	defer os.Unsetenv(key)

	got := ParseIntEnv(key, 10, logger)
	want := 42

	if got != want {
		t.Errorf("expected %d, got %d", want, got)
	}
	if logger.warnCalled {
		t.Errorf("logger.Warn should NOT be called on valid int")
	}
}

func TestParseIntEnv_InvalidInt(t *testing.T) {
	logger := &mockLogger{}
	key := "TEST_INT_ENV"
	os.Setenv(key, "not-an-int")
	defer os.Unsetenv(key)

	got := ParseIntEnv(key, 10, logger)
	want := 10

	if got != want {
		t.Errorf("expected %d, got %d", want, got)
	}
	if !logger.warnCalled {
		t.Errorf("logger.Warn SHOULD be called on invalid int")
	}
	if logger.lastMsg == "" || logger.lastMsg[:7] != "Invalid" {
		t.Errorf("logger.Warn message not set correctly: %v", logger.lastMsg)
	}
}
