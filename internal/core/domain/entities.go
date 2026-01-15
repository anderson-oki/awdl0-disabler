package domain

import "time"

// Status represents the state of a network interface
type Status string

const (
	StatusUp      Status = "UP"
	StatusDown    Status = "DOWN"
	StatusUnknown Status = "Unknown"
)

// EventType distinguishes between different system actions
type EventType string

const (
	EventDisable EventType = "Disable"
	EventEnable  EventType = "Enable"
	EventCheck   EventType = "Check"
)

// Event represents a system action occurred at a specific time
type Event struct {
	Timestamp time.Time
	Type      EventType
	Message   string
}

// Config represents the user configuration
type Config struct {
	PollingInterval time.Duration `json:"polling_interval"`
}

const (
	MinPollingInterval = 500 * time.Millisecond
	MaxPollingInterval = 60 * time.Second
)

// Clamp ensures the configuration values are within valid ranges
func (c *Config) Clamp() {
	if c.PollingInterval < MinPollingInterval {
		c.PollingInterval = MinPollingInterval
	}

	if c.PollingInterval > MaxPollingInterval {
		c.PollingInterval = MaxPollingInterval
	}
}

// Bucket represents a time slot in the histogram
type Bucket struct {
	Label string // e.g., "10:05"
	Count int
}
