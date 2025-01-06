package http

import (
	"context"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/tool/register"
	"github.com/long250038728/web/tool/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	_ "net/http/pprof"
)

var _ server.Server = &Server{}

type HandlerFunc func(*gin.Engine)

type Server struct {
	server      *http.Server
	handler     http.Handler
	address     string
	port        int
	svcInstance *register.ServiceInstance
}

// NewHttp  构造函数
func NewHttp(serverName, address string, port int, handlerFunc HandlerFunc) *Server {
	handler := gin.Default()
	svc := &Server{
		server:      &http.Server{Addr: fmt.Sprintf("%s:%d", address, port), Handler: handler},
		handler:     handler,
		address:     address,
		port:        port,
		svcInstance: register.NewServiceInstance(serverName, address, port, register.InstanceTypeHttp),
	}
	fmt.Printf("%s : %s:%d\n", svc.svcInstance.Name, svc.svcInstance.Address, svc.svcInstance.Port)

	// health 心跳检测
	handler.GET("/health", func(gin *gin.Context) {
		gin.Writer.WriteHeader(http.StatusOK)
		_, _ = gin.Writer.Write([]byte("health"))
	})

	// pprof 监控
	pprof.Register(handler, fmt.Sprintf("/%s/debug/pprof", serverName))

	// prometheus
	// Counter/CounterVec 只有增加指标 (Inc加1   Add加N)
	// Gauge/GaugeVec 增加或减少的指标 (有Set赋值    Inc Dec加减1   Add  Sub加减N)
	// Histogram/HistogramVec 统计数据分布的直方图 (定义桶Buckets: prometheus.LinearBuckets(0, 10, 5) =》(初始值，每个桶加N值, 有多少个桶)   Observe(15)方法用于记录值)
	// Summary/SummaryVec   提供更精确的百分位统计 (定义Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001} =》(初始值，每个桶加N值, 有多少个桶)   Observe(15)方法用于记录值)
	//		Objectives定义了map中3个key 0.5,0.9,0.99这三个key 代表百分之50，百分之90，百分之99有多少值/数量
	//		Objectives定义了map中key对应value  0.05, 0.01, 0.001 即允许的误差是 5% 1% 0.1%
	//		会根据Observe插入的值进行排序，
	//			此时0.5代表中位数（位置中间）的值是多少
	//			此时0.9代表百分之90的位置的值是多少
	//          此时0.9代表百分之99的位置的值是多少
	handler.GET(fmt.Sprintf("/%s/metrics", serverName), gin.WrapH(promhttp.Handler()))

	handlerFunc(handler)
	return svc
}

// Start 开始启动服务
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Stop 停止服务
func (s *Server) Stop() error {
	return s.server.Shutdown(context.Background())
}

// ServiceInstance 获取服务实例（服务注册与发现）
func (s *Server) ServiceInstance() *register.ServiceInstance {
	return s.svcInstance
}
