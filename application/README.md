
Service 层：
    Service 层负责协调应用程序中不同领域的逻辑。它们是应用程序的入口点，接收来自外部的请求并将它们传递给领域层。Service 层通常包括应用服务（Application Services）和领域服务（Domain Services）。应用服务主要负责协调领域对象以执行应用程序的用例，而领域服务则处理跨领域的业务逻辑。

Domain 层：
    Domain 层包含了应用程序的核心业务逻辑和领域对象。它们是问题域的抽象表示，包括实体（Entities）、值对象（Value Objects）、聚合根（Aggregate Roots）、领域事件（Domain Events）等。Domain 层负责实现业务规则，确保应用程序的行为符合业务需求。

Repository 层：
    Repository 层负责与数据存储进行交互，并将数据持久化到数据库或其他数据存储中。它提供了对数据的访问接口，使得领域层可以独立于具体的数据存储技术。Repository 层通常包括对领域对象进行持久化和检索的接口定义，并提供具体实现来与数据存储进行交互。