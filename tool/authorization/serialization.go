package authorization

import (
	"github.com/golang-jwt/jwt"
	"github.com/long250038728/web/tool/system_error"
)

type Serialization struct {
	SecretKey []byte
}

func (s *Serialization) SignedToken(c Claims) (token string, err error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(s.SecretKey)
}

func (s *Serialization) ParseToken(token string, c Claims, t TokenType) error {
	_, err := jwt.ParseWithClaims(token, c, func(token *jwt.Token) (interface{}, error) {
		return s.SecretKey, nil // 这里你需要提供用于签名的密钥
	})
	if err != nil {
		if validationErr, ok := err.(*jwt.ValidationError); ok && validationErr.Errors == jwt.ValidationErrorExpired {
			if t == AccessToken {
				err = system_error.AccessExpire
			}
			if t == RefreshToken {
				err = system_error.RefreshExpire
			}
		}
		return err
	}
	return nil
}