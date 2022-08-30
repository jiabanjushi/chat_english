package models

import (
	"go-fly-muti/types"
)

type VisitorBlack struct {
	Id        uint       `json:"id"`
	VisitorId string     `json:"visitor_id"`
	Name      string     `json:"name"`
	EntId     string     `json:"ent_id"`
	KefuName  string     `json:"kefu_name"`
	CreatedAt types.Time `json:"created_at"`
}

func CountVisitorBlack(query interface{}, args ...interface{}) uint {
	var v uint
	DB.Table("visitor_black").Where(query, args...).Count(&v)
	return v
}

func (this *VisitorBlack) AddVisitorBlack() error {
	return DB.Create(this).Error
}
func FindVisitorBlacks(page, pagesize int, query interface{}, args ...interface{}) []VisitorBlack {
	offset := (page - 1) * pagesize
	var res []VisitorBlack
	DB.Table("visitor_black").Where(query, args...).Order("id desc").Offset(offset).Limit(pagesize).Find(&res)
	return res
}
func FindVisitorBlack(query interface{}, args ...interface{}) VisitorBlack {
	var res VisitorBlack
	DB.Table("visitor_black").Where(query, args...).Find(&res)
	return res
}
func DelVisitorBlack(query interface{}, args ...interface{}) error {
	return DB.Where(query, args...).Delete(&VisitorBlack{}).Error
}
