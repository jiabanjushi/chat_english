package v2

import (
	"github.com/gin-gonic/gin"
	"go-fly-muti/tools"
	"go-fly-muti/types"
	"time"
)

type TokenForm struct {
	Token string `form:"token" json:"token" uri:"token" xml:"token" binding:"required"`
}

func PostRefreshToken(c *gin.Context) {
	var form TokenForm
	err := c.Bind(&form)
	if err != nil {
		c.JSON(200, gin.H{
			"code":   types.ApiCode.FAILED,
			"msg":    types.ApiCode.GetMessage(types.ApiCode.INVALID),
			"result": err.Error(),
		})
		return
	}
	c.JSON(200, RefreshToken(form))
}
func RefreshToken(form TokenForm) gin.H {
	orgToken, err := ParseToken(form.Token)
	if err != nil {
		return gin.H{
			"code": types.ApiCode.TOKEN_FAILED,
			"msg":  err.Error(),
		}
	}
	orgToken.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	orgToken.ExpiresAt = time.Now().Unix() + 24*3600
	token, err := tools.MakeCliamsToken(*orgToken)
	if err != nil {
		return gin.H{
			"code": types.ApiCode.FAILED,
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
func PostParseToken(c *gin.Context) {
	var form TokenForm
	err := c.Bind(&form)
	if err != nil {
		c.JSON(200, gin.H{
			"code":   types.ApiCode.FAILED,
			"msg":    types.ApiCode.GetMessage(types.ApiCode.INVALID),
			"result": err.Error(),
		})
		return
	}
	orgToken, err := ParseToken(form.Token)
	if err != nil {
		c.JSON(200, gin.H{
			"code": types.ApiCode.TOKEN_FAILED,
			"msg":  err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"code":   types.ApiCode.SUCCESS,
		"msg":    types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
		"result": orgToken,
	})
}
func ParseToken(token string) (*tools.UserClaims, error) {
	orgToken, err := tools.ParseCliamsToken(token, false)
	if err != nil {
		return nil, err
	}
	return orgToken, nil
}
func GetJwt(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": types.ApiCode.SUCCESS,
		"msg":  types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
	})
}
func PostRefreshTokenV1(c *gin.Context) {
	var form TokenForm
	err := c.Bind(&form)
	if err != nil {
		c.JSON(200, gin.H{
			"code":   types.ApiCode.FAILED,
			"msg":    types.ApiCode.GetMessage(types.ApiCode.INVALID),
			"result": err.Error(),
		})
		return
	}
	userinfo := tools.ParseToken(form.Token)
	if userinfo == nil || userinfo["name"] == nil || userinfo["create_time"] == nil {
		c.JSON(200, gin.H{
			"code": types.ApiCode.FAILED,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.INVALID),
		})
		return
	}
	userinfo["create_time"] = time.Now().Unix()
	token, _ := tools.MakeToken(userinfo)
	c.JSON(200, gin.H{
		"code": types.ApiCode.SUCCESS,
		"msg":  types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
		"result": gin.H{
			"token": token,
		},
	})
}
func GetJwtV1(c *gin.Context) {
	kefuName, _ := c.Get("kefu_name")
	c.JSON(200, gin.H{
		"code": types.ApiCode.SUCCESS,
		"msg":  types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
		"result": gin.H{
			"kefu_name": kefuName,
		},
	})
}
