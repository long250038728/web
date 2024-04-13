package client

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
	JenkinsToken string
	GiteeToken   string
	Features     []*feature
	Kobe         []string
}

type Pr struct {
	GiteeToken   string
	JenkinsToken string
}

func NewPrGen(giteeToken, jenkinsToken string) *Pr {
	return &Pr{GiteeToken: giteeToken, JenkinsToken: jenkinsToken}
}

func (g *Pr) GenFeature(address []string, source, target string) ([]byte, error) {
	if len(address) == 0 {
		return nil, errors.New("address num is error")
	}
	var list = make([]*feature, 0, len(address))
	for _, addr := range address {
		list = append(list, &feature{Addr: addr, Source: source, Feature: target})
	}
	return (&gen.Impl{Name: "gen feature", TmplPath: "./tmpl/gitee_feature.tmpl", Data: &features{Features: list, GiteeToken: g.GiteeToken}, IsFormat: false}).Gen()
}

func (g *Pr) GenPrCreate(address []string, source, target string) ([]byte, error) {
	if len(address) == 0 {
		return nil, errors.New("address num is error")
	}
	var list = make([]*feature, 0, len(address))
	for _, addr := range address {
		list = append(list, &feature{Addr: addr, Source: source, Feature: target})
	}
	return (&gen.Impl{Name: "gen pr create", TmplPath: "./tmpl/gitee_pr.tmpl", Data: &features{Features: list, GiteeToken: g.GiteeToken}, IsFormat: false}).Gen()
}

var tmpl = `
//1.合并kobe kobe改tag         构建
//2.合并locke                  构建
//3.其他项目合并
//3.商户管家合并

http://111.230.143.16:8081/ 用户名：admin 密码：admin@zhubaoe
https://jenkins.zhubaoe.cn/ 用户名：admin 密码：admin@zhubaoe_new

{{ $kobe := .Kobe}}{{ $giteeToken := .GiteeToken}} {{ $jenkinsToken := .JenkinsToken}}
{{- range $index,$item := .Features}}
============================ {{ $index }} {{name $item.Addr}} pr合并 ================================
curl -X POST --header 'Content-Type: application/json;charset=UTF-8' \
'{{$item.Addr}}' \
-d '{"access_token":"{{$giteeToken}}","merge_method":"merge"}'

{{- if  objectName $item.Addr "kobe"}}

!!!!!!!!!!改tag!!!!!!!!!!
https://gitee.com/zhubaoe-go/kobe

===================== 构建: =====================
{{- range $kobe_index,$kobe_item := $kobe}}
curl -X POST http://111.230.143.16:8081/job/{{$kobe_item}}/buildWithParameters \
--user {{ $jenkinsToken }} \
--data-urlencode "BRANCH=master" \
--data-urlencode "SYSTEM=root@172.16.0.34"

curl -X POST http://111.230.143.16:8081/job/{{$kobe_item}}/buildWithParameters \
--user {{ $jenkinsToken }} \
--data-urlencode "BRANCH=master" \
--data-urlencode "SYSTEM=root@172.16.0.9"

{{end}}
{{- end }}
{{- if  objectName $item.Addr "locke"}}

== 构建: ==
curl -X POST http://111.230.143.16:8081/job/locke-prod_32/build \
--user {{ $jenkinsToken }}

curl -X POST http://111.230.143.16:8081/job/locke-prod_64/build \
--user {{ $jenkinsToken }}

curl -X POST http://111.230.143.16:8081/job/locke-hot-prod-64/build \
--user {{ $jenkinsToken }}
{{ end }}
{{end}}
`

func (g *Pr) GenMerge(address []string, kobe []string) ([]byte, error) {
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
	}, Data: &features{Features: list, GiteeToken: g.GiteeToken, JenkinsToken: g.JenkinsToken, Kobe: kobe}, IsFormat: false}).Gen()
}

func (g *Pr) objectName(mainString, obj string) bool {
	str := strings.ReplaceAll(mainString, "https://", "")
	return strings.Split(str, "/")[2] == obj
}

func (g *Pr) name(mainString string) string {
	str := strings.ReplaceAll(mainString, "https://", "")
	return strings.Split(str, "/")[2]
}
