package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"go-fly-muti/lib"
	"go-fly-muti/models"
	"path"
	"strconv"
	"strings"
)

func GetConfigs(c *gin.Context) {
	configs := models.FindConfigs()
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": configs,
	})
}
func GetEntConfigs(c *gin.Context) {
	entId, _ := c.Get("ent_id")
	configs := models.FindEntConfigs(entId)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": configs,
	})
}
func GetConfig(c *gin.Context) {
	key := c.Query("key")
	config := models.FindConfig(key)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": config,
	})
}
func PostEntConfigs(c *gin.Context) {
	kefuId, _ := c.Get("kefu_id")
	name := c.PostForm("name")
	key := c.PostForm("key")
	value := c.PostForm("value")
	if key == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "error",
		})
		return
	}
	config := models.FindEntConfig(kefuId, key)
	if config.ID == 0 {
		models.CreateEntConfig(kefuId, name, key, value)
	} else {
		models.UpdateEntConfig(kefuId, name, key, value)
	}

	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": "",
	})
}

//保存微信菜单数据
func PostWechatMenu(c *gin.Context) {
	kefuId, _ := c.Get("kefu_id")
	menu := c.PostForm("menu")
	name := c.PostForm("name")
	if menu == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "error",
		})
		return
	}
	config := models.FindEntConfig(kefuId, "WechatMenu")
	if config.ID == 0 {
		models.CreateEntConfig(kefuId, name, "WechatMenu", menu)
	} else {
		models.UpdateEntConfig(kefuId, name, "WechatMenu", menu)
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}

//生成微信菜单
func GetWechatMenu(c *gin.Context) {
	entId, _ := c.Get("ent_id")
	config := models.FindEntConfig(entId, "WechatMenu")
	if config.ID == 0 || config.ConfValue == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "没有菜单数据",
		})
		return
	}
	wechatConfig, err := lib.NewWechatLib(entId.(string))
	if wechatConfig == nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	wc := wechat.NewWechat()
	cfg := &offConfig.Config{
		AppID:     wechatConfig.AppId,
		AppSecret: wechatConfig.AppSecret,
		Token:     wechatConfig.Token,
		//EncodingAESKey: "xxxx",
		Cache: memory,
	}
	officialAccount := wc.GetOfficialAccount(cfg)
	menu := officialAccount.GetMenu()
	err = menu.SetMenuByJSON(config.ConfValue)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
func PostConfig(c *gin.Context) {
	key := c.PostForm("key")
	value := c.PostForm("value")
	if key == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "error",
		})
		return
	}
	models.UpdateConfig(key, value)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}

//上传微信认证文件
func PostUploadWechatFile(c *gin.Context) {
	SendAttachment, err := strconv.ParseBool(models.FindConfig("SendAttachment"))
	if !SendAttachment || err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "禁止上传附件!",
		})
		return
	}
	f, err := c.FormFile("file")
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "上传失败!",
		})
		return
	} else {

		fileExt := strings.ToLower(path.Ext(f.Filename))
		if f.Size >= 1*1024*1024 {
			c.JSON(200, gin.H{
				"code": 400,
				"msg":  "上传失败!不允许超过1M",
			})
			return
		}
		if fileExt != ".txt" {
			c.JSON(200, gin.H{
				"code": 400,
				"msg":  "上传失败!只允许txt文件",
			})
			return
		}

		c.SaveUploadedFile(f, f.Filename)
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "上传成功!",
			"result": gin.H{
				"path": "/" + f.Filename,
			},
		})
	}
}
