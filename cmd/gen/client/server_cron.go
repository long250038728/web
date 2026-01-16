package client

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/long250038728/web/tool/gen"
	"github.com/spf13/cobra"
)

type serverValue struct {
	Server string `json:"server" yaml:"server"`
	Domain string `json:"domain" yaml:"domain"`
	Page   string `json:"page" yaml:"page"`
	Protoc string `json:"protoc" yaml:"protoc"`
}

type serverGen struct {
}

func newServerGen() *serverGen {
	return &serverGen{}
}

//go:embed tmpl/server/main.tmpl
var main string

//go:embed tmpl/server/server.tmpl
var server string

//go:embed tmpl/server/handles.tmpl
var handles string

//go:embed tmpl/server/service.tmpl
var service string

//go:embed tmpl/server/domain.tmpl
var domain string

//go:embed tmpl/server/repository.tmpl
var repository string

// genMain 生成
func (g *serverGen) genMain(data *serverValue) ([]byte, error) {
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

func (g *serverGen) genServer(data *serverValue) ([]byte, error) {
	return (&gen.Impl{
		Name:     "gen server",
		Tmpl:     server,
		Data:     data,
		IsFormat: true,
		Func: template.FuncMap{
			"serverNameFunc": g.serverName,
		},
	}).Gen()
}

func (g *serverGen) genHandles(data *serverValue) ([]byte, error) {
	return (&gen.Impl{
		Name:     "gen handles",
		Tmpl:     handles,
		Data:     data,
		IsFormat: true,
		Func: template.FuncMap{
			"serverNameFunc": g.serverName,
		},
	}).Gen()
}
func (g *serverGen) genService(data *serverValue) ([]byte, error) {
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
func (g *serverGen) genDomain(data *serverValue) ([]byte, error) {
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
func (g *serverGen) genRepository(data *serverValue) ([]byte, error) {
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

func (g *serverGen) serverName(server string) string {
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
	//go run main.go server test /Users/linlong/Desktop/web/application
	//go run main.go server test /Users/linlong/Desktop/web/application github.com/long250038728/web application protoc
	return &cobra.Command{
		Use:   "server [服务名] [输出路径] [module-path:默认github.com/long250038728/web] [项目相对路径:默认application] [proto相对路径:默认protoc]",
		Short: "创建server： 请输入 [服务名] [输出路径] [module-path:默认github.com/long250038728/web] [项目相对路径:默认application] [proto相对路径:默认protoc]",
		Long:  "创建server： 请输入 [服务名] [输出路径] [module-path:默认github.com/long250038728/web] [项目相对路径:默认application] [proto相对路径:默认protoc]",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			modulePath := "github.com/long250038728/web"
			application := "application"
			protoc := "protoc"

			if len(args) >= 4 {
				modulePath = args[3]
			}
			if len(args) >= 5 {
				application = args[4]
			}
			if len(args) >= 6 {
				protoc = args[5]
			}

			serverName := args[0]
			domainName := args[1]
			c.path = args[2]
			c.page = filepath.Join(modulePath, application)
			c.protoc = filepath.Join(modulePath, protoc)

			devops := func() error {
				var err error

				//mkdir path
				paths := []string{
					filepath.Join(c.path),
					filepath.Join(c.path, serverName),
					filepath.Join(c.path, serverName, "cmd"),
					filepath.Join(c.path, serverName, "cmd", "server"),

					filepath.Join(c.path, serverName, "internal"),
					filepath.Join(c.path, serverName, "internal", "handles"),
					filepath.Join(c.path, serverName, "internal", "service"),
					//filepath.Join(c.path, serverName, "internal", "validate"),

					filepath.Join(c.path, serverName, "internal", "biz"),
					filepath.Join(c.path, serverName, "internal", "biz", domainName),
					filepath.Join(c.path, serverName, "internal", "biz", domainName, "domain"),
					filepath.Join(c.path, serverName, "internal", "biz", domainName, "model"),
					filepath.Join(c.path, serverName, "internal", "biz", domainName, "entity"),
					filepath.Join(c.path, serverName, "internal", "biz", domainName, "repository"),
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
				var serverBytes []byte

				var handlesBytes []byte
				var serviceBytes []byte
				var domainBytes []byte
				var repositoryBytes []byte

				v := &serverValue{Server: serverName, Domain: domainName, Page: c.page, Protoc: c.protoc}

				if mainBytes, err = g.genMain(v); err != nil {
					return err
				}
				if serverBytes, err = g.genServer(v); err != nil {
					return err
				}
				if handlesBytes, err = g.genHandles(v); err != nil {
					return err
				}
				if serviceBytes, err = g.genService(v); err != nil {
					return err
				}
				if domainBytes, err = g.genDomain(v); err != nil {
					return err
				}
				if repositoryBytes, err = g.genRepository(v); err != nil {
					return err
				}

				// 辅助函数：检查文件是否存在，不存在则写入
				writeIfNotExist := func(filePath string, data []byte) error {
					_, err := os.Stat(filePath)
					if os.IsNotExist(err) {
						// 文件不存在，写入
						return os.WriteFile(filePath, data, os.ModePerm)
					}
					// 文件存在，跳过写入
					fmt.Printf("文件已存在，跳过写入: %s\n", filePath)
					return nil
				}

				// 写入各个文件，跳过已存在的文件
				if err := writeIfNotExist(filepath.Join(c.path, serverName, "cmd", "main.go"), mainBytes); err != nil {
					return err
				}
				if err := writeIfNotExist(filepath.Join(c.path, serverName, "cmd", "server", "server.go"), serverBytes); err != nil {
					return err
				}
				if err := writeIfNotExist(filepath.Join(c.path, serverName, "internal", "handles", "handles.go"), handlesBytes); err != nil {
					return err
				}
				if err := writeIfNotExist(filepath.Join(c.path, serverName, "internal", "service", serverName+".go"), serviceBytes); err != nil {
					return err
				}
				if err := writeIfNotExist(filepath.Join(c.path, serverName, "internal", "biz", domainName, "domain", domainName+".go"), domainBytes); err != nil {
					return err
				}
				if err := writeIfNotExist(filepath.Join(c.path, serverName, "internal", "biz", domainName, "repository", domainName+".go"), repositoryBytes); err != nil {
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
