package service

import (
	"context"
	"github.com/long250038728/web/application/order/ddd/domain"
	"github.com/long250038728/web/protoc/order"
	"github.com/long250038728/web/tool/server/rpc"
)

type OrderService struct {
	order.UnimplementedOrderServer
	rpc.GrpcHealth
	domain *domain.OrderDomain
}

type OrderServerOpt func(s *OrderService)

func SetOrderDomain(domain *domain.OrderDomain) OrderServerOpt {
	return func(s *OrderService) {
		s.domain = domain
	}
}

func NewService(opts ...OrderServerOpt) *OrderService {
	s := &OrderService{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *OrderService) OrderDetail(ctx context.Context, request *order.OrderDetailRequest) (*order.OrderDetailResponse, error) {
	return s.domain.OrderDetail(ctx, request)
}
