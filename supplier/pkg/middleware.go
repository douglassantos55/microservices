package pkg

import "github.com/go-kit/log"

type loggingService struct {
	next   Service
	logger log.Logger
}

func NewLoggingService(next Service, logger log.Logger) Service {
	return &loggingService{next, logger}
}

func (l *loggingService) Get(id string) (supplier *Supplier, err error) {
	defer func() {
		l.logger.Log(
			"method", "Get",
			"id", id,
			"output", supplier,
			"err", err,
		)
	}()
	return l.next.Get(id)
}

func (l *loggingService) List(page, perPage int64) (suppliers []*Supplier, total int64, err error) {
	defer func() {
		l.logger.Log(
			"method", "List",
			"page", page,
			"perPage", perPage,
			"suppliers", suppliers,
			"total", total,
			"err", err,
		)
	}()
	return l.next.List(page, perPage)
}

func (l *loggingService) Create(data Supplier) (supplier *Supplier, err error) {
	defer func() {
		l.logger.Log(
			"method", "Create",
			"input", data,
			"output", supplier,
			"err", err,
		)
	}()
	return l.next.Create(data)
}

func (l *loggingService) Update(id string, data Supplier) (supplier *Supplier, err error) {
	defer func() {
		l.logger.Log(
			"method", "Update",
			"id", id,
			"input", data,
			"output", supplier,
			"err", err,
		)
	}()
	return l.next.Update(id, data)
}

func (l *loggingService) Delete(id string) (err error) {
	defer func() {
		l.logger.Log(
			"method", "Delete",
			"id", id,
			"err", err,
		)
	}()
	return l.next.Delete(id)
}
