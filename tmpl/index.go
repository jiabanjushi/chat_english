package tmpl

import (
	"github.com/gin-gonic/gin"
	"go-fly-muti/models"
	"net/http"
)

// PageFirstIndex 匹配http://xxx/
func PageFirstIndex(c *gin.Context) {

	if c.Request.RequestURI == "/favicon.ico" {
		return
	}

	//首页默认跳转地址
	IndexJumpUrl := models.FindConfig("IndexJumpUrl")
	if IndexJumpUrl != "" {
		c.Redirect(302, IndexJumpUrl)
	}
	lang, _ := c.Get("lang")
	if lang == "en" {
		c.HTML(http.StatusOK, "index_en.html", gin.H{
			"Lang": lang,
		})
	} else {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Lang": lang,
		})
	}
}

// PageNewIndex 首页
func PageNewIndex(c *gin.Context) {
	if c.Request.RequestURI == "/favicon.ico" {
		return
	}

	lang, _ := c.Get("lang")
	if lang == "en" {
		c.HTML(http.StatusOK, "index_en.html", gin.H{
			"Lang": lang,
		})
	} else {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Lang": lang,
		})
	}

}

// PageInstallHtml 部署界面
func PageInstallHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "install.html", nil)
}

// PageDownloadHtml 下载界面
func PageDownloadHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "download.html", nil)
}

// PageContactHtml 联系界面
func PageContactHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "contact.html", nil)
}

// PageShowHtml 演示界面
func PageShowHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "show.html", nil)
}
