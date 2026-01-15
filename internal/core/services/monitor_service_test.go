package services_test

import (
	"github.com/anderson-oki/awdl0-disabler/internal/core/domain"
	"github.com/anderson-oki/awdl0-disabler/internal/core/services"
	"testing"
	"time"
)

func TestMonitorService_Tick_DisablesWhenUp(t *testing.T) {
	network := &MockNetworkPort{
		CheckFunc: func(name string) (domain.Status, error) {
			return domain.StatusUp, nil
		},
		DisableFunc: func(name string) error {
			return nil
		},
	}

	disabledCalled := false
	network.DisableFunc = func(name string) error {
		disabledCalled = true
		return nil
	}

	logger := &MockLoggerPort{}
	repo := &MockEventRepo{}

	config := &domain.Config{PollingInterval: time.Second}
	service := services.NewMonitorService(network, logger, repo, config)

	event, err := service.Tick()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !disabledCalled {
		t.Error("Expected DisableInterface to be called, but it wasn't")
	}
	if event == nil || event.Type != domain.EventDisable {
		t.Error("Expected Disable event to be returned")
	}
}

func TestMonitorService_Tick_DoNothingWhenDown(t *testing.T) {
	network := &MockNetworkPort{
		CheckFunc: func(name string) (domain.Status, error) {
			return domain.StatusDown, nil
		},
	}

	disableCalled := false
	network.DisableFunc = func(name string) error {
		disableCalled = true
		return nil
	}

	service := services.NewMonitorService(network, &MockLoggerPort{}, &MockEventRepo{}, &domain.Config{})

	event, err := service.Tick()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if disableCalled {
		t.Error("Expected DisableInterface NOT to be called when status is DOWN")
	}
	if event != nil {
		t.Errorf("Expected nil event when DOWN, got %v", event)
	}
}
