package authorization

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/long250038728/web/tool/app_error"
	"github.com/long250038728/web/tool/store"
	"time"
)

//
//type session[S any] struct {
//	store         store.Store
//	accessExpires time.Duration
//}
//
//func NewSession[S any](store store.Store, accessExpires time.Duration) Session {
//	return &session[S]{
//		store:         store,
//		accessExpires: accessExpires,
//	}
//}
//
//func (p *session[S]) GetSession(ctx context.Context, sessionId string) (session S, err error)
//	if sessionId == "" {
//		return nil, errors.New("sessionId is empty")
//	}
//	var sessionStr string
//
//	sessionStr, err = p.store.Get(ctx, sessionId)
//	if err != nil {
//		return nil, err
//	}
//
//	// 获取不到则报错
//	if len(sessionStr) == 0 {
//		return nil, app_error.SessionExpire
//	}
//	return session, json.Unmarshal([]byte(sessionStr), &session)
//}
//
//func (p *session[S]) SetSession(ctx context.Context, sessionId string, session S) (err error) {
//	var b []byte
//	var ok bool
//
//	if sessionId == "" {
//		return errors.New("sessionId is empty")
//	}
//	if b, err = json.Marshal(session); err != nil {
//		return
//	}
//
//	if ok, err = p.store.Set(ctx, sessionId, string(b), p.accessExpires); err != nil {
//		return err
//	}
//	if !ok {
//		err = errors.New("session setting is err")
//		return
//	}
//	return
//}
//
//func (p *session[T]) DeleteSession(ctx context.Context, sessionId string) error {
//	if sessionId == "" {
//		return errors.New("sessionId is empty")
//	}
//	_, err := p.store.Del(ctx, sessionId)
//	return err
//}

type SessionOpts func(s *session)

type session struct {
	store         store.Store
	accessExpires time.Duration
}

func SetAccessExpires(accessExpires time.Duration) SessionOpts {
	return func(s *session) {
		s.accessExpires = accessExpires
	}
}

func NewSession(store store.Store, opts ...SessionOpts) Session {
	sess := &session{store: store, accessExpires: time.Hour}
	for _, opt := range opts {
		opt(sess)
	}
	return sess
}

func (p *session) GetSession(ctx context.Context, sessionId string, session any) error {
	if sessionId == "" {
		return errors.New("sessionId is empty")
	}

	var sessionStr string
	sessionStr, err := p.store.Get(ctx, sessionId)
	if err != nil {
		return err
	}

	if len(sessionStr) == 0 {
		return app_error.SessionExpire
	}

	if err = json.Unmarshal([]byte(sessionStr), &session); err != nil {
		return err
	}

	return nil
}

func (p *session) SetSession(ctx context.Context, sessionId string, session any) (err error) {
	if sessionId == "" {
		err = errors.New("sessionId is empty")
		return
	}
	var b []byte
	if b, err = json.Marshal(session); err != nil {
		return
	}
	var ok bool
	ok, err = p.store.Set(ctx, sessionId, string(b), p.accessExpires)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New("session setting is err")
		return
	}

	return
}

func (p *session) DeleteSession(ctx context.Context, sessionId string) error {
	if sessionId == "" {
		return errors.New("sessionId is empty")
	}
	_, err := p.store.Del(ctx, sessionId)
	return err
}
