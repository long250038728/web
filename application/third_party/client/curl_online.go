package client

import (
	_ "embed"
	"errors"
	"github.com/long250038728/web/tool/gen"
	"strings"
	"text/template"
)

// Addr
// Source
// Feature
type feature struct {
	Addr    string `json:"addr"`
	Source  string `json:"source"`
	Feature string `json:"feature"`
}

type features struct {
	JenkinsToken string
	GiteeToken   string
	Features     []*feature
	Svc          *Svc
}

type Pr struct {
	GiteeToken   string
	JenkinsToken string
}

func NewPrGen(giteeToken, jenkinsToken string) *Pr {
	return &Pr{GiteeToken: giteeToken, JenkinsToken: jenkinsToken}
}

//go:embed script/online.tmpl
var tmpl string

type Svc struct {
	Kobe  []string `json:"kobe" yaml:"kobe"`
	Marx  []string `json:"marx" yaml:"marx"`
	Shell string   `json:"shell" yaml:"shell"`
	SQL   string   `json:"sql" yaml:"sql"`
}

func (g *Pr) GenMerge(address []string, svc *Svc) ([]byte, error) {
	if len(address) == 0 {
		return nil, errors.New("address num is error")
	}
	var list = make([]*feature, 0, len(address))
	for _, addr := range address {
		list = append(list, &feature{Addr: addr})
	}
	return (&gen.Impl{Name: "gen pr merge", Tmpl: tmpl, Func: template.FuncMap{
		"objectName": g.objectName,
		"name":       g.name,
	}, Data: &features{Features: list, GiteeToken: g.GiteeToken, JenkinsToken: g.JenkinsToken, Svc: svc}, IsFormat: false}).Gen()
}

func (g *Pr) objectName(mainString, obj string) bool {
	str := strings.ReplaceAll(mainString, "https://", "")
	return strings.Split(str, "/")[5] == obj
}

func (g *Pr) name(mainString string) string {
	str := strings.ReplaceAll(mainString, "https://", "")
	return strings.Split(str, "/")[5]
}
