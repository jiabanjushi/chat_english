package models

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-fly-muti/tools"
	"go-fly-muti/types"
	"io/ioutil"
	"log"
	"time"
)

var DB *gorm.DB

type Model struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewConnect(mysqlConfigFile string) error {
	var mysql = &types.Mysql{
		Username: "go_fly_pro",
		Password: "go_fly_pro",
		Database: "go_fly_pro",
		Server:   "127.0.0.1",
		Port:     "3306",
	}
	isExist, _ := tools.IsFileExist(mysqlConfigFile)
	if !isExist {
		panic("MYSQL配置文件不存在!" + mysqlConfigFile)
	}
	info, err := ioutil.ReadFile(mysqlConfigFile)
	if err != nil {
		panic("MYSQL配置文件读取失败!" + err.Error())
	}
	err = json.Unmarshal(info, mysql)


	if err != nil {
		panic("解析MYSQL配置文件JSON结构失败!" + err.Error())
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysql.Username, mysql.Password, mysql.Server, mysql.Port, mysql.Database)
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
		panic("数据库连接失败!")
		return err
	}
	DB.SingularTable(true)
	DB.LogMode(true)
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
	DB.DB().SetConnMaxLifetime(59 * time.Second)
	return nil
}

func Execute(sql string) error {
	db := DB.Exec(sql)
	err := db.Error
	if err != nil {
		log.Println("models.go sql execute error:" + err.Error())
		return err
	}
	return nil
}
func CloseDB() {
	defer DB.Close()
}
