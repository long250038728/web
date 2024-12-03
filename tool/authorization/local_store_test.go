package authorization

import (
	"context"
	"testing"
	"time"
)

func TestLocalStore(t *testing.T) {
	s, err := NewLocalStore(10)
	if err != nil {
		t.Error(err)
		return
	}
	evicted, err := s.SetEX(context.Background(), "hello", "world", time.Minute)
	if err != nil {
		t.Error(err)
		return
	}
	if evicted {
		t.Errorf("%s", "set is not ok")
		return
	}

	val, err := s.Get(context.Background(), "hello")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(val)

	ok, err := s.Del(context.Background(), "hello")
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
