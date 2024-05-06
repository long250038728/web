package qn

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

var config = Config{AccessKey: "VUxpsfFXUsc3UmhEkE9739o3hG_sqPnsKjLMcWd1", SecretKey: "nGZYzXEC1FRy6hA1olOIyUYokc7Z2FXV4y0kS_J_"}

func TestQn_Token(t *testing.T) {
	qn := NewQn(config)
	token := qn.Token("zhubaoe-hn")
	t.Log(token)
}

func TestQn_Upload(t *testing.T) {
	qn := NewQn(config)
	err := qn.UpLoad(context.Background(), "zhubaoe-hn", "/Users/linlong/Desktop/111.png", "goods/346.png")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("ok")
}

func TestQn_Download(t *testing.T) {
	qn := NewQn(config)
	err := qn.Download(context.Background(), "zhubaoe-hn", "/Users/linlong/Desktop/123.png", "goods/346.png")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("ok")
}

func TestQn_Delete(t *testing.T) {
	qn := NewQn(config)
	err := qn.Delete(context.Background(), "zhubaoe-hn", "goods/346.png")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("ok")
}

func TestQn_Uploads(t *testing.T) {
	//qn := NewQn(configurator)

	_ = filepath.Walk("/Users/linlong/Desktop", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		// 检查是否为文件
		if info.IsDir() {
			return nil
		}

		// 检查文件扩展名是否为png
		if strings.HasSuffix(strings.ToLower(info.Name()), ".png") {
			name := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
			fmt.Println(path)
			fmt.Println(name)
			fmt.Println("")
			//err := qn.UpLoad(context.Background(), "zhubaoe-hn", path, name)
			//if err != nil {
			//	return err
			//}
		}

		return nil
	})
}
