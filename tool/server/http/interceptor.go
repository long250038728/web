package http

import (
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/tool/authorization"
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/limiter"
	"github.com/long250038728/web/tool/server"
	"github.com/long250038728/web/tool/system_error"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"google.golang.org/grpc/metadata"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func LimitHandle(client cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		if client != nil {
			limit := limiter.NewCacheLimiter(client, limiter.SetExpiration(time.Second*10), limiter.SetTimes(10))
			if err := limit.Allow(c.Request.Context(), "http:"+c.GetHeader(server.AuthorizationKey)); err != nil {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, NewResponse(nil, system_error.TooManyRequests))
			}
		}
		c.Next()
	}
}

func BaseHandle(client cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		//链路追踪
		{
			//生成一个新的带有span的context
			// 1. c.Request.Context() 获取ctx上下文 (也可以通过context.Background()创建)
			// 2. 通过 telemetry 提取 http请求头中的参数生成一个名称为请求URI的 span (如果请求头中有traceparent 则生成一个子span，如果无则生成一个root span)
			// 3. 通过span 获取新的 ctx 以后续使用
			span := opentelemetry.NewSpan(opentelemetry.ExtractHttp(c.Request.Context(), c.Request), c.Request.RequestURI) //记录请求头
			defer span.Close()

			ctx = span.Context()

			mCarrier := map[string]string{server.AuthorizationKey: c.GetHeader("Authorization")} // mCarrier["authorization"] = authorization // 把 http 请求头中的Authorization信息写入mCarrier
			opentelemetry.InjectMap(ctx, mCarrier)                                               // 把 telemetry的id等信息写入mCarrier

			span.AddEvent(mCarrier)
			c.Header(server.TraceParentKey, mCarrier[server.TraceParentKey])

			//把所有信息写入metadata中并生成新的ctx
			ctx = metadata.NewOutgoingContext(ctx, metadata.New(mCarrier))
		}

		// 用户处理
		if client != nil {
			authSession := authorization.NewAuth(client)
			if parseCtx, err := authSession.Parse(ctx, c.GetHeader("Authorization")); err == nil {
				ctx = parseCtx
			}
		}

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func LoginCheckHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := authorization.GetClaims(c.Request.Context())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, NewResponse(nil, system_error.Unauthorized))
			return
		}
		c.Next()
	}
}

func AuthCheckHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取session对象
		sess, err := authorization.GetSession(c.Request.Context())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, NewResponse(nil, system_error.Unauthorized))
			return
		}

		//判断该url是否在session存在
		isApiAuthorized := false
		for _, url := range sess.AuthList {
			if CamelToSnake(url) == c.Request.URL.Path {
				isApiAuthorized = true
				break
			}
		}

		if !isApiAuthorized {
			c.AbortWithStatusJSON(http.StatusUnauthorized, NewResponse(nil, system_error.Unauthorized))
			return
		}

		c.Next()
	}
}

func CamelToSnake(url string) string {
	// 使用正则表达式匹配大写字母，并在前面添加下划线
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(url, "${1}_${2}")
	// 将结果转换为小写
	return strings.ToLower(snake)
}
