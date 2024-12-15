package components

import (
	"context"
	"sync"
	"time"

	"service-weaver-app/models"
)

type AnalyticsImpl struct {
	sync.Mutex
	metrics []models.Metric
}

func NewAnalytics() *AnalyticsImpl {
	return &AnalyticsImpl{metrics: []models.Metric{}}
}

func (a *AnalyticsImpl) TrackMetric(ctx context.Context, name string, value float64) error {
	a.Lock()
	defer a.Unlock()
	a.metrics = append(a.metrics, models.Metric{
		Name:  name,
		Value: value,
		Time:  time.Now().Unix(),
	})
	return nil
}
