package filesystem

import (
	"os"
	"testing"
	"time"

	"github.com/anderson-oki/awdl0-disabler/internal/core/domain"
)

func TestFileLoggerAdapter_ReadEvents(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logger_test")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(tmpDir)

	adapter := NewFileLoggerAdapter(tmpDir)
	now := time.Now()

	events := []domain.Event{
		{
			Timestamp: now.Add(-2 * time.Hour),
			Type:      domain.EventCheck,
			Message:   "First event",
		},
		{
			Timestamp: now.Add(-1 * time.Hour),
			Type:      domain.EventDisable,
			Message:   "Second event with: colons",
		},
	}

	for _, e := range events {
		if err := adapter.Log(e); err != nil {
			t.Fatalf("Failed to log event: %v", err)
		}
	}

	readEvents, err := adapter.ReadEvents(now)
	if err != nil {
		t.Fatalf("Failed to read events: %v", err)
	}

	if len(readEvents) != len(events) {
		t.Errorf("Expected %d events, got %d", len(events), len(readEvents))
	}

	for i, re := range readEvents {
		expectedTimeStr := events[i].Timestamp.Format("15:04:05")
		actualTimeStr := re.Timestamp.Format("15:04:05")

		if expectedTimeStr != actualTimeStr {
			t.Errorf("Event %d: expected time %s, got %s", i, expectedTimeStr, actualTimeStr)
		}

		if re.Type != events[i].Type {
			t.Errorf("Event %d: expected type %s, got %s", i, events[i].Type, re.Type)
		}

		if re.Message != events[i].Message {
			t.Errorf("Event %d: expected message %q, got %q", i, events[i].Message, re.Message)
		}
	}
}
