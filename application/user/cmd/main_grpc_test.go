package main

import (
	"context"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/rpc"
	"testing"
)

func TestGrpcClient(t *testing.T) {
	t.Log(userGrpcClientTest())
}

func userGrpcClientTest() (interface{}, error) {
	ctx := context.Background()

	conf, err := app.NewAppConfig("/Users/linlong/Desktop/web/config")
	if err != nil {
		return nil, err
	}

	util, err := app.NewUtilConfig(conf)
	if err != nil {
		return nil, err
	}

	//创建consul客户端
	conn, err := rpc.NewClient(rpc.Register("kobe-new", util.Register())).Dial(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()

	//请求数据
	rpcClient := user.NewUserClient(conn)
	return rpcClient.SayHello(ctx, &user.RequestHello{Name: "long"}), nil
}
