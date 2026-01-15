package ports

import (
	"time"

	"github.com/anderson-oki/awdl0-disabler/internal/core/domain"
)

// NetworkPort handles interaction with the operating system's network interfaces
type NetworkPort interface {
	CheckInterface(name string) (domain.Status, error)
	DisableInterface(name string) error
	EnableInterface(name string) error
}

// LoggerPort handles persistence of logs
type LoggerPort interface {
	Log(event domain.Event) error
}

// ConfigPort handles loading and saving configuration
type ConfigPort interface {
	Load() (*domain.Config, error)
	Save(config *domain.Config) error
}

// EventRepository (Optional/In-Memory) handles temporary storage for the graph
// We might just keep this in the service or model, but a port is cleaner for DDD
type EventRepository interface {
	Add(event domain.Event)
	GetRecent(duration time.Duration) []domain.Event
}

// SystemPort handles system-level checks and operations
type SystemPort interface {
	HasElevatedPrivileges() bool
}
