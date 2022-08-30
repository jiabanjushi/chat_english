package tmpl

import (
	"github.com/gin-gonic/gin"
	"go-fly-muti/models"
	"html/template"
	"net/http"
)

//登陆界面
func PageLogin(c *gin.Context) {

	SystemNotice := models.FindConfig("SystemNotice")
	SystemTitle := models.FindConfig("SystemTitle")
	SystemLoginTitle := models.FindConfig("SystemLoginTitle")
	SystemKeywords := models.FindConfig("SystemKeywords")
	SystemDesc := models.FindConfig("SystemDesc")
	CopyrightTxt := models.FindConfig("CopyrightTxt")
	CopyrightUrl := models.FindConfig("CopyrightUrl")
	SystemKefu := models.FindConfig("SystemKefu")
	SystemRegister := models.FindConfig("SystemRegister")
	isRegister := true
	if SystemRegister == "2" {
		isRegister = false
	}
	c.HTML(http.StatusOK, "login.html", gin.H{
		"SystemNotice":     template.HTML(SystemNotice),
		"SystemTitle":      SystemTitle,
		"SystemKeywords":   SystemKeywords,
		"SystemDesc":       SystemDesc,
		"CopyrightTxt":     CopyrightTxt,
		"CopyrightUrl":     CopyrightUrl,
		"SystemLoginTitle": SystemLoginTitle,
		"isRegister":       isRegister,
		"SystemKefu":       template.HTML(SystemKefu),
	})
}
