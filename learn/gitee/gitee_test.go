package gitee

// https://gitee.com/api/v5/swagger#/postV5ReposOwnerRepoBranches
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/server/http"
	"testing"
)

var token = "5f8aaa1e024cad5e24e86fda85c57f49"
var jenkins = "admin:11fbfc1aab366147522f497f6c7d48b2ca"

var productList = []string{
	//"zhubaoe/socrates",
	"zhubaoe/aristotle",
	//"zhubaoe/locke",
	//"fissiongeek/h5-sales",
}
var createSource = "hotfix/reshapre_20240410"
var createTarget = "master"

//============================template curl=========================

// 创建分支
func TestCreateFeatureCurl_Gen(t *testing.T) {
	gen, err := NewPrGen(token, jenkins).GenFeature(productList, createSource, createTarget)
	if err != nil {
		return
	}
	fmt.Println(string(gen))
}

// 创建PR
func TestCreatePR_Gen(t *testing.T) {
	gen, err := NewPrGen(token, jenkins).GenPrCreate(productList, createSource, createTarget)
	if err != nil {
		return
	}
	fmt.Println(string(gen))
}

// ==============================http=======================
type PrListItem struct {
	HtmlUrl string `json:"html_url"`
}

// 创建分支
func TestCreateFeature_Http(t *testing.T) {
	for _, addr := range productList {
		client := http.NewClient()
		data := map[string]any{
			"access_token": token,
			"refs":         createSource,
			"branch_name":  createTarget,
		}
		_, _, err := client.Post(context.Background(), fmt.Sprintf("https://gitee.com/api/v5/repos/%s/branches", addr), data)
		if err != nil {
			t.Error(addr, err)
			continue
		}
		t.Log(addr, "ok")
	}
}

// 创建PR
func TestCreatePR_Http(t *testing.T) {
	pr := make([]string, 0, len(productList))

	for _, addr := range productList {
		client := http.NewClient()
		data := map[string]any{
			"access_token": token,
			"title":        fmt.Sprintf("%s merge %s", createTarget, createSource),
			"head":         createSource,
			"base":         createTarget,
		}
		b, _, err := client.Post(context.Background(), fmt.Sprintf("https://gitee.com/api/v5/repos/%s/pulls", addr), data)
		if err != nil {
			t.Error(addr, err)
			continue
		}

		var item *PrListItem
		err = json.Unmarshal(b, &item)
		if err != nil {
			t.Error(err)
			return
		}
		pr = append(pr, item.HtmlUrl)
	}

	for _, p := range pr {
		fmt.Println(p)
	}
}

// 获取pr地址
func TestGetPrAddr_Http(t *testing.T) {
	pr := make([]string, 0, len(productList))

	for _, addr := range productList {
		url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/pulls?access_token=%s&state=all&head=%s&base=%s&sort=created&direction=desc&page=1&per_page=20",
			addr, token, createSource, createTarget)
		client := http.NewClient()
		b, code, err := client.Get(context.Background(), url, nil)
		if err != nil {
			t.Error(err)
			return
		}
		if code != 200 {
			t.Error(errors.New("request code is error"))
			return
		}

		var list []*PrListItem
		err = json.Unmarshal(b, &list)
		if err != nil {
			t.Error(err)
			return
		}
		if len(list) != 1 {
			t.Error(errors.New("list num is not one"))
			return
		}

		pr = append(pr, list[0].HtmlUrl)
	}

	for _, p := range pr {
		fmt.Println(p)
	}
}

//============================上线流程生产=========================

// 上线流程
func TestPrMerge_Gen(t *testing.T) {
	var address = []string{
		"https://gitee.com/zhubaoe/locke/pulls/365",
		"https://gitee.com/zhubaoe-go/kobe/pulls/913",
		"https://gitee.com/zhubaoe/aristome/pulls/365",
		"https://gitee.com/zhubaoe/socrete/pulls/365",
	}
	var kobe = []string{
		"kobe-order",
		"kobe-stock",
		"kobe-customer",
	}
	gen, err := NewPrGen(token, jenkins).GenMerge(address, kobe)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(gen))
}
