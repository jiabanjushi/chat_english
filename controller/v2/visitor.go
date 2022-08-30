package v2

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-fly-muti/models"
	v2 "go-fly-muti/models/v2"
	"go-fly-muti/tools"
	"go-fly-muti/types"
	"go-fly-muti/ws"
	"strconv"
	"time"
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
type VisitorExtra struct {
	VisitorName   string `json:"visitorName"`
	VisitorAvatar string `json:"visitorAvatar"`
	VisitorId     string `json:"visitorId"`
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
	form.UserAgent = c.GetHeader("User-Agent")
	if tools.IsMobile(form.UserAgent) {
		form.Avator = "/static/images/phone.png"
	} else {
		form.Avator = "/static/images/computer.png"

	}
	if form.VisitorId == "" {
		form.VisitorId = tools.Uuid()
	}
	ipCity := tools.ParseIpNew(c.ClientIP())
	if ipCity != nil {
		form.CityAddress = ipCity.CountryName + ipCity.RegionName + ipCity.CityName
	} else {
		form.CityAddress = "Unknown Area"
	}
	form.VisitorName = form.CityAddress
	form.ClientIp = c.ClientIP()
	c.JSON(200, VisitorLogin(form))
}
func VisitorLogin(form VisitorLoginForm) gin.H {
	makeVisitorLoginForm(&form)
	//查看企业状态
	entId, _ := strconv.Atoi(form.EntId)
	userModel := &models.User{ID: uint(entId), Name: form.ToId}
	entUserInfo := userModel.GetOneUser("*")
	ok, errCode, errMsg := entUserInfo.CheckStatusExpired()
	if !ok {
		return gin.H{
			"code": errCode,
			"msg":  errMsg,
		}
	}
	//自动分配
	serviceUser := entUserInfo
	allOffline := true
	if _, ok := ws.KefuList[form.ToId]; ok {
		serviceUser = entUserInfo
		allOffline = false
	} else {
		userModel = &models.User{Pid: uint(entId)}
		userModel.SetOrder("rec_num asc")
		kefus := userModel.GetUsers("name,nickname,avator")
		if len(kefus) == 0 {
			serviceUser = entUserInfo
			form.ToId = entUserInfo.Name
		} else {
			for _, kefu := range kefus {
				if _, ok := ws.KefuList[kefu.Name]; ok {
					serviceUser = kefu
					form.ToId = kefu.Name
					allOffline = false
					break
				}
			}
		}
	}
	visitor := createVisitor(form)
	return gin.H{
		"code": types.ApiCode.SUCCESS,
		"msg":  types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
		"result": gin.H{
			"allOffline": allOffline,
			"result":     visitor,
			"kefu": gin.H{
				"nickname": serviceUser.Nickname,
				"avatar":   serviceUser.Avator,
			},
		},
	}
}
func createVisitor(form VisitorLoginForm) *v2.Visitor {
	visitor := &v2.Visitor{
		ID:        0,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Name:      form.VisitorName,
		Avator:    form.Avator,
		SourceIp:  form.ClientIp,
		ToId:      form.ToId,
		VisitorId: form.VisitorId,
		Status:    0,
		Refer:     form.Refer,
		City:      form.CityAddress,
		ClientIp:  form.ClientIp,
		Extra:     form.Extra,
		EntId:     form.EntId,
	}
	oldVsitor := visitor.FindVisitor()
	if oldVsitor.ID != 0 {
		visitor.ID = oldVsitor.ID
		visitor.UpdateVisitor("visitor_id=? and ent_id=?", visitor.VisitorId, visitor.EntId)
		return visitor
	}
	visitor.InsertVisitor()
	return visitor
}
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
