package main

import (
	"context"
	"github.com/long250038728/web/application/user/protoc"
	"github.com/long250038728/web/tool/server/rpc"
	"testing"
)

func TestGrpcClient(t *testing.T) {
	_, _ = userGrpcClientTest()
}

func userGrpcClientTest() (interface{}, error) {
	ctx := context.Background()

	//创建consul客户端
	grpcClient := rpc.NewClient(
		//rpc.Register("User", register),    //通过服务注册与发现找到 User的实例的IP及Port
		rpc.LocalIP("192.168.1.20", 8092), //直接指定IP及Port
	)
	conn, err := grpcClient.Dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	//请求数据
	rpcClient := user.NewUserClient(conn)
	return rpcClient.SayHello(ctx, &user.RequestHello{Name: "long"})
}
