package pkg

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingService struct {
	next        Service
	reqDuration metrics.Histogram
	reqCounter  metrics.Counter
}

func NewInstrumentingService(next Service, counter metrics.Counter, duration metrics.Histogram) Service {
	return &instrumentingService{next, duration, counter}
}

func (s *instrumentingService) CreateRent(data Rent) (_ *Rent, err error) {
	defer func(begin time.Time) {
		s.reqCounter.With("method", "CreateRent", "error", fmt.Sprint(err != nil)).Add(1)
		s.reqDuration.With("method", "CreateRent").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.CreateRent(data)
}

func (s *instrumentingService) ListRents(page, perPage int64) (_ []*Rent, total int64, err error) {
	defer func(begin time.Time) {
		s.reqCounter.With("method", "ListRents", "error", fmt.Sprint(err != nil)).Add(1)
		s.reqDuration.With("method", "ListRents").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.ListRents(page, perPage)
}

func (s *instrumentingService) UpdateRent(id string, data Rent) (_ *Rent, err error) {
	defer func(begin time.Time) {
		s.reqCounter.With("method", "UpdateRent", "error", fmt.Sprint(err != nil)).Add(1)
		s.reqDuration.With("method", "UpdateRent").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.UpdateRent(id, data)
}

func (s *instrumentingService) DeleteRent(id string) (err error) {
	defer func(begin time.Time) {
		s.reqCounter.With("method", "DeleteRent", "error", fmt.Sprint(err != nil)).Add(1)
		s.reqDuration.With("method", "DeleteRent").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.DeleteRent(id)
}

func (s *instrumentingService) GetRent(id string) (_ *Rent, err error) {
	defer func(begin time.Time) {
		s.reqCounter.With("method", "GetRent", "error", fmt.Sprint(err != nil)).Add(1)
		s.reqDuration.With("method", "GetRent").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GetRent(id)
}
