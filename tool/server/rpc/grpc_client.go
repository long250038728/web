package rpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/register"
	"github.com/long250038728/web/tool/server/rpc/tool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

// Client 客户端
type Client struct {
	serverName   string
	register     register.Register
	svcInstances []*register.ServiceInstance
	balancer     tool.Balancer
}

var clientParameters = keepalive.ClientParameters{
	Time:                10 * time.Second, // 如果没有活动，每10秒发送一次心跳
	Timeout:             time.Second,      // 等待1秒钟的心跳响应，若未收到则认为连接已断开
	PermitWithoutStream: true,             // 即使没有活动的数据流，也发送心跳
}

//=================================================================================================

// ClientOpt grpc客户端opt
type ClientOpt func(client *Client)

// LocalIP 指定IP及Port
func LocalIP(address string, port int) ClientOpt {
	return func(client *Client) {
		client.svcInstances = []*register.ServiceInstance{{Address: address, Port: port}}
	}
}

// Register 指定注册中心
func Register(serverName string, register register.Register) ClientOpt {
	return func(client *Client) {
		client.serverName = serverName
		client.register = register
	}
}

//=================================================================================================

// NewClient 构造函数
func NewClient(opts ...ClientOpt) *Client {
	c := &Client{
		balancer: tool.NewRandBalancer(), //默认随机算法
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Dial 创建conn连接
func (c *Client) Dial(ctx context.Context) (*grpc.ClientConn, error) {
	//如果注册中心，那在注册中心获取列表信息
	if c.serverName != "" && c.register != nil {
		svcInstances, err := c.register.List(ctx, register.GrpcServerName(c.serverName))
		if err != nil {
			return nil, err
		}
		c.svcInstances = svcInstances
	}

	//找不到有任何服务器实例
	if c.svcInstances == nil || len(c.svcInstances) == 0 {
		return nil, errors.New("svcInstances is null")
	}

	// 负载均衡
	svcInstance := c.balancer.Balancer(c.svcInstances)

	//创建socket 连接
	return grpc.Dial(
		fmt.Sprintf("%s:%d", svcInstance.Address, svcInstance.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(clientParameters),
		//rpc.WithResolvers() 服务发现
	)
}
