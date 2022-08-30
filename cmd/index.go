package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

var (
	indexPort   string
	indexDaemon bool
)
var indexCmd = &cobra.Command{
	Use:     "index",
	Short:   "启动官网http服务",
	Example: "go-fly index",
	Run:     indexRun,
}

func init() {
	serverCmd.PersistentFlags().StringVarP(&indexPort, "indexPort", "", "8082", "官网端口")
	serverCmd.PersistentFlags().BoolVarP(&indexDaemon, "indexDaemon", "", true, "是否官网守护")
}
func indexRun(cmd *cobra.Command, args []string) {
	baseServer := "0.0.0.0:" + indexPort
	log.Println("start index server...\r\ngo：http://" + baseServer)

	engine := gin.Default()

	engine.GET("", func(context *gin.Context) {

		context.JSON(http.StatusOK, gin.H{
			"code": 200,
		})
	})
	engine.Run(baseServer)
}
