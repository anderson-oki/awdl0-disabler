package services_test

import (
	"github.com/anderson-oki/awdl0-disabler/internal/core/domain"
	"github.com/anderson-oki/awdl0-disabler/internal/core/services"
	"testing"
	"time"
)

func TestStatsService_GetHistogram(t *testing.T) {
	now := time.Now()
	repo := &MockEventRepo{
		GetRecentFunc: func(d time.Duration) []domain.Event {
			return []domain.Event{
				{Timestamp: now.Add(-5 * time.Minute)},
				{Timestamp: now.Add(-5 * time.Minute)},
				{Timestamp: now.Add(-55 * time.Minute)},
				{Timestamp: now.Add(-2 * time.Hour)},
			}
		},
	}

	service := services.NewStatsService(repo)

	buckets := service.GetHistogramAt(now, 1*time.Hour, 60)

	if len(buckets) != 60 {
		t.Errorf("Expected 60 buckets, got %d", len(buckets))
	}

	if buckets[55].Count != 2 {
		t.Errorf("Expected 2 events in bucket 55 (approx -5m), got %d", buckets[55].Count)
	}
	if buckets[5].Count != 1 {
		t.Errorf("Expected 1 event in bucket 5 (approx -55m), got %d", buckets[5].Count)
	}
	if buckets[0].Count != 0 {
		t.Errorf("Expected 0 events in bucket 0, got %d", buckets[0].Count)
	}
}
