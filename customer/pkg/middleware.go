package pkg

import "github.com/go-kit/kit/log"

type loggingService struct {
	next   Service
	logger log.Logger
}

func NewLoggingService(svc Service, logger log.Logger) Service {
	return &loggingService{svc, logger}
}

func (l *loggingService) List(page, perPage int64) (result *ListResult, err error) {
	defer func() {
		l.logger.Log(
			"method", "List",
			"page", page,
			"perPage", perPage,
			"result", result,
			"err", err,
		)
	}()
	return l.next.List(page, perPage)
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

func (l *loggingService) Update(id string, data Customer) (customer *Customer, err error) {
	defer func() {
		l.logger.Log(
			"method", "Update",
			"id", id,
			"input", data,
			"output", customer,
			"err", err,
		)
	}()
	return l.next.Update(id, data)
}
