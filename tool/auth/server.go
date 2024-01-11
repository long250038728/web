package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var ErrorToken = errors.New("token disabled")

func NewJwtAuth() Auth {
	return &JwtAuth{}
}

type JwtAuth struct {
	white map[string]interface{}
	login map[string]interface{}
	info  map[string]interface{}
}

func (auth *JwtAuth) Allow(c *gin.Context) error {
	path := c.Request.URL.Path
	//无需校验
	if _, ok := auth.white[path]; ok {
		return nil
	}

	//校验token
	if _, ok := auth.login[path]; ok {

	}

	//需要校验权限信息
	if _, ok := auth.info[path]; ok {

	}
	return ErrorToken
}
