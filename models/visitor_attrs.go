package models

import "time"

type Visitor_attr struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	VisitorId string    `json:"visitor_id"`
	EntId     string    `json:"ent_id"`
	RealName  string    `json:"real_name"`
	Tel       string    `json:"tel"`
	Email     string    `json:"email"`
	QQ        string    `json:"qq"`
	Wechat    string    `json:"wechat"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateVisitorAttr(ent_id, visitorId, realName, tel, email, qq, wechat, comment string) *Visitor_attr {
	v := &Visitor_attr{
		EntId:     ent_id,
		VisitorId: visitorId,
		RealName:  realName,
		Tel:       tel,
		Email:     email,
		QQ:        qq,
		Wechat:    wechat,
		Comment:   comment,
		CreatedAt: time.Now(),
	}
	DB.Create(v)
	return v
}
func GetVisitorAttrByVisitorId(visitorId, entId string) Visitor_attr {
	var v Visitor_attr
	if visitorId == "" {
		return v
	}
	DB.Where("visitor_id = ? and ent_id = ?", visitorId, entId).First(&v)
	return v
}
func SaveVisitorAttrByVisitorId(visitorAttr *Visitor_attr, visitorId, entId string) {
	DB.Model(visitorAttr).Where("visitor_id = ? and ent_id = ?", visitorId, entId).Update(visitorAttr)
}
func DelVisitorAttr(query interface{}, args ...interface{}) {
	DB.Where(query, args...).Delete(&Visitor_attr{})
}