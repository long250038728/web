# Tool Library - TODO

本文档记录了 `tool` 库中存在的代码问题、功能缺失和优化建议。

---

## 🔴 严重问题 (需优先修复)

### 1. context 使用错误

**位置**: `app/app.go:90-91`

```go
<-app.ctx.Done()        //等待 app.ctx.Done 触发
time.Sleep(time.Second) //这个时候就不能用 app.ctx 应该这个 ctx 已经 cancel
return app.register.DeRegister(context.Background(), svc.ServiceInstance())
```

**问题**: 在 `app.ctx` 已经取消后，又睡眠了 1 秒，此时使用的是 `context.Background()`，这会导致：
- 如果服务注册失败，无法通过超时机制快速失败
- 失去了优雅退出的意义

**建议修复**:

```go
<-app.ctx.Done()
// 创建一个带超时的 context 用于清理操作
cleanupCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
defer cancel()
return app.register.DeRegister(cleanupCtx, svc.ServiceInstance())
```

---

### 2. 未处理的错误

**位置**: `git/client.go:63`, `mq/kafka.go:177`, `mq/kafka.go:194`, `persistence/orm/orm.go:98-127`

```go
// git/client.go:63
defer fmt.Println(string(b))

// mq/kafka.go
_ = callback(ctx, nil, err)
_ = m.reader.CommitMessages(subscribeCtx, kafkaMessage)
```

**问题**:
- 日志信息直接丢弃，没有记录到正式日志系统
- 错误的 `_` 忽略可能导致消息丢失或重复消费

**建议修复**:

```go
// 引入 log 包并记录
log.Printf("response: %s", string(b))

// 错误应该上报或记录
if err := callback(...); err != nil {
    log.Errorf("callback error: %v", err)
}
```

---

### 3. 缺少错误类型判断

**位置**: `git/client.go:116-118`

```go
var res response
if err := json.Unmarshal(resp, &res); err != nil {
    return nil  // 应该返回错误
}
```

**问题**: 解析响应失败时返回 `nil` 而不是错误，调用者无法得知失败原因

**建议修复**:

```go
if err := json.Unmarshal(resp, &res); err != nil {
    return fmt.Errorf("failed to parse merge response: %w", err)
}
```

---

### 4. ID 生成器的并发安全问题

**位置**: `id/id.go`

**问题**: Snowflake ID 生成器需要确保 workerId 和数据中心 ID 的全局唯一性，当前实现可能未充分考虑并发场景下的线程安全。

**建议**: 添加互斥锁保护 ID 生成逻辑。

---

### 5. Redis Lua 脚本安全性问题

**位置**: `locker/redis_client.go`

**问题**: Lua 脚本中使用 `KEYS` 命令模式可能存在性能问题，应使用具体键名代替通配符。

**建议修复**: 检查所有 Lua 脚本，确保不使用 `KEYS` 模式匹配。

---

## 🟡 功能缺失 (建议补充)

### 1. 配置加载器缺少环境变量支持

**位置**: `configurator/yaml.go`

**问题**: 只能加载固定路径的 YAML 配置，不支持环境变量替换

**建议新增功能**:

```go
// 支持 ${ENV_VAR} 占位符替换
func LoadWithEnv(path string) (*Config, error)

// 或者直接传入配置映射
func LoadFromMap(data map[string]interface{}) *Config
```

---

### 2. ORM 缺少分页查询支持

**位置**: `persistence/orm/query.go`

**问题**: 当前只提供了条件构建，缺少分页参数支持

**建议新增功能**:

```go
type QueryBuilder struct {
    table   string
    conditions []Query
    orders   []string
    offset   int
    limit    int
}

func (b *QueryBuilder) Offset(n int) *QueryBuilder
func (b *QueryBuilder) Limit(n int) *QueryBuilder
func (b *QueryBuilder) OrderBy(field string, desc ...bool) *QueryBuilder
func (b *QueryBuilder) Count() (int64, error)
func (b *QueryBuilder) Find(result interface{}) error
```

---

### 3. MQ 缺少事务消息支持（RocketMQ）

**位置**: `mq/rocket.go`

**问题**: RocketMQ 的事务消息功能未实现

**建议新增功能**:

```go
type TransactionProducer interface {
    ProduceTransaction(ctx context.Context, topic string, msg *Message, halfCheck Callback) error
}
```

---

### 4. Store 缺少批量操作

**位置**: `store/type.go`

**问题**: 只有单 key 的 Get/Set/Del，缺少批量操作

**建议新增接口**:

```go
type Store interface {
    Get(ctx context.Context, keys ...string) (map[string]string, error)
    SetNX(ctx context.Context, kvs map[string]string, expiration time.Duration) (map[string]bool, error)
    Del(ctx context.Context, keys ...string) (int64, error) // 返回删除数量

    // 新加批量操作
    MGet(ctx context.Context, keys []string) ([]string, error)
    MSet(ctx context.Context, pairs map[string]string) error
    DelMulti(ctx context.Context, keys []string) (int64, error)

    // 支持 pattern 匹配删除
    DelByPattern(ctx context.Context, pattern string) (int64, error)

    Close()
}
```

---

### 5. Excel 缺少写入图片支持

**位置**: `excel/write.go`

**问题**: 读取支持图片 (`HeaderTypeImage`)，但写入时缺少对应的写图片方法

**建议新增**:

```go
type Write struct {
    file *excelize.File
}

func (w *Write) WriteImage(sheet, cell string, imagePath string) error
func (w *Write) WritePic sheet, cell string, pic Pic) error
```

---

### 6. Locker 缺少锁超时重试机制

**位置**: `locker/type.go`

**问题**: `Lock()` 方法一旦失败就直接返回，不支持重试

**建议新增**:

```go
func (l *Locker) LockWithRetry(ctx context.Context, maxRetries int, retryInterval time.Duration) error
```

---

### 7. Authorization 缺少 Token 黑名单/灰名单

**位置**: `authorization/auth.go`

**问题**: JWT 过期前无法使 token 失效，缺少登出功能

**建议新增**:

```go
func (auth *auth) BlacklistToken(token string, expiration time.Duration) error
func (auth *auth) IsBlacklisted(token string) bool

// 中间件检查黑名单
func AuthWithBlacklist(auth Auth, store Store) gin.HandlerFunc
```

---

### 8. Gen 缺少枚举类型生成

**位置**: `gen/enum.go`

**问题**: 虽然有 enum.tmpl，但缺少配套的 Go 枚举类型生成功能

**建议完善**:

```go
type EnumDef struct {
    Name    string
    Values  []EnumValue
}

type EnumValue struct {
    Name  string
    Value int
}

func GenEnums(enums []*EnumDef) ([]byte, error)
```

---

## 🟢 性能优化建议

### 1. Kafka Producer 批处理优化

**位置**: `mq/kafka.go:112-136`

**问题**: BulkSend 每次都会重新创建 kafka.Message 切片，频繁的内存分配影响性能

**优化建议**:

```go
// 使用对象池减少内存分配
var messagePool = sync.Pool{
    New: func() interface{} {
        return &kafka.Message{}
    },
}

func (m *producer) BulkSend(ctx context.Context, topic string, key string, messages []*Message) error {
    list := make([]kafka.Message, 0, len(messages))
    defer func() {
        for i := range list {
            messagePool.Put(&list[i])
        }
    }()
    // ... 现有逻辑
}
```

---

### 2. 使用 bytes.Buffer 拼接 SQL

**位置**: `persistence/orm/query.go:67-97`

**问题**: strings.Join 会复制所有字符串，对于复杂查询会产生多次内存分配

**优化建议**:

```go
func (b *BoolQuery) Do() (string, []interface{}) {
    var sb strings.Builder
    var args []interface{}

    if b.IsEmpty() {
        b.MustQueries = append(b.MustQueries, Raw("1 = 1"))
    }

    if len(b.MustQueries) > 0 {
        sb.WriteString("(")
        for i, q := range b.MustQueries {
            if i > 0 {
                sb.WriteString(" AND ")
            }
            sql, a := q.Do()
            sb.WriteString(sql)
            args = append(args, a...)
        }
        sb.WriteString(")")
    }
    // ...
    return sb.String(), args
}
```

---

### 3. Channel 缓冲大小优化

**位置**: `app/app.go:50`

**问题**: `quit := make(chan os.Signal, 1)` 信号队列太小，可能在繁忙时丢失信号

**建议**: 增加缓冲区或使用专门的信号处理 goroutine

---

### 4. JWT claims 复制优化

**位置**: `authorization/auth.go:59-70`

**问题**: 每次生成 token 都重新创建 Claims 对象，可以考虑预分配

**优化建议**: 对高频调用的场景实现缓存或复用 Claims

---

## 🟢 代码质量改进

### 1. 统一错误处理

多处代码使用 `errors.New()` 创建错误，建议引入更详细的错误信息：

```go
// 替换为
return fmt.Errorf("failed to create branch: %w", err)
```

这样可以使用 `errors.Is()` 和 `errors.As()` 进行错误断言。

---

### 2. 添加单元测试覆盖率

以下模块缺少测试或测试覆盖不足：

| 模块 | 当前状态 | 优先级 |
|------|----------|--------|
| `authorization/session.go` | 几乎无测试 | P0 |
| `locker/etcd_client.go` | 无测试 | P0 |
| `mq/transaction.go` | 仅基础测试 | P1 |
| `task/job/sql_job.go` | 无测试 | P1 |
| `persistence/etcd/etcd.go` | 无测试 | P2 |

---

### 3. 添加文档注释

以下公共函数缺少 godoc 注释：

- `sliceconv/Mapper` 的所有方法
- `gen/Impl` 的私有方法
- `persistence/cache/redis.go` 的具体实现
- `hook/client.go` 的配置选项

---

### 4. 移除调试代码

文件中存在明显的调试痕迹需要清理：

```bash
# git/client.go:63
defer fmt.Println(string(b))  # 应该使用正式日志

# 其他类似调试代码...
```

---

### 5. 命名规范不一致

```go
// gen/models.go
func (g *Models) Gen(...)          // 首字母大写
func (g *Models) dbSearch(...)     // 私有方法用小写但不一致
func (g *Models) tableName(...)    // 私有工具函数

// 建议统一使用小写 + 下划线风格或者全部小写
```

---

### 6. Magic Number 抽取

多处硬编码魔法值，应该定义为常量：

| 位置 | 魔法值 | 建议常量名 |
|------|--------|------------|
| `authorization/auth.go:44` | `20 * time.Minute` | `DefaultAccessTokenTTL` |
| `authorization/auth.go:45` | `24 * 7 * time.Hour` | `DefaultRefreshTokenTTL` |
| `mq/kafka.go:171` | `3` | `CommitRetryTimes` |
| `persistence/orm/orm.go:109` | `10` | `DefaultMaxIdleConns` |

---

## 🔵 架构设计改进

### 1. Server 抽象不完整

**问题**: `server/server.go` 定义了接口，但 `http.Server` 的实现与 gRPC Server 的结构不完全一致

**建议**:

```go
type Server interface {
    Start() error
    Stop(ctx context.Context) error  // 添加 context 支持超时
    ServiceInstance() *register.ServiceInstance
    Middlewares() []gin.HandlerFunc  // 暴露中间件列表供组合使用
}
```

---

### 2. App 启动顺序耦合

**问题**: `app/app.go` 中服务启动顺序是固定的，无法灵活配置

**建议**: 支持自定义启动钩子：

```go
type AppOption func(*App)

func WithPreStartHook(hook func(context.Context) error) AppOption
func WithPostStartHook(hook func(context.Context) error) AppOption
```

---

### 3. Trace 集成点不足

**问题**: OpenTelemetry 仅在 orm 中有集成，HTTP/RPC 中间件的 trace 集成不完善

**建议**: 在 `server/http` 和 `server/rpc` 中添加统一的 trace 中间件。

---

## 📋 修复优先级

| 优先级 | 问题 | 预计工作量 |
|--------|------|------------|
| P0 | context 使用错误 | 0.5h |
| P0 | 未处理的错误 | 2h |
| P0 | ID 生成器并发安全 | 4h |
| P1 | 缺少错误类型判断 | 0.5h |
| P1 | Store 批量操作 | 4h |
| P1 | Authorization 黑名单 | 2h |
| P2 | ORM 分页查询 | 4h |
| P2 | MQ 事务消息 | 8h |
| P3 | 单元测试补充 | 16h+ |
