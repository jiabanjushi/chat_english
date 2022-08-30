package controller

import (
	"github.com/gin-gonic/gin"
	"go-fly-muti/models"
)

func GetRoleListOwn(c *gin.Context) {
	roleId, _ := c.Get("role_id")
	var roles []models.Role
	if roleId.(float64) == 1 {
		roles = models.FindRoles()
	} else {
		roles = models.FindRolesOwn(roleId)
	}

	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "获取成功",
		"result": roles,
	})
}
func GetRoleList(c *gin.Context) {
	roles := models.FindRoles()
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "获取成功",
		"result": roles,
	})
}
func PostRole(c *gin.Context) {
	roleId := c.PostForm("id")
	method := c.PostForm("method")
	name := c.PostForm("name")
	path := c.PostForm("path")
	if roleId == "" || method == "" || name == "" || path == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "参数不能为空",
		})
		return
	}
	models.SaveRole(roleId, name, method, path)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})
}
