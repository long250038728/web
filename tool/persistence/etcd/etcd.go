package etcd

import (
	"context"
	"errors"
	"fmt"
	etcdClient "go.etcd.io/etcd/client/v3"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Address string `json:"address" yaml:"address"`
	Prefix  string `json:"prefix" yaml:"prefix"`
}

type EtcdCenter struct {
	io.Closer
	client *etcdClient.Client
	prefix string
}

// NewEtcdConfigCenter   配置中心
func NewEtcdConfigCenter(config *Config) (ConfigCenter, error) {
	// 账号信息
	// 		1. 通过账号密码登录 （simpleToken/JWT）
	// 		2. 通过证书登录
	// 权限RBAC
	//  	User / Role / Permission
	//		etcdctl user add alice -- user root:root  									//创建alice用户
	//		etcdctl role add admin -- user root:root  									//创建admin角色
	//		etcdctl role grant-permission admin readwrite hello helly --user root:root  //给admin角色权限添加[hello,helly]之前的读写操作
	//      etcdctl user grant-role alice admin --user root:root 						//用户alice绑定admin角色
	client, err := etcdClient.New(etcdClient.Config{
		Endpoints:   []string{config.Address},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	if config.Prefix == "" {
		config.Prefix = Prefix
	}
	return &EtcdCenter{client: client, prefix: config.Prefix}, nil
}

func (r *EtcdCenter) Get(ctx context.Context, key string) (string, error) {
	// Linearized Read 线性读 () 把请求转发到leader节点中，此时还会通过readIndex心跳获取自己是不是依旧是leader节点 —— 默认
	// Serializable Read 串行读 () 随机一个节点读取数据，该节点可能日志还未应用，所以会读到落后的数据
	res, err := r.client.Get(ctx, r.prefix+key) //etcdClient.WithSerializable()

	if err != nil {
		return "", err
	}
	if res.Count <= 0 {
		return "", errors.New("key is not value")
	}
	return string(res.Kvs[0].Value), nil
}

func (r *EtcdCenter) Set(ctx context.Context, key, value string) error {
	// etcd会检查当前的etcd db的大小，如果超过QUOTA配额会警告并拒绝写入，变成集群只读（调大后需要发生那个额外的命令）—— 默认2G
	// etcd的MVCC跟mysql区别在于，这个是用于集群中各个节点的日志写入信息。如现在写版本号是100的信息，需要判断本地是不是上一个就是99（单调递增），如果不是则需要先补齐并应用前面的日志（状态机）
	// 写入前判断:
	// 		1.简单限速
	// 		2.包最大1.5mb
	// 写入数据时：
	// 		1. leader发起提案（leader任期，投票信息，已提交索引，日志类型等），判断超时（默认7s）
	// 		2. follower投票写入
	// 超时问题：
	//		1. 由于网络问题可能节点不通讯/丢包，投票不过半重试/leader重新选举
	//		2. 磁盘io延迟（WAL，数据库(随机写入，页分裂)等写入）
	// etcd内存占用：
	// 		1.在apply之前保存在内存中写入raftlog日志中（会导致写请求过多(value值大)时该数组内存较大）
	//		2.确认后会写入内存treeIndex(b-tree)中，还会写入到WAL与boltdb中（使用了mmap技术导致db越大内存越大）
	// 		3.使用watch时会维护一定的内存
	//		4.写请求如果还有Lease时会维护这个
	_, err := r.client.Put(ctx, r.prefix+key, value)
	return err
}

func (r *EtcdCenter) Del(ctx context.Context, key string) error {
	_, err := r.client.Delete(ctx, r.prefix+key)
	return err
}

func (r *EtcdCenter) Watch(ctx context.Context, key string, callback func(changeKey, changeVal []byte)) error {
	ch := r.client.Watch(ctx, key, etcdClient.WithRange(etcdClient.GetPrefixRangeEnd(r.prefix+key)))
	for {
		select {
		case resp, ok := <-ch:
			if !ok {
				return nil
			}
			if resp.Canceled {
				return fmt.Errorf("watch operation canceled") // 操作被取消，返回错误
			}
			callback(resp.Events[0].Kv.Key, resp.Events[0].Kv.Value)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (r *EtcdCenter) UpLoad(ctx context.Context, rootPath string, files ...string) error {
	if len(files) == 0 {
		return errors.New("file list is empty")
	}

	for _, fileName := range files {
		f := strings.Split(fileName, ".")
		if len(f) != 2 {
			return errors.New("files is error: " + fileName)
		}

		// 获取
		b, err := os.ReadFile(filepath.Join(rootPath, fileName))
		if err != nil {
			return err
		}

		key := r.prefix + f[0]

		// 先删除
		_ = r.Del(ctx, key)

		// 上传
		err = r.Set(ctx, key, string(b))
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *EtcdCenter) Close() error {
	return r.client.Close()
}
