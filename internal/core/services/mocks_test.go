package services_test

import (
	"github.com/anderson-oki/awdl0-disabler/internal/core/domain"
	"time"
)

type MockNetworkPort struct {
	CheckFunc   func(name string) (domain.Status, error)
	DisableFunc func(name string) error
	EnableFunc  func(name string) error
}

func (m *MockNetworkPort) CheckInterface(name string) (domain.Status, error) {
	return m.CheckFunc(name)
}
func (m *MockNetworkPort) DisableInterface(name string) error {
	if m.DisableFunc != nil {
		return m.DisableFunc(name)
	}
	return nil
}
func (m *MockNetworkPort) EnableInterface(name string) error {
	if m.EnableFunc != nil {
		return m.EnableFunc(name)
	}
	return nil
}

type MockLoggerPort struct {
	LogFunc func(event domain.Event) error
}

func (m *MockLoggerPort) Log(event domain.Event) error {
	if m.LogFunc != nil {
		return m.LogFunc(event)
	}
	return nil
}

type MockEventRepo struct {
	AddFunc       func(event domain.Event)
	GetRecentFunc func(duration time.Duration) []domain.Event
}

func (m *MockEventRepo) Add(event domain.Event) {
	if m.AddFunc != nil {
		m.AddFunc(event)
	}
}
func (m *MockEventRepo) GetRecent(duration time.Duration) []domain.Event {
	if m.GetRecentFunc != nil {
		return m.GetRecentFunc(duration)
	}
	return nil
}
