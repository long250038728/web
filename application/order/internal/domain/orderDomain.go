package domain

import (
	"context"
	"github.com/long250038728/web/application/order/internal/repository"
	"github.com/long250038728/web/protoc/order"
)

type OrderDomain struct {
	repository *repository.OrderRepository
}

func NewOrderDomain(repository *repository.OrderRepository) *OrderDomain {
	return &OrderDomain{repository: repository}
}

func (d *OrderDomain) OrderDetail(ctx context.Context, request *order.OrderDetailRequest) (*order.OrderDetailResponse, error) {
	// 创建rpc客户端
	orderSn, err := d.repository.OrderDetail(ctx, request)
	if err != nil {
		return nil, err
	}
	return &order.OrderDetailResponse{Id: request.Id, OrderSn: orderSn}, nil
}
