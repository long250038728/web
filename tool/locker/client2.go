package locker

import (
	"context"
	etcdClient "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"time"
)

type etcd struct {
	client *etcdClient.Client
	key,
	identification string
	time    time.Duration
	session *concurrency.Session
}

func NewEtcdLocker(client *etcdClient.Client, key string) Locker {
	return &etcd{
		client: client,
		key:    key,
	}
}

func (l *etcd) Lock(ctx context.Context) error {
	session, err := concurrency.NewSession(l.client)
	if err != nil {
		return err
	}
	l.session = session

	locker := concurrency.NewMutex(l.session, "/"+l.key)
	err = locker.TryLock(ctx)
	if err != nil {
		_ = session.Close()
	}

	if err == nil {
		//自动续
		if _, err = l.client.KeepAlive(ctx, session.Lease()); err != nil {
			_ = session.Close()
		}
	}
	return err
}

func (l *etcd) UnLock(ctx context.Context) error {
	locker := concurrency.NewMutex(l.session, "/"+l.key)
	return locker.TryLock(ctx)
}

func (l *etcd) Refresh(ctx context.Context) error {
	return nil
}

func (l *etcd) AutoRefresh(ctx context.Context) error {
	return nil
}

func (l *etcd) Close() error {
	return l.session.Close()
}
