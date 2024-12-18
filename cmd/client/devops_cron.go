package client

import (
	_ "embed"
	"fmt"
	"github.com/long250038728/web/tool/gen"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

type DevopsValue struct {
	Server string `json:"server" yaml:"server"`
	Http   string `json:"http" yaml:"http"`
	Grpc   string `json:"grpc" yaml:"grpc"`
}
type KubernetesValue struct {
	*DevopsValue
	Version string `json:"version" yaml:"version"`
	Hub     string `json:"hub" yaml:"hub"`
}

type Devops struct {
}

func NewDevopsGen() *Devops {
	return &Devops{}
}

//go:embed tmpl/devops/dockerfile.tmpl
var dockerfileTmpl string

//go:embed tmpl/devops/kubernetes.tmpl
var kubernetesTmpl string

// GenDockerfile 生成
func (g *Devops) GenDockerfile(data *DevopsValue) ([]byte, error) {
	return (&gen.Impl{
		Name: "gen dockerfile",
		Tmpl: dockerfileTmpl,
		Data: data,
	}).Gen()
}
func (g *Devops) GenKubernetes(data *KubernetesValue) ([]byte, error) {
	return (&gen.Impl{
		Name: "gen kubernetes",
		Tmpl: kubernetesTmpl,
		Data: data,
	}).Gen()
}

type DevopsCorn struct {
	hub  string
	path string
}

func NewDevopsCorn(hub, path string) *DevopsCorn {
	if len(hub) == 0 {
		hub = "ccr.ccs.tencentyun.com/linl"
	}
	if len(path) == 0 {
		path = "./devops"
	}
	return &DevopsCorn{hub: hub, path: path}
}

func (c *DevopsCorn) Devops() *cobra.Command {
	return &cobra.Command{
		Use:   "devops [服务名] [http端口] [grpc端口] [版本号]",
		Short: "创建devops： 请输入 [服务名] [http端口] [grpc端口] [版本号]",
		Long:  "创建devops： 请输入 [服务名] [http端口] [grpc端口] [版本号]",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			server := args[0]
			http := args[1]
			grpc := args[2]

			version := "v1"
			if len(args) >= 4 {
				version = args[3]
			}

			devops := func() error {
				var dockerfileBytes []byte
				var kubernetesBytes []byte
				var err error

				g := NewDevopsGen()

				if dockerfileBytes, err = g.GenDockerfile(&DevopsValue{Server: server, Http: http, Grpc: grpc}); err != nil {
					return err
				}
				if kubernetesBytes, err = g.GenKubernetes(&KubernetesValue{DevopsValue: &DevopsValue{Server: server, Http: http, Grpc: grpc}, Version: version, Hub: c.hub}); err != nil {
					return err
				}

				//mkdir path
				_, err = os.Stat(filepath.Join(c.path))
				if os.IsNotExist(err) {
					if err := os.Mkdir(filepath.Join(c.path), os.ModePerm); err != nil {
						return err
					}
				}
				_, err = os.Stat(filepath.Join(c.path, server))
				if os.IsNotExist(err) {
					if err := os.Mkdir(filepath.Join(c.path, server), os.ModePerm); err != nil {
						return err
					}
				}

				//write file
				if err := os.WriteFile(filepath.Join("./", c.path, server, "dockerfile"), dockerfileBytes, os.ModePerm); err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join("./", c.path, server, "k8s.yaml"), kubernetesBytes, os.ModePerm); err != nil {
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
