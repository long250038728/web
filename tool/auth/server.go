package auth

import (
	"github.com/gin-gonic/gin"
)

var Token = "token"

// HandlerFunc 权限认证  ———— http
func HandlerFunc(auth Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		query := c.Request.URL.Query()
		token := c.Request.Header.Get(Token)

		////1.判断路径是否是无需校验的
		if auth.Neglect(path) {
			c.Next()
			return
		}

		//2.判断是否登录校验的
		_, err := auth.Token(token)
		if err != nil {
			c.Writer.Write([]byte("无权限访问"))
			c.Abort()
			return
		}

		// 是登录即可访问
		if auth.Login(token) {
			c.Next()
			return
		}

		// 路径是该角色可以访问
		if auth.Rule(token, path, query) {
			c.Next()
			return
		}

		// 不能访问
		c.Writer.Write([]byte("无权限访问"))
		c.Abort()
	}
}
