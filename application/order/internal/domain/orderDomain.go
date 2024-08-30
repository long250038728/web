package domain

import (
	"context"
	"github.com/long250038728/web/application/order/internal/repository"
	"github.com/long250038728/web/protoc"
	"github.com/long250038728/web/protoc/order"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/rpc"
)

type OrderDomain struct {
	orderRepository *repository.OrderRepository
}

func NewDomain(userRepository *repository.OrderRepository) *OrderDomain {
	return &OrderDomain{
		orderRepository: userRepository,
	}
}

func (d *OrderDomain) OrderDetail(ctx context.Context, request *order.OrderDetailRequest) (*order.OrderDetailResponse, error) {
	register, err := app.NewUtil().Register()
	if err != nil {
		return nil, err
	}

	// 创建consul客户端
	//conn, err := rpc.NewClient(rpc.LocalIP("172.40.0.3", 9001)).Dial(ctx)
	conn, err := rpc.NewClient(rpc.Register(protoc.UserService, register)).Dial(ctx)
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
