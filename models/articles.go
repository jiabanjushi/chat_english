package models

import (
	"github.com/jinzhu/gorm"
)

type ArticleCate struct {
	Id      uint   `json:"id"`
	CatName string `json:"cat_name"`
	UserId  string `json:"user_id"`
	EntId   string `json:"ent_id"`
	IsTop   uint   `json:"is_top"`
}
type Article struct {
	Id      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	CatId   uint   `json:"cat_id"`
	UserId  string `json:"user_id"`
	EntId   string `json:"ent_id"`
}

func FindArticleCates(userId interface{}) []ArticleCate {
	var list []ArticleCate
	DB.Table("article_cate").Where("user_id = ?", userId).Order("id desc").Find(&list)
	return list
}
func CountArticleList(query interface{}, args ...interface{}) uint {
	var v uint
	DB.Table("article").Where(query, args...).Count(&v)
	return v
}
func FindArticleList(page uint, pagesize uint, query interface{}, args ...interface{}) []Article {
	offset := (page - 1) * pagesize
	if offset < 0 {
		offset = 0
	}
	var list []Article
	DB.Table("article").Where(query, args...).Offset(offset).Limit(pagesize).Order("id desc").Find(&list)
	return list
}
func FindArticleCatesByEnt(userId interface{}) []ArticleCate {
	var list []ArticleCate
	DB.Table("article_cate").Where("ent_id = ?", userId).Order("id desc").Find(&list)
	return list
}

func FindArticleRow(query interface{}, args ...interface{}) Article {
	var res Article
	DB.Table("article").Where(query, args...).Order("id desc").Find(&res)
	return res
}
func DelArticles(query interface{}, args ...interface{}) {
	DB.Where(query, args...).Delete(&Article{})
}
func DelArticleCate(query interface{}, args ...interface{}) {
	DB.Where(query, args...).Delete(&ArticleCate{})
}
func (a *Article) SaveArticle(query interface{}, args ...interface{}) {
	DB.Model(&Article{}).Where(query, args...).Update(a)
}

func (a *Article) AddArticle() uint {
	DB.Create(a)
	return a.Id
}
func (a *ArticleCate) AddArticleCate() uint {
	DB.Create(a)
	return a.Id
}

//查询构造
func (a *Article) buildQuery() *gorm.DB {
	myDB := DB
	myDB.Model(a)
	if a.Id != 0 {
		myDB = myDB.Where("id = ?", a.Id)
	}
	return myDB
}
