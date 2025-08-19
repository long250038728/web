package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/tool/authorization"
	"github.com/long250038728/web/tool/server"
	"github.com/long250038728/web/tool/server/http/gateway"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"google.golang.org/grpc/metadata"
	"net/http"
	"runtime/debug"
)

// middle 中间件处理规范
// 1. 只抽离公共的基础逻辑，不应该与实际业务中有联系
// 2. 中间件中由于只提供next() 及 Abort() 方法，只控制继续或停止，无法对实际处理的handle进行捕获。
// 3. 不依赖与request中的参数，response的响应。 所以对接口的缓存，防抖的操作不在middle中处理
// BaseHandle 基本中间件（创建链路及jwt解析，生成新的ctx替换到c.Request.Context中）

// BaseHandle 基本中间件（带上链路及jwt数据）
func BaseHandle(parse authorization.Parse) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.AbortWithStatusJSON(http.StatusOK, gateway.NewResponse(nil, errors.New(fmt.Sprintf("%v", r))))
				debug.PrintStack() //打印panic报错信息
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

		// 用户处理
		if parse != nil {
			if parseCtx, err := parse.Parse(ctx, c.GetHeader("Authorization")); err == nil {
				ctx = parseCtx
			}
		}

		//把所有信息写入metadata中并生成新的ctx
		if c, err := authorization.GetClaims(ctx); err == nil {
			if b, err := json.Marshal(c); err == nil {
				mCarrier["claims"] = string(b)
			}
		}
		if s, err := authorization.GetSession(ctx); err == nil {
			if b, err := json.Marshal(s); err == nil {
				mCarrier["session"] = string(b)
			}
		}
		ctx = metadata.NewOutgoingContext(ctx, metadata.New(mCarrier))

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

//=================================================================================
