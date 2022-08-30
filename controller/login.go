package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-fly-muti/models"
	"go-fly-muti/tools"
	"go-fly-muti/types"
	"time"
)

type LoginForm struct {
	Username string `form:"username" json:"username" uri:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
}

// LoginCheckPass @Summary 登陆验证接口
// @Produce  json
// @Accept multipart/form-data
// @Param username formData   string true "用户名"
// @Param password formData   string true "密码"
// @Param type formData   string true "类型"
// @Success 200 {object} controller.Response
// @Failure 200 {object} controller.Response
// @Router /check [post]
//验证接口
func LoginCheckPass(c *gin.Context) {
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








	fmt.Println("用户名:"+form.Username)
	fmt.Println("密码:"+form.Password)

	info, ok := CheckKefuPass(form.Username, form.Password)
	if !ok {
		c.JSON(200, gin.H{
			"code": types.ApiCode.LOGIN_FAILED,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.LOGIN_FAILED),
		})
		return
	}
	if info.Status == 1 {
		c.JSON(200, gin.H{
			"code": types.ApiCode.ACCOUNT_FORBIDDEN,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.ACCOUNT_FORBIDDEN),
		})
		return
	}

	//获取远程结果
	err, ok = CheckServerAddress()
	if !ok {
		c.JSON(200, gin.H{
			"code": types.ApiCode.LOGIN_FAILED,
			"msg":  err.Error(),
		})
		return
	}


	token := GenUserToken(info)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "验证成功,正在跳转",
		"result": gin.H{
			"token": token,
		},
	})
}
func GenUserToken(info models.User) string {
	userinfo := make(map[string]interface{})
	userinfo["name"] = info.Name
	userinfo["kefu_id"] = info.ID
	userinfo["type"] = "kefu"
	uRole := models.FindRoleByUserId(info.ID)
	if uRole.RoleId != 0 {
		userinfo["role_id"] = uRole.RoleId
	} else {
		userinfo["role_id"] = 2
	}
	userinfo["create_time"] = time.Now().Unix()
	userinfo["pid"] = info.Pid

	token, _ := tools.MakeToken(userinfo)
	return token
}
