package tools

import (
	"errors"
	"github.com/golang-jwt/jwt"
)

const SECRET = "WyLoveQiXiYl"

type UserClaims struct {
	Id         uint   `json:"id"`
	Pid        uint   `json:"pid"`
	Username   string `json:"username"`
	RoleName   string `json:"role_name"`
	RoleId     uint   `json:"role_id"`
	CreateTime string `json:"create_time"`
	jwt.StandardClaims
}

func MakeToken(obj map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(obj))
	tokenString, err := token.SignedString([]byte(SECRET))
	return tokenString, err
}
func ParseToken(tokenStr string) map[string]interface{} {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(SECRET), nil
	})

	if err != nil {
		return nil
	}
	finToken := token.Claims.(jwt.MapClaims)
	return finToken
}

/**
生成jwt
*/
func MakeCliamsToken(obj UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, obj)
	tokenString, err := token.SignedString([]byte(SECRET))
	return tokenString, err
}

/**
解析jwt token
*/
func ParseCliamsToken(token string, validExpired bool) (*UserClaims, error) {
	if token == "" {
		return nil, errors.New("token failed")
	}
	tokenClaims, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})

	if tokenClaims == nil {
		return nil, err
	}
	claims, ok := tokenClaims.Claims.(*UserClaims)
	if !ok {
		return nil, err
	}
	if validExpired && !tokenClaims.Valid {
		return nil, err
	}
	return claims, nil
}
