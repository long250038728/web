package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/long250038728/web/tool/auth"
)

type Session struct {
	Store auth.Store
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
		return nil, errors.New("sessionId is empty")
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
