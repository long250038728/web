package authorization

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/long250038728/web/tool/system_error"
)

type Session struct {
	Store Store
}

func (p *Session) GetSession(ctx context.Context, sessionId string) (session *UserSession, err error) {
	if sessionId == "" {
		return nil, errors.New("sessionId is empty")
	}
	sessionStr, err := p.Store.Get(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	if len(sessionStr) == 0 {
		return nil, system_error.SessionExpire
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
	if ok, err = p.Store.Set(ctx, sessionId, string(b)); err != nil {
		return
	}
	if !ok {
		err = errors.New("session setting is err")
	}
	return
}

func (p *Session) DeleteSession(ctx context.Context, sessionId string) error {
	if sessionId == "" {
		return errors.New("sessionId is empty")
	}
	_, err := p.Store.Del(ctx, sessionId)
	return err
}
