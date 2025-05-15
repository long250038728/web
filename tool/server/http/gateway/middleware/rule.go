package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/tool/app_error"
	"github.com/long250038728/web/tool/authorization"
	"github.com/long250038728/web/tool/server/http/gateway"
	"net/url"
	"regexp"
	"strings"
)

func checkQuery(u *url.URL, requestInfo map[string]any) bool {
	values := u.Query()
	if len(values) == 0 {
		return true
	}

	for key, value := range values {
		requestValue, ok := requestInfo[key]
		if !ok {
			return false
		}
		if fmt.Sprintf("%v", requestValue) != value[0] {
			return false
		}
	}
	return true
}

// Rule 校验接口是否有权限
func Rule(c *gin.Context) gateway.ServerInterceptor {
	return func(ctx context.Context, requestInfo map[string]any, request any, handler gateway.Handler) (resp any, err error) {
		//获取session对象(session对象默认是有本地store及分布式store的，为了解决频繁获取分布式session的问题)
		sess, err := authorization.GetSession(ctx)
		if err != nil {
			return nil, err
		}
		isApiAuthorized := false
		path := c.Request.URL.Path

		for _, authPath := range sess.AuthList {
			u, err := url.Parse(authPath)
			if err != nil {
				continue
			}

			//校验path路径
			if CamelToSnake(u.Path) != CamelToSnake(path) {
				continue
			}

			// 如果参数校验成功代表通过
			if checkQuery(u, requestInfo) {
				isApiAuthorized = true
				break
			}
		}

		if !isApiAuthorized {
			return nil, app_error.Unauthorized
		}
		return handler(ctx, request)

	}
}

func CamelToSnake(url string) string {
	// 使用正则表达式匹配大写字母，并在前面添加下划线
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(url, "${1}_${2}")
	// 将结果转换为小写
	return strings.ToLower(snake)
}
