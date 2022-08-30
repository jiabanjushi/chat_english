package models

import (
	"go-fly-muti/types"
	"time"
)

type IpAuth struct {
	ID         uint       `gorm:"primary_key" json:"id"`
	IpAddress  string     `json:"ip_address"`
	ExpireTime string     `json:"expire_time"`
	Content    string     `json:"content"`
	NowTime    string     `gorm:"-" json:"now_time"`
	status     uint       `json:"status"`
	CreatedAt  types.Time `json:"created_at"`
}

func CreateIpAuth(ip, expireTime string) IpAuth {
	model := IpAuth{
		IpAddress:  ip,
		ExpireTime: expireTime,
		CreatedAt: types.Time{
			time.Now(),
		},
	}
	DB.Create(&model)
	return model
}
func FindServerIpAddress(ip string) IpAuth {
	var ipAuth IpAuth
	DB.Where("ip_address = ?", ip).First(&ipAuth)
	return ipAuth
}
