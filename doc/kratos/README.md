# Kratos框架的学习

## 1.项目新建
```
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
kratos new helloworld
```

## 2.代码分析
### App && Option
通过option类保存基本的数据信息(如id，name，timeout等)
```
type Option func(o *options)

// options is an application options.
type options struct {
	id        string
	name      string
	version   string
	metadata  map[string]string
	endpoints []*url.URL

	ctx  context.Context
	sigs []os.Signal

	logger           log.Logger
	registrar        registry.Registrar
	registrarTimeout time.Duration
	stopTimeout      time.Duration
	servers          []transport.Server

	// Before and After funcs
	beforeStart []func(context.Context) error
	beforeStop  []func(context.Context) error
	afterStart  []func(context.Context) error
	afterStop   []func(context.Context) error
}
```
定义了app类，主要存放是整个服务需要的东西(如ctx，cancel及 options等)
```
type App struct {
	opts     options
	ctx      context.Context
	cancel   context.CancelFunc
	mu       sync.Mutex
	instance *registry.ServiceInstance
}

func New(opts ...Option) *App {}
func (a *App) Run() error {}
func (a *App) Stop() (err error) {}
```
提供的方法
1. New方法实现的功能（创建*App，把基本信息赋值给App中）
   * 通过 ...Option 方式可以init opts参数
   * 创建带有取消的context的ctx,cancel 赋值给app
   * 创建App对象，把上面的对象详细保存到App中
2. Run方法实现的功能(运行)
   * 创建*registry.ServiceInstance服务注册实例赋值给*App (记录app.opts.id ,name ,version，endpoints等信息) 
   * 创建一个sctx的上下文 (把*App的信息放入到sctx上下文中)
   * 通过sctx 创建errorGroup => eg,ctx 
       1. 遍历*App.servers , 每个server中 eg调用两个go协程 
          * 第一是使用<-ctx.Done()阻塞等待，当释放时可以对server进行Stop方法的操作，此时的是新的Ctx资源(用于释放资源)
          * 第二是对server进行Start操作，此时传入的是*App 的Ctx
       2. 监听 signal.Notify(c, a.opts.sigs...) 退出信号，监听到则调用*App.Stop
   * 如果opts.registrar有值是需要进行服务注册 (根据*registry.ServiceInstance信息)
3. Stop方法实现的功能（停止）
   * 如果opts.registrar有值是需要进行服务下线 (根据*registry.ServiceInstance信息)
   * 调用*App.Cancel ，此时 eg中server中阻塞的就会唤起调用server.Stop方法进行资源释放

#### 总结
1. 使用了options的方式把所有的信息放入到option对象中
2. app中的run方法使用了errorGroup的巧妙指出用于一个协程运行server，一个协程`阻塞等待`server的停止操作
3. 当signal.Notify退出信息时调用App的Stop方法，服务注册下线及执行*App的cancel()。此时唤醒`阻塞等待`的server停止操作

---

### transport.Server
定义了接口，只要实现了这个接口的都运行当成server服务（有http，grpc）
```
type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}
```
http/grpc
1. NewServer同样都是使用了Options对Server进行初始化(address,router,timeout等信息)
2. Start/Stop方法进行对应的启动及服务停止

---

### 服务注册与发现 Registrar
1. 定义了接口，只要实现了这个接口的都运行当成服务注册与发现
```
type Registrar interface {
	// Register the registration.
	Register(ctx context.Context, service *ServiceInstance) error
	// Deregister the registration.
	Deregister(ctx context.Context, service *ServiceInstance) error
}
```
在kratos/contrib/registry目录下分装了各个服务注册与发现的工具类

---

### middleware.Middleware
与常规的http的middleware使用方式一样
```
type Handler func(ctx context.Context, req any) (any, error)

func XXX(args xxx) middleware.Middleware {
    return  func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req any) (any, error) {
            
            // 根据中间件传入的变量进行处理 （可以是工具，可以是函数）
            // 如把数据放入ctx中
            // 如做一些前置的逻辑判断
            
            resp,err := handler(ctx,req)
            
            // 可以做一些逻辑处理
            // 如计算耗时等，记录链路等信息
            return resp,err
        }
    }
}

```

### Config
1. 定义了接口，只要实现了这个接口的都运行当成config工具
2. 通过config.WithSource的方式可以把多个source进行填充，获取配置时可以根据多个渠道获取配置
```
type Config interface {  //在加载配置的工具类前再封装一层，用于存放基本数据信息
	Load() error
	Scan(v interface{}) error
	Value(key string) Value
	Watch(key string, o Observer) error
	Close() error
}

type Source interface {  //负责加载配置的工具类
	Load() ([]*KeyValue, error)
	Watch() (Watcher, error)
}
ype Watcher interface {
	Next() ([]*KeyValue, error)
	Stop() error
}
type KeyValue struct {
	Key    string  
	Value  []byte  //***数据
	Format string
}

config.New(config.WithSource(file.NewSource(flagconf)))
```
1. Load方法遍历所有的sources，调研source.Load获取数据，然后汇总到一个，同时开启协程进行watch每个source.watch
   * file
     * Load方法是记载文件
     * Watch是用fsnotify.NewWatcher()进行监听文件变化
   * consul
     * Load方法是调用consul.client.KV().List(PATH)
     * Watch是用github.com/hashicorp/consul/api/watch包的watch.Parse方法返回的wp,执行 wp.RunWithClientAndHclog(s.client)返回chan监听
2. Scan 把刚才汇总的绑定。

---

### Log
1. 定义了接口，只要实现了这个接口的都运行当成log工具
```
type Logger interface {
	Log(level Level, keyvals ...interface{}) error
}
```

---

## 3.工具分析
### MakeFile
```
.PHONY: echo
echo:
   echo "hello world" 
```

### protoc
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest                  // protoc.go生成工具
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest                 // protoc_grpc.go生成工具
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest                     // kratos cmd生成工具
go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest         // protoc_http.go生成工具
go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest              // openai 文档生成工具
go install github.com/google/wire/cmd/wire@latest                               // wire 依赖注入生成工具

protoc --proto_path=./api \                                                     //proto文件路径 （项目api）
	--proto_path=./third_party \                                                //protoc文件路径（外部依赖）
 	--go_out=paths=source_relative:./api \                                      //go out 使用source_relative,生成路径在./api
 	--go-http_out=paths=source_relative:./api \                                 //http out 使用source_relative,生成路径在./api
 	--go-grpc_out=paths=source_relative:./api \                                 //grpc out 使用source_relative,生成路径在./api
	--openapi_out=fq_schema_naming=true,default_response=false:. \              //openai文档生成 ,生成路径在.
	$(API_PROTO_FILES)	
```

## 3.总结
1. 项目中各个功能都使用的是interface方式，这样好处就是解耦，同时只要实现这些方法就可以（duck方法）
2. 在初始化时都使用options的方法可以根据需求进行参数化的配置(创建是先指定默认参数，然后通过option方法赋值)
3. 由于go语言的特性无法使用动态代理，所以会使用各种gen方法及Makefile方法