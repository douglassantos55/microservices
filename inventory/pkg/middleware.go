package pkg

import "github.com/go-kit/log"

type loggingService struct {
	next   Service
	logger log.Logger
}

func NewLoggingService(next Service, logger log.Logger) Service {
	return &loggingService{next, logger}
}

func (l *loggingService) GetEquipment(id string) (equipment *Equipment, err error) {
	defer func() {
		l.logger.Log(
			"method", "GetEquipment",
			"id", id,
			"equipment", equipment,
			"err", err,
		)
	}()
	return l.next.GetEquipment(id)
}

func (l *loggingService) CreateEquipment(data Equipment) (equipment *Equipment, err error) {
	defer func() {
		l.logger.Log(
			"method", "CreateEquipment",
			"input", data,
			"output", equipment,
			"err", err,
		)
	}()
	return l.next.CreateEquipment(data)
}

func (l *loggingService) UpdateEquipment(id string, data Equipment) (equipment *Equipment, err error) {
	defer func() {
		l.logger.Log(
			"method", "UpdateEquipment",
			"id", id,
			"input", data,
			"output", equipment,
			"err", err,
		)
	}()
	return l.next.UpdateEquipment(id, data)
}

func (l *loggingService) DeleteEquipment(id string) (err error) {
	defer func() {
		l.logger.Log(
			"method", "DeleteEquipment",
			"id", id,
			"err", err,
		)
	}()
	return l.next.DeleteEquipment(id)
}

func (l *loggingService) ListEquipment(page, perPage int) (equipment []*Equipment, total int, err error) {
	defer func() {
		l.logger.Log(
			"method", "ListEquipment",
			"page", page,
			"perPage", perPage,
			"items", equipment,
			"total", total,
			"err", err,
		)
	}()
	return l.next.ListEquipment(page, perPage)
}

func (l *loggingService) ReduceStock(id string, qty int64) (err error) {
	defer func() {
		l.logger.Log(
			"method", "ReduceStock",
			"id", id,
			"qty", qty,
			"err", err,
		)
	}()
	return l.next.ReduceStock(id, qty)
}
