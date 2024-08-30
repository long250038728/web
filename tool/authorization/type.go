package authorization

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) (bool, error)
}

type Claims interface {
	Valid() error
}

func GetSessionId(id int32) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("id:%d", id))) // 向哈希计算对象中写入字符串数据
	return hex.EncodeToString(hash.Sum(nil))
}
