package v2

import (
	"go-fly-muti/models"
	"time"
)

type Tag struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Kefu      string    `json:"kefu"`
	EntId     uint      `json:"ent_id"`
	IsTaged   uint      `json:"is_taged" sql:"-"`
}

func (this *Tag) InsertTag() *Tag {
	models.DB.Create(this)
	return this
}

func GetTag(query interface{}, args ...interface{}) Tag {
	info := Tag{}
	models.DB.Where(query, args...).First(&info)
	return info
}

func GetTags(query interface{}, args ...interface{}) []*Tag {
	var list []*Tag
	models.DB.Where(query, args...).Find(&list)
	return list
}
