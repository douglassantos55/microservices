package pkg

import (
	"context"
	"errors"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"reconcip.com.br/microservices/inventory/proto"
)

func CreateAuthEndpoints(endpoints Set, cc *grpc.ClientConn) Set {
	verify := verifyMiddleware(cc)

	return Set{
		Get:          verify(endpoints.Get),
		List:         verify(endpoints.List),
		Create:       verify(endpoints.Create),
		Update:       verify(endpoints.Update),
		Delete:       verify(endpoints.Delete),
		ReduceStock:  endpoints.ReduceStock,
		RestoreStock: endpoints.RestoreStock,
	}
}

func verifyMiddleware(cc *grpc.ClientConn) endpoint.Middleware {
	verify := createVerifyEndpoint(cc)

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r any) (any, error) {
			if _, err := verify(ctx, r); err != nil {
				return nil, err
			}
			return next(ctx, r)
		}
	}
}

func createVerifyEndpoint(cc *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		cc,
		"proto.Auth",
		"Verify",
		encodeVerifyRequest,
		decodeVerifyResponse,
		&proto.VerifyReply{},
		grpctransport.ClientBefore(jwt.ContextToGRPC()),
	).Endpoint()
}

func encodeVerifyRequest(ctx context.Context, r any) (any, error) {
	return nil, nil
}

func decodeVerifyResponse(ctx context.Context, r any) (any, error) {
	reply := r.(*proto.VerifyReply)
	if reply.Err != nil {
		return nil, NewErrorFromReply(reply.Err)
	}

	var user struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	user.ID = reply.User.Id
	user.Name = reply.User.Name

	return user, nil
}

func FetchSupplierEndpoints(endpoints Set, cc *grpc.ClientConn) Set {
	fetchSupplier := fetchSupplierMiddleware(cc)
	fetchSuppliers := fetchSuppliersMiddleware(cc)

	return Set{
		Get:          fetchSupplier(endpoints.Get),
		Create:       fetchSupplier(endpoints.Create),
		List:         fetchSuppliers(endpoints.List),
		Update:       fetchSupplier(endpoints.Update),
		Delete:       endpoints.Delete,
		ReduceStock:  endpoints.ReduceStock,
		RestoreStock: endpoints.RestoreStock,
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
	return &proto.GetRequest{Id: supplierID}, nil
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
