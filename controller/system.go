package controller

import (
	"github.com/gin-gonic/gin"
)

var StopSign = make(chan int)

func GetStop(c *gin.Context) {

	//StopSign <- 1

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "沙比不讲诚信死全家",
	})
}
