package router

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	user "github.com/long250038728/web/application/user/protoc"
	"github.com/long250038728/web/tool/server/http"
	"github.com/long250038728/web/tool/struct_map"
	"sync"
)

func RegisterUserServerServer(engine *gin.Engine, srv user.UserServerServer) {
	//设置错误
	opts := []http.MiddlewareOpt{
		http.SetErrData(map[error]http.Err{errors.New("hello"): http.Err{Code: "1233333", Message: "错误"}}),
	}
	//设置限流

	//设置权限

	//池化
	pool := sync.Pool{New: func() any {
		return http.NewMiddleware(opts...)
	}}

	engine.GET("/", func(gin *gin.Context) {
		//线程池
		middle := pool.Get().(*http.Middleware)
		defer func() {
			middle.Reset()
			pool.Put(middle)
		}()

		//实际操作
		req := &user.RequestHello{}

		type AA struct {
			Name string `json:"name" form:"name"`
		}
		reqT := &AA{}
		middle.GinContext(gin).Bind(reqT).Do(func(ctx context.Context) (interface{}, error) {
			//TODO 因为struct AA中的字段有tag，后续考虑在protoc生成的结构体有form tag
			err := struct_map.Map(reqT, req)
			if err != nil {
				return nil, err
			}
			return srv.SayHello(ctx, req)
		})
	})
}
