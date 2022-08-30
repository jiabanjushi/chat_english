package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/tidwall/gjson"
	"go-fly-muti/common"
	"go-fly-muti/models"
	"go-fly-muti/tools"
	"go-fly-muti/ws"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func PostInstall(c *gin.Context) {
	notExist, _ := tools.IsFileNotExist("./install.lock")
	if !notExist {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "系统已经安装过了",
		})
		return
	}
	server := c.PostForm("server")
	port := c.PostForm("port")
	database := c.PostForm("database")
	username := c.PostForm("username")
	password := c.PostForm("password")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, server, port, database)
	_, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "数据库连接失败:" + err.Error(),
		})
		return
	}
	isExist, _ := tools.IsFileExist(common.Dir)
	if !isExist {
		os.Mkdir(common.Dir, os.ModePerm)
	}
	fileConfig := common.MysqlConf
	file, _ := os.OpenFile(fileConfig, os.O_RDWR|os.O_CREATE, os.ModePerm)

	format := `{
	"Server":"%s",
	"Port":"%s",
	"Database":"%s",
	"Username":"%s",
	"Password":"%s"
}
`
	data := fmt.Sprintf(format, server, port, database, username, password)
	file.WriteString(data)
	//models.Connect()
	installFile, _ := os.OpenFile("./install.lock", os.O_RDWR|os.O_CREATE, os.ModePerm)
	installFile.WriteString("gofly live chat")
	ok, err := install()
	if !ok {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "安装成功",
	})
}
func install() (bool, error) {
	sqlFile := common.Dir + "go-fly.sql"
	isExit, _ := tools.IsFileExist(common.MysqlConf)
	dataExit, _ := tools.IsFileExist(sqlFile)
	if !isExit || !dataExit {
		return false, errors.New("config/mysql.json 数据库配置文件或者数据库文件go-fly.sql不存在")
	}
	sqls, _ := ioutil.ReadFile(sqlFile)
	sqlArr := strings.Split(string(sqls), "|")
	for _, sql := range sqlArr {
		if sql == "" {
			continue
		}
		err := models.Execute(sql)
		if err == nil {
			log.Println(sql, "\t success!")
		} else {
			log.Println(sql, err, "\t failed!")
		}
	}
	return true, nil
}
func MainCheckAuth(c *gin.Context) {
	id, _ := c.Get("kefu_id")
	userinfo := models.FindUserRole("user.nickname,user.avator,user.name,user.id, role.name role_name", id)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "验证成功",
		"result": gin.H{
			"avator":    userinfo.Avator,
			"name":      userinfo.Name,
			"role_name": userinfo.RoleName,
			"nick_name": userinfo.Nickname,
		},
	})
}
func GetStatistics(c *gin.Context) {
	kefuName, _ := c.Get("kefu_name")
	//kefuId, _ := c.Get("kefu_id")
	entId, _ := c.Get("ent_id")
	//今日访客数
	todayStart := time.Now().Format("2006-01-02")
	todayEnd := fmt.Sprintf("%s 23:59:59", todayStart)
	toadyVisitors := models.CountVisitors("to_id= ? and updated_at>= ? and updated_at<= ?", kefuName.(string), todayStart, todayEnd)
	visitors := models.CountVisitorsByKefuId(kefuName.(string))

	message := models.CountMessage("kefu_id=?", kefuName)
	todayMessages := models.CountMessage("kefu_id= ? and created_at>= ? and created_at<= ?", kefuName.(string), todayStart, todayEnd)
	visitorSession := 0
	for _, c := range ws.ClientList {
		if c.Ent_id == fmt.Sprintf("%v", entId) {
			visitorSession++
		}
	}
	kefuSession := 0
	for _, kefus := range ws.KefuList {
		for _, c := range kefus {
			if c.Ent_id == fmt.Sprintf("%v", entId) {
				kefuSession++
			}
		}
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"result": gin.H{
			"visitors":        visitors,
			"toady_visitors":  toadyVisitors,
			"today_messages":  todayMessages,
			"message":         message,
			"visitor_session": visitorSession,
			"kefu_session":    kefuSession,
		},
	})
}
func GetVersion(c *gin.Context) {

	versionName := "商务版"
	ipAuth := tools.Get(IP_SERVER_URL)

	ipAddress := gjson.Get(ipAuth, "result.ip_address").String()
	endTimeStr := gjson.Get(ipAuth, "result.expire_time").String()
	content := gjson.Get(ipAuth, "result.content").String()
	//if common.IsTry {
	//	versionName = "试用版"
	//	installTimeByte, _ := ioutil.ReadFile("./install.lock")
	//	installTimeByte, _ = tools.AesDecrypt(installTimeByte, []byte(common.AesKey))
	//	installTime := tools.ByteToInt64(installTimeByte)
	//	endTime := installTime + common.TryDeadline
	//	endTimeStr = time.Unix(endTime, 0).Format("2006-01-02")
	//}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"result": gin.H{
			"ip_address":   ipAddress,
			"last_time":    endTimeStr,
			"version":      versionName,
			"content":      content,
			"version_code": common.Version,
		},
	})
}
func GetOtherVersion(c *gin.Context) {
	versionCode := models.FindConfig("SystemVersion")
	versionName := models.FindConfig("SystemVersionName")
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"result": gin.H{
			"version_name": versionName,
			"version_code": versionCode,
		},
	})
}
