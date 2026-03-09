package authorization

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
)

//
//import (
//	"errors"
//	"github.com/golang-jwt/jwt/v4"
//	"github.com/long250038728/web/tool/app_error"
//)
//
//type Serialization struct {
//	SecretKey []byte
//}
//
//func (s *Serialization) SignedToken(c Claims) (token string, err error) {
//	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(s.SecretKey)
//}
//
//func (s *Serialization) ParseToken(token string, c Claims, t TokenType) error {
//	_, err := jwt.ParseWithClaims(token, c, func(token *jwt.Token) (interface{}, error) {
//		return s.SecretKey, nil // 这里你需要提供用于签名的密钥
//	})
//	if err != nil {
//		var validationErr *jwt.ValidationError
//		if errors.As(err, &validationErr) && validationErr.Errors == jwt.ValidationErrorExpired {
//			if t == AccessToken {
//				err = app_error.AccessExpire
//			}
//			if t == RefreshToken {
//				err = app_error.RefreshExpire
//			}
//		}
//		return err
//	}
//	return nil
//}

type token struct {
	SecretKey []byte
}

type TokenOpts func(s *token)

func SetSecretKey(secretKey []byte) TokenOpts {
	return func(s *token) {
		s.SecretKey = secretKey
	}
}

func NewToken(opts ...TokenOpts) Token {
	t := &token{SecretKey: []byte("temp")}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

//===============================通过Claims 生成token =============================

// Signed token生成
func (a *token) Signed(ctx context.Context, access Claims) (token string, err error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, access).SignedString(a.SecretKey)
}

//===============================通过Claims 生成token =============================

func (a *token) Parse(ctx context.Context, token string, claims Claims) error {
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return a.SecretKey, nil // 这里你需要提供用于签名的密钥
	})
	return err
}
