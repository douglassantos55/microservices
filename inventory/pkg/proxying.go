package pkg

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"reconcip.com.br/microservices/inventory/proto"
)

func FetchSupplierEndpoints(endpoints Set, cc *grpc.ClientConn) Set {
	fetchSupplier := fetchSupplierMiddleware(cc)
	fetchSuppliers := fetchSuppliersMiddleware(cc)

	return Set{
		Get:    fetchSupplier(endpoints.Get),
		Create: fetchSupplier(endpoints.Create),
		List:   fetchSuppliers(endpoints.List),
		Update: fetchSupplier(endpoints.Update),
		Delete: endpoints.Delete,
	}
}

func fetchSupplierMiddleware(cc *grpc.ClientConn) endpoint.Middleware {
	fetchSupplier := makeFetchSupplierEndpoint(cc)

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			response, err := next(ctx, r)
			if err != nil {
				return nil, err
			}

			equipment := response.(*Equipment)
			supplier, err := fetchSupplier(ctx, equipment.SupplierID)
			if err == nil {
				equipment.Supplier = supplier.(*Supplier)
			}

			return equipment, nil
		}
	}
}

func fetchSuppliersMiddleware(cc *grpc.ClientConn) endpoint.Middleware {
	fetchSupplier := makeFetchSupplierEndpoint(cc)

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			response, err := next(ctx, r)
			if err != nil {
				return nil, err
			}

			result := response.(ListResult)
			for _, item := range result.Items {
				equipment := item.(*Equipment)
				supplier, err := fetchSupplier(ctx, equipment.SupplierID)
				if err == nil {
					equipment.Supplier = supplier.(*Supplier)
				}
			}

			return result, nil
		}
	}
}

func makeFetchSupplierEndpoint(cc *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		cc,
		"proto.SupplierService",
		"Get",
		encodeRequest,
		decodeResponse,
		&proto.Supplier{},
	).Endpoint()
}

func encodeRequest(ctx context.Context, r any) (any, error) {
	supplierID, ok := r.(string)
	if !ok || supplierID == "" {
		return nil, errors.New("invalid supplier id")
	}
	return &proto.GetRequest{SupplierID: supplierID}, nil
}

func decodeResponse(ctx context.Context, r any) (any, error) {
	reply := r.(*proto.Supplier)

	return &Supplier{
		ID:         reply.Id,
		SocialName: reply.SocialName,
		LegalName:  reply.LegalName,
		Email:      reply.Email,
		Website:    reply.Email,
		Cnpj:       reply.Cnpj,
		InscEst:    reply.InscEst,
		Phone:      reply.Phone,
		Address: Address{
			Street:       reply.Address.Street,
			Number:       reply.Address.Number,
			Complement:   reply.Address.Complement,
			Neighborhood: reply.Address.Neighborhood,
			City:         reply.Address.City,
			State:        reply.Address.State,
			Postcode:     reply.Address.Postcode,
		},
	}, nil
}
