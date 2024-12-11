# Go 资源分析总结

## 一、资源分析文件生成

### 1. 使用 `go test` 生成分析文件
通过命令行生成3个文件：`v1.test`、`cpu.profile`、`mem.profile`：
```
# -bench 表示需要benchmark运行的方法,.表示运行本目录所有Benchmark开头的方法
# -benchmem 显示与内存分配相关的详细信息
# -benchtime 设定每个基准测试用例的运行时间
# -cpuprofile 生成 CPU 性能分析文件
# -memprofile 生成内存性能分析文件
go test -bench='.*' -benchmem -benchtime=10s -cpuprofile='cpu.prof' -memprofile='mem.prof'
```

### 2. 业务埋点生成分析文件
在代码中通过pprof库生成CPU和内存分析文件：
```
import "runtime/pprof"
import "os"

cpuOut, _ := os.Create("cpu.out")
defer cpuOut.Close()
memOut, _ := os.Create("mem.out")
defer memOut.Close()

pprof.StartCPUProfile(cpuOut)
defer pprof.StopCPUProfile()
defer pprof.WriteHeapProfile(memOut)

// ... 执行业务逻辑 ...
```


### 3. 通过HTTP生成分析文件
在代码中引入pprof包，启动HTTP服务，然后通过HTTP接口访问分析数据：
```
import (
    _ "net/http/pprof"
    "net/http"
    "log"
)

go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()

//go tool pprof http://192.168.1.2:8002/user/debug/pprof/profile  //获取CPU性能数据(进入pprof工具内)
//go tool pprof http://192.168.1.2:8002/user/debug/pprof/heap     //获取堆内存使用情况(进入pprof工具内)
```


### 4. Gin中间件方式
```
go get github.com/gin-contrib/pprof
import "github.com/gin-contrib/pprof"

ginRouter := gin.Default()
pprof.Register(ginRouter, "user/debug/pprof") // 默认路径为 "debug/pprof"
```



## 二、资源分析文件的使用

### 1. 通过HTTP获取并分析
```
curl http://127.0.0.1:8080/debug/pprof/profile -o cpu.profile
go tool pprof cpu.profile //等同于go tool pprof http://127.0.0.1:8080/debug/pprof/profile
```

### 2. 使用go tool pprof工具
```
# 查看服务器的分析数据
go tool pprof http://192.168.1.2:8002/user/debug/pprof/heap

# 启动交互模式分析文件
go tool pprof v1.test cpu.profile
```
交互模式下常用指令：
* top [N]：列出前N条耗时/内存占用最多的函数。 
* list <regexp>：显示符合正则表达式的代码。 
* peek <regexp>：查看匹配的函数调用关系。

### 3. 启动Web服务展示分析结果
```
go tool pprof -http :8889 v1.test cpu.profile
```

### 4. 导出分析结果为图表
```
sudo yum -y install graphviz.x86_64
go tool pprof -svg cpu.profile > cpu.svg   # 导出为 SVG 格式
go tool pprof -pdf cpu.profile > cpu.pdf   # 导出为 PDF 格式
go tool pprof -png cpu.profile > cpu.png   # 导出为 PNG 格式
```

### 5. trace分析
```
# trace查看
curl 'http://192.168.1.2:8002/user/debug/pprof/trace?seconds=30' >trace.out
go tool trace trace.out
```


## 三、Go项目依赖生成图
```
//安装graphviz
//安装godepgraph
//生成依赖图

brew install graphviz 
go install github.com/kisielk/godepgraph@latest
godepgraph -s ./application/user/cmd/ | dot -Tpng -o godepgraph.png
```

## 四、检测数据竞争
使用-race参数检测数据竞争：
```
go run -race main.go
```


## 五、查看内存逃逸情况
通过编译参数-gcflags "-m"检测内存逃逸情况：
```
go build -gcflags "-m" main.go
```
返回结果解读：
内联优化：inlining call表示进行了内联优化，将函数调用替换为函数实际代码。
栈分配：does not escape表示变量未逃逸，分配在栈上。
堆分配：escapes to heap表示变量逃逸，分配到堆上。
示例：

```
//app.Servers被内联优化。
//函数内变量在栈上分配。
//...arg在函数外被使用，导致变量逃逸到堆上。

./main.go:35:14: inlining call to app.Servers
./main.go:35:14: func literal does not escape
./main.go:35:14: ... argument escapes to heap
```


