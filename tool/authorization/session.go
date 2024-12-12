package authorization

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/long250038728/web/tool/app_error"
	"time"
)

type Session struct {
	LocalStore    Store
	Store         Store
	accessExpires time.Duration
}

func (p *Session) GetSession(ctx context.Context, sessionId string) (session *UserSession, err error) {
	if sessionId == "" {
		return nil, errors.New("sessionId is empty")
	}
	var sessionStr string

	// 按照优先级获取
	for _, s := range []Store{p.LocalStore, p.Store} {
		if s != nil {
			sessionStr, err = s.Get(ctx, sessionId)
			if err != nil {
				return nil, err
			}
			if len(sessionStr) > 0 {
				break
			}
		}
	}

	// 获取不到则报错
	if len(sessionStr) == 0 {
		return nil, app_error.SessionExpire
	}
	return session, json.Unmarshal([]byte(sessionStr), &session)
}

func (p *Session) SetSession(ctx context.Context, sessionId string, session *UserSession) (err error) {
	var b []byte
	var ok bool

	if sessionId == "" {
		return errors.New("sessionId is empty")
	}
	if b, err = json.Marshal(session); err != nil {
		return
	}

	// 数据添加到缓存中
	for _, s := range []Store{p.LocalStore, p.Store} {
		if s != nil {
			if ok, err = s.SetEX(ctx, sessionId, string(b), p.accessExpires); err != nil {
				return err
			}
			if !ok {
				err = errors.New("session setting is err")
				return
			}
		}
	}
	return
}

func (p *Session) DeleteSession(ctx context.Context, sessionId string) error {
	if sessionId == "" {
		return errors.New("sessionId is empty")
	}

	// 数据添加到缓存中
	for _, s := range []Store{p.LocalStore, p.Store} {
		if s != nil {
			if _, err := p.Store.Del(ctx, sessionId); err != nil {
				return err
			}
		}
	}
	return nil
}
