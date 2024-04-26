package qn

import (
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"io"
	"os"
)

type Config struct {
	AccessKey string `json:"access_key" yaml:"accessKey"`
	SecretKey string `json:"secret_key" yaml:"secretKey"`
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

func (qn *Qn) UpLoad(ctx context.Context, bucket, address string, fileName string) error {
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

	return uploader.PutFile(ctx, response, qn.Token(bucket), fileName, address, putExtra)
}

func (qn *Qn) Download(ctx context.Context, bucket, address string, fileName string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	bucketManager := storage.NewBucketManager(qn.mac, &storage.Config{})
	resp, err := bucketManager.Get(bucket, fileName, &storage.GetObjectInput{})
	if err != nil || resp == nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return os.WriteFile(address, body, os.ModePerm)
}

func (qn *Qn) Delete(ctx context.Context, bucket, fileName string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	bucketManager := storage.NewBucketManager(qn.mac, &storage.Config{})
	return bucketManager.Delete(bucket, fileName)
}
