package authorization

import (
	"errors"
	"github.com/golang-jwt/jwt"
)

type Token struct {
	SecretKey []byte
}

func (p *Token) SignedToken(c Claims) (token string, err error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(p.SecretKey)
}

func (p *Token) ParseToken(token string, c Claims) error {
	_, err := jwt.ParseWithClaims(token, c, func(token *jwt.Token) (interface{}, error) {
		return p.SecretKey, nil // 这里你需要提供用于签名的密钥
	})
	if err != nil {
		if validationErr, ok := err.(*jwt.ValidationError); ok && validationErr.Errors == jwt.ValidationErrorExpired {
			err = errors.New("token is Disabled")
		}
		return err
	}
	return nil
}
