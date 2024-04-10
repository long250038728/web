package gitee

import (
	"fmt"
	"testing"
)

var token = "5f8aaa1e024cad5e24e86fda85c57f49"

func TestFeature_Gen(t *testing.T) {
	var address = []string{
		"zhubaoe/socrates",
		"zhubaoe/aristotle",
		"zhubaoe/locke",
		"fissiongeek/h5-sales",
	}
	var source = "master"
	var target = "release/v3.5.60"

	gen, err := NewPrGen(token).GenFeature(address, source, target)
	if err != nil {
		return
	}
	fmt.Println(string(gen))
}

func TestPr_Gen(t *testing.T) {
	var address = []string{
		"zhubaoe/socrates",
		"zhubaoe/aristotle",
		"zhubaoe/locke",
		"fissiongeek/h5-sales",
		"zhubaoe-go/kobe",
	}
	var source = "feature/sm0501"
	var target = "master"

	gen, err := NewPrGen(token).GenPrCreate(address, source, target)
	if err != nil {
		return
	}
	fmt.Println(string(gen))
}

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
	gen, err := NewPrGen(token).GenMerge(address, kobe)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(gen))
}
