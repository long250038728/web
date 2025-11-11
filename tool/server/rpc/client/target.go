package client

import (
	"context"
	"fmt"
	"github.com/long250038728/web/tool/register"
	"google.golang.org/grpc"
)

type Target interface {
	Target(ctx context.Context, serverName string) (address string, DialOptions []grpc.DialOption, err error)
}

//=================================================================================================

func NewLocalTarget(ip string, ports map[string]int) Target {
	return &localTarget{ip: ip, ports: ports}
}

func NewRegisterTarget(r register.Register) Target {
	return &registerTarget{r: r}
}

func NewKubernetesTarget() Target {
	return &kubernetesTarget{}
}

//=================================================================================================

type localTarget struct {
	ip    string
	ports map[string]int
}

func (t *localTarget) Target(ctx context.Context, serverName string) (address string, DialOptions []grpc.DialOption, err error) {
	port, ok := t.ports[serverName]
	if !ok {
		return "", nil, fmt.Errorf("grpc client dial server port not find : %s", serverName)
	}
	return fmt.Sprintf("%s:%d", t.ip, port), []grpc.DialOption{}, nil
}

//=================================================================================================

type registerTarget struct {
	r register.Register
}

func (t *registerTarget) Target(ctx context.Context, serverName string) (address string, DialOptions []grpc.DialOption, err error) {
	if t.r == nil {
		return "", nil, fmt.Errorf("grpc client dial register is err : %w", err)
	}
	return fmt.Sprintf("%s:///%s", Scheme, serverName), []grpc.DialOption{grpc.WithResolvers(&MyResolversBuild{ctx: ctx, register: t.r})}, nil
}

//=================================================================================================

type kubernetesTarget struct {
}

func (t *kubernetesTarget) Target(ctx context.Context, serverName string) (address string, DialOptions []grpc.DialOption, err error) {
	return fmt.Sprintf("%s-grpc", serverName), []grpc.DialOption{}, nil
}

//=================================================================================================
