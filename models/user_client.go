package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User_client struct {
	ID         uint   `gorm:"primary_key" json:"id"`
	Kefu       string `json:"kefu"`
	Client_id  string `json:"client_id"`
	Created_at string `json:"created_at"`
}

func CreateUserClient(kefu, clientId string) uint {
	u := &User_client{
		Kefu:       kefu,
		Client_id:  clientId,
		Created_at: time.Now().Format("2006-01-02 15:04:05"),
	}
	DB.Create(u)
	return u.ID
}
func (this *User_client) FindClient() User_client {
	var info User_client
	this.buildQuery().First(&info)
	return info
}
func (this *User_client) FindClients() []User_client {
	var arr []User_client
	this.buildQuery().Find(&arr)
	return arr
}
func (this *User_client) DeleteClient() User_client {
	var info User_client
	this.buildQuery().Delete(info)
	return info
}

//查询构造
func (this *User_client) buildQuery() *gorm.DB {
	db := DB
	db.Model(this)
	if this.ID != 0 {
		db = db.Where("id = ?", this.ID)
	}
	if this.Client_id != "" {
		db = db.Where("client_id = ?", this.Client_id)
	}
	if this.Kefu != "" {
		db = db.Where("Kefu = ?", this.Kefu)
	}
	return db
}
