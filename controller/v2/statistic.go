package v2

import (
	"github.com/gin-gonic/gin"
	"go-fly-muti/models"
)

func GetChartStatistic(c *gin.Context) {
	kefuName, _ := c.Get("kefu_name")

	//今天0点日期字符串
	//todayTimeStr := time.Now().Format("2006-01-02")
	////日期字符串转时间
	//todayTime, _ := time.Parse("2006-01-02", todayTimeStr)
	////时间转时间戳
	//todayTimenum := todayTime.Unix()
	////15天前开始的时间
	//startTimenum := todayTimenum - 15*24*3600
	////时间戳转时间转日期字符串
	//startTime := time.Unix(startTimenum, 0).Format("2006-01-02")

	result := models.CountVisitorsEveryDay(kefuName.(string))
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": result,
	})

}
