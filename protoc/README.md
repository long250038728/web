### http 跟 grpc对比
    grpc
        1.可以通过protoc buffer 进行序列化/反序列化，传输性能高（proto文件双方都有，所以在序列化/反序列化的时候无需存放其他额外的信息）
        2.类似强类型语言（这个client有什么方法，方法中的参数是什么）—— protoc buffer
        3.同时也是弱点。客户端需要知道protoc文件才能调用，如果什么网站都需要这个，那无法进行下去，（常用于内部调用）
    http
        1.调用是只需要地址+ip就可以获取数据（常用于外部调用）
        2.由于无需protoc文件，导致序列化只能用xml，json等，性能差
        3.类似弱类型语言（传一个dict字段，后端需要根据自己的key获取对应的value，可能由于输错发送跟获取不是同一个key）



安装protoc及grpc插件
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```


用go gen
```
protoc \
--go_out=. \
--go_opt=paths=source_relative \
--go-grpc_out=. \
--go-grpc_opt=paths=source_relative \
--openapiv2_out=. \
-I=/Users/linlong/go/src/  -I=./ \
userServer.proto
```


| 参数                                                      | 解释                - | 其他 |
|---------------------------------------------------------|---------------------|:--:|
| --go_out                                                | go文件生成的目录           |    |
| --go-grpc_out                                           | grpc-go文件生成的目录      |    |
| --go_opt=paths=source_relative / --go_opt=paths=import  | 按源文件的目录组织输出         |    |
| --go-grpc_opt=paths=source_relative                     | 按源文件的目录组织输出         |    |


```
 protoc --go_out=./p/ --go_opt=paths=source_relative  --go-grpc_out=./p/  --go-grpc_opt=paths=source_relative  demo.proto  //根据目录组织
 protoc --go_out=./i/ --go_opt=paths=import  --go-grpc_out=./i/  --go-grpc_opt=paths=import  demo.proto  //根据protoc 中的import
```