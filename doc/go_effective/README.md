### 思想:
* 接口 
  * 为了保证知道的最小，接口应设计在客户端这边。但是我们看了很多代码一般都设计在服务端
  1. 通过接口可以解决循环依赖的问题
  2. 通用性(多态) 解耦(依赖注入)  约束能力
* interface{} Any  
  * 尽可能不要使用，只有在非常公共的函数内传入，用反射/类型断言获取具体值，因为会丢失很多信息
* 泛型  
  * 为了解决重复代码的问题，但是也会加大理解程度
* 嵌套
  * 具备这个功能且对外暴露，如果对内暴露的话用组合
  * 嵌套对象时，如果子对象实现了接口方法，由于嵌套对象就会直接使用子对象的方法导致出错。 json.Marshal时子对象已经实现了Marshaler接口导致序列化/反序列化有问题
* 方法使用指针或值接收时
  * 如果内部需要改变值用指针（数据较大避免拷贝也建议有用指针类型），内部无需改变值用值。（理论上避免对一个结构体同时使用指针及值方法）
* range 
  * 在for中val都是值拷贝修改无效(指针类型的拷贝还是原指针所以可以直接修改)，同时val都是同一个指针地址所以操作不当可能会导致有问题
* 读写时尽可能使用io.Reader/ io.writer，同时读取大文件时间可以考虑使用bufio.NewScanner(f)
* 指针
  * 如果传递的是指针就有可能会被外部改变。(ctx.Value("key") 这个值可能会跟保存时的不一样因为被改变)
* 内部函数结束
  * 通过ctx的done方法控制内部是否接受
  * 提供一个Close方法在defer进行关闭


### 基本类型使用事项:
* int float : 
  * 会根据操作系统32/64位生成int32/int64 在对数据操作时避免溢出使用
* slices    : 
  * (切片的底层是数组，当多个切片指针指向同一个数组时就可能会被修改/新增,此时就要考虑是否需要互相影响，是否扩容，是否会导致内存泄露等)
* map       : 
  * (与切片相同底层是数组，如果确定键值对的数量在make时指定数量避免扩容)，map的桶只会扩大不会缩小。避免占用过多所以定期替换新map或map中值为指针
* string    :
  * 如果有中文时不能使用[]byte来判断长度需要使用[]rune判断长度


### 包使用注意事项：
* error: 
  * 返回是否返回命名参数 (当返回多个参数无法知道参数是什么的使用用命名参数(返回经纬度两个参数，不知道哪个是先) ，同时命名参数的return时应该避免返回错对象如err,otherErr)
  * 接受err后如果想对err带上而外的信息可通过拼接或包装的方式， 判断两个错误是否相同不能用 == ，而是需要使用errors.Is 或 errors.As函数。
  * err处理时向上抛或包装一个新的err，而不是处理了一遍后再向上抛，
* defer: 
  * 我们经常会忽略defer中的错误这个是不可取的，当需要忽略时使用 _ = xxx 表明这个错误是忽略的而不是不处理, 不忽略时通过返回命名参数返回
* json:
  * interface 值是int/int32等类型。在json反序列化后会变成float32类型。是因为反序列化时无法知道之前的类型是float还是int。从兼容性的考虑转换为float64
* time:
  * time.After可能会导致内存泄露。 在select语句中使用time.After函数时，如果其他通道的数据到达，time.After的通道还没到达前内存不用释放。可使用time.NewTicker (整个函数只使用一个time channel)
  * 由于time有单调式及wall式的概念。对比相对时间会用单调式，对比绝对时间会用wall式的概念。可用Truncate(0)剥离单调式 ，用Equal判断相同(内部也会判断单调式)
* reflect:
  * 判断两个是否相同可用reflect.DeepEqual，但反射性能很差，一般建议使用自定义判断(看包是否提供equal方法否则自己实现)(反射不单是遍历，同时还需要获取reflect.type，reflect.value等信息所以性能差)
* interface:
  * 接口判断nil时会判断接口类型及值是不是都为nil才认定是nil。当使用自定义类型时值是nil但是类型是curr 所以不是nil。需要通过反射的方式进行判断 reflect.ValueOf(i).IsNil()
* context:
  * 用于控制操作的超时和取消，防止函数长时间运行，避免内存泄漏和资源浪费
* http body
  * 如果我们在没有读取的情况下关闭正文，默认的 HTTP 传输可能会关闭连接。 
  * 如果我们在读取之后关闭正文，默认的 HTTP 传输不会关闭连接；因此，它可以重复使用。


### 并发
* 并发与并行    
  * 并发: 是一个cpu通过中断切换到不同的应用。  
  * 并行: 是利用多个cup同时处理不同的应用
* GMP   
  * G:可以执行对象
  * M:执行的线程
  * P:一个存放各种信息的对象（G队列等）   
  >每一个 OS 线程（M）被调度到 P 上执行，然后每一个 G 运行在 M 上。
* goroutine: 
  * 虽然创建协程的开销很小，但是如果执行的内容的开销比创建协程的开销小，反而是会更加慢
  * 并发原语使用规则:  (mutex最简单性能差，atomic最底层性能最好但是处理最麻烦，channel更直观)
    * 在g1，g2处理逻辑返回到主/其他g推荐使用chanel。
    * 如果对一个段/行代码加锁后其他协程需要等待则使用mutex，
    * 如果发现可原子操作用atomic 
* channel: 
  * chan如果不主动关闭时for data := range chan 是不会退出导致内存泄露。在对内部不知道是否会close的情况下使用val,ok := <-ch 更保险
  * 经常情况下只是传递一个信号(值是什么不重要)，使用空结构体struct{}不占用额外存储空间
  * 是否需要缓冲大小应该根据业务，当插入数据时需要判断是否已满需要等待的场景
  * 如果主goroutine退出，子goroutine由于主协程的结束而被迫终止。没有被优雅关闭导致内存泄露。所以提供Close方法进行优雅关闭（主协程在return之前，主动关闭实现优雅退出）
* sync.Mutex:
  * sync.Mutex 锁
  * sync.RWMutex 读写锁 （指定读锁/写锁性能会更好——读锁的力度更小）
* atomic: 提供了简单类型的原子性操作(int,int32,int64,float,float32,float64)
```
atomic.LoadInt32(&count)    //读
atomic.SwapInt32(&count, 2) //替换/赋值
atomic.CompareAndSwapInt32(&count, 2, 4))//先对比成功再替换/赋值 
atomic.AddInt32(&count, 2)   // 加
atomic.StoreInt32(&count, 4) //减
```
* sync.WaitGroup: 当多个子goroutine执行时，此时当前goroutine如果不等待则会直接退出，所以需要有机制去让当前goroutine等待子goroutine都处理完才退出
```
var wg sync.WaitGroup
wg.Add(1)
go func() {
  defer wg.Done()
  // do something
}()
wg.Wait()
```
* SingleFight:
  * 内部有个map string *call{}对象(*call内部存放sync.WaitGroup 及 interface{}值)，调用判断传入的key是否存在。
    1. 不存在: 创建上面说的*call对象，并写入到map中。sync.WaitGroup.Add(1), 执行fn中方法返回的结果后写入interface{}中，并调用sync.WaitGroup.Done()解除阻塞。delete map中的key
    2. 存在: 利用sync.WaitGroup的Wait进行阻塞等待。等到解除阻塞后。可通过持有对象获取*struct{}对象的值。
  * 设计巧妙之处在于
    1. 通过map判断key是否存在。如果存在通过阻塞的方式进行等待操作(使用sync.WaitGroup)，
    2. 由于map的value是对象指针，在fn执行完成之后delete map中的key。 此时对象指针并未销毁回收只是解除与map的绑定，其他等待的依旧可以通过对象指针获取值
    3. fn方法返回后进行map key的删除，删除后新的相同的key进入就会当做不存在。这样可以不用主动去删除key及是否需要再去fn获取最新的值 （相同key在fn执行期间只保证只有一个请求。如果fn执行结束后又有相同的key就当做新的与之前不相关）







### 测试用例:
* 测试时可以加上以下参数: 
  * -race 检测是否竞态检测     
  * -shuffle 随机    
  * -short 指定本次为short类型    
  * -tags=test1 指定go:build test1 的用例文件
  * -coverprofile=coverage.out 代码覆盖
* 单元测试函数:		     
  * t.Skip("this is t.Skip") 跳过该用例		
  * t.Parallel() 允许并发       
  * testing.Short() 判断当前是不是short类型
* 表格驱动测试
  * 单元测试用例建议使用表格驱动测试,可以覆盖到不同的场景
* 其他
  * 测试提供了httptest包及iotest包



### 优化
在编写可读清晰的代码是第一要素
* CPU
  * L1和L2缓存是在片上他意味与处理器是在同一块硅片。L3是在片外所以性能会差一点，然后就是RAM速度是L1的50到100倍
  * L1缓存有两个子缓存 1. 用于缓存数据 2. 用于缓存命令上下文
  * 高速缓存行概念，当访问特定的内存位置时不久将来。(当CPU想要访问一个特定的内存位置时，它首先检查 L1，然后是 L2，然后是 L3，最后，如果该位置不在这些缓存中，则检查主内存)
    1. 相同的位置会被再次引用 (时间局部性)
    2. 会同时引用附近的存储位置（空间局部性） sum([]int32)的性能一定比sum([]struct.num)快的多 
  * 当一个缓存行在多个内核之间共享，并且至少有一个 goroutine 是写线程时，整个缓存行都会失效
* 堆栈
  * 一个变量既可以分配在栈上，也可以分配在堆上。这两种类型的内存有着根本的不同
    * 栈是默认内存；它是一种后进先出(LIFO)的数据结构，存储特定 goroutine 的所有局部变量
* pprof
  * 通过net/http/pprof包开启HTTP(应用通过SIGPROF信号要求操作系统每隔 10 ms 中断一次。当应用接收到一个SIGPROF时，它会挂起当前的活动，并将执行转移到分析器。分析器收集数据)
    * /debug/pprof/profile端点来激活 CPU 分析
    * /debug/pprof/heap/?debug=0 heap分析
    * /debug/pprof/heap?gc=1 heap分析
    * /debug/pprof/goroutine/?debug=0 goroutine分析
    * /debug/pprof/goroutine/?debug=2
    * /debug/pprof/block
* 通过基准测试 
  * go test -bench=. -cpuprofile profile.out  && go tool pprof -http=:8080 <file>  ||  go tool pprof -http=:8080 -diff_base <file2> <file1>
* 跟踪器
  * go test -bench=. -v -trace=trace.out  && go tool trace trace.out