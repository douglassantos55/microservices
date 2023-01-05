package pkg

import "github.com/go-kit/kit/log"

type logging struct {
	next   Service
	logger log.Logger
}

func LoggingMiddleware(svc Service, logger log.Logger) Service {
	return &logging{svc, logger}
}

func (l *logging) Create(data Customer) (customer *Customer, err error) {
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
