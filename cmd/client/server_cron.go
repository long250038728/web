package client

import (
	_ "embed"
	"fmt"
	"github.com/long250038728/web/tool/gen"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

type ServerValue struct {
	Server string `json:"server" yaml:"server"`
	Page   string `json:"page" yaml:"page"`
}

type Server struct {
}

func NewServerGen() *Server {
	return &Server{}
}

//go:embed tmpl/server/main.tmpl
var main string

//go:embed tmpl/server/router.tmpl
var router string

//go:embed tmpl/server/service.tmpl
var service string

//go:embed tmpl/server/domain.tmpl
var domain string

//go:embed tmpl/server/repository.tmpl
var repository string

// GenMain 生成
func (g *Server) GenMain(data *ServerValue) ([]byte, error) {
	return (&gen.Impl{
		Name:     "gen main",
		Tmpl:     main,
		Data:     data,
		IsFormat: true,
	}).Gen()
}
func (g *Server) GenRouter(data *ServerValue) ([]byte, error) {
	return (&gen.Impl{
		Name:     "gen router",
		Tmpl:     router,
		Data:     data,
		IsFormat: true,
	}).Gen()
}

func (g *Server) GenService(data *ServerValue) ([]byte, error) {
	return (&gen.Impl{
		Name:     "gen service",
		Tmpl:     service,
		Data:     data,
		IsFormat: true,
	}).Gen()
}

func (g *Server) GenDomain(data *ServerValue) ([]byte, error) {
	return (&gen.Impl{
		Name:     "gen domain",
		Tmpl:     domain,
		Data:     data,
		IsFormat: true,
	}).Gen()
}

func (g *Server) GenRepository(data *ServerValue) ([]byte, error) {
	return (&gen.Impl{
		Name:     "gen repository",
		Tmpl:     repository,
		Data:     data,
		IsFormat: true,
	}).Gen()
}

type ServerCorn struct {
	path string
	page string
}

func NewServerCornCorn(path, page string) *ServerCorn {
	if len(path) == 0 {
		path = "/Users/linlong/Desktop/web/application"
	}

	if len(page) == 0 {
		page = "github.com/long250038728/web/application"
	}
	return &ServerCorn{path: path, page: page}
}

func (c *ServerCorn) Server() *cobra.Command {
	return &cobra.Command{
		Use:   "server [服务名] [http端口] [grpc端口] [版本号]",
		Short: "创建server： 请输入 [服务名] [http端口] [grpc端口] [版本号]",
		Long:  "创建server： 请输入 [服务名] [http端口] [grpc端口] [版本号]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			server := args[0]
			devops := func() error {
				var err error

				//mkdir path
				paths := []string{
					filepath.Join(c.path),
					filepath.Join(c.path, server),
					filepath.Join(c.path, server, "cmd"),
					filepath.Join(c.path, server, "internal"),
					filepath.Join(c.path, server, "internal", "domain"),
					filepath.Join(c.path, server, "internal", "model"),
					filepath.Join(c.path, server, "internal", "repository"),
					filepath.Join(c.path, server, "internal", "router"),
					filepath.Join(c.path, server, "internal", "service"),
				}
				for _, path := range paths {
					_, err = os.Stat(path)
					if os.IsNotExist(err) {
						if err := os.Mkdir(path, os.ModePerm); err != nil {
							return err
						}
					}
				}

				g := NewServerGen()
				var mainBytes []byte
				var routerBytes []byte

				var serverBytes []byte
				var domianBytes []byte
				var repositoryBytes []byte

				v := &ServerValue{Server: server, Page: c.page}

				if mainBytes, err = g.GenMain(v); err != nil {
					return err
				}
				if routerBytes, err = g.GenRouter(v); err != nil {
					return err
				}

				if serverBytes, err = g.GenService(v); err != nil {
					return err
				}
				if domianBytes, err = g.GenDomain(v); err != nil {
					return err
				}
				if repositoryBytes, err = g.GenRepository(v); err != nil {
					return err
				}

				//write file
				if err := os.WriteFile(filepath.Join(c.path, server, "cmd", "main.go"), mainBytes, os.ModePerm); err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(c.path, server, "internal", "router", "router.go"), routerBytes, os.ModePerm); err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(c.path, server, "internal", "service", "service.go"), serverBytes, os.ModePerm); err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(c.path, server, "internal", "domain", "domain.go"), domianBytes, os.ModePerm); err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(c.path, server, "internal", "repository", "repository.go"), repositoryBytes, os.ModePerm); err != nil {
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
