package client

import (
	_ "embed"
	"fmt"
	"github.com/long250038728/web/tool/gen"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type serverValue struct {
	Server string `json:"server" yaml:"server"`
	Page   string `json:"page" yaml:"page"`
	Protoc string `json:"protoc" yaml:"protoc"`
}

type server struct {
}

func newServerGen() *server {
	return &server{}
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

// genMain 生成
func (g *server) genMain(data *serverValue) ([]byte, error) {
	return (&gen.Impl{
		Name:     "gen main",
		Tmpl:     main,
		Data:     data,
		IsFormat: true,
		Func: template.FuncMap{
			"serverNameFunc": g.serverName,
		},
	}).Gen()
}
func (g *server) genRouter(data *serverValue) ([]byte, error) {
	return (&gen.Impl{
		Name:     "gen router",
		Tmpl:     router,
		Data:     data,
		IsFormat: true,
		Func: template.FuncMap{
			"serverNameFunc": g.serverName,
		},
	}).Gen()
}
func (g *server) genService(data *serverValue) ([]byte, error) {
	return (&gen.Impl{
		Name:     "gen service",
		Tmpl:     service,
		Data:     data,
		IsFormat: true,
		Func: template.FuncMap{
			"serverNameFunc": g.serverName,
		},
	}).Gen()
}
func (g *server) genDomain(data *serverValue) ([]byte, error) {
	return (&gen.Impl{
		Name:     "gen domain",
		Tmpl:     domain,
		Data:     data,
		IsFormat: true,
		Func: template.FuncMap{
			"serverNameFunc": g.serverName,
		},
	}).Gen()
}
func (g *server) genRepository(data *serverValue) ([]byte, error) {
	return (&gen.Impl{
		Name:     "gen repository",
		Tmpl:     repository,
		Data:     data,
		IsFormat: true,
		Func: template.FuncMap{
			"serverNameFunc": g.serverName,
		},
	}).Gen()
}

func (g *server) serverName(server string) string {
	//对server字符串第一个转换为大写
	return fmt.Sprintf("%s%s", strings.ToUpper(server[:1]), server[1:])
}

type ServerCorn struct {
	path   string
	page   string
	protoc string
}

func NewServerCornCorn() *ServerCorn {
	return &ServerCorn{}
}

func (c *ServerCorn) Server() *cobra.Command {
	return &cobra.Command{
		Use:   "server [服务名]",
		Short: "创建server： 请输入 [服务名] [输出路径] [项目包名] [项目相对路径] [proto相对路径]",
		Long:  "创建server： 请输入 [服务名] [输出路径] [项目包名] [项目相对路径] [proto相对路径]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			server := args[0]
			if len(args) >= 2 {
				c.path = args[1]
			}

			webBase := "github.com/long250038728/web"
			application := "application"
			protoc := "protoc"

			if len(args) >= 3 {
				webBase = args[2]
			}
			if len(args) >= 4 {
				application = args[3]
			}
			if len(args) >= 5 {
				protoc = args[4]
			}
			c.page = filepath.Join(webBase, application)
			c.protoc = filepath.Join(webBase, protoc)

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

				g := newServerGen()
				var mainBytes []byte
				var routerBytes []byte

				var serverBytes []byte
				var domainBytes []byte
				var repositoryBytes []byte

				v := &serverValue{Server: server, Page: c.page, Protoc: c.protoc}

				if mainBytes, err = g.genMain(v); err != nil {
					return err
				}
				if routerBytes, err = g.genRouter(v); err != nil {
					return err
				}

				if serverBytes, err = g.genService(v); err != nil {
					return err
				}
				if domainBytes, err = g.genDomain(v); err != nil {
					return err
				}
				if repositoryBytes, err = g.genRepository(v); err != nil {
					return err
				}

				//write file
				if err := os.WriteFile(filepath.Join(c.path, server, "cmd", "main.go"), mainBytes, os.ModePerm); err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(c.path, server, "internal", "router", "router.go"), routerBytes, os.ModePerm); err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(c.path, server, "internal", "service", server+".go"), serverBytes, os.ModePerm); err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(c.path, server, "internal", "domain", server+"_domain.go"), domainBytes, os.ModePerm); err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(c.path, server, "internal", "repository", server+"_repository.go"), repositoryBytes, os.ModePerm); err != nil {
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
