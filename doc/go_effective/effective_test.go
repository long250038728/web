package go_effective

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"testing/iotest"
	"time"
)

func TestNumber(t *testing.T) {
	a := int64(100)
	b := int64(100)
	if a > math.MaxInt64-b {
		t.Log("溢出")
	}
	t.Log(a, b)
}

func TestRange(t *testing.T) {
	t.Run("Update_Val_Ptr", func(t *testing.T) {
		s := []curr{{Code: 1}, {Code: 2}, {Code: 3}}
		s2 := []*curr{{Code: 1}, {Code: 2}, {Code: 3}}
		for _, val := range s {
			val.Code = val.Code + 1 //由于是[]curr 值修改不会影响源数据
		}
		for _, val := range s2 {
			val.Code = val.Code + 1 //由于是[]*curr 指针所以修改有效
		}
	})

	t.Run("One_Ptr", func(t *testing.T) {
		data := []int32{1, 2, 3, 4, 5}
		for _, val := range data {
			go func() {
				t.Log(val) //由于val是同个指针地址，导致val打印出来都是同个值
			}()
			value := val
			go func() {
				t.Log(value) // 循环中新增value指针，把val值赋值给value，所以数据为具体的val值 (把变量传入func里面也可以，因为此时也会是值拷贝)
			}()
		}
		time.Sleep(time.Second)
	})
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
	t.Run("two err compute", func(t *testing.T) {
		sourceErr := errors.New("this is source err")
		err1 := fmt.Errorf("this is v err : %v", sourceErr) //拼接(原err信息丢失)
		err2 := fmt.Errorf("this is w err : %w", sourceErr) //包装(原err信息不丢失)

		//不能用 == 来判断两个err是否相等， 用Is及As来判断
		t.Log(err1, errors.Is(err1, sourceErr))
		t.Log(err2, errors.Is(err2, sourceErr))

		//通过递归方式进行取包
		t.Log(errors.Unwrap(err2))

	})

	//错误只应处理一次(直接返回)
	t.Run("once err", func(t *testing.T) {
		err := func() error {
			err := errors.New("once err")
			if err != nil {
				return err
			}
			return nil
		}()
		t.Log(err)
	})

	//错误只应处理一次(通过包装添加额外信息)
	t.Run("once err with append", func(t *testing.T) {
		err := func() error {
			err := errors.New("once err")
			if err != nil {
				//fmt.Println(err)  //这个是错误的，如果在这里print ，在上层也print ,如果并发print就会导致数据不在一起难以排查
				return fmt.Errorf("this is once err %w", err)
			}
			return nil
		}()
		t.Log(err)
	})

	//自定义error类型(因为error是一个接口，实现这个接口就代表可以成为error,使用errors.As进行赋值，赋值成功则代码err是该类型)
	t.Run("curr error", func(t *testing.T) {
		var err error = curr{Msg: "this is err", Code: 200}
		var err2 curr
		if errors.As(err, &err2) {
			t.Log(err2.Code)
		}
	})
}

func TestDefer(t *testing.T) {
	// 一般在程序运行中业务错误不使用panic，一般用于初始化时报错后续无法继续则通过panic终止服务
	// 可通过recover捕获异常，只能捕获同一个协程的，如果不同协程会导致无法捕获
	t.Run("Panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Log(r)
			}
		}()
		panic("this is panic")
	})

	// 我们经常在defer 调用某个函数后会忽略错误值， 可以通过 _ = xxx来进行忽略而非是 无返回值接受(表明是忽略错误而不是不处理错误)
	t.Run("IgnoreErr", func(t *testing.T) {
		defer func() {
			_ = errors.New("this is ignore error")
		}()
		return
	})

	// 我们经常在defer 调用某个函数后会忽略错误值，但是这个错误有时候是必要的，可以通过返回命名变量处理
	t.Run("ReturnErrVal", func(t *testing.T) {
		err := func() (err error) {
			defer func() {
				deferErr := errors.New("this is defer error")
				if err == nil { //我们一般希望业务报错时函数返回业务的报错，如果业务没报错再返回defer中的报错
					err = deferErr
				}
			}()
			err = nil
			return
		}()
		t.Log(err)
	})

	// defer中 return会在defer前执行，导致defer赋值到函数中变量的无效
	t.Run("ReturnErrValFromVal", func(t *testing.T) {
		err := func() error {
			var err error
			defer func() {
				deferErr := errors.New("this is defer error")
				if err == nil { //我们一般希望业务报错时函数返回业务的报错，如果业务没报错再返回defer中的报错
					err = deferErr
				}
			}()
			return err
		}()
		t.Log(err)
	})
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

func TestGoroutine(t *testing.T) {
	t.Run("chan", func(t *testing.T) {
		ch := make(chan int, 0) //分为有缓冲区及无缓冲区 （当缓冲区满了后插入数据就会等待直到可以插入）

		go func() { //这里的wait需要开启一个协程。如果不开启协程的话，wg.Wait()在主协程导致一直等待，无法写入到chan
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				ch <- 1
			}()
			go func() {
				defer wg.Done()
				ch <- 2
			}()
			wg.Wait()
			close(ch) // 关闭通道，通知接收者没有更多数据
		}()

		for val := range ch { //当chan中的数据为空会等待直到有数据到达
			t.Log(val)
			time.Sleep(time.Second)
		}
	})

	t.Run("chan leak", func(t *testing.T) {
		ch := make(chan int, 1) //如果buffer为0时，time.After优先到达，此时chan没有消费者，此时往里面插数据会阻塞导致ch<-1 这里会内存泄露

		go func() {
			time.Sleep(time.Second * 2)
			ch <- 1
		}()

		select {
		case data := <-ch:
			t.Log(data)
		case <-time.After(time.Second):
			t.Log("time after")
		}
	})

	t.Run("chan close", func(t *testing.T) {
		ch := make(chan int)
		closeCh := make(chan struct{})

		// 用另外的chan来判断当前的chan是否被关闭
		go func() {
			num := 0
			for {
				num += 1
				select {
				case <-closeCh:
					close(ch)
					return
				case ch <- num:
					if num == 10000 { //当num == 10000 时关闭closeCh的信号，退出for循环，此时关闭ch
						close(closeCh) // 不能使用num >= 10000 ，当10000发生close(closeCh)时，由于竞争问题case <-closeCh还没执行而是执行了ch <- num ,10001时又再close一次就报错了，
					}
				}
			}
		}()

		for val := range ch {
			fmt.Println(val)
		}
	})

	t.Run("atomic", func(t *testing.T) {
		var count int32 = 1
		t.Log(atomic.LoadInt32(&count)) //读

		t.Log(atomic.SwapInt32(&count, 2)) //替换/赋值
		t.Log(count)

		t.Log(atomic.CompareAndSwapInt32(&count, 2, 4)) //先对比成功再替换/赋值
		t.Log(count)

		t.Log(atomic.AddInt32(&count, 2)) // 加
		t.Log(count)

		atomic.StoreInt32(&count, 4) //减
		t.Log(count)
	})

	t.Run("mutex", func(t *testing.T) {
		rwMutex := sync.RWMutex{}

		count := 0
		var wg sync.WaitGroup
		wg.Add(5)

		for i := 0; i < 5; i++ {
			go func() {
				rwMutex.Lock() //多个协程同时竞争这个锁，如果抢到执行，抢不到等待（原子性atomic就需要自己处理锁等待跟争抢问题）
				defer func() {
					rwMutex.Unlock()
					wg.Done()
				}()
				count += 1
				time.Sleep(time.Second)
			}()
		}

		wg.Wait()
		t.Log(count)
	})

	t.Run("deadlock mutex", func(t *testing.T) {
		c := &curr{}
		c.mu.Lock() //加锁
		defer c.mu.Unlock()
		t.Log(fmt.Sprintf("xxxx %v", c)) //会调用c的String方法，此时由于c的String方法加锁，导致死锁，此时应该不应该用同一个锁

		//解决方案
		c2 := &curr{}
		str := c2.String() //此时已经加锁并解锁

		c2.mu.Lock()                       //加锁
		defer c2.mu.Unlock()               //解锁
		t.Log(fmt.Sprintf("xxxx %v", str)) //此时不会有死锁

	})

	t.Run("context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel() //记得执行cancel函数（避免内存泄露，如果没执行时select其他chan到达，会等到ctx的chan到达才会销毁释放，加上后会立刻销毁释放）

		ctx = context.WithValue(ctx, curr{}, "value") //用curr{}的原因是因为这个结构体是不对外暴露，可以让其他包无法访问。
		t.Log(ctx.Value(curr{}))

		select {
		case <-ctx.Done(): //监听ctx是否已经done了
			t.Log(ctx.Err()) //查看ctx的err信息
			return
		default:
			//当有default，不会阻塞等待，当其他chan没到达就直接到default，无default时就会阻塞等待其中一个chan到达
			t.Log("select不会阻塞等待其他chan")
			return
		}
	})

}

func TestTime(t *testing.T) {
	// time
	t.Run("time use", func(t *testing.T) {
		t1 := time.Now()
		t2, _ := time.Parse(time.DateTime, "2024-07-01 12:00:00")
		t.Log(t1) //2024-07-27 11:05:44.185951 +0800 CST m=+0.000843600 (单调式类型，多了m=+0.000843600，用于比较相对时间)
		t.Log(t2) // 2024-07-01 12:00:00 +0000 UTC (wall类型， 用于比较绝对时间)

		t3 := t1.Truncate(0) //剥离了单调式时间
		t3.Equal(t1)         //判断两个时间用Equal 而不是使用 ==  (会对比单调式类型)
	})

	// 检查内存情况
	menPrint := func() {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		fmt.Printf("Allocated memory: %d bytes\n", memStats.Alloc)
		fmt.Printf("Total memory allocated and not yet freed: %d bytes\n", memStats.TotalAlloc)
		fmt.Printf("Memory obtained from system: %d bytes\n", memStats.Sys)
		fmt.Printf("Number of heap objects: %d\n", memStats.HeapObjects)
		fmt.Println("=============================================")
	}

	//time.After内存泄露
	t.Run("time leak", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		consumer := make(chan int32, 1)
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					consumer <- 1
				}
			}
		}()
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					time.Sleep(time.Second) //每秒检查内存情况
					menPrint()
				}
			}
		}()
		for {
			select {
			//直到time.After的chan数据到达才会释放。如果每次都是其他chan先到达，则就会导致内存泄露
			case <-time.After(time.Second):
				t.Log("time after")
			case <-consumer:

			case <-ctx.Done():
				return
			}
		}
	})

	//解决time.After内存泄露问题
	t.Run("time ticker", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		consumer := make(chan int32, 1)
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					consumer <- 1
				}
			}
		}()
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					time.Sleep(time.Second) //每秒检查内存情况
					menPrint()
				}

			}
		}()
		//不重复创建一个time.After而是只有一个，通过reset重置时间，同时也能通过stop停止进行销毁
		ticker := time.NewTicker(time.Second)
		defer func() {
			ticker.Stop() //执行stop时就会把chan关闭
		}()

		for {
			ticker.Reset(time.Second)
			select {
			case <-ticker.C:
				t.Log("time ticker")
				return
			case <-consumer:

			case <-ctx.Done():
				return
			}
		}
	})
}

func TestJson(t *testing.T) {
	// interface{} 值int类型json序列化后再反序列化就变成值float64
	t.Run("interface int to float", func(t *testing.T) {
		var value int32 = 2
		data := map[string]interface{}{
			"hello": value,
		}
		var unmarshalData map[string]interface{}

		b, _ := json.Marshal(data)
		_ = json.Unmarshal(b, &unmarshalData)

		t.Log(reflect.TypeOf(unmarshalData["hello"]).Kind() == reflect.Float64)
	})

	// 嵌套类型，子类型实现了Marshaler接口后导致序列化/反序列化问题
	t.Run("Marshaler interface", func(t *testing.T) {
		//curr实现了  json.Marshaler接口  {"Msg":"ok","Code":200,"Time":"2024-07-28T21:26:23.97965+08:00"}
		//curr未实现了json.Marshaler接口  "2024-07-28T21:26:45.669756+08:00"  (用的是time.Time中的MarshalJSON方法)
		c := curr{Msg: "ok", Code: 200, Time: time.Now()}
		b, _ := json.Marshal(&c)
		t.Log(string(b))
	})
}

func TestReader(t *testing.T) {
	// bufio包Scanner的使用（减少一次性读取大文件的问题）
	t.Run("bufio", func(t *testing.T) {
		f, _ := os.Open("./effective_test.go")
		scanner := bufio.NewScanner(f)
		builder := bytes.Buffer{}

		// 每次读取4096字节的数据，如果遇到切割符号就返回，如果没有就继续for读取 4096 * X 的数据  存储到buf字段中
		// 当数据已经包含了切割符号，取开始到切割符号的位置 存储到token字段中 ，scanner.Bytes读取的是token的数据
		//	下次scan时继续把数据存储到buf字段中(当开始==结束 && buf >  maxInt  代表缓冲区满了会执行清空再读)。 直到读不到数据
		for scanner.Scan() {
			builder.Write(append(scanner.Bytes(), []byte("\n")...))
		}
		if err := scanner.Err(); err != nil {
			t.Error(err)
		}
		t.Log(builder.String())
	})

	t.Run("strings", func(t *testing.T) {
		reader := strings.NewReader("hello")

		writeWriter := strings.Builder{} //把读到的写入io.writer
		copyWriter := strings.Builder{}  //把读到的写入io.writer

		//每次读写到writeWriter 中
		data := make([]byte, 2, 2) //每次读2byte的内容
		for {
			n, err := reader.Read(data)
			if err != nil && err != io.EOF {
				t.Error(err)
				return
			}

			if err != nil || n == 0 {
				break
			}
			writeWriter.Write(data[:n]) //由于data是复用的，如果不使用[:n]就会导致n后面的内容是上一次的内容或是0
		}
		t.Log(writeWriter.String())

		// 一次性从io.reader 写到 io.writer
		_, _ = io.Copy(&copyWriter, strings.NewReader(writeWriter.String()))
		t.Log(copyWriter.String())
	})
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
	t.Run("race", func(t *testing.T) {
		count := 0
		for i := 0; i < 1000; i++ {
			go func() {
				count += 1
			}()
		}
		t.Log(count)
	})

	t.Run("method", func(t *testing.T) {
		//标识该测试用例跳过
		t.Skip("this is t.Skip")

		// 并行执行多个单元测试用例
		t.Parallel()

		// 执行单元测试时指定-short ，表示只执行短测试
		// go test -short -v .
		if testing.Short() {
			t.Log("this is testing.Short")
		}
	})

	//使用表格驱动测试
	t.Run("table test", func(t *testing.T) {
		tests := map[string]string{
			"hello": "world",
		}
		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				t.Log(tt)
			})
		}
	})

	//httptest包 及 iotest包
	t.Run("ohter test library", func(t *testing.T) {
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
	})
}

func BenchmarkBenchmark(b *testing.B) {
	// 默认1s  可通过 -benchmark=10s 设置时间  或 -count=10
	b.Log(b.N)
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

// 数据传递
func TestChanTransfer(t *testing.T) {
	// 1,2,3,4 ,1,2,3,4, 1,2,3,4 ...
	newWorker := func(id int, ch chan struct{}, nextChan chan struct{}) {
		for {
			tk := <-ch
			t.Log(id + 1)
			time.Sleep(time.Second)
			nextChan <- tk
		}
	}
	chs := []chan struct{}{
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
	}
	/*
		此时有四个协程newWorker方法
			newWorker1  newWorker2  newWorker3  newWorker4
		他们全部阻塞等待数据。此时chs[0] <- struct{}{} 往第一个插入数据
			newWorker1 当前阻塞解除执行逻辑后. 往下一个chan插入数据 。 下一个对应的函数解除阻塞处理逻辑。以此类推
			   ... ...
			newWorker4 当前阻塞解除执行逻辑后给下一个然后当前的阻塞  (chs[(i+1)%4])
	*/
	for i := 0; i < 4; i++ {
		go newWorker(i, chs[i], chs[(i+1)%4])
	}
	chs[0] <- struct{}{}
	select {}
}

//=======

type curr struct {
	mu   sync.Mutex
	Msg  string
	Code int
	time.Time
}

func (c *curr) Bar() string {
	return "bar"
}
func (c curr) Error() string { return c.Msg }
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

func (c *curr) String() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return "this is curr"
}
