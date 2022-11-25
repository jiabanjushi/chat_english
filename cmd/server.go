package cmd

import (
	"context"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zh-five/xdaemon"
	"go-fly-muti/common"
	"go-fly-muti/controller"
	"go-fly-muti/middleware"
	"go-fly-muti/models"
	"go-fly-muti/router"
	"go-fly-muti/setting"
	"go-fly-muti/static"
	"go-fly-muti/tools"
	"go-fly-muti/ws"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	port     string
	daemon   bool
	rootPath string
)
var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "启动客服http服务",
	Example: "go-fly server",
	Run:     run,
}

func init() {
	serverCmd.PersistentFlags().StringVarP(&rootPath, "rootPath", "r", "", "程序根目录")
	serverCmd.PersistentFlags().StringVarP(&port, "port", "p", "8081", "监听端口号")
	serverCmd.PersistentFlags().BoolVarP(&daemon, "daemon", "d", false, "是否为守护进程模式")
}
func run(cmd *cobra.Command, args []string) {

	//初始化目录
	initDir()
	//初始化守护进程
	initDaemon()

	baseServer := "0.0.0.0:" + port

	//if common.RpcStatus {
	//	go frpc.NewRpcServer(common.RpcServer)
	//	log.Println("start rpc server...\r\ngo：tcp://" + common.RpcServer)
	//}
	//加载配置
	if err := setting.Init(); err != nil {
		fmt.Println("配置文件初始化事变", err)
		return
	}
	//设置时区
	loc, err := time.LoadLocation(viper.GetString("app.  timeZone"))
	if err == nil {
		time.Local = loc // -> this is setting the global timezone
		fmt.Println(time.Now().Format("2006-01-02 15:04:05 "))
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	//是否编译模板
	if common.IsCompireTemplate {
		templ := template.Must(template.New("").ParseFS(static.TemplatesEmbed, "templates/**/*.html"))
		engine.SetHTMLTemplate(templ)
	} else {
		engine.LoadHTMLGlob(common.StaticDirPath + "templates/**/*")
	}
	engine.Static("/static", common.StaticDirPath)
	engine.Use(tools.Session("gofly"))
	//跨域设置
	engine.Use(middleware.CrossSite)

	router.InitViewRouter(engine)
	router.InitApiRouter(engine)

	//限流类
	tools.NewLimitQueue()
	//清理
	//ws.CleanVisitorExpire()
	//后端websocket
	go ws.WsServerBackend()
	//初始化数据
	//logger := lib.NewLogger()
	models.NewConnect(common.ConfigDirPath + "/mysql.json")
	//初始化配置数据
	models.InitConfig()
	//后端定时客服
	go ws.UpdateVisitorStatusCron()
	log.Println("GOFLY服务开始运行:" + baseServer)
	//性能监控
	pprof.Register(engine)
	//engine.Run(baseServer)

	srv := &http.Server{
		Addr:    baseServer,
		Handler: engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("GOFLY服务监听: %s\n", err)
		}
	}()

	<-controller.StopSign
	log.Println("关闭服务...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务关闭失败:", err)
	}
	log.Println("服务已关闭")

}

//初始化目录
func initDir() {
	if rootPath == "" {
		rootPath = tools.GetRootPath()
	}
	log.Println("GOFLY服务运行路径:" + rootPath)
	common.RootPath = rootPath
	common.LogDirPath = rootPath + "/logs/"
	common.ConfigDirPath = rootPath + "/config/"
	common.StaticDirPath = rootPath + "/static/"
	common.UploadDirPath = rootPath + "/static/upload/"

	common.KFMYArray = strings.Split(common.KFMY, "\n")

	if noExist, _ := tools.IsFileNotExist(common.RootPath + "/install.lock"); noExist {
		panic("未检测到" + common.RootPath + "/install.lock,请先安装服务!")

	}

	if noExist, _ := tools.IsFileNotExist(common.LogDirPath); noExist {
		if err := os.MkdirAll(common.LogDirPath, 0777); err != nil {
			log.Println(err.Error())
		}
	}
	isMainUploadExist, _ := tools.IsFileExist(common.UploadDirPath)
	if !isMainUploadExist {
		os.Mkdir(common.UploadDirPath, os.ModePerm)
	}
}

//初始化守护进程
func initDaemon() {
	//启动进程之前要先杀死之前的金额

	pid, err := ioutil.ReadFile("Project.sock")
	if err != nil {
		return
	}
	pidSlice := strings.Split(string(pid), ",")
	var command *exec.Cmd
	for _, pid := range pidSlice {
		if runtime.GOOS == "windows" {
			command = exec.Command("taskkill.exe", "/f", "/pid", pid)
		} else {
			fmt.Println("成功结束进程:", pid)
			command = exec.Command("kill", pid)
		}
		command.Start()
	}

	if daemon == true {
		d := xdaemon.NewDaemon(common.LogDirPath + "gofly.log")
		d.MaxError = 10
		d.Run()
	}
	//记录pid
	ioutil.WriteFile(common.RootPath+"/gofly.sock", []byte(fmt.Sprintf("%d,%d", os.Getppid(), os.Getpid())), 0666)
}
