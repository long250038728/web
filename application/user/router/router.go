package router

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	user "github.com/long250038728/web/application/user/protoc"
	"github.com/long250038728/web/tool/server/http"
)

func RegisterUserServerServer(engine *gin.Engine, srv user.UserServerServer) {
	//设置错误
	opts := []http.MiddlewareOpt{
		http.SetErrData(map[error]http.Err{errors.New("hello"): http.Err{Code: "1233333", Message: "错误"}}),
	}
	//设置链路

	//设置限流

	//设置权限

	engine.GET("/help", func(gin *gin.Context) {
		req := &user.RequestHello{}
		http.NewMiddleware(gin, opts...).Bind(req).Do(func(ctx context.Context) (interface{}, error) {
			return srv.SayHello(ctx, req)
		})
	})
}
