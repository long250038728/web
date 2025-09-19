package handles

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/user/internal/service"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http/gateway"
	"github.com/long250038728/web/tool/server/http/gateway/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"time"
)

type Handles struct {
	util *app.Util
}

func NewHandles(util *app.Util) *Handles {
	return &Handles{util: util}
}

func (r *Handles) RegisterHTTPServer(engine *gin.Engine, srv *service.User) {
	authorized, err := r.util.Auth()
	if err != nil {
		panic(err)
	}
	cache, err := r.util.Cache()
	if err != nil {
		panic(err)
	}

	// 服务/领域/接口
	userGroup := engine.Group("/user/user/").Use(middleware.BaseHandle(authorized))
	{
		userGroup.GET("json", func(c *gin.Context) {
			gateway.Json(c, &user.RequestHello{}).Use(
				middleware.Login(),                    //需要登录
				middleware.Validate([]string{"name"}), //参数必填项
				middleware.Rule(c),                    //权限判断
				middleware.Locker(c, cache, []string{"name"}, middleware.Claims(true), middleware.Expiration(time.Second*3)), //分布式锁处理
				middleware.Cache(c, cache, []string{"name"}, middleware.Claims(true)),                                        //缓存处理
			).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.SayHello(ctx, req.(*user.RequestHello))
			})
		})

		userGroup.POST("xml", func(c *gin.Context) {
			gateway.Xml(c, &user.RequestHello{}).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.SayHello(ctx, req.(*user.RequestHello))
			})
		})

		userGroup.GET("file", func(c *gin.Context) {
			gateway.File(c, &user.RequestHello{}).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.File(ctx, req.(*user.RequestHello))
			})
		})

		userGroup.GET("sse", func(c *gin.Context) {
			gateway.SSE(c, &user.RequestHello{}).Handle(func(ctx context.Context, req any) (any, error) {
				return srv.SendSSE(ctx, req.(*user.RequestHello))
			})
		})
	}
}

func (r *Handles) RegisterGRPCServer(engine *grpc.Server, srv *service.User) {
	user.RegisterUserServer(engine, srv)
	grpc_health_v1.RegisterHealthServer(engine, srv)
}
