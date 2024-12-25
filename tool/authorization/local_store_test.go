package authorization

import (
	"context"
	"testing"
	"time"
)

func TestLocalStore(t *testing.T) {
	s := NewLocalStore(5 * 1024 * 1024)
	ok, err := s.SetEX(context.Background(), "hello", "world", time.Minute)
	if err != nil {
		t.Error(err)
		return
	}
	if !ok {
		t.Errorf("%s", "set is not ok")
		return
	}

	val, err := s.Get(context.Background(), "hello")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(val)

	ok, err = s.Del(context.Background(), "hello")
	if err != nil {
		t.Error(err)
		return
	}
	if !ok {
		t.Errorf("%s", "del is not ok")
		return
	}
	t.Log("ok")
}
