package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "go-fly-pro",
	Short: "go-fly-pro",
	Long:  `简洁快速的GO语言在线客服系统GOFLY`,
}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(indexCmd)
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println("执行命令参数错误:", err)
		os.Exit(1)
	}
}
