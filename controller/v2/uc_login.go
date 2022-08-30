package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go-fly-muti/models"
	"go-fly-muti/tools"
	"go-fly-muti/types"
	"strconv"
	"time"
)

type LoginForm struct {
	Username string `form:"username" json:"username" uri:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
}

func PostUcLogin(c *gin.Context) {
	var form LoginForm
	err := c.Bind(&form)

	if err != nil {
		c.JSON(200, gin.H{
			"code":   types.ApiCode.FAILED,
			"msg":    types.ApiCode.GetMessage(types.ApiCode.INVALID),
			"result": err.Error(),
		})
		return
	}
	c.JSON(200, UcLogin(form))
}
func UcLogin(form LoginForm) gin.H {
	var user *models.User
	user = &models.User{
		Name: form.Username,
	}
	*user = user.GetOneUser("*")
	md5Pass := tools.Md5(form.Password)
	if user.ID == 0 || user.Password != md5Pass {
		return gin.H{
			"code": types.ApiCode.LOGIN_FAILED,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.LOGIN_FAILED),
		}
	}
	ok, errCode, errMsg := user.CheckStatusExpired()
	if !ok {
		return gin.H{
			"code": errCode,
			"msg":  errMsg,
		}
	}
	roleId, _ := strconv.Atoi(user.RoleId)
	tokenCliams := tools.UserClaims{
		Id:         user.ID,
		Username:   user.Name,
		RoleId:     uint(roleId),
		Pid:        user.Pid,
		RoleName:   user.RoleName,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 24*3600,
		},
	}
	token, err := tools.MakeCliamsToken(tokenCliams)
	if err != nil {
		return gin.H{
			"code": types.ApiCode.LOGIN_FAILED,
			"msg":  err.Error(),
		}
	}
	return gin.H{
		"code": types.ApiCode.SUCCESS,
		"msg":  types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
		"result": gin.H{
			"token": token,
		},
	}
}
