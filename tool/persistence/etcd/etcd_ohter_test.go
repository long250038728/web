package etcd

import (
	"context"
	"fmt"
	"github.com/long250038728/web/tool/configurator"
	etcdClient "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
	"log"
	"sync"
	"testing"
	"time"
)

func TestEtcd(t *testing.T) {
	ctx := context.Background()
	var centerConfig Config
	configurator.NewYaml().MustLoadConfigPath("center.yaml", &centerConfig)
	client, err := etcdClient.New(etcdClient.Config{
		Endpoints:   []string{centerConfig.Address},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		_ = client.Close()
	}()

	// 创建watch监听
	t.Run("watch", func(t *testing.T) {
		go func() {
			ch := client.Watch(ctx, "/my-election", etcdClient.WithPrefix()) //监听key的变化(包括leader)
			for msg := range ch {
				for _, even := range msg.Events {
					t.Log(even)
				}
			}
		}()
	})

	t.Run("Lease", func(t *testing.T) {
		// 创建一个租约，有效期为1秒
		lease, _ := client.Grant(ctx, 1)

		// 查看租约的ttl及过期时间
		t.Log(client.TimeToLive(ctx, lease.ID))

		// 修改key 并对这个key添加租约（ttl）如果无租约则代表永久有效
		_, _ = client.Put(ctx, "key", "val", etcdClient.WithLease(lease.ID))

		// 保持租约活跃 —————— 继续使用租约进行其他操作
		resCh, _ := client.KeepAlive(ctx, lease.ID)
		go func() {
			for chResp := range resCh {
				t.Log(fmt.Sprintf("KeepAlive response: %v", chResp))
			}
		}()

		//移除租约(此时key已经没有ttl了，key会被删除)
		t.Log(client.Revoke(ctx, lease.ID))
	})

	session, _ := concurrency.NewSession(client)
	defer func() {
		_ = session.Close() //当session close时Campaign就会取消成为leader（租约（lease）会被撤销）
	}()

	t.Run("election", func(t *testing.T) {
		election := concurrency.NewElection(session, "/my-election")

		// 	txn := client.Txn(ctx).If(v3.Compare(v3.CreateRevision(k), "=", 0))
		//	txn = txn.Then(v3.OpPut(k, val, v3.WithLease(s.Lease())))
		//	txn = txn.Else(v3.OpGet(k))
		//	resp, err := txn.Commit()
		//
		//  封装了一层，通过txn事务。如果这个k(pfx)的Revision == 0 则修改为自己 （最先能写入的则代表获取成功）
		if err := election.Campaign(ctx, "candidate-1"); err != nil { //竞选leader （如果竞选不到则会阻塞等待）
			t.Log(err)
			return
		}

		//	cmp := v3.Compare(v3.CreateRevision(e.leaderKey), "=", e.leaderRev)
		//	txn := client.Txn(ctx).If(cmp)
		//	txn = txn.Then(v3.OpPut(e.leaderKey, val, v3.WithLease(e.leaderSession.Lease())))
		//	resp, terr := txn.Commit()

		//  封装了一层，通过txn事务。如果如果这个k(pfx)的Revision == 自己 则修改
		if err := election.Proclaim(ctx, "candidate-2"); err != nil { //更新leader中的val（当前节点不是领导者，调用它会导致错误）
			t.Log(err)
			return
		}

		t.Log(election.Leader(ctx))
		t.Log(election.Key(), election.Rev())

		//  cmp := v3.Compare(v3.CreateRevision(e.leaderKey), "=", e.leaderRev)
		//	resp, err := client.Txn(ctx).If(cmp).Then(v3.OpDelete(e.leaderKey)).Commit()
		//	封装了一层，通过txn事务。如果如果这个k(pfx)的Revision == 自己 则删除
		if err := election.Resign(ctx); err != nil { //	放弃leader（当前节点不是领导者，调用它会导致错误）
			log.Fatal(err)
		}
		t.Log(election.Leader(ctx))
	})

	// 锁：
	// redis 常用的方式是 set key value EX 10 NX
	//  	1.通过NX确保不存在才新增
	//      2.通过EX保证如果程序意外退出没有delete导致该key永远存在
	//      3.value是需要设置一个自己的随机值，这边删除时需要指定value是自己生成的（lua原子性）
	// redis的缺点
	//		1.由于没有续约的机制，导致了可能还没有执行完，EX设定的时间已经超过，别的抢到了锁
	//      2.redis 主备同步切换异步同步问题导致数据不一致（备还没同步），同时redis可能会出现脑裂
	//      3.redis需要使用redlock解决
	// etcd
	//		1.etcd使用Raft原生支持避免脑裂，主备切换问题
	//		2.etcd使用线性读确保数据准确性
	//		3.支持事务满足原子性问题
	//		4.天生支持租约可以续租
	//      5.可以通过watch机制保证client crash时，其他client快速感知
	//		6.社区通过concurrency包封装了锁，选举的问题
	//			架构问题			   |     		  使用问题
	//		主从异步 (主从切换)	   |		lua原子性不好用(事务)
	//		  主从切换 (脑裂)		   |		  需要自己实现续约
	//							   |	   	    无watch功能
	t.Run("lock", func(t *testing.T) {
		t.Run("mutex", func(t *testing.T) {
			locker := concurrency.NewMutex(session, "/mutex")
			//_ = locker.TryLock(context.Background()) //如果获取锁失败会返回err无阻塞等待
			_ = locker.Lock(ctx) //如果获取锁失败会进行锁等待
			_ = locker.Unlock(ctx)
			t.Log("this is mutex", locker.Key())
		})
		t.Run("lock", func(t *testing.T) {
			locker := concurrency.NewLocker(session, "/locker") //NewMutex 的一个包装,只提供Lock跟UnLock方法
			locker.Lock()                                       //如果获取锁失败会进行锁等待
			defer locker.Unlock()
			t.Log("this is locker")
		})
	})

	//分布式消息队列
	t.Run("queue", func(t *testing.T) {
		t.Run("queue", func(t *testing.T) {
			queue := recipe.NewQueue(client, "/queue")

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer func() {
					wg.Done()
				}()
				data, err := queue.Dequeue()
				if err != nil {
					t.Error(err)
					return
				}
				t.Log(data)
			}()

			time.Sleep(time.Second)
			err := queue.Enqueue("hello")
			if err != nil {
				t.Error(err)
			}

			wg.Wait()
		})

		//分布式优先消息队列 （当在消费队列中已经通过优先级排序然后消费）
		t.Run("priority_query", func(t *testing.T) {
			queue := recipe.NewPriorityQueue(client, "/priority_query")

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer func() {
					wg.Done()
				}()
				data, err := queue.Dequeue()
				if err != nil {
					t.Error(err)
					return
				}
				t.Log(data)
			}()

			time.Sleep(time.Second)
			err := queue.Enqueue("hello", 10)
			if err != nil {
				t.Error(err)
			}

			wg.Wait()
		})
	})

	//事务 (原子性操作)
	t.Run("tx", func(t *testing.T) {
		t.Run("tx", func(t *testing.T) {
			res, err := client.Txn(ctx).
				Then(etcdClient.OpPut("hello", "world"), etcdClient.OpPut("say", "hello")).
				Commit() //提交事务
			if err != nil {
				t.Error(err)
				return
			}
			t.Log(res)
		})
		t.Run("stm", func(t *testing.T) {
			res, err := concurrency.NewSTM(client, func(stm concurrency.STM) error {
				stm.Put("hello", "world")
				t.Log(stm.Get("hello"))
				return nil //提交事务
			})
			t.Log(res, err)
		})
	})

	// 栅栏，目的多个消费者barrier.Wait此时阻塞，有一个发起barrier.Release则全部释放
	t.Run("barrier", func(t *testing.T) {
		t.Run("barrier", func(t *testing.T) {
			barrier := recipe.NewBarrier(client, "/barrier")
			t.Log(barrier.Hold())
			go func() {
				time.Sleep(time.Second)
				t.Log(barrier.Release())
			}()
			t.Log(barrier.Wait())
		})

		//当count的数量到达时才会全部解除阻塞一起执行之后的内容
		t.Run("double_barrier", func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(3)
			for i := 0; i < 3; i++ {
				go func(num int) {
					defer wg.Done()
					barrier := recipe.NewDoubleBarrier(session, "/double_barrier1", 3)
					t.Log(num)
					_ = barrier.Enter() //进入逻辑处理，结束执行Leave进入阻塞等待，
					if num == 1 {
						time.Sleep(time.Second * 5)
						t.Log(num, time.Now().Local())
					} else {
						t.Log(num, time.Now().Local())
					}
					_ = barrier.Leave() //需要等count的数量到达时才会全部解除阻塞一起执行之后的内容
					t.Log("end ", num, time.Now().Local())
				}(i)
			}

			wg.Wait()
		})
	})
}
