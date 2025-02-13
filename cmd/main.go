package main

import (
	"fmt"
	client "github.com/long250038728/web/cmd/client"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/git"
	"github.com/long250038728/web/tool/hook"
	"github.com/long250038728/web/tool/jenkins"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/ssh"
	"github.com/spf13/cobra"
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

//var services = &client.Svc{Kobe: make([]string, 0, 0), Marx: make([]string, 0, 0)}

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

func main() {
	gitCron := client.NewGitCron(gitClient)
	olineCron := client.NewOnlineCron(gitClient, jenkinsClient, ormClient, sshClient, hookClient, tels)
	devopsCron := client.NewDevopsCorn("", "")
	serverCron := client.NewServerCornCorn()
	chatCron := client.NewChatCorn()

	rootCmd := &cobra.Command{
		Use:   "linl",
		Short: "快速生成工具",
	}

	rootCmd.AddCommand(gitCron.Branch())
	rootCmd.AddCommand(gitCron.Pr())
	rootCmd.AddCommand(gitCron.List())
	rootCmd.AddCommand(gitCron.Completion())

	rootCmd.AddCommand(olineCron.Json())
	rootCmd.AddCommand(olineCron.Action())
	rootCmd.AddCommand(olineCron.Cron())
	rootCmd.AddCommand(devopsCron.Devops())
	rootCmd.AddCommand(serverCron.Server())
	rootCmd.AddCommand(chatCron.Chat())

	rootCmd.AddCommand()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
	}
}
