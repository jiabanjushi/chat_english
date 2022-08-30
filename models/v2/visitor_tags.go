package v2

import (
	"go-fly-muti/models"
	"time"
)

type VisitorTag struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	VisitorId string    `json:"visitor_id"`
	TagId     uint      `json:"tag_id"`
	Kefu      string    `json:"kefu"`
	EntId     uint      `json:"ent_id"`
}
type VisitorTagList struct {
	VisitorId string `json:"visitor_id"`
	TagId     uint   `json:"tag_id"`
	TagName   string `json:"tag_name"`
}

func (this *VisitorTag) InsertVisitorTag() *VisitorTag {
	models.DB.Create(this)
	return this
}
func GetVisitorTags(query interface{}, args ...interface{}) []VisitorTag {
	list := []VisitorTag{}
	models.DB.Where(query, args...).Find(&list)
	return list
}
func GetVisitorTagsByVisitorId(visitorId string, entId uint) []VisitorTagList {
	var list []VisitorTagList
	models.DB.Table("visitor_tag").Select("visitor_tag.*,tag.name as tag_name").Joins("left join tag on visitor_tag.tag_id=tag.id").
		Where("visitor_tag.visitor_id = ? and visitor_tag.ent_id = ?", visitorId, entId).Find(&list)
	return list
}
func DelVisitorTags(query interface{}, args ...interface{}) {
	models.DB.Where(query, args...).Delete(&VisitorTag{})
}
