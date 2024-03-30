package etcd

//
//import (
//	"context"
//	"fmt"
//	"github.com/coreos/etcd/clientv3"
//	"github.com/long250038728/web/tool/register"
//	"strconv"
//	"strings"
//	"time"
//)
//
//var leaseTTLSeconds int64 = 10
//
//type Register struct {
//	client *clientv3.Client
//}
//
//// NewEtcdRegister   创建etcd服务注册类
//func NewEtcdRegister(addr string) (*Register, error) {
//	client, err := clientv3.New(clientv3.Config{
//		Endpoints:   []string{"http://localhost:2379"},
//		DialTimeout: 5 * time.Second,
//	})
//	if err != nil {
//		return nil, err
//	}
//	return &Register{client: client}, nil
//}
//
//func (r Register) Register(ctx context.Context, serviceInstance *register.ServiceInstance) error {
//	res, err := r.client.Grant(ctx, leaseTTLSeconds)
//	if err != nil {
//		return err
//	}
//	_, err = r.client.Put(ctx, serviceInstance.TableName, fmt.Sprintf("%s:%d", serviceInstance.Address, serviceInstance.Port), clientv3.WithLease(res.ID))
//	if err != nil {
//		return err
//	}
//
//	select {
//	case <-ctx.Done():
//		return nil
//	default:
//		go func() {
//			time.Sleep(time.Duration(leaseTTLSeconds/2) * time.Second)
//			_, _ = r.client.KeepAliveOnce(ctx, res.ID)
//		}()
//	}
//	return nil
//}
//
//func (r Register) DeRegister(ctx context.Context, serviceInstance *register.ServiceInstance) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (r Register) List(ctx context.Context, serviceName string) ([]*register.ServiceInstance, error) {
//	res, err := r.client.Get(ctx, serviceName)
//	if err != nil {
//		return nil, err
//	}
//
//	var list = make([]*register.ServiceInstance, 0, len(res.Kvs))
//	for _, svc := range res.Kvs {
//		val := strings.Split(":", string(svc.Value))
//		port, _ := strconv.ParseInt(val[1], 10, 0)
//		list = append(list, &register.ServiceInstance{
//			TableName:    serviceName,
//			Address: val[0],
//			Port:    int(port),
//		})
//	}
//
//	return list, nil
//}
//
//func (r Register) Subscribe(ctx context.Context, serviceName string) (<-chan *register.ServiceInstance, error) {
//	watch := r.client.Watch(ctx, serviceName)
//	for res := range watch {
//		for _, event := range res.Events {
//			fmt.Printf("Event Type: %v, Key: %s, Value: %s\n", event.Type, event.Kv.Key, event.Kv.Value)
//		}
//	}
//	return nil, nil
//}
