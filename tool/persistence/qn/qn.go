package qn

import (
	"context"
	"fmt"
	"github.com/long250038728/web/tool/server/http"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"io"
	"math/rand"
	"os"
	"strings"
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

func (qn *Qn) UpLoad(ctx context.Context, bucket, path, fileName string) error {
	response := &storage.PutRet{} //响应
	return storage.NewFormUploader(&storage.Config{}).PutFile(ctx, response, qn.Token(bucket), fileName, path, nil)
}

func (qn *Qn) UpLoadUrl(ctx context.Context, bucket, url, fileName string) error {
	filePaths := strings.Split(url, "/")
	tmpPath := fmt.Sprintf("./%d_%s", rand.Int(), filePaths[len(filePaths)-1])

	b, _, err := http.NewClient().Get(ctx, url, nil)
	if err != nil {
		return err
	}
	if err = os.WriteFile(tmpPath, b, os.ModePerm); err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(tmpPath)
	}()
	return qn.UpLoad(ctx, bucket, tmpPath, fileName)
}

func (qn *Qn) Download(ctx context.Context, bucket, path, fileName string) error {
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
	return os.WriteFile(path, body, os.ModePerm)
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
