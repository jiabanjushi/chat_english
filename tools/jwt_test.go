package tools

import (
	"github.com/golang-jwt/jwt"
	"testing"
	"time"
)

func TestJwt(t *testing.T) {
	tokenCliams := UserClaims{
		Id:         1,
		Username:   "kefu2",
		RoleId:     2,
		Pid:        1,
		CreateTime: time.Now(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 24*3600,
		},
	}
	token, err := MakeCliamsToken(tokenCliams)
	t.Log(token, err)

	orgToken, err := ParseCliamsToken(token, true)
	t.Logf("%+v,%+v", orgToken, err)
}
