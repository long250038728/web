package client

// https://gitee.com/api/v5/swagger#/postV5ReposOwnerRepoBranches
import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

var token = "5f8aaa1e024cad5e24e86fda85c57f49"
var jenkins = "admin:11fbfc1aab366147522f497f6c7d48b2ca"

var productList = []string{
	"zhubaoe/locke",
	//"zhubaoe-go/kobe",
	//"zhubaoe/socrates",
	"zhubaoe/aristotle",
	"fissiongeek/h5-sales",
	//"zhubaoe/hume",
}

// 创建feature/release分支 （创建分支）
//var createSource = "master"
//var createTarget = "release/v3.5.53"

// //提pr （创建pr）
//var createSource = "feature/sm0403"
//var createTarget = "release/v3.5.53"

// //合到check （创建pr）
//var createSource = "release/v3.5.53"
//var createTarget = "check"

// //合到master （创建pr）
var createSource = "release/v3.5.53"
var createTarget = "master"

// ==============================http=======================
// 创建分支
func TestCreateFeature_Http(t *testing.T) {
	client := NewGiteeClinet(token)
	ctx := context.Background()
	for _, addr := range productList {
		if err := client.CreateFeature(ctx, addr, createSource, createTarget); err != nil {
			t.Error(addr, err)
			continue
		}
		t.Log(addr, "ok")
	}
}

// 创建PR
func TestCreatePR_Http(t *testing.T) {
	pr := make([]string, 0, len(productList))
	errs := make([]string, 0, len(productList))

	client := NewGiteeClinet(token)
	ctx := context.Background()

	for _, addr := range productList {
		if createTarget == "check" && addr == "zhubaoe/locke" { //locke 没有check
			continue
		}
		info, err := client.CreatePR(ctx, addr, createSource, createTarget)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		pr = append(pr, info.HtmlUrl)
	}

	if len(pr) == 0 {
		t.Log("没有pr成功")
	}

	prText := strings.Join(pr, "\n")
	errText := strings.Join(errs, "\n")
	text := fmt.Sprintf("(%s 合并到 %s)\nPR创建成功(%d):\n%s\n", createSource, createTarget, len(pr), prText)
	if len(errs) > 0 {
		text = fmt.Sprintf("%sPR创建失败(%d):\n%s\n", text, len(errs), errText)
	}

	hookClient := NewQyHookClient("991cbde3-6963-4adc-a25c-7a6402ab7d38")
	err := hookClient.sendHook(context.Background(), text, []string{"18575538087"})
	t.Log(err, text)
}

// 获取pr地址
func TestGetPrAddr_Http(t *testing.T) {
	pr := make([]string, 0, len(productList))

	client := NewGiteeClinet(token)
	ctx := context.Background()
	for _, addr := range productList {
		list, err := client.GetPR(ctx, addr, createSource, createTarget)
		if err != nil {
			t.Error(err, addr)
			continue
		}
		if len(list) != 1 {
			t.Error(errors.New("list num is not one"), len(list), addr)
			continue
		}
		pr = append(pr, list[0].HtmlUrl)
	}

	for _, p := range pr {
		fmt.Println(p)
	}
}

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

//============================上线流程生产=========================

// 上线流程
func TestPrMerge_Gen(t *testing.T) {
	var address = []string{
		"https://gitee.com/zhubaoe/locke/pulls/121",
		"https://gitee.com/zhubaoe/aristotle/pulls/914",
		"https://gitee.com/fissiongeek/h5-sales/pulls/385",
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
	fmt.Println(os.WriteFile("./out/online.md", gen, os.ModePerm))
}
