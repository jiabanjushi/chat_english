package router

import (
	"github.com/gin-gonic/gin"
	"go-fly-muti/middleware"
	"go-fly-muti/tmpl"
	"net/http"
)

func InitViewRouter(engine *gin.Engine) {
	//首页
	engine.GET("/", tmpl.PageFirstIndex)
	engine.GET("/index", middleware.SetLanguage, tmpl.PageNewIndex)
	engine.GET("/index_:lang", middleware.SetLanguage, tmpl.PageNewIndex)
	engine.GET("/install.html", middleware.SetLanguage, tmpl.PageInstallHtml)
	engine.GET("/download.html", middleware.SetLanguage, tmpl.PageDownloadHtml)
	engine.GET("/contact.html", middleware.SetLanguage, tmpl.PageContactHtml)
	engine.GET("/show.html", middleware.SetLanguage, tmpl.PageShowHtml)
	engine.GET("/aboutus.html", middleware.SetLanguage, tmpl.PageContactHtml)
	engine.GET("/document.html", middleware.SetLanguage, tmpl.PageShowHtml)
	engine.GET("/kfxieyi.html", middleware.SetLanguage, tmpl.PageShowHtml)
	engine.GET("/yhxieyi.html", middleware.SetLanguage, tmpl.PageShowHtml)
	engine.GET("/zxkfxt.html", middleware.SetLanguage, tmpl.PageNewIndex)
	engine.GET("/wzkfxt.html", middleware.SetLanguage, tmpl.PageNewIndex)
	engine.GET("/wyltgj.html", middleware.SetLanguage, tmpl.PageNewIndex)
	engine.GET("/mfzxkf.html", middleware.SetLanguage, tmpl.PageNewIndex)

	engine.GET("/install", tmpl.PageInstall)
	engine.GET("/login", middleware.DomainLimitMiddleware, middleware.TryUseLimitMiddleware, tmpl.PageLogin)
	engine.GET("/chat_page", middleware.SetLanguage, tmpl.PageChat)
	engine.GET("/chatIndex", middleware.SetLanguage, tmpl.PageChat)
	engine.GET("/wechatIndex", middleware.SetLanguage, tmpl.PageWechatChat)
	engine.GET("/chatRoom", tmpl.PageChatRoom)
	//后台主界面
	engine.GET("/mainGuide", func(c *gin.Context) {
		c.HTML(http.StatusOK, "main_guide.html", gin.H{})
	})
	engine.GET("/main", middleware.DomainLimitMiddleware, middleware.TryUseLimitMiddleware, middleware.JwtPageMiddleware, tmpl.PageMain)
	engine.GET("/chat_main", middleware.JwtPageMiddleware, tmpl.PageChatMain)
	engine.GET("/setting", tmpl.PageSetting)
	engine.GET("/setting_statistics", tmpl.PageSettingStatis)
	engine.GET("/setting_indexpage", tmpl.PageSettingIndexPage)
	engine.GET("/setting_mysql", tmpl.PageSettingMysql)
	engine.GET("/setting_welcome", tmpl.PageSettingWelcome)
	engine.GET("/setting_deploy", tmpl.PageSettingDeploy)
	engine.GET("/setting_kefu_list", tmpl.PageKefuList)
	engine.GET("/setting_visitor_list", tmpl.PageVisitorList)
	engine.GET("/setting_visitor_message", tmpl.PageVisitorMessage)
	engine.GET("/setting_user_all", tmpl.PageUserAllList)
	engine.GET("/setting_kefu_message", tmpl.PageKefuMessage)
	engine.GET("/setting_avator", tmpl.PageAvator)
	engine.GET("/setting_modifypass", tmpl.PageModifypass)
	engine.GET("/setting_ipblack", tmpl.PageIpblack)
	engine.GET("/setting_config", tmpl.PageConfig)
	engine.GET("/setting_configs", tmpl.PageConfigs)
	engine.GET("/setting_articles", tmpl.PageSettingArticles)
	engine.GET("/setting_wechat_menu", tmpl.PageSettingWechatMenu)
	engine.GET("/setting_news", tmpl.PageSettingNews)
	engine.GET("/mail_list", tmpl.PageMailList)
	engine.GET("/roles_list", tmpl.PageRoleList)
	//演示页面
	engine.GET("/deploy", func(c *gin.Context) {
		kefuName := c.Query("kefuName")
		entId := c.Query("entId")
		siteUrl := c.Query("siteUrl")
		c.HTML(http.StatusOK, "deploy.html", gin.H{
			"kefuName": kefuName,
			"entId":    entId,
			"siteUrl":  siteUrl,
		})
	})
	//测试
	engine.GET("/test/:entId/:kefuName", func(c *gin.Context) {
		kefuName := c.Param("kefuName")
		entId := c.Param("entId")
		siteUrl := c.Query("siteUrl")
		c.HTML(http.StatusOK, "deploy.html", gin.H{
			"kefuName": kefuName,
			"entId":    entId,
			"siteUrl":  siteUrl,
		})
	})
	//sitemap
	engine.GET("/sitemap.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sitemap.html", gin.H{})
	})
}
