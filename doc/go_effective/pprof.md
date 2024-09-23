## 生成文件
通过命令行方式(会生成3个文件 v1.test cpu.profile men.profile)
```
go test -bench=".*" -cpuprofile cpu.profile -memprofile men.profile
```

通过代码生成
```
import "runtime/pprof"
    
cpuOut,_ := os.Create("cpu.out")
defer cpuOut.Close()
memOut,_ := os.Create("men.out")
defer memOut.Close()

pprof.StartCPUProfile(cpuOut)     

defer pprof.StopCPUProfile()
defer pprof.WriteHeapProfile(memOut)

....... do something.......
```
通过http生成
```
go get github.com/gin-contrib/pprof

import "github.com/gin-contrib/pprof"

ginRouter := gin.Default()
pprof.Register(ginRouter, "dev/pprof")  // default is "debug/pprof"
```


## 使用
通过http 生成 cpu.profile文件
```
curl http://127.0.0.1:8080/debug/pprof/profile -o cpu.profile
go tool pprof cpu.profile
```

go tool pprof使用
```
// 服务器查看
go tool pprof http://localhost:8080/debug/pprof/heap 
(heap profile block trace)


// go tool pprof交互模式
go tool pprof v1.test cpu.profile
* top[N]  列出top N 条数据
* list <regexp> 符合正则的代码
* peek <regexp> 符合正则的调用函数及被调用函数


// go tool pprof web模式  本地开启8081端口加载文件
go tool pprof -http="0.0.0.0:8081" v1.test cpu.profile

//go tool pprof 图片生成
sudo yum -y install graphviz.x86_64
go tool pprof -svg cpu.profile > cpu.svg  # svg 格式
go tool pprof -pdf cpu.profile > cpu.pdf # pdf 格式
go tool pprof -png cpu.profile > cpu.png # png 格式
```


