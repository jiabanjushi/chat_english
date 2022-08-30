package v2

import (
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-fly-muti/common"
	"go-fly-muti/models"
	"go-fly-muti/tools"
	"go-fly-muti/types"
	"time"
)

type RegisterForm struct {
	Username   string `form:"username" json:"username" uri:"username" xml:"username" binding:"required"`
	Password   string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
	RePassword string `form:"rePassword" json:"rePassword" uri:"rePassword" xml:"rePassword" binding:"required"`
	Nickname   string `form:"nickname" json:"nickname" uri:"nickname" xml:"nickname" binding:"required"`
	Captcha    string `form:"captcha" json:"captcha" uri:"captcha" xml:"captcha" binding:"required"`
}

func PostUcRegister(c *gin.Context) {
	var form RegisterForm
	err := c.Bind(&form)

	if err != nil {
		c.JSON(200, gin.H{
			"code":   types.ApiCode.FAILED,
			"msg":    types.ApiCode.GetMessage(types.ApiCode.INVALID),
			"result": err.Error(),
		})
		return
	}
	//验证码
	session := sessions.Default(c)
	if captchaId := session.Get("captcha"); captchaId != nil {
		session.Delete("captcha")
		_ = session.Save()
		if !captcha.VerifyString(captchaId.(string), form.Captcha) {
			c.JSON(200, gin.H{
				"code": types.ApiCode.CAPTCHA_FAILED,
				"msg":  types.ApiCode.GetMessage(types.ApiCode.CAPTCHA_FAILED),
			})
			return
		}
	} else {
		c.JSON(200, gin.H{
			"code": types.ApiCode.CAPTCHA_FAILED,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.CAPTCHA_FAILED),
		})
		return
	}

	c.JSON(200, UcRegister(form))
}
func UcRegister(form RegisterForm) gin.H {
	//重复密码
	if form.Password != form.RePassword {
		return gin.H{
			"code": types.ApiCode.INVALID_PASSWORD,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.INVALID_PASSWORD),
		}
	}
	//账户是否存在
	var user *models.User
	user = &models.User{
		Name: form.Username,
	}
	*user = user.GetOneUser("*")
	if user.ID != 0 {
		return gin.H{
			"code": types.ApiCode.ACCOUNT_EXIST,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.ACCOUNT_EXIST),
		}
	}

	//插入用户
	mStr := fmt.Sprintf("%ds", common.TryDeadline)
	duration, _ := time.ParseDuration(mStr)
	expired := time.Now().Add(duration)
	user = &models.User{
		Name:      form.Username,
		Password:  tools.Md5(form.Password),
		Avator:    "/static/images/4.jpg",
		Nickname:  form.Nickname,
		Pid:       1,
		UpdatedAt: time.Now(),
		ExpiredAt: types.Time{
			expired,
		},
		RecNum: 0,
		Status: 0,
	}
	userId := user.AddUser()
	if userId == 0 {
		return gin.H{
			"code": types.ApiCode.FAILED,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.FAILED),
		}
	}
	models.CreateUserRole(userId, types.Constant.EntRoleId)
	return gin.H{
		"code": types.ApiCode.SUCCESS,
		"msg":  types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
	}
}
