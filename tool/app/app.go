package app

import (
	"context"
	"github.com/long250038728/web/tool/auth"
	"github.com/long250038728/web/tool/register"
	"github.com/long250038728/web/tool/server"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var _ Application = &App{}

type App struct {
	Application
	ctx    context.Context
	cancel context.CancelFunc

	servers  []server.Server
	register register.Register
	auth     auth.Auth
}

func NewApp(opts ...Option) Application {
	app := &App{}
	for _, opt := range opts {
		opt(app)
	}
	return app
}

func (app *App) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	group, ctx := errgroup.WithContext(ctx)
	app.cancel = cancel
	app.ctx = ctx

	//优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	group.Go(func() error {
		select {
		case <-app.ctx.Done():
			return app.ctx.Err()
		case <-quit: //此时阻塞，收到指令 ctx.Done触发
			app.Stop()
		}
		return nil
	})

	//遍历服务
	for _, s := range app.servers {
		svc := s

		//启动服务
		group.Go(func() error {
			err := svc.Start() //此时阻塞，其中有一个报错则 ctx.Done触发
			return err
		})

		//关闭服务
		group.Go(func() error {
			<-app.ctx.Done() //此时阻塞，等待 ctx.Done触发 ，去关闭服务
			err := svc.Stop()
			return err
		})

		//服务注册 && 取消
		if app.register != nil {
			group.Go(func() error {
				select {
				case <-app.ctx.Done():
					return nil
				default:
					return app.register.Register(context.Background(), svc.ServiceInstance())
				}
			})

			group.Go(func() error {
				<-app.ctx.Done() //此时阻塞，等待 ctx.Done触发 ，去取消注册
				time.Sleep(time.Second * 2)
				return app.register.DeRegister(context.Background(), svc.ServiceInstance())
			})
		}
	}

	//监听所有err
	return group.Wait()
}

func (app *App) Stop() {
	app.cancel()
}
