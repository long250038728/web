package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/application/user/ddd/service"
	user "github.com/long250038728/web/application/user/protoc"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/auth"
	"github.com/long250038728/web/tool/limiter"
	"github.com/long250038728/web/tool/server/http/tool"
	_ "net/http/pprof"
	"time"
)

//go func() {
//	log.Println(http.ListenAndServe("localhost:6060", nil))  //"net/http/pprof"
//}()

func RegisterUserServer(engine *gin.Engine, srv *service.UserService, util *app.Util) {
	//设置错误
	//设置限流
	//设置权限
	opts := []tool.MiddlewareOpt{
		tool.Limiter(limiter.NewRedisLimiter(util.Cache(), time.Second, 10)),
		tool.Auth(auth.NewRedisAuth(util.Cache(), auth.WhiteList([]string{"/"}))),
	}

	middleware := tool.NewMiddle(opts...)

	engine.GET("/", func(gin *gin.Context) {
		//请求参数处理
		req := &user.RequestHello{Name: "HELLO"}

		//请求处理
		middleware.JSON(gin, func(ctx context.Context) (interface{}, error) {
			return srv.SayHello(ctx, req)
		})
	})
}
