package mu

import (
	"sync"
)

// GMP 还有 饥饿模式，正常模式

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
