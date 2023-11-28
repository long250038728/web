package mu

import (
	"context"
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync"
)

type mu struct {
}

//============================== 读写锁 =====================================

func ReadWriteLock(key string) string {
	//读写锁
	var rwm sync.RWMutex

	//加锁 / 解锁
	rwm.RLocker()
	defer rwm.RUnlock()

	//返回
	return "hello"
}

//============================== 单例 =====================================

// 单例
var once sync.Once
var m *mu

func singleton() *mu {
	once.Do(func() {
		m = &mu{}
	})
	return m
}

//============================== 池化 =====================================

func pool() {
	syncPool := sync.Pool{
		New: func() any {
			return "str" //减少内存分配及垃圾回收
		},
	}
	data := syncPool.Get().(string)
	syncPool.Put(data)
}

//============================== 协程等待 =====================================

func wait() {
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		//处理
	}()
	go func() {
		defer wg.Done()
		//处理
	}()
	go func() {
		defer wg.Done()
		//处理
	}()
	wg.Wait()
}

// ============== chan 一般都用select default(解决队列满等待) ================

// 实现一个简单的消息队列
func channel(ctx context.Context, data bool) {
	var c = make(chan bool)

	//插入数据如果要遍历多个消费者，这时候就加上读锁
	//加入消费者要往数组添加，，这时候就加上写锁
	//ctx是解决两个问题  1.写入超时的问题  2.消费时for退出的问题
	//close(c) 是为了把所有监听的go协程都退出，如果往里面只发一个的话，多个协程只有一个关闭

	//生产者
	select {
	case c <- data:
	case <-ctx.Done():
		//这个时候可以解决<-chan 阻塞一直释放不了的问题
	}

	//消费者
	for {
		select {
		case <-ctx.Done():
			return //这个可以解决外部退出的问题
		case chData, ok := <-c:
			if !ok {
				return //这个可以解决外部退出的问题
			}
			fmt.Println(chData)
		}
	}

	//结束close
}

func sf() {
	sf := singleflight.Group{}
	data, err, share := sf.Do("hello", func() (interface{}, error) {
		return "", nil
	})

	fmt.Println(data)
	fmt.Println(err)
	fmt.Println(share)
}
