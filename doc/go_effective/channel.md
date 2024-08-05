### Chan的概念
chan是使用CSP模型(通信顺序进程)。表示goroutine不要通过共享内存的方式通信，而是要通过 Channel 通信的方式分享数据。
> 虽然chan是一个非常重要的思想，但是对于不熟悉底层原理会非常容易产生panic及内存泄露。所以使用起来需要额外注意

### channel使用场景:
关键字是传递（从A传递到B）
* 数据交流: 多个goroutine可以当成生产者和消费者的关系
* 数据传递: 一个goroutine把数据交给另外一个
* 信号号通知: 一个goroutine把信号交给另外一个
* 任务编排: 一组goroutine按照一定顺序并发或串行执行
* 锁: 利用channel实现互斥锁

### channel panic的情况:
* close 为nil的chan
* send 已经close的chan
* close 已经close的chan

### channel 内存泄露的情况:
* 生产者
  * chan如果设置为空buffer时由于已经没有消费者，(生产者就会一直阻塞无法写入导致生产者内存泄露 )
* 消费者（select 多个chan，当其中有一个chan到达，其他chan未到达时，就会导致内存泄露）
  * time.After需要等到时间到了才会释放
  * 生产者已经结束了没有进行close (消费者就会一直阻塞无法等待数据导致消费者内存泄露 )

### channel使用注意事项:
* close方法应该在生产者处理（在明确知道已经不会往chan写入数据进行close）
* 生产者及消费者的内存泄露