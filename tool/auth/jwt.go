package auth

import (
	errors2 "errors"
	"github.com/golang-jwt/jwt"
	"github.com/long250038728/web/tool/struct_map"
	"time"
)

//github.com/golang-jwt/jwt

type userClaims struct {
	jwt.StandardClaims
	UserClaims
}

var secretKey = []byte("secret_key")

func Claims(c *UserClaims) (string, error) {
	//外部不带有 jwt.StandardClaims 的对象转换为带有 jwt.StandardClaims 的对象
	claims := &userClaims{}
	err := struct_map.Map(c, claims)
	if err != nil {
		return "", err
	}

	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // 创建JWT对象
	return token.SignedString(secretKey)                       // 使用密钥进行签名
}

func Parse(signedString string) (*UserClaims, error) {
	c := &UserClaims{}
	if len(signedString) == 0 {
		return c, nil
	}

	// 解析JWT字符串
	token, err := jwt.ParseWithClaims(signedString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil // 这里你需要提供用于签名的密钥
	})
	if err != nil {
		if validationErr, ok := err.(*jwt.ValidationError); ok && validationErr.Errors == jwt.ValidationErrorExpired {
			return nil, errors2.New("Token is Disabled")
		}
		return nil, err
	}

	//获取Claims对象
	claims := token.Claims.(*userClaims)

	//带有jwt.StandardClaims 的对象 转换为 外部不带有 jwt.StandardClaims 的对象
	return c, struct_map.Map(claims.UserClaims, c)
}
