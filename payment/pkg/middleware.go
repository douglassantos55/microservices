package pkg

import "github.com/go-kit/log"

type loggingService struct {
	next   Service
	logger log.Logger
}

func NewLoggingService(next Service, logger log.Logger) *loggingService {
	return &loggingService{next, logger}
}

func (l *loggingService) GetPaymentMethod(id string) (method *PaymentMethod, err error) {
	defer func() {
		l.logger.Log(
			"method", "GetPaymentMethod",
			"id", id,
			"method", method,
			"err", err,
		)
	}()
	return l.next.GetPaymentMethod(id)
}

func (l *loggingService) ListPaymentMethods() (methods []*PaymentMethod, err error) {
	defer func() {
		l.logger.Log(
			"method", "ListPaymentMethods",
			"methods", methods,
			"err", err,
		)
	}()
	return l.next.ListPaymentMethods()
}

func (l *loggingService) CreatePaymentMethod(data PaymentMethod) (method *PaymentMethod, err error) {
	defer func() {
		l.logger.Log(
			"method", "CreatePaymentMethod",
			"data", data,
			"method", method,
			"err", err,
		)
	}()
	return l.next.CreatePaymentMethod(data)
}

func (l *loggingService) UpdatePaymentMethod(id string, data PaymentMethod) (method *PaymentMethod, err error) {
	defer func() {
		l.logger.Log(
			"method", "UpdatePaymentMethod",
			"id", id,
			"data", data,
			"method", method,
			"err", err,
		)
	}()
	return l.next.UpdatePaymentMethod(id, data)
}

func (l *loggingService) DeletePaymentMethod(id string) (err error) {
	defer func() {
		l.logger.Log(
			"method", "DeletePaymentMethod",
			"id", id,
			"err", err,
		)
	}()
	return l.next.DeletePaymentMethod(id)
}

func (l *loggingService) CreatePaymentType(data PaymentType) (paymentType *PaymentType, err error) {
	defer func() {
		l.logger.Log(
			"method", "CreatePaymentType",
			"data", data,
			"paymentType", paymentType,
			"err", err,
		)
	}()
	return l.next.CreatePaymentType(data)
}

func (l *loggingService) ListPaymentTypes() (types []*PaymentType, err error) {
	defer func() {
		l.logger.Log(
			"method", "ListPaymentTypes",
			"types", types,
			"err", err,
		)
	}()
	return l.next.ListPaymentTypes()
}

func (l *loggingService) UpdatePaymentType(id string, data PaymentType) (paymentType *PaymentType, err error) {
	defer func() {
		l.logger.Log(
			"method", "UpdatePaymentType",
			"id", id,
			"data", data,
			"paymentType", paymentType,
			"err", err,
		)
	}()
	return l.next.UpdatePaymentType(id, data)
}
