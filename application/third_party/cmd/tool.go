package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/long250038728/web/application/third_party/client"
	"github.com/long250038728/web/tool/config"
	"github.com/spf13/cobra"
	"os"
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

var source string
var target string
var products []string
var services = &client.Svc{
	Kobe: make([]string, 0, 0), Marx: make([]string, 0, 0),
}

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

var gitToken = "5f8aaa1e024cad5e24e86fda85c57f49"

var name = "admin"
var password = "11fbfc1aab366147522f497f6c7d48b2ca"
var jenkinsToken = name + ":" + password

//var jenkinsToken = "admin:11fbfc1aab366147522f497f6c7d48b2ca"

var git = client.NewGiteeClinet(gitToken)
var jenkins = client.NewJenkinsClient("http://111.230.143.16:8081", "admin", "11fbfc1aab366147522f497f6c7d48b2ca")
var ctx = context.Background()

func main2() {
	rootCmd := &cobra.Command{
		Use:   "linl",
		Short: "上线快速生成工具",
	}

	rootCmd.AddCommand(branch())
	rootCmd.AddCommand(pr())
	rootCmd.AddCommand(list())
	rootCmd.AddCommand(online())
	rootCmd.AddCommand(request())

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

func pr() *cobra.Command {
	return &cobra.Command{
		Use:   "pr [来源分支] [目标分支] [项目缩写名...]",
		Short: "PR创建： 请输入【来源分支】【目标分支】【项目缩写名...】",
		Long:  "PR创建： 请输入【来源分支】【目标分支】【项目缩写名...】\nlocke\nkobe\nhume\nari\nh5\nsoc\nplato\nmarx",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			source = args[0]
			target = args[1]
			products = args[2:]

			//通过快捷名称找对应的项目地址
			if err := checkProduct(products); err != nil {
				fmt.Println(err.Error())
				return
			}

			for _, product := range products {
				addr := productHash[product]
				if _, err := git.CreatePR(ctx, addr, source, target); err != nil {
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
			source = args[0]
			target = args[1]
			products = args[2:]

			//通过快捷名称找对应的项目地址
			if err := checkProduct(products); err != nil {
				fmt.Println(err.Error())
				return
			}

			for _, product := range products {
				addr := productHash[product]
				if err := git.CreateFeature(ctx, addr, source, target); err != nil {
					fmt.Println("生成失败:" + err.Error())
					continue
				}
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
			source = args[0]
			target = args[1]

			var address = make([]string, 0, len(productList))

			for _, addr := range productList {
				list, err := git.GetPR(ctx, addr, source, target)
				if err != nil {
					continue
				}
				if len(list) != 1 {
					continue
				}
				address = append(address, list[0].HtmlUrl)
			}
			fmt.Println("有", len(address), "个PR项目")
			fmt.Println(address)
		},
	}
}

// online
func online() *cobra.Command {
	return &cobra.Command{
		Use:   "online [来源分支] [目标分支] [kobe/marx列表(.yaml)]",
		Short: "shell生成： 请输入【来源分支】【目标分支】【项目列表文件】",
		Long:  "shell生成： 请输入【来源分支】【目标分支】【项目列表文件】",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			source = args[0]
			target = args[1]
			conf := ""
			if len(args) == 3 {
				conf = args[2]
			}

			if len(conf) > 0 {
				if err := (&config.Yaml{}).Load(conf, &services); err != nil {
					fmt.Println(err)
					return
				}
			}

			var address = make([]string, 0, len(productList))

			for _, addr := range productList {
				list, err := git.GetPR(ctx, addr, source, target)
				if err != nil {
					continue
				}
				if len(list) != 1 {
					continue
				}

				if addr == "zhubaoe-go/kobe" && len(services.Kobe) == 0 {
					fmt.Println("有kobe项目，但是未添加服务")
					return
				}
				if addr == "zhubaoe/marx" && len(services.Marx) == 0 {
					fmt.Println("有marx项目，但是未添加服务")
					return
				}
				address = append(address, list[0].Url)
			}

			b, err := client.NewPrGen(gitToken, jenkinsToken).GenMerge(address, services)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			err = os.WriteFile("./online.md", b, os.ModePerm)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("文件生成./online.md")
			fmt.Println(string(b))
		},
	}
}

type RequestType struct {
	Type    int
	Address string
	Params  map[string]any
	Num     int32
}

func request() *cobra.Command {
	return &cobra.Command{
		Use:   "request [来源分支] [目标分支] [kobe/marx列表(.yaml)]",
		Short: "shell生成： 请输入【来源分支】【目标分支】【项目列表文件】",
		Long:  "shell生成： 请输入【来源分支】【目标分支】【项目列表文件】",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			source = args[0]
			target = args[1]
			if len(args) == 3 {
				if err := (&config.Yaml{}).Load(args[2], &services); err != nil {
					fmt.Println(err)
					return
				}
			}

			var address = make([]*RequestType, 0, 100)

			for _, addr := range productList {
				list, err := git.GetPR(ctx, addr, source, target)
				if err != nil {
					continue
				}
				if len(list) != 1 {
					continue
				}
				if addr == "zhubaoe-go/kobe" && len(services.Kobe) == 0 {
					fmt.Println("有kobe项目，但是未添加服务")
					return
				}
				if addr == "zhubaoe/marx" && len(services.Marx) == 0 {
					fmt.Println("有marx项目，但是未添加服务")
					return
				}

				//调用合并分支
				address = append(address, &RequestType{Type: 1, Address: list[0].Url, Num: list[0].Number})

				//两台服务器
				if addr == "zhubaoe-go/kobe" {
					for _, svc := range services.Kobe {
						address = append(address, &RequestType{Type: 2, Address: svc, Params: map[string]any{"BRANCH": "origin/master", "SYSTEM": "root@172.16.0.34"}})
						address = append(address, &RequestType{Type: 2, Address: svc, Params: map[string]any{"BRANCH": "origin/master", "SYSTEM": "root@172.16.0.9"}})
					}
				}

				// 一台服务器
				if addr == "zhubaoe/marx" {
					for _, svc := range services.Marx {
						address = append(address, &RequestType{Type: 2, Address: svc})
					}
				}

				if addr == "zhubaoe/plato" {
					address = append(address, &RequestType{Type: 2, Address: "plato-prod", Params: map[string]any{"BRANCH": "origin/master"}})
				}

				// 三个服务
				if addr == "zhubaoe-go/locke" {
					address = append(address, &RequestType{Type: 2, Address: "locke-prod_32"})
					address = append(address, &RequestType{Type: 2, Address: "locke-prod_64"})
					address = append(address, &RequestType{Type: 2, Address: "locke-hot-prod-64"})
				}
			}

			for _, request := range address {
				switch request.Type {
				case 1:
					if err := git.Merge(ctx, request.Address, request.Num); err != nil {
						fmt.Println(request.Address, "pr merge", err)
						return
					}
				case 2:
					if err := jenkins.BlockBuild(ctx, request.Address, request.Params); err != nil {
						fmt.Println(request.Address, "block build", err)
						return
					}
				default:
					fmt.Println("type is err")
					return
				}
			}

		},
	}
}
