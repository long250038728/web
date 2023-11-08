package qn

import (
	"context"
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
	err := qn.UpLoad(context.Background(), "zhubaoe-hn", "/Users/linlong/Desktop/111.png", "goods/345.png")
	t.Log(err)
}
