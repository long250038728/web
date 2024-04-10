package gitee

import (
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
	Token    string
	Features []*feature
	Kobe     []string
}

type Pr struct {
	Token string
}

func NewPrGen(token string) *Pr {
	return &Pr{token}
}

func (g *Pr) GenFeature(address []string, source, target string) ([]byte, error) {
	if len(address) == 0 {
		return nil, errors.New("address num is error")
	}
	var list = make([]*feature, 0, len(address))
	for _, addr := range address {
		list = append(list, &feature{Addr: addr, Source: source, Feature: target})
	}
	return (&gen.Impl{Name: "gen feature", TmplPath: "./gitee_feature.tmpl", Data: &features{Features: list, Token: g.Token}, IsFormat: false}).Gen()
}

func (g *Pr) GenPrCreate(address []string, source, target string) ([]byte, error) {
	if len(address) == 0 {
		return nil, errors.New("address num is error")
	}
	var list = make([]*feature, 0, len(address))
	for _, addr := range address {
		list = append(list, &feature{Addr: addr, Source: source, Feature: target})
	}
	return (&gen.Impl{Name: "gen pr create", TmplPath: "./gitee_pr.tmpl", Data: &features{Features: list, Token: g.Token}, IsFormat: false}).Gen()
}

func (g *Pr) GenMerge(address []string, kobe []string) ([]byte, error) {
	if len(address) == 0 {
		return nil, errors.New("address num is error")
	}
	var list = make([]*feature, 0, len(address))
	for _, addr := range address {
		list = append(list, &feature{Addr: addr})
	}
	return (&gen.Impl{Name: "gen pr merge", TmplPath: "./gitee_pr_merge.tmpl", Func: template.FuncMap{
		"objectName": g.objectName,
	}, Data: &features{Features: list, Token: g.Token, Kobe: kobe}, IsFormat: false}).Gen()
}

func (g *Pr) objectName(mainString, obj string) bool {
	str := strings.ReplaceAll(mainString, "https://", "")
	return strings.Split(str, "/")[2] == obj
}
