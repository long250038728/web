package rpc

import (
	"context"
	"fmt"
	"github.com/long250038728/web/tool/server/rpc/client"
	"github.com/long250038728/web/tool/server/rpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/resolver"
	"sync"
	"time"
)

// circuitList 设置需要熔断的接口
var circuitList []string

// SetCircuits 设置需要熔断的接口(不需要append是因为列表是确定的不会运行中添加，所以无需动态添加)
func SetCircuits(circuits []string) {
	circuitList = circuits
}

// =================================================================================================

type Conn struct {
	*grpc.ClientConn
}

func (c *Conn) Close() error {
	return c.ClientConn.Close()
}

//=================================================================================================

// Client 客户端
type Client struct {
	target client.Target
	once   sync.Once
}

// Opt grpc客户端opt
type Opt func(client *Client)

var clientParameters = keepalive.ClientParameters{
	Time:                10 * time.Second, // 如果没有活动，每10秒发送一次心跳
	Timeout:             time.Second,      // 等待1秒钟的心跳响应，若未收到则认为连接已断开
	PermitWithoutStream: true,             // 即使没有活动的数据流，也发送心跳
}

// NewClient 构造函数
func NewClient(target client.Target, opts ...Opt) *Client {
	c := &Client{
		target: target,
	}
	for _, opt := range opts {
		opt(c)
	}

	c.once.Do(func() {
		resolver.Register(&client.MyResolversBuild{})
	})
	return c
}

func (c *Client) Dial(ctx context.Context, serverName string) (conn *Conn, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("grpc client dial err: %v", r)
		}
	}()

	address, dialOptions, err := c.target.Target(ctx, serverName)
	if err != nil {
		return nil, err
	}

	//获取target 信息
	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(server.ServerCircuitInterceptor(circuitList)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(clientParameters),
	}
	opts = append(opts, dialOptions...)

	//创建socket 连接
	clientConn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, err
	}
	return &Conn{clientConn}, err
}
