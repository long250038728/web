package client

import (
	_ "embed"
	"fmt"
	"github.com/long250038728/web/tool/gen"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

type devopsValue struct {
	Server string `json:"server" yaml:"server"`
	Http   string `json:"http" yaml:"http"`
	Grpc   string `json:"grpc" yaml:"grpc"`
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

//go:embed tmpl/devops/build.tmp
var bashTmpl string

// genDockerfile 生成
func (g *devops) genDockerfile(data *devopsValue) ([]byte, error) {
	return (&gen.Impl{
		Name: "gen dockerfile",
		Tmpl: dockerfileTmpl,
		Data: data,
	}).Gen()
}
func (g *devops) genKubernetes(data *devopsValue) ([]byte, error) {
	return (&gen.Impl{
		Name: "gen kubernetes",
		Tmpl: kubernetesTmpl,
		Data: data,
	}).Gen()
}

// genBashFile 生成
func (g *devops) genBashFile() ([]byte, error) {
	return (&gen.Impl{
		Name: "gen bashFile",
		Tmpl: bashTmpl,
	}).SetDelims("[[", "]]").Gen()
}

type DevopsCorn struct {
}

func NewDevopsCorn() *DevopsCorn {
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

			devops := func() error {
				var dockerfileBytes []byte
				var kubernetesBytes []byte
				var bashBytes []byte
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
				if kubernetesBytes, err = g.genKubernetes(&devopsValue{Server: server, Http: http, Grpc: grpc}); err != nil {
					return err
				}
				if bashBytes, err = g.genBashFile(); err != nil {
					return err
				}

				//write file
				if err := os.WriteFile(filepath.Join(path, server, "dockerfile"), dockerfileBytes, os.ModePerm); err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(path, server, "kubernetes.yaml"), kubernetesBytes, os.ModePerm); err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(path, "build.sh"), bashBytes, os.ModePerm); err != nil {
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
