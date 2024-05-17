package domain

import (
	"context"
	"github.com/long250038728/web/application/order/ddd/repository"
	"github.com/long250038728/web/protoc"
	"github.com/long250038728/web/protoc/order"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/rpc"
)

type OrderDomain struct {
	orderRepository *repository.OrderRepository
}

func NewOrderDomain(userRepository *repository.OrderRepository) *OrderDomain {
	return &OrderDomain{
		orderRepository: userRepository,
	}
}

func (d *OrderDomain) OrderDetail(ctx context.Context, request *order.OrderDetailRequest) (*order.OrderDetailResponse, error) {
	util := app.NewUtil()

	// 创建consul客户端
	conn, err := rpc.NewClient(rpc.Register(protoc.UserService, util.Register())).Dial(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()

	// grpc获取数据
	resp, err := user.NewUserClient(conn).SayHello(ctx, &user.RequestHello{Name: "long"})
	if err != nil {
		return nil, err
	}
	return &order.OrderDetailResponse{Id: request.Id, OrderSn: resp.Str}, nil
}
