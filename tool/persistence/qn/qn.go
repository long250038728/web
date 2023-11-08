package qn

import (
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Config struct {
	AccessKey string
	SecretKey string
}

type Qn struct {
	mac *qbox.Mac
}

func NewQn(config Config) *Qn {
	return &Qn{
		qbox.NewMac(config.AccessKey, config.SecretKey),
	}
}

func (qn *Qn) Token(bucket string) string {
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	return putPolicy.UploadToken(qn.mac)
}

func (qn *Qn) UpLoad(context context.Context, bucket, address string, fileName string) error {
	//配置
	uploader := storage.NewFormUploader(&storage.Config{})

	//额外信息
	putExtra := &storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	//响应
	response := &storage.PutRet{}

	return uploader.PutFile(context, response, qn.Token(bucket), fileName, address, putExtra)
}
