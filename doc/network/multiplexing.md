## 多路复用：
(场景:现在假设有一个标准数据及读取套接字，由于fget跟read是阻塞读取，如果执行了fget的话进行阻塞，要执行read就需要等待fget阻塞完成后才能调用，在多线程的场景下，由于cpu资源优先不可能盲目的开启n个线程去单独阻塞处理数据。那么就需要一个可以同时处理多个阻塞的东西。此时就用到了多路复用)

多路复用的设计初衷就是把标准输入、套接字等都看做I/O的一路。
多路复用的意思就是任何一路I/O有事件发生情况下，通知应用程序去处理对应I/O事件。仿佛在同一时刻可以同时处理多个I/O事件。
如果有标准输入有数据，立即从标准输入读入数据，通过套接字发送出去
如果有套接字有数据可以读，立即读出数据。

I/O事件的类型
1. 标准输入文件描述符准备好可以读
2. 监听套接字准备好，新的连接已经建立连接
3. 已连接的套接字准备好可以写
4. I/O事件超过10秒发生超时


### select
使用select函数，通知内核挂起进程，当一个或多个I/O事件发生后，控制权还给应用程序，有应用程序进行I/O事件的处理

select函数（通过返回值判断是否成功）
1. maxfd则表示有maxfd-1需要处理。
2. 三个集合描述符集合
    * 用一个整型数组表示一个描述集合，一个32位整型数组可以表示32个描述字
    * readset读描述符集合 、 writeset写描述符集合   exceptset异常描述符集合
    * 三个集合描述符集合每一个都可以设置为空，表示不需要内核进行相关检测。
3. timeval
   1. 设置为空，表示没有I/O事件select一直等待下去
   2. 设置一个非0的值，表示在固定一段时间内从阻塞中返回
   3. rv_sec,tv_usec设置为0，表示不等待直接返回（这种情况相对较小）

* FD_ZERO将向量元素设置为0
* FD_SET用来把对应的套接字fd的元素设置为1
* FD_CLR用来把对应的套接字fd的元素设置为0
* FD_ISSET用来把对应的套接字fd的set集合中是否存在(0代表不需要处理，1代表需要处理)



### poll
由于select的最大值是1024（缺点），无法在高并发下处理多个fd，所以设置了poll函数来处理多路复用
>int poll (struct pollfd * fds,unsigned long nfds,int timeout)

1. fds：是一个pollfd类型的结构体数组 struct pollfd { int fd; short  events;  //监听事件（读写） short  revents'  //返回事件}
2. nfds：有多少个fd
3. 阻塞时间：与select类似



### epoll
poll的缺点是返回的数据是全部fd数组的数据，遍历该数组判断revents有值代表有返回，


创建
```
int epoll_create(int size) //size设置一个大于0的整数
int epoll_create1(int flags)
表示epoll实例，返回-1表示出错返回值大于0
```

增加/删除监听事件
```
epoll_ctl
int epoll_ctl(int epfd, int op , int fd ,struct epoll_event* event)
成功返回0，-1失败
epfd：刚刚调用的epoll_create创建的epoll实例描述字
op：操作
    EPOLL_CTL_ADD //新增
    EPOLL_CRL_DEL //删除
    EPOLL_CTL_MOD //修改
fd:需要操作的文件描述符（套接字）
event：注册的事件类型，可以在这个结构体设置用户需要的数据，
    其中最为常见的是使用联合结构里的fd字段，表示事件对应的文件描述符
    typedef union epoll_data {
        void *ptr
        int fd
        uint32_t u32
        uint64_t u64
    }epoll_data_t;

    struct epoll_enve {
        uint32_t events    //event事件
        epoll_data_t data //用户的数据
    }
```

监听（等待）
```
表示调用程序挂起，等待内核io事件分发
int epoll_wati(int epfd,struct epoll_event * events,int maxevents , int timeout)
返回值大于0表示事件个数，返回0表示超时，-1表示出错

epfd：刚刚调用的epoll_create创建的epoll实例描述字
events：返回用户空间需要处理的io事件，数组，数组的大小由epoll_wait返回值决定，
    这个数组每个元素都是需要待处理的io事件。其中events表示具体的事件类型，与epoll_ctl相同，这个epoll_event结构体里
    的data就是epoll_ctl中设置的data，也是用户空间和内核空间调用需要的数据
maxevents：大于0的整数，表示epoll_wait可以返回的最大事件值
timeout：超时时间

```


1. 事件集合。
在每次使用poll或select之前，都需要准备一个刚兴趣的事件集合内核拿到该集合进行分析并在内核控件构建对应的数据结构来完成事件集合的注册。
epoll维护的是一个全局的事件集合，通过epoll句柄可以操纵这个时间的集合，增加跟删除或修改这个时间集合里的某个元素。要知道绝大多数情况下，事件集合的
变化没有那么大，这样操纵系统内核就不需要每次扫描事件集合，构建内核空间数据结构
2. 就绪列表。
每次使用poll或select之后，程序都需要扫描整个感兴趣的事件集合。找出正在活动的事件，每次扫描需要花费系统资源。事实上很多情况扫描完一圈真正活跃的事件只有几个。而epoll只返回活动事件的列表，应用程序减少大量扫描时间。

----

### fd大小的查询及设置
ulimit -n
su vi /etc/sysctl.conf

### 查看每个tcp的发送缓冲区和接收缓冲区
* cat /proc/sys/net/ipv4/tcp_wmem //写
* cat /proc/sys/net/ipv4/tcp_rmem //读
* 4096：最小分配值	16384：默认分配值	4194304：最大分配值


### 各种IO模型
1. 阻塞IO + 进程
每个连接通过fork派生一个子进程进行处理，即便阻塞IO也不会影响到其他IO(效率不高，扩展性差，资源占用高)
accept connections
fork for conneced connection fd
proess_run(fd)


2. 阻塞IO + 线程
通过pthred_create创建单独线程，达到上面进程的效果
accept connections
prhread_create for conneced connection fd
thread_run(fd)


3. 非阻塞IO + readiness notification + 单线程
每次遍历fdset，循环判断是否有时间相应，cpu消耗资源，让操作系统告诉我们哪个套接字准备好可读/写,
在结果返回前，把CPU的控制权交出去，让操作系统把cpu时间调度给其他需要的进程，这就是select、poll。
这样的方法就需要每次dispatch之后，对所有注册的套接字进行排查效率不是那么的高，如果dispatch调用返回的
是提供有IO事件或IO变化的套接字，那效率就高很多，那就是epoll


4. 非阻塞IO + readiness notification + 多线程
所有IO都在一个线程处理，利用cpu多核的能力，让每个核都可以作为一个IO分发器进行IO事件并发
这就是reactor事件。基于epoll/poll/select的IO分发器可以叫reactor，叫做事件驱动/事件轮训


5. 异步IO + 多线程
当调用结束后，请求立即返回，由操作系统后台完成对应的操作。当最终操作完成，产生一个信号或是通过回调处理


### 各IO流程
```
阻塞              读   阻塞   阻塞   阻塞         读到了    
非阻塞            读   轮训   轮训   轮训         读到了  
非阻塞+多路复用    调用多路复用等到数据来了    读    读到了
异步             读   做其他事  数据到了通知      读到了  
非阻塞IO的问题就是需要轮训去查询时间是否到来，结合了多路复用，可以实现数据进行非阻塞同时无需轮训去处理，而是通过内核直到数据到来后去去取
```


### 总结：
由于fget，read，write都是阻塞操作，一个线程只能处理一个命令直到该命令阻塞完成后才能进行下一个命令，那么就把这几个命令通过select，poll函数把多个操作放入该函数中，还函数会在内核中处理，此时该函数会出现阻塞直到至少有一个返回的时候返回。这样就可以不用通过一个命令一个线程/进程，避免线程的频繁切换上下文。而且阻塞等待，大大浪费了cpu资源。多路复用是提高cpu有效资源。

### select:
1. 3个集合，最多只有1024个，遍历所有socket_id，fd少的话遍历的数据少(set)
2. select是传入set集合后返回set集合进行响应，有readset，writeset，exceptset事件，
3. 通过FD_ISSET(socket_fd,&集合))函数判断，该fd是否响应

### poll：
一个数组，没有长度限制,遍历poll数组, fd少的话依旧需要遍历整个数组(list)
poll是传入一个数组，数组中每个对象是结构体(fd,events,revent),通过poll函数返回的后，遍历数组，
如果该该对象中fd >0 且 revent是某个操作(读写等)，对该数据进行操作。

### epoll：
通过改进接口的设计，避免用户态-内核态的频繁数据拷贝。
1. poll与select每次都是传一个集合进去函数中，函数需要遍历传入的集合。epoll则是传一个需要处理的fd
2. poll与select返回的是所有集合，epoll返回的是已经完成io操作的集合



### select
>select(int nfds, fd_set *r, fd_set *w, fd_set *e, struct timeval *timeout)

select它存在4个问题：
1. 这3个bitmap有限制大小（FD_SETSIZE，通常为1024）；
2. 因为这3个集合在返回时会被内核改动，因此我们每次调用时都须要又一次设置；
3. 我们在调用完毕后须要扫描这3个集合才干知道哪些fd的读/写事件发生了
4. 内核在每次调用都须要扫描这3个fd集合，然后查看哪些fd的事件实际发生.


### poll
poll(struct pollfd *fds, int nfds, int timeout)
因为存在这些问题，于是人们对select进行了改进。从而有了poll。
struct pollfd { int fd; short events; short revents; }
poll调用须要传递的是一个pollfd结构的数组。调用返回时结果信息也存放在这个数组里面。
pollfd的结构中存放着fd、
我们对该fd感兴趣的事件(events)
该fd实际发生的事件(revents)。

### poll对select做的改进
1. poll传递的不是固定大小的bitmap，因此select的问题1攻克了。
2. poll将感兴趣事件和实际发生事件分开了，因此select的问题2也攻克了。
3. 但select的问题3和问题4仍然没有解决。

对于select的问题4，我们为什么须要每次调用都传递全量的fd呢？内核可不能够在第一次调用的时候记录这些fd，然后我们在以后的调用中不须要再传这些fd呢？
对于每一次系统调用，内核不会记录下不论什么信息。所以每次调用都须要反复传递同样信息。


### epoll和kqueue。
上帝说要有状态。所以我们有了epoll和kqueue。
```
int epoll_create(int size);
int epoll_ctl(int epfd, int op, int fd, struct epoll_event *event);
int epoll_wait(int epfd, struct epoll_event *events, int maxevents, int timeout);
epoll_create的作用是创建一个context，这个context相当于状态保存者的概念。
epoll_ctl的作用是，当你对一个新的fd的读/写事件感兴趣时，通过该调用将fd与对应的感兴趣事件更新到context中。
epoll_wait的作用是，等待context中fd的事件发生。
```

epoll是Linux中的实现，kqueue则是在FreeBSD的实现。
```
int kqueue(void);
int kevent(int kq, const struct kevent *changelist, int nchanges, struct kevent *eventlist, int nevents, const struct timespec *timeout);
与epoll同样的是，kqueue创建一个context。与epoll不同的是。kqueue用kevent取代了epoll_ctl和epoll_wait。
epoll和kqueue攻克了select存在的问题。通过它们，我们能够高效的通过系统调用来获取多个套接字的读/写事件，从而解决一个线程处理多个连接的问题。
```
