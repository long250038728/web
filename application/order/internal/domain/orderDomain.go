package domain

import (
	"context"
	"github.com/long250038728/web/application/order/internal/repository"
	"github.com/long250038728/web/protoc"
	"github.com/long250038728/web/protoc/order"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/server/rpc"
)

type Domain struct {
	repository *repository.Repository
}

func NewDomain(repository *repository.Repository) *Domain {
	return &Domain{
		repository: repository,
	}
}

func (d *Domain) OrderDetail(ctx context.Context, request *order.OrderDetailRequest) (*order.OrderDetailResponse, error) {
	// 创建rpc客户端
	conn, err := rpc.NewClient().Dial(ctx, protoc.UserService)
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
