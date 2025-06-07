package main

import (
	"fmt"
	"github.com/long250038728/web/cmd/gen/client"
	client2 "github.com/long250038728/web/cmd/release/client"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/git"
	"github.com/long250038728/web/tool/hook"
	"github.com/long250038728/web/tool/jenkins"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/ssh"
	"github.com/spf13/cobra"
)

// go get -u github.com/spf13/cobra

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
	configLoad.MustLoadConfigPath("other/gitee.yaml", &gitConfig)
	configLoad.MustLoadConfigPath("other/jenkins.yaml", &jenkinsConfig)
	configLoad.MustLoadConfigPath("online/db.yaml", &ormConfig)
	configLoad.MustLoadConfigPath("other/ssh.yaml", &sshConfig)

	if gitClient, err = git.NewGiteeClient(&gitConfig); err != nil {
		panic(err)
	}
	if jenkinsClient, err = jenkins.NewJenkinsClient(&jenkinsConfig); err != nil {
		panic(err)
	}
	if ormClient, err = orm.NewMySQLGorm(&ormConfig); err != nil {
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
	gitCron := client2.NewGitCron(gitClient)
	olineCron := client2.NewReleaseCron(gitClient, jenkinsClient, ormClient, sshClient, hookClient, tels)
	devopsCron := client.NewDevopsCorn("", "")
	serverCron := client.NewServerCornCorn()

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
	rootCmd.AddCommand()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
	}
}
