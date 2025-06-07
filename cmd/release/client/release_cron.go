package client

import (
	"context"
	"fmt"
	"github.com/long250038728/web/tool/git"
	"github.com/long250038728/web/tool/hook"
	"github.com/long250038728/web/tool/jenkins"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/ssh"
	"github.com/long250038728/web/tool/task/cron_job/robfig"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type ReleaseCorn struct {
	gitClient     git.Git
	jenkinsClient *jenkins.Client
	ormClient     *orm.Gorm
	sshClient     ssh.SSH
	hookClient    hook.Hook
	tels          []string
}

func NewReleaseCron(gitClient git.Git, jenkinsClient *jenkins.Client, ormClient *orm.Gorm, sshClient ssh.SSH, hookClient hook.Hook, tels []string) *ReleaseCorn {
	return &ReleaseCorn{
		gitClient:     gitClient,
		jenkinsClient: jenkinsClient,
		ormClient:     ormClient,
		sshClient:     sshClient,
		hookClient:    hookClient,
		tels:          tels,
	}
}

func (c *ReleaseCorn) Json() *cobra.Command {
	return &cobra.Command{
		Use:   "json [来源分支] [目标分支] [kobe/marx列表(.yaml)]",
		Short: "shell生成： 请输入【来源分支】【目标分支】【项目列表文件】",
		Long:  "shell生成： 请输入【来源分支】【目标分支】【项目列表文件】",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			source := args[0]
			target := args[1]
			svcPath := "./svc.yaml"

			if len(args) == 3 {
				svcPath = args[2]
			}

			if err := NewTaskClient(
				SetGit(c.gitClient),
				SetJenkins(c.jenkinsClient),
				SetOrm(c.ormClient),
				SetRemoteShell(c.sshClient),
				SetQyHook(c.hookClient, c.tels),
			).Build(ctx, source, target, svcPath); err != nil {
				fmt.Println("error :", err)
			}
			fmt.Println("ok")
		},
	}
}

func (c *ReleaseCorn) Action() *cobra.Command {
	return &cobra.Command{
		Use:   "action",
		Short: "上线操作",
		Long:  "上线操作",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			path := "./"
			if len(args) >= 1 {
				path = args[0]

				// 使用 os.Stat 获取文件信息
				info, err := os.Stat(path)
				if err != nil {
					fmt.Println("the path is error :", path)
					return
				}
				if !info.IsDir() {
					fmt.Println("the path is not dir :", path)
					return
				}
			}

			ctx := context.Background()
			if err := NewTaskClient(
				SetGit(c.gitClient),
				SetJenkins(c.jenkinsClient),
				SetOrm(c.ormClient),
				SetRemoteShell(c.sshClient),
				SetQyHook(c.hookClient, c.tels),
				SetOutPath(path),
			).Request(ctx); err != nil {
				fmt.Println("error :", err)
			}
			fmt.Println("ok")
		},
	}
}

func (c *ReleaseCorn) Cron() *cobra.Command {
	return &cobra.Command{
		Use:   "cron [执行时] [执行分]",
		Short: "cron： 执行请输入时间",
		Long:  "cron： 执行请输入时间",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			//信号退出
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			defer signal.Stop(quit)

			h := args[0]
			m := args[1]

			path := "./"
			if len(args) >= 3 {
				path = args[2]

				// 使用 os.Stat 获取文件信息
				info, err := os.Stat(path)
				if err != nil {
					fmt.Println("the path is error :", path)
					return
				}
				if !info.IsDir() {
					fmt.Println("the path is not dir :", path)
					return
				}
			}

			if _, err := strconv.Atoi(h); err != nil {
				fmt.Println("hour cron is error :", err)
				return
			}
			if _, err := strconv.Atoi(m); err != nil {
				fmt.Println("minute cron is error :", err)
				return
			}
			spec := fmt.Sprintf("%s %s * * *", m, h)
			fmt.Println(spec)

			taskClient := NewTaskClient(
				SetGit(c.gitClient),
				SetJenkins(c.jenkinsClient),
				SetOrm(c.ormClient),
				SetRemoteShell(c.sshClient),
				SetQyHook(c.hookClient, c.tels),
				SetOutPath(path),
			)
			ctx := context.Background()
			_ = taskClient.HookSend(ctx, spec)

			//创建任务
			job := robfig.NewCronJob()
			job.Start()
			defer func() {
				fmt.Println("=========")
				job.Close()
			}()

			//添加任务
			ch := make(chan error)
			_, _ = job.AddFunc(spec, func() {
				fmt.Println("执行了")
				ch <- taskClient.Request(ctx)
			})

			select {
			case err := <-ch: //等待执行
				fmt.Println(err)

			case s := <-quit: //监听信号
				fmt.Println(s.String())
			}
		},
	}
}
