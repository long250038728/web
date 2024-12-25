package bak

//
//import (
//	"errors"
//	"github.com/gin-gonic/gin"
//	"sync"
//)
//
//type Func func() (interface{}, error)
//
//type FileInterface interface {
//	FileName() string
//	FileData() []byte
//}
//
//type Function interface {
//	Run(middleware *GatewayMiddleware, function Func)
//}
//
//type Gateway struct {
//	pool sync.Pool
//}
//
//func NewGateway() *Gateway {
//	return &Gateway{
//		pool: sync.Pool{
//			New: func() any {
//				return NewGatewayMiddleware()
//			},
//		},
//	}
//}
//
//func (g *Gateway) getGatewayMiddleware(gin *gin.Context) *GatewayMiddleware {
//	return g.pool.Get().(*GatewayMiddleware).SetGin(gin)
//}
//func (g *Gateway) putGatewayMiddleware(middleware *GatewayMiddleware) {
//	middleware.Reset()
//	g.pool.Put(middleware)
//}
//
////==========================================================================================================
//
//func (g *Gateway) JSON(gin *gin.Context, request any, function Func) {
//	g.do(&JSON{}, gin, request, function)
//}
//
//func (g *Gateway) File(gin *gin.Context, request any, function Func) {
//	g.do(&File{}, gin, request, function)
//}
//
//func (g *Gateway) SSE(gin *gin.Context, request any, function Func) {
//	g.do(&SSE{}, gin, request, function)
//}
//
//func (g *Gateway) do(f Function, gin *gin.Context, request any, function Func) {
//	middleware := g.getGatewayMiddleware(gin)
//	defer g.putGatewayMiddleware(middleware)
//
//	if err := middleware.Bind(request); err != nil {
//		middleware.WriteErr(err)
//		return
//	}
//
//	f.Run(middleware, function)
//}
//
////==========================================================================================================
//
//type JSON struct {
//}
//
//func (r *JSON) Run(middleware *GatewayMiddleware, function Func) {
//	resp, err := function()
//	if err != nil {
//		middleware.WriteErr(err)
//		return
//	}
//	middleware.WriteData(resp)
//}
//
//type File struct {
//}
//
//func (r *File) Run(middleware *GatewayMiddleware, function Func) {
//	//处理业务
//	res, err := function()
//	if err != nil {
//		middleware.WriteErr(err)
//		return
//	}
//	if file, ok := res.(FileInterface); !ok {
//		middleware.WriteErr(errors.New("the File is not interface to FileInterface"))
//	} else {
//		middleware.WriteFile(file.FileName(), file.FileData())
//	}
//}
//
//type SSE struct {
//}
//
//func (r *SSE) Run(middleware *GatewayMiddleware, function Func) {
//	// 检查是否支持SSE
//	if err := middleware.CheckSupportSEE(); err != nil {
//		middleware.WriteErr(err)
//		return
//	}
//	//处理业务
//	res, err := function()
//	if err != nil {
//		middleware.WriteErr(err)
//		return
//	}
//
//	if ch, ok := res.(<-chan string); !ok {
//		middleware.WriteErr(errors.New("the response is not chan string"))
//	} else {
//		middleware.WriteSSE(ch)
//	}
//}
