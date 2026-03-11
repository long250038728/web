# Tool Library

Go 微服务基础工具库，提供构建微服务所需的常用功能组件。

## 目录结构

```
tool/
├── app/                    # 应用框架（HTTP/RPC 服务器管理、优雅退出）
├── app_const/              # 常量定义
├── app_error/              # 错误定义
├── authorization/          # JWT 认证和会话管理
├── configurator/           # YAML 配置加载器
├── excel/                  # Excel 文件读写
├── gen/                    # 代码生成器（Model/Proto）
├── git/                    # Gitee API 客户端
├── hook/                   # Webhook 通知（企业微信）
├── id/                     # Snowflake 分布式 ID 生成
├── jenkins/                # Jenkins CI/CD 客户端
├── locker/                 # 分布式锁（Redis/Etcd）
├── mq/                     # 消息队列（Kafka/RocketMQ）
├── paths/                  # 配置路径解析
├── persistence/            # 数据持久化层
│   ├── cache/             # Redis 缓存
│   ├── qn/                # 七牛云存储
│   ├── etcd/              # Etcd KV 存储
│   ├── orm/               # GORM 封装（MySQL/ClickHouse）
│   └── es/                # Elasticsearch 配置
├── register/               # 服务注册与发现（Consul）
├── server/                 # 服务器抽象
│   ├── http/              # Gin HTTP 服务器
│   └── rpc/               # gRPC 服务器/客户端
├── sliceconv/              # 数据结构转换工具
├── ssh/                    # SSH 命令执行
├── store/                  # 多层存储（内存+Redis）
├── task/                   # 任务调度
│   ├── cron_job/          # Cron 定时任务
│   └── job/               # SQL 定时任务
└── tracing/                # 分布式追踪（OpenTelemetry）
```

---

## 模块使用说明

### 1. app - 应用框架

**功能**: 管理 HTTP/RPC 服务器、服务注册、优雅退出

**关键接口**:
- `Start()` - 启动服务
- `Stop()` - 停止服务

**调用示例**:

```go
package main

import (
    "context"
    "github.com/long250038728/web/tool/app"
    "github.com/long250038728/web/tool/server/http"
)

func main() {
    // 创建 HTTP 服务
    httpServer := http.NewHttp("my-service", "0.0.0.0", 8080, func(r *gin.Engine) {
        r.GET("/health", func(c *gin.Context) {
            c.String(200, "ok")
        })
        r.GET("/api/data", getDataHandler)
    })

    // 创建应用并启动
    application, err := app.NewApp(
        app.WithServer(httpServer),
        app.WithRegister(registerClient),
        app.WithTracer(tracer),
    )
    if err != nil {
        panic(err)
    }

    if err := application.Start(); err != nil {
        panic(err)
    }
}
```

---

### 2. authorization - JWT 认证

**功能**: JWT token 生成、解析、刷新

**关键接口**:
- `Signed()` - 生成 access_token 和 refresh_token
- `Parse()` - 解析 token 获取用户信息
- `Refresh()` - 刷新 token

**调用示例**:

```go
// 初始化认证
auth := authorization.NewAuth(store,
    authorization.SecretKey([]byte("your-secret-key")),
    authorization.AccessExpires(20*time.Minute),
    authorization.RefreshExpires(7*24*time.Hour),
)

// 登录时生成 token
userClaims := &authorization.UserInfo{
    Id:   "user123",
    Name: "张三",
}
accessToken, refreshToken, err := auth.Signed(ctx, userClaims)

// 中间件解析 token
func AuthMiddleware(auth authentication.Auth) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        ctx, err := auth.Parse(c.Request.Context(), token)
        if err != nil {
            c.JSON(401, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}

// 获取当前用户
func getCurrentUser(c *gin.Context) *authorization.UserInfo {
    claims := authorization.GetClaims(c.Request.Context())
    return claims.UserInfo
}
```

---

### 3. mq - 消息队列

**功能**: Kafka/RocketMQ 统一接口，支持发布订阅、事务消息

**配置示例**:

```yaml
mq:
  kafka:
    address: "192.168.1.100:9092"
    env: "dev"
    batch_size: 100
    batch_bytes: 1048576
    batch_timeout: 1s
```

**生产者调用示例**:

```go
config := &mq.Config{
    Address:      "192.168.1.100:9092",
    Env:          "dev",
    BatchSize:    100,
    BatchBytes:   1048576,
    BatchTimeout: time.Second,
}

producer := mq.NewKafkaProducer(config)
defer producer.Close()

msg := &mq.Message{
    Headers: []mq.Header{
        {Key: "trace_id", Value: []byte("abc123")},
    },
    Data: []byte(`{"userId": 123, "action": "order_create"}`),
}

err := producer.Send(ctx, "topic-order", "key123", msg)
// 批量发送
err := producer.BulkSend(ctx, "topic-order", "key123", []*Message{msg1, msg2})
```

**消费者调用示例**:

```go
consumer := mq.NewKafkaConsumer(config, "topic-order", "consumer-group-1")
defer consumer.Close()

err := consumer.Subscribe(ctx, func(ctx context.Context, c *mq.Message, err error) error {
    if err != nil {
        log.Println("receive error:", err)
        return err
    }

    var data OrderEvent
    json.Unmarshal(c.Data, &data)

    // 处理业务逻辑
    processOrder(data)
    return nil
})
```

---

### 4. persistence/orm - 数据库 ORM

**功能**: GORM 封装，支持 MySQL/ClickHouse，内置 OpenTelemetry 追踪

**配置示例**:

```yaml
persistence:
  orm:
    mysql:
      address: "127.0.0.1"
      port: 3306
      database: "mydb"
      user: "root"
      password: "password"
      table_prefix: "app_"
      read_only: false
```

**调用示例**:

```go
// 创建 MySQL 连接
mysqlConfig := &orm.Config{
    Address:     "127.0.0.1",
    Port:        3306,
    Database:    "mydb",
    User:        "root",
    Password:    "password",
    TablePrefix: "app_",
}
db, err := orm.NewMySQLGorm(mysqlConfig)
if err != nil {
    panic(err)
}

// 查询
type User struct {
    ID   uint
    Name string
    Age  int
}

var users []User
db.Debug().Where("age > ?", 18).Find(&users)

// 使用条件查询
conditions := orm.NewBoolQuery().
    Must(orm.Eq("age", 18), orm.Neq("status", 0)).
    Should(orm.Gt("score", 80), orm.Lte("score", 100))
sql, args := conditions.Do()
db.Raw("SELECT * FROM user WHERE "+sql, args...).Find(&users)

// 插入
db.Create(&User{Name: "张三", Age: 25})

// 更新
db.Model(&User{}).Where("id = ?", 1).Update("age", 26)

// 删除
db.Delete(&User{}, 1)
```

---

### 5. locker - 分布式锁

**功能**: 基于 Redis/Etcd 的分布式锁，支持自动续期

**接口**:
- `Lock()` - 加锁
- `UnLock()` - 解锁
- `Refresh()` - 手动续期
- `AutoRefresh()` - 自动续期
- `Close()` - 关闭

**调用示例**:

```go
// 创建分布式锁
locker := locker.NewRedisLocker(redisStore, "lock_key_123")
defer locker.Close()

ctx := context.Background()
if err := locker.Lock(ctx); err != nil {
    log.Fatal("failed to acquire lock:", err)
}
defer locker.UnLock(ctx)

// 自动续期模式
if err := locker.AutoRefresh(ctx); err != nil {
    log.Fatal("auto refresh failed:", err)
}
// 此时会自动在后台续期，记得最后 Close()

// 业务逻辑
doBusinessWork()
```

---

### 6. register - 服务注册与发现

**功能**: 基于 Consul 的服务注册与发现

**接口**:
- `Register()` - 注册服务
- `DeRegister()` - 注销服务
- `List()` - 获取服务列表
- `Subscribe()` - 订阅服务变化

**类型定义**:

```go
type ServiceInstance struct {
    ID      string
    Name    string
    Address string
    Port    int
    Type    string // "HTTP" or "GRPC"
}
```

**调用示例**:

```go
// 创建服务实例
instance := register.NewServiceInstance("user-service", "192.168.1.100", 8080, register.InstanceTypeHttp)

// 注册服务
err := consul.Register(ctx, instance)
defer consul.DeRegister(ctx, instance)

// 获取服务列表
instances, _ := consul.List(ctx, "user-service")
for _, inst := range instances {
    log.Printf("%s:%d", inst.Address, inst.Port)
}

// 订阅服务变化
ch, _ := consul.Subscribe(ctx, "user-service")
for {
    select {
    case ins := <-ch:
        log.Println("service changed:", ins)
    case <-ctx.Done():
        return
    }
}
```

---

### 7. server/http - HTTP 服务器

**功能**: 基于 Gin 的 HTTP 服务器，内置健康检查、性能监控、Prometheus 指标

**特性**:
- /health - 健康检查
- /{server_name}/pprof - 性能分析
- /{server_name}/metrics/metrics - Prometheus 指标

**调用示例**:

```go
server := http.NewHttp("my-api", "0.0.0.0", 8080, func(r *gin.Engine) {
    // 路由定义
    api := r.Group("/api")
    {
        api.GET("/users", listUsers)
        api.POST("/users", createUser)
        api.GET("/users/:id", getUser)
        api.PUT("/users/:id", updateUser)
        api.DELETE("/users/:id", deleteUser)
    }
})

// 启动
application.AddServer(server)
```

---

### 8. excel - Excel 文件处理

**功能**: 读取/写入 Excel 文件，支持字符串、数字、图片等类型

**调用示例**:

```go
// 读取 Excel
headers := []excel.Header{
    {Name: "name", Key: "Name", Type: excel.HeaderTypeString},
    {Name: "age", Key: "Age", Type: excel.HeaderTypeInt},
    {Name: "salary", Key: "Salary", Type: excel.HeaderTypeFloat},
}

read := excel.NewRead("users.xlsx")
var users []User
if err := read.Read("Sheet1", headers, &users); err != nil {
    panic(err)
}
defer read.Close()

// 写入 Excel
write := excel.NewWrite("output.xlsx")
write.Write("Sheet1", headers, users)
write.Close()
```

---

### 9. gen - 代码生成器

**功能**: 根据数据库表结构生成 Go Model 和 Protobuf 定义

**调用示例**:

```go
// 从数据库生成 Model
db, _ := orm.NewMySQLGorm(mysqlConfig)
gen := gen.NewModelsGen(db)

tables := []string{"user_info", "order_info"}
code, err := gen.Gen("mydb", tables)
if err != nil {
    panic(err)
}
os.WriteFile("models.go", code, 0644)

// 生成 Proto
protoCode, err := gen.GenProto("mydb", tables)
os.WriteFile("model.proto", protoCode, 0644)
```

---

### 10. git - Gitee API 客户端

**功能**: Gitee/Git 仓库自动化操作

**调用示例**:

```go
gitClient, err := git.NewGiteeClient(&git.Config{
    Token: "your-access-token",
})

// 创建分支
err = gitClient.CreateFeature(ctx, "org/repo", "main", "feature/login")

// 创建 PR
pr, err := gitClient.CreatePR(ctx, "org/repo", "feature/login", "main")

// 合并 PR
err = gitClient.MergePR(ctx, "org/repo", "feature/login", "main")
```

---

### 11. hook - Webhook 通知

**功能**: 发送企业微信 webhook 通知

**调用示例**:

```go
client := hook.NewHookClient(hook.WebhookURL("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx"))

msg := &hook.Message{
    MsgType: "text",
    Text: &hook.Text{
        Content: "部署成功！版本 v1.2.3",
        MentionedList: []string{"13800138000"}, // 提及的成员手机号
    },
}

err := client.SendHook(ctx, msg)
```

---

### 12. id - 分布式 ID 生成

**功能**: Snowflake 风格的分布式 ID 生成器

**调用示例**:

```go
idGen := id.NewSnowflakeId(1, 1) // workerId=1, dataCenterId=1

// 生成全局唯一 ID
id := idGen.GenerateId()

// 生成时间有序 ID
tsId := idGen.Generate()
```

---

### 13. jenkins - Jenkins CI/CD

**功能**: Jenkins 构建触发和状态监控

**调用示例**:

```go
jenkins := jenkins.NewJenkinsClient(&jenkins.Config{
    URL:    "http://jenkins.example.com",
    Token:  "api-token",
})

// 触发构建
buildNumber, err := jenkins.Build(ctx, "job-name", map[string]string{
    "branch": "main",
    "version": "1.0.0",
})

// 等待构建完成
err = jenkins.Block(ctx, "job-name", buildNumber)

// 等待构建并开始
number, err := jenkins.BlockBuild(ctx, "job-name")
```

---

### 14. task - 任务调度

**CronJob 调用示例**:

```go
cronJob := task_cron_job.NewCronJob()

// 添加定时任务
id, err := cronJob.AddFunc("*/5 * * * *", func() {
    log.Println("每 5 分钟执行一次")
})

// 添加复杂调度
id2, _ := cronJob.AddFunc("0 0 12 * * MON-FRI", func() {
    log.Println("工作日中午 12 点执行")
})

// 启动（会阻塞）
cronJob.Start()
defer cronJob.Close()
```

---

### 15. tracing/opentelemetry - 分布式追踪

**功能**: OpenTelemetry 集成 Jaeger 进行全链路追踪

**调用示例**:

```go
// 创建 Jaeger Exporter
exporter, _ := jaeger.NewExporter(
    jaeger.Endpoint("http://jaeger:14268/api/traces"),
    jaeger.ServiceName("user-service"),
)

// 创建 Tracer
tracer, err := opentelemetry.NewTrace(ctx, exporter, "user-service")

// 在 Handler 中使用
func handler(c *gin.Context) {
    ctx := c.Request.Context()

    // 创建 Span
    span := opentelemetry.NewSpan(ctx, "process_request")
    defer span.Close()

    span.AddEvent(map[string]any{
        "step": "start_processing",
    })

    // 业务逻辑...
}
```

---

### 16. sliceconv - 数据结构转换

**功能**: 在不同数据类型之间转换

**调用示例**:

```go
// 切片转切片
ints := []int{1, 2, 3}
strings := sliceconv.Map[int, string](ints, func(v int) string {
    return strconv.Itoa(v)
}) // ["1", "2", "3"]

// 结构体字段映射
users := []User{{ID: 1, Name: "张三"}}
maps := sliceconv.Map[User, map[string]interface{}](users, func(u User) map[string]interface{} {
    return map[string]interface{}{
        "id":   u.ID,
        "name": u.Name,
    }
})

// 带 Tag 映射
m := sliceconv.Mapper[User, DTO]{}
dtos := m.Change(users) // 根据 tag 自动映射字段
```

---

### 17. ssh - SSH 命令执行

**功能**: 本地和远程 SSH 命令执行

**调用示例**:

```go
// 本地执行
local := ssh.NewLocalClient()
output, err := local.Run("ls -la")

// 上传文件
err := local.RunFile("echo hello > test.txt")

// 远程执行
remote := ssh.NewRemoteClient(ssh.Host("192.168.1.100"), ssh.User("admin"), ssh.Password("password"))
output, err := remote.Run("docker ps")

// 远程执行脚本
err := remote.RunFile("./deploy.sh")
```

---

### 18. store - 多层存储

**功能**: 内存 + Redis 双层存储，支持 Pub/Sub

**调用示例**:

```go
// 创建 Store
store := store.NewMultiStore(
    store.NewMemoryStore(),
    store.NewRedisStore(&redis.Config{Addr: "localhost:6379"}),
)

// 基本操作
store.Set(ctx, "key", "value", 5*time.Minute)
val, _ := store.Get(ctx, "key")
store.Del(ctx, "key")

// Pub/Sub
pub, _ := store.Pub()
sub, _ := store.Sub("channel")

go func() {
    pub.Publish(ctx, "channel", "message")
}()

msg := <-sub.Receive(ctx)
```

---

## 架构模式

本工具库采用以下设计模式:

1. **Interface-based design** - 所有核心组件都通过接口定义，便于替换实现
2. **Option pattern** - 使用 `Opt func(*T)` 模式进行灵活配置
3. **Context-first design** - 所有操作都接受 context，支持超时和取消
4. **Multi-store strategy** - Store 层支持多策略组合（内存 +Redis）
5. **Unified error handling** - 统一的错误码体系

## 依赖清单

```go
module github.com/long250038728/web/tool

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-jwt/jwt/v4 v4.5.0
    github.com/robfig/cron/v3 v3.0.1
    github.com/segmentio/kafka-go v0.4.47
    github.com/xuri/excelize/v2 v2.8.0
    go.opentelemetry.io/otel v1.21.0
    go.opentelemetry.io/otel/sdk v1.21.0
    gorm.io/driver/clickhouse v0.5.0
    gorm.io/driver/mysql v1.5.2
    gorm.io/gorm v1.25.5
)
```
