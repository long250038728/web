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

	////创建consul客户端
	//register, err := consul.NewConsulRegister("192.168.0.89:8500")
	//if err != nil {
	//	return nil, err
	//}

	//tracing.OpentracingGlobalTracer("http://link.zhubaoe.cn:14268/api/traces", "Aclient")

	//创建consul客户端
	grpcClient := rpc.NewClient(
		//rpc.Register("AUser", register),
		rpc.LocalIP("192.168.1.20", 8092),
	)
	conn, err := grpcClient.Dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	//请求数据
	rpcClient := user.NewUserServerClient(conn)
	return rpcClient.SayHello(ctx, &user.RequestHello{Name: "long"})
}
