package router

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	user "github.com/long250038728/web/application/user/protoc"
)

func RegisterUserServerServer(engine *gin.Engine, srv user.UserServerServer) {
	engine.GET("/help", func(gin *gin.Context) {
		resp, err := srv.SayHello(gin.Request.Context(), &user.RequestHello{
			Name: "hello",
		})

		if err != nil {
			gin.Writer.Write([]byte("获取数据错误"))
			return
		}

		data, err := json.Marshal(resp)
		if err != nil {
			gin.Writer.Write([]byte("数据压缩失败"))
			return
		}

		gin.Writer.Write(data)
	})
}
