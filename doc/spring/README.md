### 基础组件
* `Nacos`: 服务注册与发现 && 配置中心
* `OpenFeign`: 远程调用
* `Sentinel`: 容错
* `Seata`: 分布式事务


### 环境准备
1. java安装: OpenJDK
2. Maven安装
3. IntelliJ IDEA
4. Lombok(提高编程效率和代码可读性)在插件中


### spring boot / spring cloud 总结
1. IOC: 控制反转
   * 通过注解的方式在项目启动时把对应注解的加入到bean容器中统一管理（无需在代码中使用new）
   * 通过bean容器进行管理对象的生命周期
   * 解耦：调用方无需关心依赖的创建细节。
   * 优势：无需管理各个对象之间的关系，在bean中会自动new对应依赖关系的对象（可以通过注解指定是普通对象还是单例对象）
   * 优势：在改变对象中的依赖时无需手动的去添加依赖的创建
2. AOP: 面向切面编程
    * 把一些常用的功能封装成注解统一对外提供服务
    * 通过注解的方式让原先的类\方法添加额外的功能
    * 添加功能不影响原先业务的代码，只需添加对应注解（从业务逻辑中分离出来）
    * 把服务注册发现/日志/事务/远程调用等通过注解的方式添加，对函数内部无任何修改及影响
3. 自动装配
    * 通过自动装配功能，给对象/函数添加注解，就可以自动实现对应的功能。让业务方无需过多考虑项目的配置文件加载,跨服务函数调用，服务注册与发现等


### spring cloud 使用
* `Nacos`: 服务注册与发现
    1. 对`pom.xml`添加依赖（nacos-discovery包）
    2. 在`配置文件`application.yaml添加`nacos`的配置信息
    3. 在项目中添加注解`@EnableDiscoveryClient`（老版本需要、新版本可以不添加）
    4. 程序启动时会根据application.yaml中的`spring.application.name`参数及`nacos相关`参数把服务注册到服务注册中心
* `Nacos`: 配置中心
    1. 对`pom.xml`添加依赖（nacos-config包）
    2. 创建bootstrap.yaml配置添加nacos配置信息（比application.yaml 加载更前）
    3. 代码无需添加启动时就会自动加载（如加载到application.yaml效果一样）
    4. 代表变量中可通过@Value("${xxxxxx:默认值}")注入。同时需要对对应的类添加@RefreshScope注解
* `OpenFeign`: 远程调用
    1. 对`pom.xml`添加依赖
    2. 服务端：无需有任何代码的修改
    3. 客户端： 创建一个接口添加注解`@FeignClient(value = "coupon-calculation-serv", path = "/calculator")`
        * value: application.yaml中的`spring.application.name`
        * path: `@RequestMapping("calculator")`
    4. 对客户端Application添加注解`@EnableFeignClients(basePackages = {"com.xxxx"})`
        * basePackages: 启动项目的时找到所有位于com.xxx包下加载FeignClient修饰的接口
    5. 其他
        * 添加日志 (application.yaml中配置）
        * 设置超时  (application.yaml中配置）
        * 服务降级hystrix组件（pom.xml添加依赖。创建fallback 或 fallback工厂，添加到FeignClient中）
* `Sleuth`链路（生成traceId和在其他组件中自动埋点）——主要对接Zipkin
    1. 对`pom.xml`添加依赖（sleuth包）
    2. 在`配置文件`application.yaml添加`sleuth`的配置信息（采样率及每秒最多采样多少）
    3. 在项目中会自动埋点（rabbitMQ,Kafka,OpenFeign,Controller等） 
* `Micrometer Tracing`新版本链路——打通spring组件（否则就要在各个组件写拦截器）
    1. 对`pom.xml`添加依赖（micrometer、opentelemetry包）
    2. 对application.yaml添加配置 
    3. 在项目中会自动埋点
    4. 在业务中埋点如mybatis-plus埋点需要添加拦截器，在业务中埋点可以通过`Span span = tracer.nextSpan().name("createOrder").start(); span.End()`
* `Gateway`
    1. 对`pom.xml`添加依赖（gateway包）
    2. 在`配置文件`application.yaml中添加服务路由（路由地址、断言、过滤器（全局和自定义过滤器））
* `Stream`消息队列
    1. 依赖 (pom.xml spring-cloud-starter-stream-rabbit包)
    2. 生产者`@autowired private StreamBridge streamBridge;  streamBridge.send("binding名称", "data")`
        * binding名称 格式 topicName-out-index   topicName主题消息，out生产者，index编号(默认是0) 这个主要用于配置中有绑定bander
    3. 消费者`@Bean public topicName<String> input() { return message -> { // doSomething(); }}`
        * 方法名就是主题名称
    4. 配置文件application.yaml
        * spring.cloud.stream.binders.rabbit.type=rabbit     //添加名为rabbit的binder
        * spring.cloud.stream.binders.rabbit.environment.spring.rabbitmq.host
        * spring.cloud.stream.binders.rabbit.environment.spring.rabbitmq.port
        * spring.cloud.stream.binders.rabbit.environment.spring.rabbitmq.username
        * spring.cloud.bindings.topic-out-0.destination=topicName  //生产者
        * spring.cloud.bindings.topic-out-0.binder=rabbit
        * spring.cloud.bindings.topic-in-0.destination=topicName  //消费者
        * spring.cloud.bindings.topic-in-0.binder=rabbit
        * spring.cloud.bindings.topic-in-0.group=consumerGroupA 
* `Seata`: 分布式事务
    1. 对`pom.xml`添加依赖（seata包）
    2. 搭配mysql数据库
    3. 配置文件application.yaml添加seata的相关配置
    4. nacos配置文件中添加seata配置
    5. 代码中无需增加任何代码
    6. 项目中引入seata包（pom.xml中添加）后 会创建两个bean SeataDataSourceProxy 和 SeataSqlSessionFactoryBean
    7. 基本原理 Seata会创建代理事务处理器，通过代理事务处理器对sql进行改造。
        * 框架会根据配置好的规则进行扫描，并创建全局锁。



### 补充
分布式事务可以简化为“全局事务”和“分支事务”，分支事务全部完成就代表全局事务完成，分支事务只要有一个回滚，代表全局事务失败需要全部回滚（最终一致性）
seata 担任协调者的角色manager，记录各个“分支事务的执行状态”，"锁状态"，同时记录“undo log用于回滚”（AT模式自动生成）（SAGE则是手动添加undo log语句）




### domo
OpenFeign
```客户端java
@FeignClient(value = "coupon-calculation-serv", path = "/calculator")
public interface CalculationService {

    // 优惠券结算
    @PostMapping("/checkout")
    ShoppingCart checkout(ShoppingCart settlement);

    // 优惠券列表挨个试算
    // 给客户提示每个可用券的优惠额度，帮助挑选
    @PostMapping("/simulate")
    SimulationResponse simulate(SimulationOrder simulator);
}

 @Autowired
 private CalculationService templateService;  //感觉像调用本地方法一样调用openFeign

 templateService.checkout(xxxxx)
```



### spring boot/Cloud 与其他语言的对比
* 服务注册
  * spring boot/Cloud： 可通过注解(声明式)的方式进行服务的注册与发现。
  * 在其他语言：手动创建一个服务于注册的对象，然后手动调用服务注册的方法。
* 远程调用
  * OpenFeign 创建一个接口并添加@FeignClient注解，此时调用该接口方法即可实现远程调用（Feign 会自动生成动态代理，完成 HTTP 请求发送、参数序列化、响应反序列化等）
  * 在其他语言：手动创建一个http客户端或grpc客户端，虽然满足了远程调用，但是在代码中并不像调用本地函数一样（需要显式处理服务地址、路径、参数编码、请求发送、响应解析等步骤）
  * 在spring中使用OpenFeign进行远程调用（http），而不是使用GRPC。
* Sleuth/Micrometer Tracing/Stream
  * 目的：在实际的库前加一层中间层。简化引入困难
  * Sleuth/Micrometer Tracing 无需关心使用什么链路（通过application.yaml配置 + 添加依赖），同时对原生组件进行自动埋点，无需再添加代码
  * Stream 无需关心使用的是哪个消息队列（通过application.yaml配置 + 添加依赖），同时消费者监听提交等逻辑都封装在内部
* spring boot/Cloud 劣势
  * 由于是自动装配 + 隐式注入，很多逻辑是 Spring 帮你做的，初学者看不到过程，不知道“是在哪里调用的”、“数据怎么传的”、“链路是怎么创建的”。想做自定义改造时没入口。无法像其他语言那样 new 出来（显式调用），从调用入口推到实现细节。
  * 接入简单，少写代码，但是对某个功能需要自定义/扩展时就无从下手
