长短token

1.登录时服务器返回长短两个 token:
    客户端在登录成功后，服务器返回长短两个 token。长 token 通常具有较长的有效期，而短 token 则具有较短的有效期。
2.客户端每次请求时携带短 token:
    客户端在每次请求时将短 token 放置在请求的 Header 或其他适当的位置上，以便服务器进行验证。
3.服务器判断短 token 失效时返回特殊错误码:
    当服务器检测到短 token 失效时，它返回一个特殊的错误码（如 401 Unauthorized），提示客户端该 token 已失效。
4.客户端中间件拦截错误码，请求长 token 获取短 token:
    客户端可以设置一个拦截器或中间件来检查每个请求的响应。当客户端收到特殊的错误码时，它会触发请求长 token 获取新的短 token 的接口。
    此时服务器也会返回长token（如果长token失效会返回新的，如果未失效则与之前相同） 客户端需要保存着两个token
5.客户端保存短 token:
    客户端收到新的短 token 后，将其保存起来，以备将来的请求使用。
6.再次发起之前的接口请求:
    当客户端收到新的短 token 后，它可以再次发起之前因为 token 失效而失败的接口请求。这次请求将会携带新的有效 token，以确保请求能够成功执行。


为什么要使用JWT

    JWT数据存储于客户端： 
        JSON Web Tokens（JWT）是一种轻量级的身份验证和授权方法，它将用户信息编码为 JSON 格式并使用数字签名进行验证。
        JWT 的一个主要优势是，它将用户的认证状态和相关信息存储在客户端，而不是在服务器端。这意味着服务器不需要存储任何会话信息或用户状态，
        从而减轻了服务器的负担，并且消除了在集群或多台服务器之间同步会话信息的需求。
    JWT长短token的风险
        存放在客户端就有可能遭到窃取加大了安全风险，设置一个较短的有效期限，避免他们将有更多时间来滥用该令牌（短token）
        长token只会在短token失效的时候才会传输，减少了窃取的机会。同时也避免短token失效后需要重新登陆的问题
    短token存储的内容
        用户 ID、用户名等，以及一些访问控制信息。可以避免频繁的数据库查询，从而提高了性能和效率。
    长token存储的内容
        包含一些敏感且相对稳定的信息，例如用于生成新 token 的密钥或令牌刷新令牌（refresh token)所以需要较高的安全性
    短token的时间不宜设置过长     
        除了传输的风险外，服务器端获取一些基本的信息都在短token中获取，当用户信息发生变化时例如用户修改了个人资料。如果过期较长则会有较长的时间获取旧数据。



JWT（无状态，客户端存储）

    优点：
        减少服务器存储需求：因为token存储在客户端，服务器不需要存储任何会话数据。
        水平扩展性：JWT适用于分布式系统，无需在服务器间同步状态。
        跨域支持：JWT可以轻松地在不同域之间传递，适合单点登录（SSO）场景。
        自包含：token包含了所有必要的信息，不需要额外查询数据库。
    缺点：
        安全性：token若被截获，未加密或签名不正确的情况下，信息可能会被篡改。
        性能：大型token可能会增加请求的大小，影响性能。
        状态管理：无法在服务器端控制token中的数据，除非用户登出。
Redis（有状态，服务器存储）

    优点：
        安全性：所有敏感信息存储在服务器端，减少了token被截获的风险。
        灵活性：可以灵活地修改用户的状态，如更新权限或用户信息。
        性能：Redis是一个高性能的内存数据库，适合存储频繁访问的数据。
        控制：可以更容易地实现会话超时、强制登出等控制逻辑。
    缺点：
        服务器负载：需要在服务器端存储会话信息，可能会增加服务器的负载。
        复杂性：需要管理Redis实例，包括部署、扩展和维护。
        分布式挑战：在分布式系统中，需要额外的逻辑来处理会话共享或复制。
决定因素：

    安全性需求：如果安全性是首要考虑，可能倾向于使用Redis。
    分布式系统：如果你的应用程序是分布式的，或者打算使用微服务架构，JWT可能更合适。
    性能要求：如果会话信息非常大，使用JWT可能会影响性能；使用Redis则可以减轻HTTP请求的大小。
    开发和维护成本：使用Redis需要额外的基础设施管理，而JWT更简单。
    更新频率：如果用户信息经常变化，Redis提供了更方便的方式来更新这些信息。
    单点登录：如果需要跨域认证，JWT可能是更好的选择。（多个下游服务一般不共享第三方存储）
                

OAuth2.0

    需要授权第三方应用有对应的访问权限。假如：使用一个第三方软件对淘宝/抖音平台进行数据获取分析等。但不可能直接把账号密码给到第三方，需要通过一个访问令牌的方式。
    第三方获取淘宝/抖音数据时需要客户进行授权给第三方。第三方才有权限去获取数据
```
    第一次交互： （客户与第三方的交互）
    客户-->第三方    :  我这边需要你去调用淘宝/抖音获取数据                                                     
    第三方-->客户    :  好的没有问题，我的身份证号是xxxxx, 我会引导你到淘宝/抖音那边去，然后点击授权就可以了         
    
    第二次交互： （客户与淘宝/抖音的交互）
    客户-->淘宝/抖音  :  （点击授权按钮）我现在需要授权给身份证号是xxxxx的这个人。我已确认授权
    淘宝/抖音-->客户  :  收到，我现在生成一个授权码给你，你发给刚才提供第三方的地址上 （授权码在回调地址上）           

    第三次交互： （第三方、淘宝/抖音的交互）
    淘宝/抖音-->第三方  :  这是授权码，我发给你，你拿这个授权码去申请访问令牌 （回调单向）[96973a6f5637fb3d1049f6d456702932.webp](..%2F..%2F..%2F96973a6f5637fb3d1049f6d456702932.webp)
    第三方-->淘宝/抖音  :   这是给我的授权码，我现在需要申请访问令牌
    淘宝/抖音-->第三方  :   好的没问题，这个访问令牌个给你
    第三方-->客户       :   已经完成了授权

```

分析

    * 第一次交互:  返回三个重要信息。（你去哪里申请，申请授权给谁,授权成功如何获取）
        1. 淘宝/抖音的授权页面
        2. 第三方的信息(身份证号是xxxxx)
        3. 授权后回调的地址（不直接http响应是为了验证地址是否第三方注册时提供的合法不）
    * 第二次交互:  
        1. 客户确定授权
    * 第三次交互
        1. 第三方前端通过回调拿到授权码
        2. 第三方把授权码去淘宝/抖音获取访问令牌
        3. 淘宝/抖音返回给第三方，第三方告诉客户申请通过

目的:

    1. 授权的发起及同意是在客户手上决定的（第三方只是提供了授权页面，第三方信息及回调地址）
    2. 为什么需要回调地址（当客户在授权页面点击授权后，如果没有回调重定向的话就一直停留在这个页面。如果只是返回上一页也不合理）所以重定向顺便带上授权码
    2. 返回授权码给第三方不直接返回访问令牌是为了需要让第三方主动的获取。从安全的方面考虑（只认第三方主动调的，以后也是第三方调用，同时令牌不应该暴露在浏览器）
    3. 第三方拿到授权码去申请访问令牌，淘宝/抖音会校验授权码跟第三方的信息，他们一致我才会颁发授权令牌

其他：

    1. 第三方需要提前在淘宝/抖音进行注册（回调域名，获取权限列表），获取app_id及app_secret
    2. 淘宝/抖音授权页需要带上几个信息 app_id, open_id,回调地址，获取权限范围，通过生成code
        * app_id及回调地址是为了验证第三方是否合法
        * open_id获取客户信息
    3. 第三方通过code，app_id及app_secret 去获取访问令牌，会校验上面的三个值，然后获取code内容生成，同时生成刷新令牌
    4. 刷新令牌refresh_token是为了访问令牌一般过期时间为10分钟，超过后如果再次让客户授权不合理则可以通过刷新令牌再次获取
    5. 刷新令牌使用后是否有效，是否有有效期，超过有效期是续还是重新授权这个是根据不同平台的规则



微信第三方授权流程：
```
    1. 获取component_verify_ticket（通过创建应用时的回到地址，每10分钟会主动推送，有效期12h）
    2. 获取component_access_token （有效期2h）
        POST https://api.weixin.qq.com/cgi-bin/component/api_component_token
            component_appid	string	是	第三方平台 appid
            component_appsecret	string	是	第三方平台 appsecret
            component_verify_ticket	string	是	微信后台推送的 ticket
    3. 获取pre_auth_code （获取生成授权码获连接需要pre_auth_code）及前端/应用跳到微信授权页参数
            POST https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token=COMPONENT_ACCESS_TOKEN
                component_access_token	string	是	第三方平台component_access_token，不是authorizer_access_token
                component_appid	string	是	第三方平台 appid
    4. 生成微信授权页参数返回到前端（或前端自己拼装url）
            component_appid	是	第三方平台方 appid
            pre_auth_code	是	预授权码
            redirect_uri	是	- 授权回调 URI(填写格式为https://xxx)。（插件版无该参数）
            auth_type	    是	- 要授权的账号类型，即商家点击授权链接或者扫了授权码之后，展示在用户手机端的授权账号类型。（小程序/微信公众号等）
            biz_appid	    否	- 指定授权唯一的小程序或公众号 。
            category_id_list	否	- 指定的权限集id列表，如果不指定，则默认拉取当前第三方账号已经全网发布的权限集列表。
    5. 用户点击授权后会回调传递authorization_code参数，通过authorization_code换authorizer_access_token/authorizer_refresh_token
            POST https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token=COMPONENT_ACCESS_TOKEN
                component_access_token	string	是	第三方平台component_access_token，不是authorizer_access_token
                component_appid	string	是	第三方平台 appid
                authorization_code	string	是	授权码, 会在授权成功时返回给第三方平台，详见第三方平台授权流程说明
    6. 当authorizer_access_token过期需要authorizer_refresh_token获取最新的authorizer_access_token
            POST https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token?component_access_token=COMPONENT_ACCESS_TOKEN
                component_access_token	string	是	第三方平台component_access_token
                component_appid	string	是	第三方平台 appid
                authorizer_appid	string	是	授权方 appid
                authorizer_refresh_token	string	是	刷新令牌，获取授权信息时得到
```
微信第三方授权流程总结：

    1. 获取第三方相关的每个接口都需要带上component_access_token、component_appid参数
    2. 第1，2第三方应用前准备工作     第3,4为授权前的准备工作   第5，6是获取authorizer_access_token及更新authorizer_refresh_token
    3. component中 component_verify_ticket在回调获取，component_access_token两个小时有效接口获取一次
    4. authorizer中 authorizer_refresh_token不会过期,authorizer_access_token两个小时有效接口获取一次
    component_access_token， 过期通过component_verify_ticket请求接口获取 （第一次通过也是component_verify_ticket请求接口获取）
    authorizer_access_token，过期通过authorizer_refresh_token请求接口获取 （第一次通过authorization_code获取）






















