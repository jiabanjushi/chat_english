package lib

import (
	"github.com/sirupsen/logrus"
	"go-fly-muti/common"
	"log"
	"os"
	"path"
	"time"
)

var logrusObj *logrus.Logger

func NewLogger() *logrus.Logger {

	if logrusObj != nil {
		src, _ := setOutputFile(common.LogDirPath)
		//设置输出
		logrusObj.Out = src
		return logrusObj
	}

	//实例化
	logger := logrus.New()
	src, _ := setOutputFile(common.LogDirPath)
	//设置输出
	logger.Out = src
	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrusObj = logger
	return logger
}
func setOutputFile(logFilePath string) (*os.File, error) {
	now := time.Now()
	logFileName := now.Format("2006-01-02") + ".log"
	//日志文件
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	//写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return src, nil
}
