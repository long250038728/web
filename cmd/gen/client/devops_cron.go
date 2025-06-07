package client

import (
	_ "embed"
	"fmt"
	"github.com/long250038728/web/tool/gen"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)

type devopsValue struct {
	Server string `json:"server" yaml:"server"`
	Http   string `json:"http" yaml:"http"`
	Grpc   string `json:"grpc" yaml:"grpc"`
}
type kubernetesValue struct {
	*devopsValue
	Version string `json:"version" yaml:"version"`
	Hub     string `json:"hub" yaml:"hub"`
}

type devops struct {
}

func newDevopsGen() *devops {
	return &devops{}
}

//go:embed tmpl/devops/dockerfile.tmpl
var dockerfileTmpl string

//go:embed tmpl/devops/kubernetes.tmpl
var kubernetesTmpl string

// genDockerfile 生成
func (g *devops) genDockerfile(data *devopsValue) ([]byte, error) {
	return (&gen.Impl{
		Name: "gen dockerfile",
		Tmpl: dockerfileTmpl,
		Data: data,
	}).Gen()
}
func (g *devops) genKubernetes(data *kubernetesValue) ([]byte, error) {
	return (&gen.Impl{
		Name: "gen kubernetes",
		Tmpl: kubernetesTmpl,
		Data: data,
	}).Gen()
}

type DevopsCorn struct {
	//hub  string
	//path string
}

func NewDevopsCorn() *DevopsCorn {
	//if len(hub) == 0 {
	//	hub = "ccr.ccs.tencentyun.com/linl"
	//}
	//if len(path) == 0 {
	//	path = "./devops"
	//}
	return &DevopsCorn{}
}

func (c *DevopsCorn) Devops() *cobra.Command {
	//go run main.go devops /Users/linlong/Desktop/web/application  test 100081 100082
	//go run main.go devops /Users/linlong/Desktop/web/application  test 100081 100082
	return &cobra.Command{
		Use:   "devops  [输出路径] [服务名] [http端口号:默认10000] [rpc端口号:默认10001] ",
		Short: "创建devops： 请输入  [输出路径] [服务名] [http端口号:默认10000] [rpc端口号:默认10001] ",
		Long:  "创建devops： 请输入  [输出路径] [服务名] [http端口号:默认10000] [rpc端口号:默认10001] ",
		Args:  cobra.MinimumNArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			server := args[1]
			http := args[2]
			grpc := args[3]

			version := time.Now().Format("20060102_15_04_05")
			hub := "ccr.ccs.tencentyun.com/linl"

			if len(args) >= 5 {
				version = args[4]
			}
			if len(args) >= 6 {
				hub = args[4]
			}

			devops := func() error {
				var dockerfileBytes []byte
				var kubernetesBytes []byte
				var err error

				_, err = os.Stat(filepath.Join(path, server))
				if os.IsNotExist(err) {
					if err := os.MkdirAll(filepath.Join(path, server), os.ModePerm); err != nil {
						return err
					}
				}

				g := newDevopsGen()

				if dockerfileBytes, err = g.genDockerfile(&devopsValue{Server: server, Http: http, Grpc: grpc}); err != nil {
					return err
				}
				if kubernetesBytes, err = g.genKubernetes(&kubernetesValue{devopsValue: &devopsValue{Server: server, Http: http, Grpc: grpc}, Version: version, Hub: hub}); err != nil {
					return err
				}

				//write file
				if err := os.WriteFile(filepath.Join(path, server, "dockerfile"), dockerfileBytes, os.ModePerm); err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(path, server, "kubernetes.yaml"), kubernetesBytes, os.ModePerm); err != nil {
					return err
				}
				return nil
			}

			if err := devops(); err != nil {
				fmt.Println("执行出错", err.Error())
			}

			fmt.Println("全部执行完成")
			return
		},
	}
}
