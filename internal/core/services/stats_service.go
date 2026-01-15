package services

import (
	"math"
	"time"

	"github.com/anderson-oki/awdl0-disabler/internal/core/domain"
	"github.com/anderson-oki/awdl0-disabler/internal/core/ports"
)

type StatsService struct {
	repo ports.EventRepository
}

func NewStatsService(r ports.EventRepository) *StatsService {
	return &StatsService{repo: r}
}

// GetHistogram generates a fixed number of buckets for the given duration
func (s *StatsService) GetHistogram(duration time.Duration, numBuckets int) []domain.Bucket {
	return s.GetHistogramAt(time.Now(), duration, numBuckets)
}

// GetHistogramAt generates histogram relative to a specific time (useful for testing)
func (s *StatsService) GetHistogramAt(now time.Time, duration time.Duration, numBuckets int) []domain.Bucket {
	events := s.repo.GetRecent(duration)

	buckets := make([]domain.Bucket, numBuckets)

	startTime := now.Add(-duration)
	slotDuration := duration / time.Duration(numBuckets)

	for i := 0; i < numBuckets; i++ {
		buckets[i] = domain.Bucket{
			Count: 0,
			Label: "",
		}
	}

	for _, evt := range events {
		if evt.Timestamp.Before(startTime) || evt.Timestamp.After(now) {
			continue
		}

		offset := evt.Timestamp.Sub(startTime)
		index := int(math.Floor(float64(offset) / float64(slotDuration)))

		if index >= 0 && index < numBuckets {
			buckets[index].Count++
		}
	}

	return buckets
}

// GetRecentEvents returns raw events from the repository for the given duration
func (s *StatsService) GetRecentEvents(duration time.Duration) []domain.Event {
	return s.repo.GetRecent(duration)
}
