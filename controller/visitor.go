package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-fly-muti/common"
	"go-fly-muti/models"
	v2 "go-fly-muti/models/v2"
	"go-fly-muti/tools"
	"go-fly-muti/types"
	"go-fly-muti/ws"
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type VisitorLoginForm struct {
	VisitorId   string `form:"visitor_id" json:"visitor_id" uri:"visitor_id" xml:"visitor_id"`
	Refer       string `form:"refer" json:"refer" uri:"refer" xml:"refer"`
	ReferUrl    string `form:"refer_url" json:"refer" uri:"refer" xml:"refer"`
	Url         string `form:"url" json:"url" uri:"url" xml:"url"`
	ToId        string `form:"to_id" json:"to_id" uri:"to_id" xml:"to_id"  binding:"required"`
	EntId       string `form:"ent_id" json:"ent_id" uri:"ent_id" xml:"ent_id" binding:"required"`
	Avator      string `form:"avator" json:"avator" uri:"avator" xml:"avator"`
	UserAgent   string `form:"user_agent" json:"user_agent" uri:"user_agent" xml:"user_agent"`
	Extra       string `form:"extra" json:"extra" uri:"extra" xml:"extra"`
	ClientIp    string `form:"client_ip" json:"client_ip" uri:"client_ip" xml:"client_ip"`
	CityAddress string `form:"city_address" json:"city_address" uri:"city_address" xml:"city_address"`
	VisitorName string `form:"visitor_name" json:"visitor_name" uri:"visitor_name" xml:"visitor_name"`
}

func PostVisitorLogin(c *gin.Context) {
	var form VisitorLoginForm
	err := c.Bind(&form)
	if err != nil {
		c.JSON(200, gin.H{
			"code":   types.ApiCode.FAILED,
			"msg":    types.ApiCode.GetMessage(types.ApiCode.INVALID),
			"result": err.Error(),
		})
		return
	}

	//限流
	if !tools.LimitFreqSingle("visitor_login:"+c.ClientIP(), 1, 2) {
		c.JSON(200, gin.H{
			"code": types.ApiCode.FREQ_LIMIT,
			"msg":  c.ClientIP() + types.ApiCode.GetMessage(types.ApiCode.FREQ_LIMIT),
		})
		return
	}

	isWechat := false
	if form.VisitorId == "" {
		form.VisitorId = tools.Uuid()
		form.VisitorId = fmt.Sprintf("%s|%s", form.EntId, form.VisitorId)
	} else {
		//验证访客黑名单
		if !CheckVisitorBlack(form.VisitorId) {
			c.JSON(200, gin.H{
				"code": types.ApiCode.VISITOR_BAN,
				"msg":  types.ApiCode.GetMessage(types.ApiCode.VISITOR_BAN),
			})
			return
		}
		visitorIds := strings.Split(form.VisitorId, "|")
		if len(visitorIds) > 1 {
			if visitorIds[0] != "wx" {
				form.VisitorId = strings.Join(visitorIds[1:], "|")
				form.VisitorId = fmt.Sprintf("%s|%s", form.EntId, form.VisitorId)
			} else {
				isWechat = true
			}
		} else {
			form.VisitorId = fmt.Sprintf("%s|%s", form.EntId, form.VisitorId)
		}
	}

	form.UserAgent = c.GetHeader("User-Agent")
	if form.Avator == "" && !isWechat {
		if tools.IsMobile(form.UserAgent) {
			form.Avator = "/static/images/phone.png"
		} else {
			form.Avator = "/static/images/computer.png"
		}
	} else {
		form.Avator, _ = url.QueryUnescape(form.Avator)
	}

	ipCity := tools.ParseIpNew(c.ClientIP())
	if ipCity != nil {
		form.CityAddress = ipCity.CountryName + ipCity.RegionName + ipCity.CityName
	} else {
		form.CityAddress = "Unknown Area"
	}
	if form.VisitorName == "" && !isWechat {
		form.VisitorName = form.CityAddress
	} else {
		form.VisitorName, _ = url.QueryUnescape(form.VisitorName)
	}

	form.ClientIp = c.ClientIP()

	makeVisitorLoginForm(&form)
	allOffline := true

	entKefuInfo := models.FindUserByUid(form.EntId)
	if entKefuInfo.ID == 0 {
		c.JSON(200, gin.H{
			"code": types.ApiCode.ENT_ERROR,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.ENT_ERROR),
		})
		return
	}
	if entKefuInfo.Status == 1 {
		c.JSON(200, gin.H{
			"code": types.ApiCode.ACCOUNT_FORBIDDEN,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.ACCOUNT_FORBIDDEN),
		})
		return
	}
	dstKefu := models.FindUser(form.ToId)
	//判断是否在线
	_, ok := ws.KefuList[form.ToId]
	if dstKefu.OnlineStatus == 1 && ok {
		allOffline = false
	} else {
		kefus := models.FindUsersWhere("(pid = ? or id=?) and online_status=1", form.EntId, form.EntId)
		//kefus := models.FindUsersByPid(form.EntId)
		if len(kefus) == 0 {
			form.ToId = entKefuInfo.Name
		} else {
			for _, kefu := range kefus {
				if _, ok := ws.KefuList[kefu.Name]; ok {
					form.ToId = kefu.Name
					allOffline = false
					dstKefu = kefu
					break
				}
			}
		}

	}

	visitor := models.FindVisitorByVistorId(form.VisitorId)
	visitor.ToId = form.ToId
	if visitor.Name != "" {
		if form.Avator == "" && isWechat {
			form.Avator = visitor.Avator
		}
		if form.VisitorName == "" && isWechat {
			form.VisitorName = visitor.Name
		}
		if visitor.RealName != "" {
			form.VisitorName = visitor.RealName
		}
		//更新状态上线
		models.UpdateVisitor(form.EntId, form.VisitorName, form.Avator, form.VisitorId, form.ToId, visitor.Status, c.ClientIP(), c.ClientIP(), form.Refer, form.Extra)
	} else {
		visitor = *models.CreateVisitor(form.VisitorName, form.Avator, c.ClientIP(), form.ToId, form.VisitorId, form.Refer, form.CityAddress, form.ClientIp, form.EntId, form.Extra)
	}

	if form.VisitorName != "" {
		visitor.Name = form.VisitorName
	}
	go SendVisitorLoginNotice(form.ToId, visitor.Name, visitor.Avator, visitor.Name+"来了", visitor.VisitorId)
	go models.AddVisitorExt(visitor.VisitorId, Address, form.UserAgent, form.Url, form.Refer, form.ClientIp)
	go SendWechatVisitorTemplate(form.ToId, visitor.Name, "上线", visitor.EntId)
	//go SendWechatKefuNotice(form.ToId, "[访客]"+visitor.Name+",访问："+form.Refer, visitor.EntId)
	go SendNoticeEmail(visitor.Name, "[访客]"+visitor.Name, form.EntId, "访问："+form.Refer)
	go SendAppGetuiPush(dstKefu.Name, "[访客]"+visitor.Name, visitor.Name+"来了")
	c.JSON(200, gin.H{
		"code":       200,
		"msg":        "ok",
		"alloffline": allOffline,
		"result":     visitor,
		"kefu": gin.H{
			"username": dstKefu.Nickname,
			"avatar":   dstKefu.Avator,
		},
	})
}

//组合扩展信息
func makeVisitorLoginForm(form *VisitorLoginForm) {
	//扩展信息
	extraJson := tools.Base64Decode(form.Extra)
	if extraJson != "" {
		var extraObj VisitorExtra
		err := json.Unmarshal([]byte(extraJson), &extraObj)
		if err == nil {
			if extraObj.VisitorName != "" {
				form.VisitorName = extraObj.VisitorName
			}
			if extraObj.VisitorAvatar != "" {
				form.Avator = extraObj.VisitorAvatar
			}
			if extraObj.VisitorId != "" {
				form.VisitorId = extraObj.VisitorId
			}
		}
	}
}
func GetVisitor(c *gin.Context) {
	visitorId := c.Query("visitorId")
	vistor := models.FindVisitorByVistorId(visitorId)
	exts := models.FindVisitorExtByVistorId(visitorId, 1, 1)
	osVersion := ""
	browser := ""
	if len(exts) != 0 {
		ext := exts[0]
		uaParser := tools.NewUaParser(ext.Ua)
		osVersion = uaParser.OsVersion
		browser = uaParser.Browser
	}
	c.JSON(200, gin.H{
		"code":        200,
		"msg":         "ok",
		"create_time": vistor.CreatedAt.Format("2006-01-02 15:04:05"),
		"last_time":   vistor.UpdatedAt.Format("2006-01-02 15:04:05"),
		"os_version":  osVersion,
		"browser":     browser,
		"result":      vistor,
	})
}

/**
获取访客访问动态
*/
func GetVisitorExt(c *gin.Context) {
	visitorId := c.Query("visitor_id")
	page, pagesize := HandlePagePageSize(c)
	list := make([]VisitorExtend, 0)
	count := models.CountVisitorExtByVistorId(visitorId)
	exts := models.FindVisitorExtByVistorId(visitorId, page, pagesize)
	for _, ext := range exts {
		uaParser := tools.NewUaParser(ext.Ua)
		item := VisitorExtend{
			CreatedAt: ext.CreatedAt.Format("2006-01-02 15:04:05"),
			ID:        ext.ID,
			VisitorId: ext.VisitorId,
			Url:       ext.Url,
			Ua:        ext.Ua,
			Title:     ext.Title,
			ClientIp:  ext.ClientIp,
			OsVersion: uaParser.OsVersion,
			Browser:   uaParser.Browser,
		}
		list = append(list, item)
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"result": gin.H{
			"list":     list,
			"count":    count,
			"pagesize": pagesize,
		},
	})
}

// @Summary 获取访客列表接口
// @Produce  json
// @Accept multipart/form-data
// @Param page query   string true "分页"
// @Param token header string true "认证token"
// @Success 200 {object} controller.Response
// @Failure 200 {object} controller.Response
// @Router /visitors [get]
func GetVisitors(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	var pagesize uint
	myPagesize, _ := strconv.Atoi(c.Query("pagesize"))
	if myPagesize != 0 {
		pagesize = uint(myPagesize)
	} else {
		pagesize = common.VisitorPageSize
	}
	kefuId, _ := c.Get("kefu_name")
	vistors := models.FindVisitorsByKefuId(uint(page), pagesize, kefuId.(string))
	visitorIds := make([]string, 0)
	for _, visitor := range vistors {
		visitorIds = append(visitorIds, visitor.VisitorId)
	}
	messagesMap := models.FindLastMessageMap(visitorIds)
	unreadMap := models.FindUnreadMessageNumByVisitorIds(visitorIds, "visitor")
	//log.Println(unreadMap)
	users := make([]VisitorOnline, 0)
	for _, visitor := range vistors {
		var unreadNum uint32
		if num, ok := unreadMap[visitor.VisitorId]; ok {
			unreadNum = num
		}
		user := VisitorOnline{
			Id:          visitor.ID,
			VisitorId:   visitor.VisitorId,
			Avator:      visitor.Avator,
			Ip:          visitor.SourceIp,
			Username:    fmt.Sprintf("#%d %s", visitor.ID, visitor.Name),
			LastMessage: messagesMap[visitor.VisitorId],
			UpdatedAt:   visitor.UpdatedAt,
			UnreadNum:   unreadNum,
			Status:      visitor.Status,
		}
		users = append(users, user)
	}
	count := models.CountVisitorsByKefuId(kefuId.(string))
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"result": gin.H{
			"list":     users,
			"count":    count,
			"pagesize": pagesize,
		},
	})
}
func GetVisitorsList(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	pagesize, _ := strconv.Atoi(c.Query("pagesize"))
	entId, _ := c.Get("ent_id")
	search := ""
	args := []interface{}{}
	//通过访客名搜索
	visitorName := c.Query("visitorName")
	if visitorName != "" {
		search += " and (name like ? or visitor_id = ?)"
		args = append(args, visitorName+"%")
		args = append(args, fmt.Sprintf("%v_%s", entId, visitorName))
	}
	//通过客服名搜索
	kefuName := c.Query("kefuName")
	if kefuName != "" {
		search += " and to_id = ? "
		args = append(args, kefuName)
	}
	//根据tag找出visitor
	visitorTag := c.Query("visitorTag")
	var visitorIdsArr []string
	if visitorTag != "" {
		visitorTags := v2.GetVisitorTags("tag_id = ? ", visitorTag)
		for _, visitorTag := range visitorTags {
			visitorIdsArr = append(visitorIdsArr, visitorTag.VisitorId)
		}
	}
	if len(visitorIdsArr) != 0 {
		search += " and visitor_id in (?) "
		args = append(args, visitorIdsArr)
	}
	//排序参数
	orderBy := c.Query("orderBy")
	if orderBy == "" {
		orderBy = "updated_at desc"
	} else {
		orderBy = orderBy + " desc"
	}
	vistors := models.FindVisitorsByEntId(uint(page), uint(pagesize), fmt.Sprintf("%v", entId), orderBy, search, args...)
	count := models.CountVisitorsByEntid(fmt.Sprintf("%v", entId), search, args...)

	//获取最后一条消息内容
	visitorIds := make([]string, 0)
	for _, visitor := range vistors {
		visitorIds = append(visitorIds, visitor.VisitorId)
	}
	messagesMap := models.FindLastMessageMap(visitorIds)
	//获取访客未读数
	unreadMap := models.FindUnreadMessageNumByVisitorIds(visitorIds, "visitor")
	//log.Println(unreadMap)
	users := make([]VisitorOnline, 0)
	for _, visitor := range vistors {
		var unreadNum uint32
		if num, ok := unreadMap[visitor.VisitorId]; ok {
			unreadNum = num
		}
		username := fmt.Sprintf("#%d %s", visitor.ID, visitor.Name)
		if visitor.RealName != "" {
			username = visitor.RealName
		}
		user := VisitorOnline{
			Id:          visitor.ID,
			VisitorId:   visitor.VisitorId,
			City:        visitor.City,
			Avator:      visitor.Avator,
			Ip:          visitor.SourceIp,
			Username:    username,
			LastMessage: messagesMap[visitor.VisitorId],
			UpdatedAt:   visitor.UpdatedAt,
			UnreadNum:   unreadNum,
			Status:      visitor.Status,
		}
		users = append(users, user)
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"result": gin.H{
			"list":     users,
			"count":    count,
			"pagesize": uint(pagesize),
			"page":     page,
		},
	})
}

func GetVisitorMessageByKefu(c *gin.Context) {
	entId, _ := c.Get("ent_id")
	visitorId := c.Query("visitor_id")
	log.Println(entId, visitorId)
	messages := models.FindMessageByWhere("message.visitor_id=? and message.ent_id=?", visitorId, entId)
	chatMessages := make([]ChatMessage, 0)

	for _, message := range messages {
		var chatMessage ChatMessage
		chatMessage.Time = message.CreatedAt.Format("2006-01-02 15:04:05")
		chatMessage.Content = message.Content
		chatMessage.MesType = message.MesType
		if message.MesType == "kefu" {
			chatMessage.Name = message.KefuName
			chatMessage.Avator = message.KefuAvator
		} else {
			chatMessage.Name = message.VisitorName
			chatMessage.Avator = message.VisitorAvator
		}
		chatMessages = append(chatMessages, chatMessage)
	}
	models.ReadMessageByVisitorId(visitorId, "visitor")
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": chatMessages,
	})
}

// @Summary 获取在线访客列表接口
// @Produce  json
// @Success 200 {object} controller.Response
// @Failure 200 {object} controller.Response
// @Router /visitors_online [get]
func GetVisitorOnlines(c *gin.Context) {
	users := make([]map[string]string, 0)
	visitorIds := make([]string, 0)
	for uid, visitor := range ws.ClientList {
		userInfo := make(map[string]string)
		userInfo["uid"] = uid
		userInfo["name"] = visitor.Name
		userInfo["avator"] = visitor.Avator
		users = append(users, userInfo)
		visitorIds = append(visitorIds, visitor.Id)
	}

	//查询最新消息
	messages := models.FindLastMessage(visitorIds)
	temp := make(map[string]string, 0)
	for _, mes := range messages {
		temp[mes.VisitorId] = mes.Content
	}
	for _, user := range users {
		user["last_message"] = temp[user["uid"]]
	}

	tcps := make([]string, 0)
	for ip, _ := range clientTcpList {
		tcps = append(tcps, ip)
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"result": gin.H{
			"ws":  users,
			"tcp": tcps,
		},
	})
}

/**
设置访客属性
*/
func PostVisitorAttrs(c *gin.Context) {
	entId, _ := c.Get("ent_id")
	var param VisitorAttrParams
	if err := c.BindJSON(&param); err != nil {
		c.String(200, err.Error())
	}
	attrInfo := models.GetVisitorAttrByVisitorId(param.VisitorId, fmt.Sprintf("%v", entId))

	if attrInfo.ID == 0 {
		models.CreateVisitorAttr(
			fmt.Sprintf("%v", entId),
			param.VisitorId,
			param.VisitorAttr.RealName,
			param.VisitorAttr.Tel,
			param.VisitorAttr.Email,
			param.VisitorAttr.QQ,
			param.VisitorAttr.Wechat,
			param.VisitorAttr.Comment)
	} else {
		var attr = &models.Visitor_attr{
			RealName: param.VisitorAttr.RealName,
			Tel:      param.VisitorAttr.Tel,
			Email:    param.VisitorAttr.Email,
			QQ:       param.VisitorAttr.QQ,
			Wechat:   param.VisitorAttr.Wechat,
			Comment:  param.VisitorAttr.Comment,
		}
		models.SaveVisitorAttrByVisitorId(attr, param.VisitorId, fmt.Sprintf("%v", entId))
	}
	if param.VisitorAttr.RealName != "" {
		models.UpdateVisitorRealName(param.VisitorAttr.RealName,
			fmt.Sprintf("%v", entId),
			param.VisitorId)
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}

/**
设置访客属性
*/
func GetVisitorAttr(c *gin.Context) {
	entId, _ := c.Get("ent_id")
	visitorId := c.Query("visitor_id")
	attrInfo := models.GetVisitorAttrByVisitorId(visitorId, fmt.Sprintf("%v", entId))

	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": attrInfo,
	})
}

// @Summary 获取客服的在线访客列表接口
// @Produce  json
// @Success 200 {object} controller.Response
// @Failure 200 {object} controller.Response
// @Router /visitors_kefu_online [get]
func GetKefusVisitorOnlines(c *gin.Context) {
	kefuName, _ := c.Get("kefu_name")
	entId, _ := c.Get("ent_id")
	kefuInfo := models.FindUserByUid(entId)
	if kefuInfo.Status != 0 && kefuInfo.Status == 1 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "账户禁用",
		})
		return
	}
	users := make([]*VisitorOnline, 0)
	visitorIds := make([]string, 0)
	clientList := sortMapToSlice(ws.ClientList)
	for _, visitor := range clientList {
		if visitor.To_id != kefuName {
			continue
		}
		userInfo := new(VisitorOnline)
		userInfo.UpdatedAt = visitor.UpdateTime
		userInfo.VisitorId = visitor.Id
		userInfo.Username = visitor.Name
		userInfo.Avator = visitor.Avator
		users = append(users, userInfo)
		visitorIds = append(visitorIds, visitor.Id)
	}

	//查询最新消息
	messages := models.FindLastMessageMap(visitorIds)
	//查未读数
	unreadMap := models.FindUnreadMessageNumByVisitorIds(visitorIds, "visitor")
	for _, user := range users {
		user.LastMessage = messages[user.VisitorId]
		if user.LastMessage == "" {
			user.LastMessage = "new visitor"
		}
		var unreadNum uint32
		if num, ok := unreadMap[user.VisitorId]; ok {
			unreadNum = num
		}
		user.UnreadNum = unreadNum
	}

	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": users,
	})
}

func DelVisitor(c *gin.Context) {
	visitorId := c.Query("visitor_id")
	entIdStr, _ := c.Get("ent_id")
	entId, _ := strconv.Atoi(entIdStr.(string))

	models.DelVisitor("ent_id = ? and visitor_id = ?", entId, visitorId)
	models.DelMessage("ent_id = ? and visitor_id = ?", entId, visitorId)
	models.DelVisitorAttr("ent_id = ? and visitor_id = ?", entId, visitorId)
	c.JSON(200, gin.H{
		"code": types.ApiCode.SUCCESS,
		"msg":  types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
	})
}

func sortMapToSlice(youMap map[string]*ws.User) []*ws.User {
	keys := make([]string, 0)
	for k, _ := range youMap {
		keys = append(keys, k)
	}
	myMap := make([]*ws.User, 0)
	sort.Strings(keys)
	for _, k := range keys {
		myMap = append(myMap, youMap[k])
	}
	return myMap
}
