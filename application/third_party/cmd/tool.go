package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/long250038728/web/application/third_party/client"
	"github.com/long250038728/web/tool/configurator"
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

// go get -u github.com/spf13/cobra
//
// 打包成工具
// go build tool
// mv tool linl
//
// 使用
// ./linl -h
// ./linl branch master feature/sm100 hume locke

var services = &client.Svc{Kobe: make([]string, 0, 0), Marx: make([]string, 0, 0)}

var productHash = map[string]string{
	"locke": "zhubaoe/locke",
	"kobe":  "zhubaoe-go/kobe",
	"hume":  "zhubaoe/hume",
	"ari":   "zhubaoe/aristotle",
	"h5":    "fissiongeek/h5-sales",
	"soc":   "zhubaoe/socrates",
	"plato": "zhubaoe/plato",
	"marx":  "zhubaoe/marx",
}
var productList = []string{
	"zhubaoe/locke",
	"zhubaoe-go/kobe",
	"zhubaoe/hume",
	"zhubaoe/socrates",
	"zhubaoe/aristotle",
	"fissiongeek/h5-sales",
	"zhubaoe/plato",
	"zhubaoe/marx",
}

var gitConfig git.Config
var jenkinsConfig jenkins.Config
var ormConfig orm.Config
var sshConfig ssh.Config

var gitClient git.Git
var jenkinsClient *jenkins.Client
var ormClient *orm.Gorm
var sshClient ssh.SSH
var hookClient hook.Hook

var hookToken = "bb3f6f61-04b8-4b46-a167-08a2c91d408d"
var tels = []string{"18575538087"}

func init() {
	var err error
	configLoad := configurator.NewYaml()
	configLoad.MustLoadConfigPath("gitee.yaml", &gitConfig)
	configLoad.MustLoadConfigPath("jenkins.yaml", &jenkinsConfig)
	configLoad.MustLoadConfigPath("online/db.yaml", &ormConfig)
	configLoad.MustLoadConfigPath("ssh.yaml", &sshConfig)

	if gitClient, err = git.NewGiteeClient(&gitConfig); err != nil {
		panic(err)
	}
	if jenkinsClient, err = jenkins.NewJenkinsClient(&jenkinsConfig); err != nil {
		panic(err)
	}
	if ormClient, err = orm.NewGorm(&ormConfig); err != nil {
		panic(err)
	}
	if sshClient, err = ssh.NewRemoteSSH(&sshConfig); err != nil {
		panic(err)
	}
	if hookClient, err = hook.NewQyHookClient(&hook.Config{Token: hookToken}); err != nil {
		panic(err)
	}
}

var ctx = context.Background()

func main() {
	rootCmd := &cobra.Command{
		Use:   "linl",
		Short: "上线快速生成工具",
	}

	rootCmd.AddCommand(branch())
	rootCmd.AddCommand(pr())
	rootCmd.AddCommand(list())
	rootCmd.AddCommand(json())
	rootCmd.AddCommand(action())
	rootCmd.AddCommand(completion())
	rootCmd.AddCommand(cron())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
	}
}

// checkProduct 通过快捷名称找对应的项目地址
func checkProduct(products []string) error {
	for _, product := range products {
		if _, ok := productHash[product]; !ok {
			return errors.New("输入有误:" + product + "找不到对应的项目")
		}
	}
	return nil
}

func completion() *cobra.Command {
	return &cobra.Command{
		Use:   "completion",
		Short: "无",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("该方法无实际命令")
		},
	}
}

func pr() *cobra.Command {
	return &cobra.Command{
		Use:   "pr [来源分支] [目标分支] [项目缩写名...]",
		Short: "PR创建： 请输入【来源分支】【目标分支】【项目缩写名...】",
		Long:  "PR创建： 请输入【来源分支】【目标分支】【项目缩写名...】\nlocke\nkobe\nhume\nari\nh5\nsoc\nplato\nmarx",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			source := args[0]
			target := args[1]
			products := args[2:]

			//通过快捷名称找对应的项目地址
			if err := checkProduct(products); err != nil {
				fmt.Println(err.Error())
				return
			}

			for _, product := range products {
				addr := productHash[product]
				if _, err := gitClient.CreatePR(ctx, addr, source, target); err != nil {
					fmt.Println("生成失败:" + err.Error())
					continue
				}
			}
			fmt.Println("全部执行完成")
			return
		},
	}
}

func branch() *cobra.Command {
	return &cobra.Command{
		Use:   "branch [来源分支] [目标分支] [项目缩写名...]",
		Short: "分支创建： 请输入【来源分支】【目标分支】【项目缩写名...】",
		Long:  "分支创建： 请输入【来源分支】【目标分支】【项目缩写名...】\nlocke\nkobe\nhume\nari\nh5\nsoc\nplato\nmarx",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			source := args[0]
			target := args[1]
			products := args[2:]

			//通过快捷名称找对应的项目地址
			if err := checkProduct(products); err != nil {
				fmt.Println(err.Error())
				return
			}
			for _, product := range products {
				addr := productHash[product]
				if err := gitClient.CreateFeature(ctx, addr, source, target); err != nil {
					fmt.Println("生成失败:" + err.Error())
					continue
				}
				fmt.Println(product + ": branch ok")
			}
			fmt.Println("全部执行完成")
			return
		},
	}
}

// list 上线检查脚本
func list() *cobra.Command {
	return &cobra.Command{
		Use:   "list [来源分支] [目标分支]",
		Short: "获取PR列表： 请输入【来源分支】【目标分支】",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			source := args[0]
			target := args[1]

			var address = make([]string, 0, len(productList))

			for _, addr := range productList {
				list, err := gitClient.GetPR(ctx, addr, source, target)
				if err != nil {
					continue
				}
				if len(list) == 0 {
					continue
				}
				address = append(address, list[0].HtmlUrl)
			}
			fmt.Println("有", len(address), "个PR项目")
			fmt.Println(address)
		},
	}
}

func json() *cobra.Command {
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

			if err := client.NewTaskClient(
				client.SetGit(gitClient),
				client.SetJenkins(jenkinsClient),
				client.SetOrm(ormClient),
				client.SetRemoteShell(sshClient),
				client.SetQyHook(hookClient, tels),
			).Build(ctx, source, target, svcPath); err != nil {
				fmt.Println("error :", err)
			}
			fmt.Println("ok")
		},
	}
}

func action() *cobra.Command {
	return &cobra.Command{
		Use:   "action",
		Short: "上线操作",
		Long:  "上线操作",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			if err := client.NewTaskClient(
				client.SetGit(gitClient),
				client.SetJenkins(jenkinsClient),
				client.SetOrm(ormClient),
				client.SetRemoteShell(sshClient),
				client.SetQyHook(hookClient, tels),
			).Request(ctx); err != nil {
				fmt.Println("error :", err)
			}
			fmt.Println("ok")
		},
	}
}

func cron() *cobra.Command {
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
				ctx := context.Background()
				ch <- client.NewTaskClient(
					client.SetGit(gitClient),
					client.SetJenkins(jenkinsClient),
					client.SetOrm(ormClient),
					client.SetRemoteShell(sshClient),
					client.SetQyHook(hookClient, tels),
				).Request(ctx)
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
