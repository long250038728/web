package service

import (
	"context"
	"github.com/long250038728/web/application/order/internal/domain"
	"github.com/long250038728/web/protoc/order"
	"github.com/long250038728/web/tool/server/rpc/server"
)

type Order struct {
	order.UnimplementedOrderServer
	server.GrpcHealth
	domain *domain.Order
}

type OrderServerOpt func(s *Order)

func SetDomain(domain *domain.Order) OrderServerOpt {
	return func(s *Order) {
		s.domain = domain
	}
}

func NewService(opts ...OrderServerOpt) *Order {
	s := &Order{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Order) OrderDetail(ctx context.Context, request *order.OrderDetailRequest) (*order.OrderDetailResponse, error) {
	return s.domain.OrderDetail(ctx, request)
}
