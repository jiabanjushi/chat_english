package models

import (
	"go-fly-muti/types"
)

type New struct {
	Id        uint       `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Tag       string     `json:"tag"`
	status    uint       `json:"status"`
	CreatedAt types.Time `json:"created_at"`
}

func CountNews(query interface{}, args ...interface{}) uint {
	var v uint
	DB.Table("new").Where(query, args...).Count(&v)
	return v
}
func (this *New) SaveNews(query interface{}, args ...interface{}) error {
	db := DB.Table("new").Where(query, args...).Update(this)
	return db.Error
}
func (this *New) AddNews() error {
	return DB.Create(this).Error
}
func FindNews(page, pagesize int, query interface{}, args ...interface{}) []New {
	offset := (page - 1) * pagesize
	var res []New
	DB.Table("new").Where(query, args...).Order("id desc").Offset(offset).Limit(pagesize).Find(&res)
	return res
}
func DelNews(query interface{}, args ...interface{}) error {
	return DB.Where(query, args...).Delete(&New{}).Error
}
