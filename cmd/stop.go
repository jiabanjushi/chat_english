package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-fly-muti/tools"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止客服http服务",
	Run: func(cmd *cobra.Command, args []string) {

		pids, err := ioutil.ReadFile("gofly.sock")
		rootPath=tools.GetRootPath()
		if err != nil {
			fmt.Sprintf(err.Error())
			return
		}
		pidSlice := strings.Split(string(pids), ",")
		var command *exec.Cmd
		for _, pid := range pidSlice {
			fmt.Println(pid)
			if runtime.GOOS == "windows" {
				command = exec.Command("taskkill.exe", "/f", "/pid", pid)
			} else {
				command = exec.Command("kill", pid)
			}
			err := command.Start()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	},
}
