package controller

import (
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-fly-muti/common"
	"go-fly-muti/models"
	"go-fly-muti/tools"
	"go-fly-muti/ws"
	"strconv"
)

func GetKefuInfo(c *gin.Context) {
	kefuId, _ := c.Get("kefu_id")
	entIdStr, _ := c.Get("ent_id")
	entId, _ := strconv.Atoi(entIdStr.(string))
	//user := models.FindUserById(kefuId)
	//info := make(map[string]interface{})
	//info["name"] = user.Nickname
	//info["id"] = user.Name
	//info["avator"] = user.Avator
	//info["role_id"] = user.RoleId
	user := models.User{ID: uint(kefuId.(float64))}
	result := user.GetOneUser("*")
	if result.ID == 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "客服不存在",
		})
		return
	}
	result.Password = ""
	result.EntId = entId
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": result,
	})
}
func GetKefuInfoAll(c *gin.Context) {
	id, _ := c.Get("kefu_id")
	userinfo := models.FindUserRole("user.avator,user.name,user.id, role.name role_name", id)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "验证成功",
		"result": userinfo,
	})
}
func GetOtherKefuList(c *gin.Context) {
	idInterface, _ := c.Get("kefu_id")
	entId, _ := c.Get("ent_id")
	id := idInterface.(float64)
	result := make([]interface{}, 0)
	ws.SendPingToKefuClient()
	kefus := models.FindUsersByEntId(entId)
	for _, kefu := range kefus {
		if uint(id) == kefu.ID {
			continue
		}

		item := make(map[string]interface{})
		item["name"] = kefu.Name
		item["nickname"] = kefu.Nickname
		item["avator"] = kefu.Avator
		item["status"] = "offline"
		kefus, ok := ws.KefuList[kefu.Name]
		if ok && len(kefus) != 0 {
			item["status"] = "online"
		}
		result = append(result, item)
	}
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": result,
	})
}
func PostTransKefu(c *gin.Context) {
	kefuId := c.Query("kefu_id")
	visitorId := c.Query("visitor_id")
	curKefuId, _ := c.Get("kefu_name")
	user := models.FindUser(kefuId)
	visitor := models.FindVisitorByVistorId(visitorId)
	if user.Name == "" || visitor.Name == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "访客或客服不存在",
		})
		return
	}
	models.UpdateVisitorKefu(visitorId, kefuId)
	ws.UpdateVisitorUser(visitorId, kefuId)
	go ws.VisitorOnline(kefuId, visitor)
	go ws.VisitorOffline(curKefuId.(string), visitor.VisitorId, visitor.Name)
	ws.VisitorNotice(visitor.VisitorId, "客服转接到"+user.Nickname)
	ws.VisitorTransfer(visitor.VisitorId, kefuId)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "转移成功",
	})
}
func GetKefuInfoSettingOwn(c *gin.Context) {
	kefuId := c.Query("kefu_id")
	pid, _ := c.Get("kefu_id")
	roleId, _ := c.Get("role_id")
	var user models.User
	if roleId.(float64) == 1 {
		user = models.FindUserById(kefuId)
	} else {
		user = models.FindUserByIdPid(pid, kefuId)

	}
	user.Password = ""
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": user,
	})
}
func GetKefuInfoSetting(c *gin.Context) {
	kefuId := c.Query("kefu_id")
	user := models.FindUserById(kefuId)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": user,
	})
}
func PostKefuRegister(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	rePassword := c.PostForm("rePassword")
	avator := "/static/images/4.jpg"
	nickname := c.PostForm("nickname")
	captchaCode := c.PostForm("captcha")
	roleId := 2
	if name == "" || password == "" || rePassword == "" || nickname == "" || captchaCode == "" {
		c.JSON(200, gin.H{
			"code":   400,
			"msg":    "参数不能为空",
			"result": "",
		})
		return
	}
	if password != rePassword {
		c.JSON(200, gin.H{
			"code":   400,
			"msg":    "密码不一致",
			"result": "",
		})
		return
	}
	oldUser := models.FindUser(name)
	if oldUser.Name != "" {
		c.JSON(200, gin.H{
			"code":   400,
			"msg":    "用户名已经存在",
			"result": "",
		})
		return
	}
	session := sessions.Default(c)
	if captchaId := session.Get("captcha"); captchaId != nil {
		session.Delete("captcha")
		_ = session.Save()
		if !captcha.VerifyString(captchaId.(string), captchaCode) {
			c.JSON(200, gin.H{
				"code":   400,
				"msg":    "验证码验证失败",
				"result": "",
			})
			return
		}
	} else {
		c.JSON(200, gin.H{
			"code":   400,
			"msg":    "验证码失效",
			"result": "",
		})
		return
	}
	//插入新用户
	uid := models.CreateUser(name, tools.Md5(password), avator, nickname, 1, 10)
	if uid == 0 {
		c.JSON(200, gin.H{
			"code":   400,
			"msg":    "增加用户失败",
			"result": "",
		})
		return
	}
	models.CreateUserRole(uid, uint(roleId))

	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "注册完成",
		"result": "",
	})
}
func PostKefuAvator(c *gin.Context) {

	avator := c.PostForm("avator")
	if avator == "" {
		c.JSON(200, gin.H{
			"code":   400,
			"msg":    "不能为空",
			"result": "",
		})
		return
	}
	kefuName, _ := c.Get("kefu_name")
	models.UpdateUserAvator(kefuName.(string), avator)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": "",
	})
}
func PostKefuinfo(c *gin.Context) {

	avator := c.PostForm("avator")
	nickname := c.PostForm("nickname")
	if avator == "" || nickname == "" {
		c.JSON(200, gin.H{
			"code":   400,
			"msg":    "参数不能为空",
			"result": "",
		})
		return
	}
	kefuName, _ := c.Get("kefu_name")
	models.UpdateKefuInfoByName(kefuName.(string), avator, nickname)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": "",
	})
}
func PostKefuPass(c *gin.Context) {
	kefuName, _ := c.Get("kefu_name")
	newPass := c.PostForm("new_pass")
	confirmNewPass := c.PostForm("confirm_new_pass")
	old_pass := c.PostForm("old_pass")
	if newPass != confirmNewPass {
		c.JSON(200, gin.H{
			"code":   400,
			"msg":    "密码不一致",
			"result": "",
		})
		return
	}
	user := models.FindUser(kefuName.(string))
	if user.Password != tools.Md5(old_pass) {
		c.JSON(200, gin.H{
			"code":   400,
			"msg":    "旧密码不正确",
			"result": "",
		})
		return
	}
	models.UpdateUserPass(kefuName.(string), tools.Md5(newPass))
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": "",
	})
}
func PostKefuInfoStatus(c *gin.Context) {
	kefuId, _ := c.Get("kefu_id")
	roleId, _ := c.Get("role_id")
	id := c.PostForm("id")
	status := c.PostForm("status")
	statusInt, _ := strconv.Atoi(status)
	query := " id = ? "
	arg := []interface{}{
		id,
	}
	if roleId.(float64) != 1 {
		query = query + " and pid = ? "
		arg = append(arg, kefuId)
	}

	models.UpdateUserStatusWhere(uint(statusInt), query, arg...)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
func PostKefuInfo(c *gin.Context) {
	kefuName, _ := c.Get("kefu_name")
	kefuId, _ := c.Get("kefu_id")
	mRoleId, _ := c.Get("role_id")

	kefuInfo := models.FindUser(kefuName.(string))
	var query string
	var arg = []interface{}{}
	if mRoleId.(float64) != 1 {
		query = "pid=?"
		arg = append(arg, kefuId)
		count := models.CountUsersWhere(query, arg)
		if kefuInfo.AgentNum <= count {
			c.JSON(200, gin.H{
				"code": 400,
				"msg":  fmt.Sprintf("子账号个数超过上限数: %d!", kefuInfo.AgentNum),
			})
			return
		}
	}

	id := c.PostForm("id")
	name := c.PostForm("name")
	password := c.PostForm("password")
	avator := c.PostForm("avator")
	nickname := c.PostForm("nickname")
	roleId := c.PostForm("role_id")
	agentNum, _ := strconv.Atoi(c.PostForm("agent_num"))
	if roleId == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "请选择角色!",
		})
		return
	}
	roleIdInt, _ := strconv.Atoi(roleId)
	if mRoleId.(float64) != 1 && roleIdInt <= int(mRoleId.(float64)) {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "修改角色无权限!",
		})
		return
	}
	//插入新用户
	if id == "" {
		oldUser := models.FindUser(name)
		if oldUser.Name != "" {
			c.JSON(200, gin.H{
				"code":   400,
				"msg":    "用户名已经存在",
				"result": "",
			})
			return
		}
		pid, _ := c.Get("kefu_id")
		uid := models.CreateUser(name, tools.Md5(password), avator, nickname, uint(pid.(float64)), uint(agentNum))
		if uid == 0 {
			c.JSON(200, gin.H{
				"code":   400,
				"msg":    "增加用户失败",
				"result": "",
			})
			return
		}
		roleIdInt, _ := strconv.Atoi(roleId)
		models.CreateUserRole(uid, uint(roleIdInt))
	} else {
		//更新用户
		if password != "" {
			password = tools.Md5(password)
		}
		oldUser := models.FindUser(name)
		if oldUser.Name != "" && id != fmt.Sprintf("%d", oldUser.ID) {
			c.JSON(200, gin.H{
				"code":   400,
				"msg":    "用户名已经存在",
				"result": "",
			})
			return
		}
		models.UpdateUser(id, name, password, avator, nickname, uint(agentNum))
		roleIdInt, _ := strconv.Atoi(roleId)
		uid, _ := strconv.Atoi(id)
		models.DeleteRoleByUserId(uid)
		models.CreateUserRole(uint(uid), uint(roleIdInt))
	}

	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": "",
	})
}

func GetKefuList(c *gin.Context) {
	roleId, _ := c.Get("role_id")
	if roleId.(float64) != 1 {
		c.JSON(200, gin.H{
			"code":   200,
			"msg":    "无权限",
			"result": "",
		})
	}
	page, _ := strconv.Atoi(c.Query("page"))
	count := models.CountUsers()
	list := models.FindUsersPages(uint(page), common.PageSize)
	//users := models.FindUsers()
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"result": gin.H{
			"count":    count,
			"page":     page,
			"list":     list,
			"pagesize": common.PageSize,
		},
	})
}
func GetKefuListOwn(c *gin.Context) {
	kefuId, _ := c.Get("kefu_id")
	roleId, _ := c.Get("role_id")
	query := "1=1 "
	var arg = []interface{}{}
	if roleId.(float64) != 1 {
		query += "and pid = ? "
		arg = append(arg, kefuId)
	}

	//通过客服名搜索
	kefuName := c.Query("kefu_name")
	if kefuName != "" {
		query += "and user.name like ? "
		arg = append(arg, kefuName+"%")
	}

	pagesize, _ := strconv.Atoi(c.Query("pagesize"))
	page, _ := strconv.Atoi(c.Query("page"))
	count := models.CountUsersWhere(query, arg...)
	list := models.FindUsersOwn(uint(page), uint(pagesize), query, arg...)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"result": gin.H{
			"count":    count,
			"page":     page,
			"list":     list,
			"pagesize": pagesize,
		},
	})
}
func GetKefuListMessages(c *gin.Context) {
	kefuId, _ := c.Get("kefu_id")
	roleId, _ := c.Get("role_id")
	id := c.Query("id")
	page, _ := strconv.Atoi(c.Query("page"))
	user := models.FindUserById(id)
	fmt.Printf("%T\n", kefuId)
	fmt.Printf("%T\n", user.Pid)
	if fmt.Sprintf("%v", user.Pid) != fmt.Sprintf("%v", kefuId) {
		if roleId.(float64) != 1 {
			c.JSON(200, gin.H{
				"code":   400,
				"msg":    "无权限",
				"result": "",
			})
			return
		}
	}
	count := models.CountMessage("message.kefu_id=?", user.Name)
	list := models.FindMessageByPage(uint(page), common.PageSize, "message.kefu_id=?", user.Name)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"result": gin.H{
			"count":    count,
			"page":     page,
			"list":     list,
			"pagesize": common.PageSize,
		},
	})
}
func GetVisitorListMessages(c *gin.Context) {
	visitorId := c.Query("visitor_id")
	entId, _ := c.Get("ent_id")
	page, _ := strconv.Atoi(c.Query("page"))

	count := models.CountMessage("message.ent_id= ? and message.visitor_id=?", entId, visitorId)
	list := models.FindMessageByPage(uint(page), common.PageSize, "message.ent_id= ? and message.visitor_id=?", entId, visitorId)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"result": gin.H{
			"count":    count,
			"page":     page,
			"list":     list,
			"pagesize": common.PageSize,
		},
	})
}
func DeleteKefuInfo(c *gin.Context) {
	kefuId := c.Query("id")
	models.DeleteUserById(kefuId)
	models.DeleteRoleByUserId(kefuId)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "删除成功",
		"result": "",
	})
}
func DeleteKefuInfoOwn(c *gin.Context) {
	kefuId := c.Query("id")
	pid, _ := c.Get("kefu_id")

	roleId, _ := c.Get("role_id")

	if kefuId == fmt.Sprintf("%v", pid) {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "不能删除自己",
		})
		return
	}
	if roleId.(float64) == 1 {
		models.DeleteUserById(kefuId)
		models.DeleteRoleByUserId(kefuId)
	} else {
		models.DeleteUserByIdPid(kefuId, pid)
		models.DeleteRoleByUserId(kefuId)
	}

	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "删除成功",
		"result": "",
	})
}

//更新用户的在线状态
func GetUpdateOnlineStatus(c *gin.Context) {
	status := c.Query("status")
	statusInt, _ := strconv.Atoi(status)
	kefuId, _ := c.Get("kefu_id")
	user := &models.User{
		ID:           uint(kefuId.(float64)),
		OnlineStatus: uint(statusInt),
	}
	user.UpdateUser()
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "成功",
	})
}
