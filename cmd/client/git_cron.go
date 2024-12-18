package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/git"
	"github.com/spf13/cobra"
)

type GitCorn struct {
	gitClient git.Git
}

func NewGitCron(gitClient git.Git) *GitCorn {
	return &GitCorn{
		gitClient: gitClient,
	}
}

// checkProduct 通过快捷名称找对应的项目地址
func (c *GitCorn) checkProduct(products []string) error {
	for _, product := range products {
		if _, ok := ProductHash[product]; !ok {
			return errors.New("输入有误:" + product + "找不到对应的项目")
		}
	}
	return nil
}

func (c *GitCorn) Completion() *cobra.Command {
	return &cobra.Command{
		Use:   "completion",
		Short: "无",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("该方法无实际命令")
		},
	}
}

func (c *GitCorn) Pr() *cobra.Command {
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
			if err := c.checkProduct(products); err != nil {
				fmt.Println(err.Error())
				return
			}

			ctx := context.Background()
			for _, product := range products {
				addr := ProductHash[product]
				if _, err := c.gitClient.CreatePR(ctx, addr, source, target); err != nil {
					fmt.Println("生成失败:" + err.Error())
					continue
				}
			}
			fmt.Println("全部执行完成")
			return
		},
	}
}

func (c *GitCorn) Branch() *cobra.Command {
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
			if err := c.checkProduct(products); err != nil {
				fmt.Println(err.Error())
				return
			}

			ctx := context.Background()
			for _, product := range products {
				addr := ProductHash[product]
				if err := c.gitClient.CreateFeature(ctx, addr, source, target); err != nil {
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

func (c *GitCorn) List() *cobra.Command {
	return &cobra.Command{
		Use:   "list [来源分支] [目标分支]",
		Short: "获取PR列表： 请输入【来源分支】【目标分支】",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			source := args[0]
			target := args[1]

			var address = make([]string, 0, len(ProductList))
			ctx := context.Background()
			for _, addr := range ProductList {
				list, err := c.gitClient.GetPR(ctx, addr, source, target)
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
