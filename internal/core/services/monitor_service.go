package services

import (
	"time"

	"github.com/anderson-oki/awdl0-disabler/internal/core/domain"
	"github.com/anderson-oki/awdl0-disabler/internal/core/ports"
)

type MonitorService struct {
	network ports.NetworkPort
	logger  ports.LoggerPort
	repo    ports.EventRepository
	config  *domain.Config
}

func NewMonitorService(n ports.NetworkPort, l ports.LoggerPort, r ports.EventRepository, c *domain.Config) *MonitorService {
	return &MonitorService{
		network: n,
		logger:  l,
		repo:    r,
		config:  c,
	}
}

func (s *MonitorService) Tick() (*domain.Event, error) {
	status, err := s.network.CheckInterface("awdl0")
	if err != nil {
		return nil, err
	}

	if status != domain.StatusUp {
		return nil, nil
	}

	if err := s.network.DisableInterface("awdl0"); err != nil {
		return nil, err
	}

	evt := domain.Event{
		Timestamp: time.Now(),
		Type:      domain.EventDisable,
		Message:   "awdl0 detected UP. Disabling...",
	}

	_ = s.logger.Log(evt)

	s.repo.Add(evt)

	return &evt, nil
}

func (s *MonitorService) GetConfig() *domain.Config {
	return s.config
}

func (s *MonitorService) Restore() error {
	return s.network.EnableInterface("awdl0")
}
