package go_effective

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"testing/iotest"
	"time"
)

// 接口  为了保证知道的最小，接口应设计在客户端这边。但是我们看了很多代码一般都设计在服务端
//       通过接口可以解决循环依赖的问题
//		 通用性(多态) 解耦(依赖注入)  约束能力
// interface any  尽可能不要使用，只有在非常公共的函数内传入，用反射/类型断言获取具体值，因为会丢失很多信息
// 泛型  为了解决重复代码的问题，但是也会加大理解程度
// 嵌套 具备这个功能且对外暴露，如果对内暴露的话用组合

// int float会根据操作系统32/64位生成int32/int64 在对数据操作时避免溢出使用
// 切片(切片的底层是数组，当多个切片指针指向同一个数组时就可能会被修改/新增,此时就要考虑是否需要互相影响，是否扩容，是否会导致内存泄露等)
// map(与切片相同底层是数组，如果确定键值对的数量在make时指定数量避免扩容)，map的桶只会扩大不会缩小。所以定期替换新map或map中值为指针
// string 如果有中文时不能使用[]byte来判断长度需要使用[]rune判断长度

// 方法使用指针或值接收时，如果内部需要改变值用指针（数据较大避免拷贝也建议有用指针类型），内部无需改变值用值。（理论上避免对一个结构体同时使用指针及值方法）
// 返回是否返回命名参数 (当返回多个参数无法知道参数是什么的使用用命名参数(返回经纬度两个参数，不知道哪个是先) ，同时命名参数的return时应该避免返回错对象如err,otherErr)
// error 接受err后如果想对err带上而外的信息可通过拼接或包装的方式， 判断两个错误是否相同不能用 == ，而是需要使用errors.Is 或 errors.As函数。同时err处理时向上抛或包装一个新的err，而不是处理了一遍后再向上抛，
// defer 我们经常会忽略defer中的错误这个是不可取的，当需要忽略时使用 _ = xxx 表明这个错误是忽略的而不是不处理, 不忽略时通过返回命名参数返回
// interface 接口判断nil时会判断接口类型及值是不是都为nil才认定是nil。当使用自定义类型时值是nil但是类型是curr 所以不是nil。需要通过反射的方式进行判断 reflect.ValueOf(i).IsNil()
// 判断两个是否相同可用reflect.DeepEqual，但反射性能很差，一般建议使用自定义判断(反射不单是遍历，同时还需要获取reflect.type，reflect.value等信息所以性能差)
// range

//
//
//

// time.After内存泄露。 在select语句中使用time.After函数时，如果其他通道的数据到达，time.After的通道还没到达前内存不用释放。
// 如果interface的值是int/int32等类型。在json反序列化后会变成float32类型。是因为反序列化时无法知道之前的类型是float还是int。从兼容性的考虑转换为float64
// 嵌套对象时，json.Marshal 有可能子对象已经实现了Marshaler接口导致序列化/反序列化有问题
// 读写时尽可能使用io.Reader/ io.writer，同时读取大文件时间可以考虑使用bufio.NewScanner(f)
//

func TestNumber(t *testing.T) {
	a := int64(100)
	b := int64(100)
	if a > math.MaxInt64-b {
		t.Log("溢出")
	}
	t.Log(a, b)
}

func TestSlices(t *testing.T) {
	a := []string{"1", "2", "3", "4", "5", "6", "7"}

	// b此时是从坐标2取到到坐标4(不包含4)，同时限制了长度为2，此时如果append，就会开辟一个新的空间，此时修改不会影响到a（因为扩容了）
	// c此时是从坐标2取到到坐标4(不包含4)，此时不限制长度，此时如果append，发现底层a数组无需扩容，只是就会把append替换原先的值，如果需要扩容就会开辟一个新的空间，此时修改不会影响到a
	b := a[2:4:4]
	b = append(b, "6")
	b[1] = "aaa"
	c := a[2:4]
	c = append(c, "append")

	// 使用copy 找最小的长度 (是长度而不是容量，如果长度为0时copy后为依旧为空)
	d := make([]string, 2, 2)
	copy(d, a)

	//由于切割是指针指向原有的数组上面操作，所以可能会内存泄露(如果数组a有1G，此时a退出作用域应该被销毁，但是切割引用了a，导致a无法被销毁)
	//可以用copy 或 新增一个新的数组append
	e := a[:2]                     //内存泄露
	f := make([]string, 0, len(a)) //解决内存泄露
	f = append(f, a...)

	//空数组及数组长度为空
	//所以一般不使用  g == nil (true)      h == nil (false)
	//应该使用len(g)  len(h)
	var g []string            //空数字 null
	h := make([]string, 0, 0) //数组长度为空

	t.Log(a, b, c, d, e, f, g == nil, h == nil, len(g), len(h))
}

func TestMap(t *testing.T) {
	hash := make(map[string]string, 1)
	t.Log(len(hash), fmt.Sprintf("%p %p", hash, &hash))
	hash["hello"] = "hello"
	hash["world"] = "world"
	hash["this"] = "this"
	hash["is"] = "is"
	hash["a"] = "a"
	hash["hash"] = "hash"
	hash["map"] = "map"
	t.Log(len(hash), fmt.Sprintf("%p %p", hash, &hash))

	//copy 只允许切片所以这里通过遍历
	hash2 := make(map[string]string, 10)
	for key, val := range hash {
		hash2[key] = val
	}
	t.Log(hash2)

	//判断是否相同 (用反射的方法判断，但是性能会非常差，一般非生产环境下使用)
	isEqual := reflect.DeepEqual(hash, hash2)
	t.Log(isEqual)

	//通过遍历的方式判断是否相同（使用自定义的方式）
	forIsEqual := func() bool {
		for key, val := range hash {
			if hash2[key] != val {
				return false
			}
		}
		return true
	}()
	t.Log(forIsEqual)
}

func TestString(t *testing.T) {
	str := "my name is lin 中"
	strRune := []rune(str) //int32  一个中文占1个int32 (1到4个字节，一个字节8位,即32位)
	strByte := []byte(str) //uint8  一个中文占3个uint8 导致打印byte时不是一个完整的字
	t.Log(len(str), len(strRune), len(strByte))

	for _, s := range strRune {
		t.Log(string(s))
	}
	t.Log("===============\n")
	for _, s := range strByte {
		t.Log(string(s))
	}

	//字符串拼接
	//由于string是一个空间，每次拼接会生成一个新的空间
	appendStr := "hello" + "world" + "my" + "brother"
	//[]byte是通过数组/切片的方式，在某个情况下会扩容，扩容的大小如果能容下之后的数据就无需多次扩容
	builderStr := strings.Builder{}
	builderStr.Grow(len(appendStr)) //如果已知长度设置大小避免扩容
	builderStr.Write([]byte("hello"))
	builderStr.Write([]byte("world"))
	builderStr.Write([]byte("my"))
	builderStr.Write([]byte("brother"))

	t.Log(appendStr, builderStr.String())

	//内存泄露
	strLong := "hello sister and brother , very happy to see you , i meet you"
	strSmart := strLong[:5]
	t.Log(strLong)
	t.Log(strSmart) //跟切片一样内部引用相同的数组

	//通过copy的方式
	strCopy := strings.Clone(strLong[:5])
	t.Log(strCopy)
}

func TestError(t *testing.T) {
	//哨兵err
	sourceErr := errors.New("this is source err")

	err1 := fmt.Errorf("this is v err : %v", sourceErr) //拼接(原err信息丢失)
	err2 := fmt.Errorf("this is w err : %w", sourceErr) //包装(原err信息不丢失)

	//不能用 == 来判断两个err是否相等， 用Is及As来判断
	t.Log(err1, errors.Is(err1, sourceErr))
	t.Log(err2, errors.Is(err2, sourceErr))

	//错误只应处理一次(直接返回)
	onceErr := func() error {
		err := errors.New("once err")
		if err != nil {
			return err
		}
		return nil
	}

	//错误只应处理一次(通过包装添加额外信息)
	once2Err := func() error {
		err := errors.New("once err")
		if err != nil {
			//fmt.Println(err)  //这个是错误的，如果在这里print ，在上层也print ,如果并发print就会导致数据不在一起难以排查
			return fmt.Errorf("this is once err %w", err)
		}
		return nil
	}
	t.Log(onceErr(), once2Err())

	//自定义error类型(因为error是一个接口，实现这个接口就代表可以成为error,使用errors.As进行赋值，赋值成功则代码err是该类型)
	currErr := func() {
		var err error = curr{Msg: "this is err", Code: 400}
		var err2 curr
		if errors.As(err, &err2) {
			t.Log(err2.Code)
		}
	}
	currErr()
}

func TestDefer(t *testing.T) {
	// 一般在程序运行中业务错误不使用panic，一般用于初始化时报错后续无法继续则通过panic终止服务
	// 可通过recover捕获异常，只能捕获同一个协程的，如果不同协程会导致无法捕获
	panicFunc := func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
		go func() {
			panic("this is panic")
		}()
	}
	panicFunc()

	deferFunc := func() (err error) {
		// 我们经常在defer 调用某个函数后会忽略错误值， 可以通过 _ = xxx来进行忽略而非是 无返回值接受(表明是忽略错误而不是不处理错误)
		defer func() {
			_ = errors.New("this is ignore error")
		}()
		return
	}
	t.Log(deferFunc())

	// 我们经常在defer 调用某个函数后会忽略错误值，但是这个错误有时候是必要的，可以通过返回命名变量处理
	deferFunc2 := func() (err error) {
		defer func() {
			deferErr := errors.New("this is defer error")
			if err == nil { //我们一般希望业务报错时函数返回业务的报错，如果业务没报错再返回defer中的报错
				err = deferErr
			}
		}()
		err = nil
		return
	}
	t.Log(deferFunc2())

	// defer中 return会在defer前执行，导致defer赋值到函数中变量的无效
	deferFunc3 := func() error {
		var err error
		defer func() {
			deferErr := errors.New("this is defer error")
			if err == nil { //我们一般希望业务报错时函数返回业务的报错，如果业务没报错再返回defer中的报错
				err = deferErr
			}
		}()
		return err
	}
	t.Log(deferFunc3())
}

func TestInterface(t *testing.T) {
	// 语法糖： func Bar(c *curr) string 当传入一个nil值不会有什么问题 ，但是如果func Bar(c curr) string 当传入一个nil值就会报错
	var c *curr
	t.Log(c.Bar())

	//由于这里是自定义error类型(实现了Error)。接口判断nil时会判断接口类型及值是不是都为nil才认定是nil。这个例子值是nil但是类型是curr 所以ta不是nil
	//需要通过反射的方式进行判断 reflect.ValueOf(i).IsNil()
	xx := func() error {
		var c *curr
		return c
	}
	t.Log(xx() == nil, c)
}

func TestChan(t *testing.T) {
}

func TestTime(t *testing.T) {
	//time.After内存泄露
	_ = func() {
		consumer := make(chan int32, 1)
		go func() {
			for {
				consumer <- 1
			}
		}()
		go func() {
			for {
				//每秒检查内存情况
				time.Sleep(time.Second)
				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)

				fmt.Printf("Allocated memory: %d bytes\n", memStats.Alloc)
				fmt.Printf("Total memory allocated and not yet freed: %d bytes\n", memStats.TotalAlloc)
				fmt.Printf("Memory obtained from system: %d bytes\n", memStats.Sys)
				fmt.Printf("Number of heap objects: %d\n", memStats.HeapObjects)
				fmt.Println("=============================================")
			}
		}()

		for {
			select {
			//直到time.After的chan数据到达才会释放。如果每次都是其他chan先到达，则就会导致内存泄露
			case <-time.After(time.Second):
				fmt.Println("time after")
			case <-consumer:
			}
		}
	}

	//解决time.After内存泄露问题
	_ = func() {
		consumer := make(chan int32, 1)
		go func() {
			for {
				consumer <- 1
			}
		}()
		go func() {
			for {
				//每秒检查内存情况
				time.Sleep(time.Second)
				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)

				fmt.Printf("Allocated memory: %d bytes\n", memStats.Alloc)
				fmt.Printf("Total memory allocated and not yet freed: %d bytes\n", memStats.TotalAlloc)
				fmt.Printf("Memory obtained from system: %d bytes\n", memStats.Sys)
				fmt.Printf("Number of heap objects: %d\n", memStats.HeapObjects)
				fmt.Println("=============================================")
			}
		}()
		//不重复创建一个time.After而是只有一个，通过reset重置时间，同时也能通过stop停止进行销毁
		t := time.NewTicker(time.Second)
		defer func() {
			t.Stop() //执行stop时就会把chan关闭
		}()

		for {
			t.Reset(time.Second)
			select {
			case <-t.C:
				fmt.Println("time after")
				return
			case <-consumer:
				//fmt.Println("consumer")
			}
		}
	}

	// time
	_ = func() {
		t1 := time.Now()
		t2, _ := time.Parse(time.DateTime, "2024-07-01 12:00:00")
		t.Log(t1) //2024-07-27 11:05:44.185951 +0800 CST m=+0.000843600 (单调式类型，多了m=+0.000843600，用于比较相对时间)
		t.Log(t2) // 2024-07-01 12:00:00 +0000 UTC (wall类型， 用于比较绝对时间)

		t3 := t1.Truncate(0) //剥离了单调式时间

		t3.Equal(t1) //判断两个时间用Equal 而不是使用 ==
	}
}

func TestJson(t *testing.T) {
	// interface{} 值int类型json序列化后再反序列化就变成值float64
	_ = func() {
		var value int32 = 2
		data := map[string]interface{}{
			"hello": value,
		}
		var unmarshalData map[string]interface{}

		b, _ := json.Marshal(data)
		_ = json.Unmarshal(b, &unmarshalData)

		t.Log(reflect.TypeOf(unmarshalData["hello"]).Kind())
	}

	_ = func() {
		c := curr{Msg: "ok", Code: 200, Time: time.Now()}
		b, _ := json.Marshal(&c)
		t.Log(string(b))

		//curr实现了  json.Marshaler接口  {"Msg":"ok","Code":200,"Time":"2024-07-28T21:26:23.97965+08:00"}
		//curr未实现了json.Marshaler接口  "2024-07-28T21:26:45.669756+08:00"
	}
}

func TestReader(t *testing.T) {
	// bufio包Scanner的使用（减少一次性读取大文件的问题）
	_ = func() {
		f, _ := os.Open("./effective_test.go")
		scanner := bufio.NewScanner(f)

		// 每次读取4096字节的数据，如果遇到切割符号就返回，如果没有就继续for读取 4096 * X 的数据  存储到buf字段中
		// 当4096中的数据已经包含了切割符号，取开始到切割符号的位置 存储到token字段中 ，
		//	下次scan时继续把数据存储到buf字段中(当开始==结束 && buf >  maxInt  代表缓冲区满了会执行清空再读)。 直到读不到数据
		for scanner.Scan() {
			t.Log(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			t.Log(err)
		}
	}
}

func TestTest(t *testing.T) {
	// //go:build test1
	//  标识种类的方法
	//	go test -v .    			表示值运行不含build标识的方法
	//  go test --tags=test1 -v .   表示值运行不含build标识及test1标识的方法

	// 随机执行单元测试用例
	// go test -shuffle=on -v .
	// 希望以之前的随机顺序执行，指定上次的执行返回的随机值
	// go test -shuffle=324213412  -v .

	//代码覆盖率
	// go test -coverprofile=coverage.out ./...
	// go test -coverprofile=./... -coverprofile=coverage.out ./...
	// go tool cover -html=coverage.out

	// 执行单元测试时指定-race ，表示运行时会检查竞态检测
	// go test -race -v .

	//标识该测试用例跳过
	_ = func() {
		t.Skip("this is t.Skip")
	}

	// 并行执行多个单元测试用例
	_ = func() {
		t.Parallel()
		t.Log("hello")
	}

	// 执行单元测试时指定-short ，表示只执行短测试
	// go test -short -v .
	_ = func() {
		if testing.Short() {
			t.Log("this is testing.Short")
		}
	}

	//使用表格驱动测试
	_ = func() {
		tests := map[string]struct{}{}
		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				t.Log(tt)
			})
		}
	}

	//httptest包 及 iotest包
	func() {
		//测试服务端处理逻辑
		handle := func(http.ResponseWriter, *http.Request) {
			//do something
		}
		req := httptest.NewRequest(http.MethodGet, "https://localhost", strings.NewReader("isOpen=true"))
		w := httptest.NewRecorder()
		handle(w, req)
		t.Log(w.Result().Status)

		//测试客户端发起请求
		svc := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = w.Write([]byte("hello"))
		}))
		defer svc.Close()
		t.Log(svc.Client().Post("https://localhost", "application/json", nil))

		//测试iotest
		_ = iotest.TestReader(strings.NewReader("hello"), []byte("hello"))
	}()
}

func BenchmarkBenchmark(b *testing.B) {
	// 默认1s  可通过 -benchmark=10s 设置时间  或 -count=10

	b.ResetTimer() //重置时间
	b.StopTimer()  //停止
	b.StartTimer() //开始
	for i := 0; i < b.N; i++ {
		func() {
			c := 0
			for i := 0; i < 1000; i++ {
				c += 1
			}
		}()
	}
}

//=======

type curr struct {
	Msg  string
	Code int
	time.Time
}

func (c curr) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Msg  string
		Code int
		Time time.Time
	}{
		Msg:  c.Msg,
		Code: c.Code,
		Time: c.Time,
	})
}
func (c *curr) Bar() string {
	return "bar"
}
func (c curr) Add(Msg string, Code int) {
	c.Msg = Msg
	c.Code = Code
}
func (c curr) Error() string { return c.Msg }
