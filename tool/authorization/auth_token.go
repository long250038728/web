package authorization

import (
	"github.com/golang-jwt/jwt"
	"github.com/long250038728/web/tool/system_error"
)

type Token struct {
	SecretKey []byte
}

func (p *Token) SignedToken(c Claims) (token string, err error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(p.SecretKey)
}

func (p *Token) ParseToken(token string, c Claims, t int) error {
	_, err := jwt.ParseWithClaims(token, c, func(token *jwt.Token) (interface{}, error) {
		return p.SecretKey, nil // 这里你需要提供用于签名的密钥
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
