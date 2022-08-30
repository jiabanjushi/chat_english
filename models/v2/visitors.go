package v2

import (
	"github.com/jinzhu/gorm"
	"go-fly-muti/models"
	"time"
)

type Visitor struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Avator    string    `json:"avator"`
	SourceIp  string    `json:"source_ip"`
	ToId      string    `json:"to_id"`
	VisitorId string    `json:"visitor_id"`
	Status    uint      `json:"status"`
	Refer     string    `json:"refer"`
	City      string    `json:"city"`
	ClientIp  string    `json:"client_ip"`
	Extra     string    `json:"extra"`
	EntId     string    `json:"ent_id"`
}

func (this *Visitor) InsertVisitor() *Visitor {
	models.DB.Create(this)
	return this
}
func (this *Visitor) FindVisitor() Visitor {
	var info Visitor
	this.buildQuery().First(&info)
	return info
}
func (this *Visitor) UpdateVisitor(query interface{}, args ...interface{}) {
	models.DB.Model(this).Where(query, args...).Update(this)
}

//查询构造
func (this *Visitor) buildQuery() *gorm.DB {
	db := models.DB
	db.Model(this)
	if this.ID != 0 {
		db = db.Where("id = ?", this.ID)
	}
	if this.VisitorId != "" {
		db = db.Where("visitor_id = ?", this.VisitorId)
	}
	if this.EntId != "" {
		db = db.Where("ent_id = ?", this.EntId)
	}
	return db
}
