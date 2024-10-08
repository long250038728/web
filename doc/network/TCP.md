### 学习网络编程，掌握两个核心
1. 就是理解网络协议，并在这个基础上和操作系统内核配合，感知各种网络 I/O 事件；
2. 就是学会使用线程处理并发。

OSI & TCP/IP
阿帕网前期非常简陋致力开发下一代协议，此时OSI组织游说众多厂商一起开发一个标准，这样大家遵循这个标准，那么大家就有钱赚。但是等到1984年把这个标准发布的时候却发现满世界都在用一个叫TCP/IP的协议栈东西。
```
OSI：
OSI组织游说众多厂商我们一起定义一个网络互连的标准。把事情做大，我们就有钱赚了。搞出了一个非常强悍的标准OSI参考模型。OSI模型过于复杂，没有参考实现在一定程度上阻碍了普及，但是教科书般的层次模型。对后世影响深远
TCP/IP: 
在发布OSI标准的时候已经是1984年，OSI组织发现满世界都在用一个TCP/IP协议栈的东西。TCP/IP解决了实际上的问题，在实际中不断的完善。
```

```
OSI：           TCP/IP协议栈
应用层            应用层
表示层
会话层
传输层            传输层
网络层            网际层
数据链路层        网络接口层
物理层
```

### 套接字对(是linux具体的实现对象/函数)
一个连接可以通过客户端-服务器端的IP及端口号确定，这叫套接字对


### 区别TCP/UDP
* TCP（字节流套接字）中的连接是谁发起的 ***连接
可靠的、双向连接的通讯串流，以1-2-3输出到套接字上，另一端一定会以1-2-3的顺序抵达。通过连接管理，阻塞控制，数据流与窗口管理，超时和重传等一系列设计
* UDP（数据报套接字）中的报文是谁发送的 ***报文
使用UDP的原因是速度，通过广播及多播技术向网络中多个节点同时发送信息。在丢失一两个丢失不会造成多大问题。


### 服务器socket
1. 服务器创建socket后进行bind绑定到一个地址及端口上
2. 执行listen操作把原先的socket转化为服务器的socket
3. 服务端最后阻塞在accept上等待客户端请求


### 客户端socket
1. 初始化socket
2. 执行connect向服务器的地址及端口发现请求
3. 请求的过程使用了三次握手，握手成功后进行数据的发送、响应及四次挥手后断开


### socket请求（一旦建立连接，传输数据是双向的）
1. 客户端进程向操作系统内核发起write字节流写操作
2. 内核协议栈将字节流通过网络设备传输到服务器端
3. 服务器端从内核得到信息，将字节流从内核读取到进程中
4. 处理业务逻辑后把结果以同样的方式写给客户端

```
### 创建socket
int socket(int domain,int type,int protocal)
domain:指定是ipv4，ipv6，本地等（枚举值）
type:  指定tcp，upd,原始套接字（枚举值）
protocal：0 指定通讯协议现在已经废弃

### 创建bind(注意，这里不生成返回值)
bind(int fd,sockaddr* addr,socklen_t len)
fd: socket生成的fd
addr: ipv4，ipv6套接字结构体（结构体，里面包含地址跟端口号，端口号传0表示不指定代表生产随机数，一般用于客户端）
len: 地址长度

### 创建listen
int listen(int socket,int backlog)
socket: socket生成的fd
backlog:表示已完成且未accept的队列大小，决定了接受并发的并发数

### 建立连接accept
int accept(int listensockfd,struct sockaddr * cliaddr,socklen_t * addrlen)
listensockfd:listen创建的套接字
cliaddr：客户端的地址等信息（结构体）
addrlen：地址长度
返回一个新的fd，是因为accept进行连接后，生成新的fd之后处理的事项都是新的fd，如果用listensockfd，那此时fd在处理事项阻塞其他请求则无法接收


### 客户端连接connect
int connect(int sockfd,const struct sockaddr * servaddr,socklen addrlen)
sockfd:创建socket生成的fd
servaddr：服务器的地址等信息（结构体）
addrlen：地址长度
```


调用connect会激发TCP三次握手
客户端                      服务器

       --     发送j    ->

       <- 响应j+1，发送k --

       --    响应k+1   ->


TCP发送缓冲区
1.在建立tcp三次握手后操作系统会创建配套的基础设施如发送缓存区
发送缓冲区的大小可以通过套接字选择来改变，当我们调用write函数时，实际上做的是把数据从应用程序中拷贝到操作系统内核的发送缓冲区，并不是把数据通过套接字发送出去。
当发送缓冲区不够大，应用程序会被阻塞，缓冲区就像流水线会不断的取出数据,按照TCP/IP的语义，将数据封装成TCP的MSS包，以及IP的MTU包，最后通过数据链路层将数据发送出去。这样发送缓冲区就空了一部分，于是可以继续往发送缓冲区插入数据，直到所有数据都发送成功，write才会停止阻塞返回。注意返回成功只能表示数据全部发送到发送缓冲区，不表示对端成功接收

TCP读取数据
read函数要求操作系统内核从套接字描述子socketfd读取最多多少字节，并把结果存储到buffer中，返回值告诉我们实际的读取字节数目。对于read来说，需要循环读取数据，并且需要考虑EOF等异常条件
如果read函数返回值为0表示EOF,表示要处理断连。。
如果read函数返回时-1，表示出错。


UDP(无上下文)
到达顺序不保证，能不能收到也不保证（没有重传、确认。有序传输及拥塞控制等能力），优点在于简单，对延迟丢包等不是特别敏感在多人通信的场景都是用UDP协议
recvfrom函数等待客户端报文的发送，客户端调用sendto函数往目标地址及端口发送UDP报文


工具
ping 帮助我们进行网络连通的探测
ifconfig 显示当前系统的所有网络设备
netstat和lsof可以查看活动连接状态
tcpdump对奇怪环境的抓包，了解报文


TIME_WAIT
在即将断开时会经历四次挥手的状态，当发起方发起连接断开的一段时间会处于TIME_WAIT状态
发起方                                                         接收方

       --    发送j,我这边没有数据传了我先进入半断开状态    ->
       <-   响应j+1，好的没问题，我接收到了,你进入半断开把  --

              此时发起方进入TIME_WAIT状态，是因为在
              发起方这个时候是把数据发出去了，但是接收方
              并不一定全部接收，有可能还需要重传，如果我这个时候
              全断开了，接收方要我重传我都断开了
                        等待中
                        等待中
                        等待中
                         ...

         <-   发送k，好了没问题了我全部都接收完，不用重传了  --
         --       响应k+1好的，那我进入全断开了          ->


### 疑惑梳理
TCP: 
1. 由于http2.x支持KeepAlive多路复用的。此时第一次发起HTTP请求时建立TCP三次握手，发送完结束后一般不会立即关闭断开TCP连接所以无四次挥手（只是保持连接活跃而已，并不意味着连接不会关闭。还要看其他的条件）
2. HTTP 协议使用 Content-Length 头部或者分块传输编码来告知接收方数据的结束点，TCP 负责数据的可靠传输，但并不直接管理 HTTP 数据的结束
3. 关闭TCP的条件 
   * 客户端主动关闭
   * 服务端空闲时间超出释放连接。服务器资源不够需要释放资源
4. 四次挥手的目的：为了优雅地关闭 TCP 连接。它的核心作用是确保双方都能够完成数据的发送和接收，避免数据丢失。
5. 会把数据拆分成多个段，每个段都有一个序列号(有序的)。当接收方收到数据后发送"确认包"表示已经接收。当收到了FIN的段表示该请求最大的数据发送完毕，如果中间哪个段没收到会再发送请求重发
6. 由于在发送过程中接收方收到的段是乱序的（发送方是顺序的。网络原因可能导致接收方乱序）。收到后在发送"确认包"时告知下一个期望的序列号。
7. 发送一个包会记录发送时间，当超过超时时间也会进行重发
8. 段头包含源端口、目标端口、序列号、确认号、数据偏移、标志位、窗口大小、校验和、紧急指针等字段
9. 加入有段1，2，3，4 此时收到了段1后收到了段3，会把（段1和段3）缓存到缓冲区中。TCP协议本身会保证数据的顺序，等到全部段到达后重组交给应用层
10. 
