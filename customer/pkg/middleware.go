package pkg

import "github.com/go-kit/kit/log"

type loggingService struct {
	next   Service
	logger log.Logger
}

func NewLoggingService(svc Service, logger log.Logger) Service {
	return &loggingService{svc, logger}
}

func (l *loggingService) Create(data Customer) (customer *Customer, err error) {
	defer func() {
		l.logger.Log(
			"method", "Create",
			"input", data,
			"output", customer,
			"err", err,
		)
	}()
	return l.next.Create(data)
}
