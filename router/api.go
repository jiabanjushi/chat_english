package router

import (
	"github.com/gin-gonic/gin"
	"go-fly-muti/controller"
	controllerV2 "go-fly-muti/controller/v2"
	"go-fly-muti/middleware"
	"go-fly-muti/ws"

)

func InitApiRouter(engine *gin.Engine) {
	//用户身份接口
	ucv1 := engine.Group("/uc/v1", middleware.Ipblack, middleware.SetLanguage)
	{
		//刷新token
		ucv1.POST("/refreshToken", controllerV2.PostRefreshTokenV1)
		//刷新token
		ucv1.POST("/loginCheck", controller.LoginCheckPass)
		//访客登录
		ucv1.POST("/visitorLogin", middleware.Ipblack, controller.PostVisitorLogin)
	}
	uc := engine.Group("/uc/v2", middleware.Ipblack, middleware.SetLanguage)
	{
		//注册
		uc.POST("/register", controllerV2.PostUcRegister)
		//登录
		uc.POST("/loginCheck", controllerV2.PostUcLogin)
		//刷新token
		uc.POST("/refreshToken", controllerV2.PostRefreshToken)
		//解析token
		uc.POST("/parseToken", controllerV2.PostParseToken)
		//登录
		uc.POST("/emailCode", controllerV2.PostEmailCode)
		//获取ip地址
		uc.GET("/ipAuth", controllerV2.GetIpAuth)
	}
	//客服接口
	kefuV2 := engine.Group("/kefu/v2", middleware.Ipblack, middleware.SetLanguage, middleware.JwtApiV2Middleware)
	{
		//解析jwt
		kefuV2.GET("/parseJwt", controllerV2.GetJwt)
		kefuV2.GET("/visitorExt", controller.GetVisitorExt)
	}
	//访客接口
	visitorV2 := engine.Group("/visitor/v2", middleware.Ipblack, middleware.SetLanguage)
	{
		visitorV2.POST("/login", controllerV2.PostVisitorLogin)
	}
	//消息接口
	//微信接口
	//企业接口
	//路由分组
	v2 := engine.Group("/2", middleware.SetLanguage)
	{
		//获取欢迎信息
		//分页获取数据
		v2.GET("/notices", controller.GetNotice)
		//获取消息
		v2.GET("/messages", controller.GetMessagesV2)
		v2.GET("/messages_unread", controller.GetMessagesVisitorUnread)
		v2.POST("/messages_read", controller.PostMessagesVisitorRead)
		//发送单条信息
		v2.POST("/message", middleware.Ipblack, controller.SendMessageV2)
		//分页获取数据
		v2.GET("/messages_page", controller.GetVisitorListMessagesPage)
		v2.POST("/message_ask", controller.PostMessagesAsk)
	}
	//企业路由分组
	entGroup := engine.Group("/ent")
	{
		entGroup.POST("/email_message", controller.SendEntEmail)
		entGroup.GET("/article_cate", middleware.JwtApiMiddleware, controller.GetArticleCates)
		entGroup.GET("/article_list", middleware.JwtApiMiddleware, controller.GetArticleList)
	}
	//客服路由分组
	kefuGroup := engine.Group("/kefu")
	kefuGroup.Use(middleware.JwtApiMiddleware)
	{
		//获取客服信息
		kefuGroup.GET("/kefuinfo", controller.GetKefuInfo)
		//解析jwt
		kefuGroup.GET("/parseJwt", controllerV2.GetJwtV1)
		//app注册
		kefuGroup.POST("/appClient", controllerV2.PostAppKefuClient)
		kefuGroup.GET("/visitor", controller.GetVisitor)
		kefuGroup.GET("/visitor_attr", controller.GetVisitorAttr)
		kefuGroup.POST("/visitor_attrs", controller.PostVisitorAttrs)
		kefuGroup.POST("/message", controller.SendKefuMessage)
		//关闭连接
		kefuGroup.GET("/message_close", controller.SendCloseMessageV2)
		kefuGroup.GET("/visitor/messages", controller.GetVisitorMessageByKefu)
		//分页获取数据
		kefuGroup.GET("/messages_page", controller.GetVisitorListMessagesPageBykefu)
		kefuGroup.POST("/message_delete", controller.DeleteMessage)
		kefuGroup.POST("/messages_read", controller.PostMessagesKefuRead)
		kefuGroup.GET("/visitorExt", controller.GetVisitorExt)
		kefuGroup.GET("/onlineVisitors", controller.GetKefusVisitorOnlines)
		//统计信息
		kefuGroup.GET("/statistics", controller.GetStatistics)
		//图表统计信息
		kefuGroup.GET("/chartStatistics", controllerV2.GetChartStatistic)
		//更新在线状态
		kefuGroup.GET("/updateOnlineStatus", controller.GetUpdateOnlineStatus)
		//给访客打tag
		kefuGroup.POST("/visitorTag", controller.PostVisitorTag)
		//获取tag
		kefuGroup.GET("/tags", controller.GetTags)
		//获取访客tags
		kefuGroup.GET("/visitorAllTag", controller.GetVisitorAllTags)
		//获取访客tags
		kefuGroup.GET("/visitorTag", controller.GetVisitorTags)
		//删除访客tags
		kefuGroup.GET("/delVisitorTag", controller.DelVisitorTag)
		//增加文章分类
		kefuGroup.POST("/addArticleCate", controller.PostArticleCate)
		//增加文章
		kefuGroup.POST("/addArticle", controller.PostArticle)
		//删除文章
		kefuGroup.GET("/delArticle", controller.DelArticle)
		//删除文章分类
		kefuGroup.GET("/delArticleCate", controller.DelArticleCate)
		//搜索访客列表
		kefuGroup.GET("/visitorList", controller.GetVisitorsList)
		//获取客服用户
		kefuGroup.GET("/kefuUsers", controller.GetKefuListOwn)
		//获取二维码
		kefuGroup.GET("/qrcode", controller.GetQrcode)
		//删除访客
		kefuGroup.DELETE("/delVisitor", controller.DelVisitor)
		//自动欢迎列表
		kefuGroup.GET("/notices", controller.GetNotices)
		//添加欢迎
		kefuGroup.POST("/notice", controller.PostNotice)
		//删除欢迎
		kefuGroup.DELETE("/notice", controller.DelNotice)
		//更新欢迎
		kefuGroup.POST("/updateNotice", controller.PostNoticeSave)
		//配置企业
		kefuGroup.POST("/entConfigs", middleware.RbacAuth, controller.PostEntConfigs)
		//子账号列表
		kefuGroup.GET("/kefuList", middleware.RbacAuth, controller.GetKefuListOwn)
		//编辑账号
		kefuGroup.POST("/kefuInfo", middleware.RbacAuth, controller.PostKefuInfo)
		//删除账号
		kefuGroup.DELETE("/kefuInfo", middleware.RbacAuth, controller.DeleteKefuInfoOwn)
		//切换状态
		kefuGroup.POST("/kefuStatus", middleware.RbacAuth, controller.PostKefuInfoStatus)
		//查看别的客服
		kefuGroup.GET("/kefuSetting", middleware.RbacAuth, controller.GetKefuInfoSettingOwn)
		//生成微信菜单
		kefuGroup.POST("/wechatMenu", middleware.RbacAuth, controller.PostWechatMenu)
		//生成微信菜单
		kefuGroup.GET("/mkWechatMenu", middleware.RbacAuth, controller.GetWechatMenu)
		//上传微信认证文件
		kefuGroup.POST("/uploadWechatAuthFile", controller.PostUploadWechatFile)
		//拉黑IP
		kefuGroup.POST("/ipblack", controller.PostIpblack)
		//IP黑名单列表
		kefuGroup.GET("/ipblacks", controller.GetIpblacksByKefuId)
		//删除IP
		kefuGroup.DELETE("/ipblack", controller.DelIpblack)
		//删除访客聊天记录
		kefuGroup.GET("/delVisitorMessage", controller.DeleteVisitorMessage)
		//添加访客黑名单
		kefuGroup.POST("/visitorBlack", controller.PostVisitorBlack)
		//访客黑名单
		kefuGroup.GET("/visitorBlacks", controller.GeVisitorBlacks)
		//删除访客黑名单
		kefuGroup.GET("/delVisitorBlack", controller.DelVisitorBlack)
		//获取角色
		kefuGroup.GET("/roleList", controller.GetRoleListOwn)
	}
	//聊天室路由分组
	room := engine.Group("/room")
	{
		room.POST("/login", middleware.Ipblack, controller.PostRoomLogin)
		room.POST("/message", middleware.Ipblack, controller.PostRoomMessage)
	}
	engine.GET("/captcha", controller.GetCaptcha)
	engine.POST("/check", controller.LoginCheckPass)
	engine.POST("/JyKey",controller.JyKey)

	engine.POST("/check_auth", middleware.JwtApiMiddleware, controller.MainCheckAuth)
	engine.GET("/userinfo", middleware.JwtApiMiddleware, controller.GetKefuInfoAll)
	engine.POST("/register", middleware.Ipblack, controller.PostKefuRegister)
	//engine.POST("/install", controller.PostInstall)
	//前后聊天
	engine.GET("/ws_kefu", middleware.JwtApiMiddleware, ws.NewKefuServer)
	engine.GET("/ws_visitor", middleware.Ipblack, ws.NewVisitorServer)
	//websocket路由分组
	//wsGroup := engine.Group("/ws")
	//{
	//	wsGroup.GET("/v2/visitor", middleware.Ipblack, wsV2.VisitorWebsocketServer)
	//}

	engine.GET("/message_notice", controller.SendVisitorNotice)
	//发送单条消息
	//上传文件
	engine.POST("/uploadimg", middleware.Ipblack, controller.UploadImg)
	//上传文件
	engine.POST("/uploadfile", middleware.Ipblack, controller.UploadFile)
	//上传文件
	engine.POST("/uploadaudio", middleware.Ipblack, controller.UploadAudio)
	engine.POST("/call_kefu", controller.PostCallKefu)

	//获取客服信息
	engine.POST("/kefuinfo_peerid", middleware.JwtApiMiddleware, controller.PostKefuPeerId)
	engine.GET("/kefuinfo", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetKefuInfo)
	engine.GET("/kefuinfo_setting", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetKefuInfoSetting)

	engine.POST("/kefuinfo", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostKefuInfo)

	engine.GET("/kefulist", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetKefuList)
	engine.GET("/kefulist_own", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetKefuListOwn)
	engine.GET("/kefulist_message", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetKefuListMessages)
	engine.GET("/visitorlist_message", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetVisitorListMessages)
	engine.GET("/other_kefulist", middleware.JwtApiMiddleware, controller.GetOtherKefuList)
	engine.GET("/trans_kefu", middleware.JwtApiMiddleware, controller.PostTransKefu)
	engine.POST("/modifypass", middleware.JwtApiMiddleware, controller.PostKefuPass)
	engine.POST("/modifyavator", middleware.JwtApiMiddleware, controller.PostKefuAvator)
	engine.POST("/modify_kefuinfo", middleware.JwtApiMiddleware, controller.PostKefuinfo)
	//角色列表
	engine.GET("/roles", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetRoleList)
	engine.POST("/role", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostRole)

	engine.GET("/visitors_online", controller.GetVisitorOnlines)
	engine.GET("/visitors_kefu_online", middleware.JwtApiMiddleware, controller.GetKefusVisitorOnlines)
	engine.GET("/clear_online_tcp", controller.DeleteOnlineTcp)
	engine.POST("/visitor_login", middleware.Ipblack, controller.PostVisitorLogin)
	engine.GET("/visitor_ext", middleware.JwtApiMiddleware, controller.GetVisitorExt)

	//engine.POST("/visitor", controller.PostVisitor)
	engine.GET("/visitor", middleware.JwtApiMiddleware, controller.GetVisitor)
	engine.GET("/visitors", middleware.JwtApiMiddleware, controller.GetVisitors)
	engine.GET("/statistics", middleware.JwtApiMiddleware, controller.GetStatistics)
	//前台接口
	engine.GET("/about", controller.GetAbout)
	engine.POST("/about", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostAbout)
	engine.GET("/notice", middleware.SetLanguage, controller.GetNotice)

	engine.GET("/ipblacks_all", middleware.JwtApiMiddleware, controller.GetIpblacks)

	engine.GET("/configs", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetConfigs)
	engine.GET("/ent_configs", middleware.JwtApiMiddleware, controller.GetEntConfigs)

	engine.POST("/config", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostConfig)
	engine.GET("/config", controller.GetConfig)
	engine.GET("/autoreply", controller.GetAutoReplys)
	engine.GET("/replys", middleware.JwtApiMiddleware, controller.GetReplys)
	engine.POST("/reply", middleware.JwtApiMiddleware, controller.PostReply)
	engine.POST("/reply_content", middleware.JwtApiMiddleware, controller.PostReplyContent)
	engine.POST("/reply_content_save", middleware.JwtApiMiddleware, controller.PostReplyContentSave)
	engine.DELETE("/reply_content", middleware.JwtApiMiddleware, controller.DelReplyContent)
	engine.DELETE("/reply", middleware.JwtApiMiddleware, controller.DelReplyGroup)
	engine.POST("/reply_search", middleware.JwtApiMiddleware, controller.PostReplySearch)
	engine.POST("/translate", controller.GetTranslate)

	//微信
	wechat := engine.Group("/wechat")
	{
		wechat.GET("/checkSignature", controller.GetCheckWeixinSign)
		wechat.GET("/showQrcode", controller.GetShowQrCode)
		wechat.Any("/server/:entId/:kefuName", controller.PostWechatServer)
		//查询微信绑定表
		wechat.GET("/oauth", controller.GetWechatOauth)
	}
	//系统相关
	systemGroup := engine.Group("/system")
	systemGroup.Use(middleware.JwtApiMiddleware, middleware.AdminAuth)
	{
		systemGroup.GET("/stop", controller.GetStop)
		//systemGroup.POST("/saveNews", controller.PostNews)
	//	systemGroup.GET("/delNews", controller.DelNews)
	}
	//版本信息
	engine.GET("/version", controller.GetVersion)
	//其他接口
	otherGroup := engine.Group("/other")
	{
		//版本信息
		otherGroup.GET("/version", controller.GetOtherVersion)
		otherGroup.GET("/news", controller.GetNews)
	}
}
