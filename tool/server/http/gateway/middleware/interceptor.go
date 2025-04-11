package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/tool/app_error"
	"github.com/long250038728/web/tool/authorization"
	"github.com/long250038728/web/tool/persistence/redis"
	"github.com/long250038728/web/tool/server"
	"github.com/long250038728/web/tool/server/http/gateway"
	"github.com/long250038728/web/tool/store"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"google.golang.org/grpc/metadata"
	"net/http"
	"regexp"
	"strings"
)

// middle 中间件处理规范
// 1. 只抽离公共的基础逻辑，不应该与实际业务中有联系
// 2. 中间件中由于只提供next() 及 Abort() 方法，只控制继续或停止，无法对实际处理的handle进行捕获。
// 3. 不依赖与request中的参数，response的响应。 所以对接口的缓存，防抖的操作不在middle中处理
// BaseHandle 基本中间件（创建链路及jwt解析，生成新的ctx替换到c.Request.Context中）
// LoginCheckHandle 登录中间件检验 (通过BaseHandle生成的ctx获取Claims对象，获取不到则报错)
// AuthCheckHandle  权限中间件  (通过BaseHandle生成的ctx获取Session对象，根据session中判断是否能进行访问)
// LimitHandle 限流中间件 （对单个用户(token存在获取token，不存在则获取ip)进行限流处理）

// BaseHandle 基本中间件（带上链路及jwt数据）
func BaseHandle(client redis.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.AbortWithStatusJSON(http.StatusOK, gateway.NewResponse(nil, errors.New(fmt.Sprintf("%v", r))))
				return
			}
		}()

		//跨域处理
		{
			// 设置跨域相关的头部信息
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			// 如果是预检请求（OPTIONS 方法），则直接返回成功响应
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent) // 204 No Content
				return
			}
		}

		//前后端分离后，后端不需要处理favicon.ico
		if c.Request.URL.Path == "/favicon.ico" {
			c.AbortWithStatus(http.StatusNoContent) // 204 No Content
			return
		}

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
			authSession := authorization.NewAuth(store.NewRedisStore(client))
			if parseCtx, err := authSession.Parse(ctx, c.GetHeader("Authorization")); err == nil {
				ctx = parseCtx
			}
		}

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

//=================================================================================

func Api(path string) gateway.ServerInterceptor {
	return func(ctx context.Context, requestInfo map[string]any, request any, handler gateway.Handler) (resp any, err error) {
		//获取session对象(session对象默认是有本地store及分布式store的，为了解决频繁获取分布式session的问题)
		sess, err := authorization.GetSession(ctx)
		if err != nil {
			return nil, err
		}

		//判断该url是否在session存在
		isApiAuthorized := false
		for _, url := range sess.AuthList {
			if CamelToSnake(url) == path {
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

// Limit 示例中间件：限流拦截器
func Limit() gateway.ServerInterceptor {
	return func(ctx context.Context, requestInfo map[string]any, request any, handler gateway.Handler) (resp any, err error) {
		fmt.Println("limit")
		// 限流逻辑（省略实际实现）
		return handler(ctx, request)
	}
}

//// LoginCheckHandle 登录中间件检验（校验jwt是否有效）前提是需要执行BaseHandle 中间件
//func LoginCheckHandle() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		_, err := authorization.GetClaims(c.Request.Context())
//		if err != nil {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, NewResponse(nil, app_error.Unauthorized))
//			return
//		}
//		c.Next()
//	}
//}
//
//// AuthCheckHandle 权限中间件(校验jwt中对应的session是否有路径访问权限) 前提是需要执行BaseHandle 中间件
//func AuthCheckHandle() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		//获取session对象(session对象默认是有本地store及分布式store的，为了解决频繁获取分布式session的问题)
//		sess, err := authorization.GetSession(c.Request.Context())
//		if err != nil {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, NewResponse(nil, app_error.Unauthorized))
//			return
//		}
//
//		//判断该url是否在session存在
//		isApiAuthorized := false
//		for _, url := range sess.AuthList {
//			if CamelToSnake(url) == c.Request.URL.Path {
//				isApiAuthorized = true
//				break
//			}
//		}
//
//		if !isApiAuthorized {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, NewResponse(nil, app_error.Unauthorized))
//			return
//		}
//
//		c.Next()
//	}
//}
//
//// LimitHandle 限流中间件 (优先获取用户的token信息，如果接口无需token参数，那通过IP的方式 ---- 单个用户)
//func LimitHandle(client store.Cache) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		if client != nil {
//			identification := c.GetHeader(server.AuthorizationKey)
//			if len(identification) == 0 {
//				identification = c.ClientIP()
//			}
//
//			// 1s 10次
//			limit := limiter.NewCacheLimiter(client, limiter.SetExpiration(time.Second), limiter.SetTimes(10))
//			if err := limit.Allow(c.Request.Context(), fmt.Sprintf("http:%s", identification)); err != nil {
//				c.AbortWithStatusJSON(http.StatusTooManyRequests, NewResponse(nil, app_error.TooManyRequests))
//			}
//		}
//		c.Next()
//	}
//}

func CamelToSnake(url string) string {
	// 使用正则表达式匹配大写字母，并在前面添加下划线
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(url, "${1}_${2}")
	// 将结果转换为小写
	return strings.ToLower(snake)
}
