package authorization

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/long250038728/web/tool/store"
	"testing"
)

// UserSession 内部使用的信息
type UserSession struct {
	Id       int32  `json:"id" yaml:"id"`
	Name     string `json:"name" yaml:"name"`
	AuthList []string
}

func GetSessionKey(id int32) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("id:%d", id))) // 向哈希计算对象中写入字符串数据
	return hex.EncodeToString(hash.Sum(nil))
}

func TestSession(t *testing.T) {
	sess := NewSession(store.NewLocalStore(10))

	userSessionId := "123456"
	userSession := &UserSession{Id: 1, Name: "test"}

	if err := sess.SetSession(context.Background(), userSessionId, userSession); err != nil {
		t.Error(err)
		return
	}

	userSession = &UserSession{}
	if err := sess.GetSession(context.Background(), userSessionId, userSession); err != nil {
		t.Error(err)
		return
	}
	if err := sess.DeleteSession(context.Background(), userSessionId); err != nil {
		t.Error(err)
		return
	}
}
