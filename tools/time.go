package tools

import "time"

func Now() int64 {
	return time.Now().Unix()
}

//时间戳转时间字符串
func IntToTimeStr(sourceTime int64, formatStr string) string {
	return time.Unix(sourceTime, 0).Format(formatStr)
}

//时间字符串转时间戳
func TimeStrToInt(sourceTime string) int64 {
	timeLayout := "2006-01-02 15:04:05"                             //转化所需模板
	loc, _ := time.LoadLocation("Local")                            //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, sourceTime, loc) //使用模板在对应时区转化为time.time类型
	return theTime.Unix()
}
