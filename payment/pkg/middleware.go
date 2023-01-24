package pkg

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

func getTypeMiddleware(svc Service) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			res, err := next(ctx, r)
			if err != nil {
				return nil, err
			}

			if conditions, ok := res.([]*Condition); ok {
				for _, condition := range conditions {
					paymentType, err := svc.GetPaymentType(condition.PaymentTypeID)
					if err == nil {
						condition.PaymentType = paymentType
					}
				}
				return conditions, nil
			}

			if condition, ok := res.(*Condition); ok {
				paymentType, err := svc.GetPaymentType(condition.PaymentTypeID)
				if err == nil {
					condition.PaymentType = paymentType
				}
				return condition, nil
			}

			return res, nil
		}
	}
}

type loggingService struct {
	next   Service
	logger log.Logger
}

func NewLoggingService(next Service, logger log.Logger) *loggingService {
	return &loggingService{next, logger}
}

func (l *loggingService) GetPaymentMethod(id string) (method *Method, err error) {
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

func (l *loggingService) ListPaymentMethods() (methods []*Method, err error) {
	defer func() {
		l.logger.Log(
			"method", "ListPaymentMethods",
			"methods", methods,
			"err", err,
		)
	}()
	return l.next.ListPaymentMethods()
}

func (l *loggingService) CreatePaymentMethod(data Method) (method *Method, err error) {
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

func (l *loggingService) UpdatePaymentMethod(id string, data Method) (method *Method, err error) {
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

func (l *loggingService) CreatePaymentType(data Type) (paymentType *Type, err error) {
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

func (l *loggingService) ListPaymentTypes() (types []*Type, err error) {
	defer func() {
		l.logger.Log(
			"method", "ListPaymentTypes",
			"types", types,
			"err", err,
		)
	}()
	return l.next.ListPaymentTypes()
}

func (l *loggingService) UpdatePaymentType(id string, data Type) (paymentType *Type, err error) {
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

func (l *loggingService) DeletePaymentType(id string) (err error) {
	defer func() {
		l.logger.Log(
			"method", "DeletePaymentType",
			"id", id,
			"err", err,
		)
	}()
	return l.next.DeletePaymentType(id)
}

func (l *loggingService) GetPaymentType(id string) (paymentType *Type, err error) {
	defer func() {
		l.logger.Log(
			"method", "GetPaymentType",
			"id", id,
			"paymentType", paymentType,
			"err", err,
		)
	}()
	return l.next.GetPaymentType(id)
}

func (l *loggingService) CreatePaymentCondition(data Condition) (condition *Condition, err error) {
	defer func() {
		l.logger.Log(
			"method", "CreatePaymentCondition",
			"data", data,
			"condition", condition,
			"err", err,
		)
	}()
	return l.next.CreatePaymentCondition(data)
}

func (l *loggingService) ListPaymentConditions() (conditions []*Condition, err error) {
	defer func() {
		l.logger.Log(
			"method", "ListPaymentConditions",
			"conditions", conditions,
			"err", err,
		)
	}()
	return l.next.ListPaymentConditions()
}

func (l *loggingService) UpdatePaymentCondition(id string, data Condition) (condition *Condition, err error) {
	defer func() {
		l.logger.Log(
			"method", "UpdatePaymentCondition",
			"id", id,
			"data", data,
			"condition", condition,
			"err", err,
		)
	}()
	return l.next.UpdatePaymentCondition(id, data)
}

func (l *loggingService) DeletePaymentCondition(id string) (err error) {
	defer func() {
		l.logger.Log(
			"method", "DeletePaymentCondition",
			"id", id,
			"err", err,
		)
	}()
	return l.next.DeletePaymentCondition(id)
}

func (l *loggingService) GetPaymentCondition(id string) (condition *Condition, err error) {
	defer func() {
		l.logger.Log(
			"method", "GetPaymentCondition",
			"id", id,
			"condition", condition,
			"err", err,
		)
	}()
	return l.next.GetPaymentCondition(id)
}

func (l *loggingService) CreateInvoice(data Invoice) (invoice *Invoice, err error) {
	defer func() {
		l.logger.Log(
			"method", "CreateInvoice",
			"data", data,
			"invoice", invoice,
			"err", err,
		)
	}()
	return l.next.CreateInvoice(data)
}

func (l *loggingService) ListInvoices(page, perPage int64) (invoices []*Invoice, total int64, err error) {
	defer func() {
		l.logger.Log(
			"method", "ListInvoices",
			"page", page,
			"perPage", perPage,
			"invoices", invoices,
			"total", total,
			"err", err,
		)
	}()
	return l.next.ListInvoices(page, perPage)
}

func (l *loggingService) UpdateInvoice(id string, data Invoice) (invoice *Invoice, err error) {
	defer func() {
		l.logger.Log(
			"method", "UpdateInvoice",
			"id", id,
			"data", data,
			"invoice", invoice,
			"err", err,
		)
	}()
	return l.next.UpdateInvoice(id, data)
}

func (l *loggingService) DeleteInvoice(id string) (err error) {
	defer func() {
		l.logger.Log(
			"method", "DeleteInvoice",
			"id", id,
			"err", err,
		)
	}()
	return l.next.DeleteInvoice(id)
}
