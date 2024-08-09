package config_center

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

var client ConfigCenter
var eClient *EtcdCenter

func init() {
	var err error
	var centerConfig Config
	configurator.NewYaml().MustLoadConfigPath("center.yaml", &centerConfig)
	if client, err = NewEtcdConfigCenter(&centerConfig); err != nil {
		panic(err)
	}
	if eClient, err = NewEtcd(&centerConfig); err != nil {
		panic(err)
	}
}

func TestConfig(t *testing.T) {
	ctx := context.Background()
	t.Run("watch", func(t *testing.T) {
		t.Log(client.Watch(ctx, "hello", func(changeKey, changeVal []byte) {
			fmt.Println(string(changeKey), string(changeVal))
		}))
	})
	t.Run("set", func(t *testing.T) {
		t.Log(client.Set(ctx, "hello", "123456"))
	})
	t.Run("set", func(t *testing.T) {
		t.Log(client.Set(ctx, "hello", "4567"))
	})
	t.Run("get", func(t *testing.T) {
		t.Log(client.Get(ctx, "hello"))
	})
	t.Run("del", func(t *testing.T) {
		t.Log(client.Del(ctx, "hello"))
	})

	_ = client.Close()
}

func TestConfig_Upload(t *testing.T) {
	t.Log(client.UpLoad(context.Background(), "/Users/linlong/Desktop/web/config"))
}

func TestEtcd(t *testing.T) {
	defer func() {
		_ = eClient.Close()
	}()
	client := eClient.client

	session, err := concurrency.NewSession(client)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		t.Log(session.Lease())
		_ = session.Close() //当session close时Campaign就会取消成为leader（租约（lease）会被撤销）
	}()

	t.Run("watch", func(t *testing.T) {
		//监听key的变化(包括leader)
		go func() {
			ch := client.Watch(context.Background(), "/my-election", etcdClient.WithPrefix())
			for msg := range ch {
				for _, even := range msg.Events {
					t.Log(even)
				}
			}
		}()
	})

	t.Run("Lease", func(t *testing.T) {
		// 创建一个租约，有效期为1秒
		lease, err := client.Grant(context.Background(), 1)
		if err != nil {
			t.Error(err)
			return
		}

		_, err = client.Put(context.Background(), "key", "val", etcdClient.WithLease(lease.ID))
		if err != nil {
			t.Error(err)
			return
		}

		// 继续使用租约进行其他操作，例如保持租约活跃
		resCh, err := client.KeepAlive(context.Background(), lease.ID)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			//阻塞监听
			for chResp := range resCh {
				t.Log(fmt.Sprintf("KeepAlive response: %v", chResp))
			}
		}()
		time.Sleep(time.Second * 5)

		//移除租约
		t.Log(client.Revoke(context.Background(), lease.ID))

		//获取值（如果租约已经失效，则获取不到）
		val, err := client.Get(context.Background(), "key")
		if err != nil {
			log.Fatal(err)
		}
		t.Log(val.Kvs)
	})

	t.Run("election", func(t *testing.T) {
		election := concurrency.NewElection(session, "/my-election")

		//监听当前的election
		go func() {
			ch := election.Observe(context.Background())
			for _ = range ch {
			}
		}()

		//竞选leader
		if err := election.Campaign(context.Background(), "candidate-1"); err != nil {
			t.Log(err)
			return
		}
		t.Log(election.Leader(context.Background()))

		//更新leader中的val（当前节点不是领导者，调用它会导致错误）
		if err := election.Proclaim(context.Background(), "candidate-2"); err != nil {
			t.Log(err)
			return
		}
		t.Log(election.Leader(context.Background()))

		t.Log(election.Key(), election.Rev())

		//放弃leader（当前节点不是领导者，调用它会导致错误）
		if err := election.Resign(context.Background()); err != nil {
			log.Fatal(err)
		}
		t.Log(election.Leader(context.Background()))
	})

	t.Run("lock", func(t *testing.T) {
		locker := concurrency.NewLocker(session, "/locker")
		locker.Lock()
		defer locker.Unlock()
		t.Log("this is locker")
	})

	t.Run("mutex", func(t *testing.T) {
		locker := concurrency.NewMutex(session, "/mutex")
		_ = locker.Lock(context.Background())
		_ = locker.Unlock(context.Background())
		t.Log("this is mutex", locker.Key())
	})

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

	t.Run("tx", func(t *testing.T) {
		res, err := client.Txn(context.Background()).
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

	t.Run("barrier", func(t *testing.T) {
		barrier := recipe.NewBarrier(client, "/barrier")
		t.Log(barrier.Hold())
		go func() {
			time.Sleep(time.Second)
			t.Log(barrier.Release())
		}()

		t.Log(barrier.Wait())
	})

	t.Run("double_barrier", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(3)
		for i := 0; i < 3; i++ {
			go func(num int) {
				defer wg.Done()
				barrier := recipe.NewDoubleBarrier(session, "/double_barrier1", 3)
				t.Log(num)
				_ = barrier.Enter()
				if num == 1 {
					time.Sleep(time.Second * 5)
					t.Log(num, time.Now().Local())
				} else {
					t.Log(num, time.Now().Local())
				}
				_ = barrier.Leave()
				t.Log("end ", num, time.Now().Local())
			}(i)
		}

		wg.Wait()
	})
}
