package rpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/const_g"
	"github.com/long250038728/web/tool/register"
	"github.com/long250038728/web/tool/server/rpc/tool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/resolver"
	"strings"
	"sync"
	"time"
)

// Client 客户端
type Client struct {
	svcInstances []*register.ServiceInstance
	balancer     tool.Balancer
	once         sync.Once

	rpcType string
	rpcPort map[string]int
	r       register.Register
	ip      string
}

var clientParameters = keepalive.ClientParameters{
	Time:                10 * time.Second, // 如果没有活动，每10秒发送一次心跳
	Timeout:             time.Second,      // 等待1秒钟的心跳响应，若未收到则认为连接已断开
	PermitWithoutStream: true,             // 即使没有活动的数据流，也发送心跳
}

// =================================================================================================
// circuitList 设置需要熔断的接口
var circuitList []string

// SetCircuits 设置需要熔断的接口(不需要append是因为列表是确定的不会运行中添加，所以无需动态添加)
func SetCircuits(circuits []string) {
	circuitList = circuits
}

//=================================================================================================

// ClientOpt grpc客户端opt
type ClientOpt func(client *Client)

// Balancer 设置负载均衡
func Balancer(balancer tool.Balancer) ClientOpt {
	return func(client *Client) {
		client.balancer = balancer
	}
}

//=================================================================================================

// NewClient 构造函数
func NewClient(ip string, rpcType string, rpcPort map[string]int, r register.Register, opts ...ClientOpt) *Client {
	c := &Client{
		ip:       ip,
		balancer: tool.NewRandBalancer(), //默认随机算法
		rpcType:  rpcType,
		rpcPort:  rpcPort,
		r:        r,
	}
	for _, opt := range opts {
		opt(c)
	}

	c.once.Do(func() {
		resolver.Register(&MyResolversBuild{})
	})
	return c
}

func (c *Client) Dial(ctx context.Context, serverName string) (conn *grpc.ClientConn, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("grpc client dial err: %v", r)
		}
	}()

	//获取target 信息
	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(tool.ServerCircuitInterceptor(circuitList)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(clientParameters),
	}

	port, ok := c.rpcPort[serverName]
	if (c.rpcType == const_g.RpcLocal || c.rpcType == const_g.RpcKubernetes) && !ok {
		return nil, fmt.Errorf("grpc client dial server port not find : %s", serverName)
	}

	if c.rpcType == const_g.RpcRegister && c.r == nil {
		return nil, fmt.Errorf("grpc client dial register is err : %w", err)
	}

	target := ""
	switch c.rpcType {
	case const_g.RpcLocal:
		{
			//获取本地ip:port
			target = fmt.Sprintf("%s:%d", c.ip, port)
		}
	case const_g.RpcKubernetes:
		{
			// server-name.default.svc.cluster.local:port
			// 如果客户端和服务在同一个命名空间（例如 default），可以直接使用短地址: server-name:port
			target = fmt.Sprintf("%s:%d", serverName, port)
		}
	case const_g.RpcRegister:
		{ //服务注册与发现
			target = fmt.Sprintf("%s:///%s", Scheme, serverName)
			opts = append(opts, grpc.WithResolvers(&MyResolversBuild{ctx: ctx, register: c.r})) //服务发现
		}
	default:
		return nil, errors.New("config grpc is err")
	}

	//创建socket 连接
	return grpc.NewClient(target,
		opts...,
	)
}

const Scheme = "svc"

// 1. 创建一个resolver.Builder的一个对象
//	  1.1 连接时遍历所有的Builder对象,判断scheme与Builder.Scheme()匹配,找出对应的Builder对象
//    1.2 调用Build方法返回具体的Resolver对象
// 2. 通过Builder中build方法，可以获取target,cc,opts这几个参数存入Resolver对象中后续逻辑判断
//	  2.1 target存放是一个url,是目标服务的域名地址等信息
//    2.2 cc 是ClientConn连接信息
//    2.3 opts 一些额外的参数信息
// 3. 调用Resolver中的ResolveNow方法更新ClientConn连接的地址列表

type MyResolversBuild struct {
	ctx      context.Context
	register register.Register
}

func (m *MyResolversBuild) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &MyResolver{
		target:   target,
		cc:       cc,
		opts:     opts,
		register: m.register,
		ctx:      m.ctx,
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

func (m *MyResolversBuild) Scheme() string {
	return Scheme
}

//=============================

type MyResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
	opts   resolver.BuildOptions

	register register.Register
	ctx      context.Context
}

func (r *MyResolver) ResolveNow(options resolver.ResolveNowOptions) {
	svcInstances, err := r.register.List(r.ctx, register.GetServerName(strings.Replace(r.target.URL.Path, "/", "", 1), register.InstanceTypeGRPC))
	adders := make([]resolver.Address, 0, 0)

	_ = r.cc.UpdateState(resolver.State{Addresses: adders})
	if err != nil {
		return
	}

	//adders := make([]resolver.Address, 0, len(svcInstances))
	for _, instance := range svcInstances {
		adders = append(adders, resolver.Address{Addr: fmt.Sprintf("%s:%d", instance.Address, instance.Port)})
	}
	_ = r.cc.UpdateState(resolver.State{Addresses: adders})
}

// Close 当resolver不再被使用时，需要调用这个方法来关闭resolver，释放任何它持有的资源
func (r *MyResolver) Close() {
}
