## net/rpc库 v1版本
net/rpc库 rpc的service
```
type helloService struct{}
func (s *helloService) Hello(request string, reply *string) error {
    *reply = "hello :" + request
    return nil
}

func main() {
    rpc.RegisterName("HelloService", new(helloService)) // 修正：helloService 小写
    listener, err := net.Listen("tcp", ":1234") 
    if err != nil {
        log.Fatal("Listen error:", err)
        return
    }
    for { // 建议添加循环，否则只能处理一个连接
        conn, err := listener.Accept() 
        if err != nil {
            log.Fatal("Accept error:", err)
            return
        }
        go rpc.ServeConn(conn)
    }
}
```
net/rpc库 rpc的client
```
func main() {
    client,err := rpc.Dial("tcp","localhost:1234")
    if err != nil {
        return
    }
    var reply string
    err := client.Call("HelloService.Hello","this is request",&reply)
    if err != nil {
        return
    }
    fmt.Println(reply)
}
```

## net/rpc库 v2版本
把一些东西进行抽离,把rpc的库包装在类中，对外只有接口（service完成接口中定义的函数，client端只能调用接口定义的函数）
1. 把ServiceName独立处理
2. 定义这Server有什么方法改为接口类型
3. 提供一个注册类的方法（服务端进行注册，客户端进行调用，他们都遵循Service定义的接口）
```
const HelloServiceName = "HelloService"

type HelloServiceInterface interface {
    Hello(request string, reply *string) error 
}
```
service
```
type helloService struct{}

func RegisterHelloService(svc HelloServiceInterface) error {
    return rpc.RegisterName(HelloServiceName, svc) // 建议返回 error
}

func (s *helloService) Hello(request string, reply *string) error {
    *reply = "hello: " + request // 修正拼写和格式
    return nil
}

func main() {
    RegisterHelloService(new(helloService))
    listener, err := net.Listen("tcp", ":1234")
    if err != nil {
        log.Fatal("Listen error:", err)
    }

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal("Accept error:", err)
            continue
        }
        go rpc.ServeConn(conn)
    }
}
```

client
```
type HelloServiceClient struct {
    client *rpc.Client
}

func (c *HelloServiceClient) Hello(request string, reply *string) error {
    return c.client.Call(HelloServiceName+".Hello", request, reply)
}

func DialHelloService(network, address string) (*HelloServiceClient, error) {
    conn, err := net.Dial(network, address)
    if err != nil {
        return nil, err
    }
    return &HelloServiceClient{client: rpc.NewClient(conn)}, nil
}

func main() {
    client, err := DialHelloService("tcp", "localhost:1234")
    if err != nil {
        log.Fatal("Dial error:", err)
        return
    }
    var reply string
    err = client.Hello("this is request", &reply)
    if err != nil {
        log.Fatal("Call error:", err)
        return
    }
    fmt.Println(reply)
}
```

## net/rpc库 v3版本
为了保持其他语言可以互相调用，所以使用了插件net/rpc/jsonrpc扩展.用有一个普通的TCP服务代替go的RPC版本
```service
func main() {
    rpc.RegisterName(HelloServiceName, new(helloService))
    listener, err := net.Listen("tcp", ":1234")
    if err != nil {
        log.Fatal("Listen error:", err)
    }
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal("Accept error:", err)
            continue
        }
        // 关键修改：使用 jsonrpc 的编解码器
        go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
    }
}

在调用过程中变成了{"method":"HelloService.Hello","params":["this is request"],"id":0}
```

```client
type HelloServiceClient struct {
    client *rpc.Client
}

func (c *HelloServiceClient) Hello(request string, reply *string) error {
    return c.client.Call(HelloServiceName+".Hello", request, reply)
}


func DialHelloService(network, address string) (*HelloServiceClient, error) {
    conn, err := net.Dial(network, address)
    if err != nil {
        return nil, err
    }
    // 将 TCP 连接包装成 JSON-RPC 客户端
    client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
    return &HelloServiceClient{client: client}, nil
}

func (c *HelloServiceClient) Close() error {
    return c.client.Close()
}

func main() {
    client, err := DialHelloService("tcp", "localhost:1234")
    if err != nil {
        log.Fatal("Connect error:", err)
    }
    defer client.Close()

    var reply string
    err = client.Hello("this is request", &reply)
    if err != nil {
        log.Fatal("Call error:", err)
    }
    fmt.Println(reply)
}

响应变成了 {"id"：1,"result":"hello:hello","error":null}
```


## 总结
1. 可以使用net/rpc 进行rpc的调用
2. 把net/rpc + 在jsonrpc库插件 (实现跨语言的json rpc) , 把rpc的细节封装在service及client类中 
3. protoc-gen-go (实现跨语言的protobuf)工具自动生成了service及client的基础代码。
   * 自动生成定义server、client接口
   * 自动生成定义了service的创建方法，（传入conn，svc（满足定义的server接口结构体） 两个参数，然后调用grpc的库） 
   * 自动生成定义了client的创建方法（传入conn），返回client接口。 每个接口内部即完成了类似 c.client.Call(HelloServiceName+".Hello", request, reply) 的实现（c.cc.Invoke(ctx, Auth_Login_FullMethodName, in, out, cOpts...)）
4. grpc库 (把通用逻辑/复杂逻辑放到grpc库中，把重复且易出错的交给代码自动生成)
   * 服务端：通过已经实现了该服务定义接口（自动生成）的结构体对象 + "tcp" + "address:port" 即可实现服务端的创建
   * 客户端：通过 "tcp" + "address:port"  + 初始化函数（自动生成），即可实现客户端，通过调用客户端的方法(自动生成)即可对远程服务进行调用
   * 流程如下： 客户端创建（自动生成的初始化函数） =》 调用客户端接口（自动生成的接口） =》  接口内部实现了c.cc.Invoke(ctx, Auth_Login_FullMethodName, in, out, cOpts...)（自动生成） =》 发起请求到服务器（grpc库实现） =》 服务器接收到消息后反序列化后到执行具体方法（grpc库实现）  =》 方法执行返回(业务实现方)
