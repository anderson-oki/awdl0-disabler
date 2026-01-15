package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anderson-oki/awdl0-disabler/internal/core/domain"
)

type FileLoggerAdapter struct {
	LogDir string
}

func NewFileLoggerAdapter(dir string) *FileLoggerAdapter {
	_ = os.MkdirAll(dir, 0755)
	return &FileLoggerAdapter{LogDir: dir}
}

func (l *FileLoggerAdapter) Log(event domain.Event) error {
	filename := event.Timestamp.Format("2006-01-02") + ".log"
	path := filepath.Join(l.LogDir, filename)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	line := fmt.Sprintf("[%s] %s: %s\n",
		event.Timestamp.Format("15:04:05"),
		event.Type,
		event.Message,
	)

	_, err = file.WriteString(line)
	return err
}

// ReadEvents reads log events for a specific date
func (l *FileLoggerAdapter) ReadEvents(date time.Time) ([]domain.Event, error) {
	filename := date.Format("2006-01-02") + ".log"
	path := filepath.Join(l.LogDir, filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []domain.Event{}, nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var events []domain.Event
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, "] ", 2)
		if len(parts) != 2 {
			continue
		}

		timeStr := strings.TrimPrefix(parts[0], "[")
		rest := parts[1]

		parsedTime, err := time.Parse("15:04:05", timeStr)
		if err != nil {
			continue
		}

		fullTime := time.Date(
			date.Year(), date.Month(), date.Day(),
			parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0,
			date.Location(),
		)

		msgParts := strings.SplitN(rest, ": ", 2)
		if len(msgParts) != 2 {
			continue
		}
		eventType := domain.EventType(msgParts[0])
		message := msgParts[1]

		events = append(events, domain.Event{
			Timestamp: fullTime,
			Type:      eventType,
			Message:   message,
		})
	}

	return events, nil
}
