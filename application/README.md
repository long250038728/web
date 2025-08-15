## DDD思想
### Service 层：
Service 层负责协调应用程序中不同领域的逻辑。它们是应用程序的入口点，接收来自外部的请求并将它们传递给领域层。Service 层通常包括应用服务（Application Services）和领域服务（Domain Services）。应用服务主要负责协调领域对象以执行应用程序的用例，而领域服务则处理跨领域的业务逻辑。

### Domain 层：
Domain 层包含了应用程序的核心业务逻辑和领域对象。它们是问题域的抽象表示，包括实体（Entities）、值对象（Value Objects）、聚合根（Aggregate Roots）、领域事件（Domain Events）等。Domain 层负责实现业务规则，确保应用程序的行为符合业务需求。

### Repository 层：
Repository 层负责与数据存储进行交互，并将数据持久化到数据库或其他数据存储中。它提供了对数据的访问接口，使得领域层可以独立于具体的数据存储技术。Repository 层通常包括对领域对象进行持久化和检索的接口定义，并提供具体实现来与数据存储进行交互。


## 流程
1. 命令行解析:  获取配置文件根路径
2. 加载配置/初始化基础构件工具:  app.NewUtil()
3. 创建业务组件:  repository.NewRepository => domain.NewDomain => service.NewService => handles.NewHandles
4. 创建App:  组装App实例，提供Start,Stop方法，通过opts把http/grpc(依赖业务组件)，服务注册与发现，链路等应用周期的实例加入其中
5. 运行App:  通过waitGroup方式运行多个server（http/grpc），并对服务注册发现，链路等的实例进行应用
6. 监听退出:  监听退出信号后进行数据的关闭

## 重要思想
1. util工具: 把常规的mysql/redis/mq等在util中创建并持有(单例模式)，依赖注入于repository，保证只有repository才能调用基础工具(domain/service不允许直接调用)
2. app工具: 负责整个应用的生命周期，通过waitGroup并行运行多个server（http/grpc）,服务注册与关闭,openTelemetry链路初始化,同时监听退出信号回收/关闭系统资源。
3. 其他: 对于ES,Mysql由于提供的方法多所以无法使用interface接口的方式进行抽离，其他的尽可能使用接口暴露及使用（为了解耦及可替换性）
4. 依赖注入： 尽可能使用依赖注入的方式，如果链路复杂可使用wire工具进行静态生成