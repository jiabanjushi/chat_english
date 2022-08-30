package models

import (
	"time"
)

type Visitor struct {
	Model
	Name      string `json:"name"`
	RealName  string `json:"real_name"`
	Avator    string `json:"avator"`
	SourceIp  string `json:"source_ip"`
	ToId      string `json:"to_id"`
	VisitorId string `json:"visitor_id"`
	Status    uint   `json:"status"`
	Refer     string `json:"refer"`
	City      string `json:"city"`
	ClientIp  string `json:"client_ip"`
	Extra     string `json:"extra"`
	EntId     string `json:"ent_id"`
}
type VisitorExt struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	VisitorId string    `json:"visitor_id"`
	Url       string    `json:"url"`
	ServerIp  string    `json:"server_ip"`
	ClientIp  string    `json:"client_ip"`
	Title     string    `json:"title"`
	Ua        string    `json:"ua"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateVisitor(name, avator, sourceIp, toId, visitorId, refer, city, clientIp, entId, extra string) *Visitor {
	v := &Visitor{
		Name:      name,
		Avator:    avator,
		SourceIp:  sourceIp,
		ToId:      toId,
		VisitorId: visitorId,
		Status:    1,
		Refer:     refer,
		City:      city,
		ClientIp:  clientIp,
		Extra:     extra,
		EntId:     entId,
	}
	v.UpdatedAt = time.Now()
	DB.Create(v)
	return v
}
func AddVisitorExt(visitorId, serverIp, ua, url, title, clientIp string) {
	v := &VisitorExt{
		VisitorId: visitorId,
		ServerIp:  serverIp,
		Ua:        ua,
		Url:       url,
		CreatedAt: time.Now(),
		Title:     title,
		ClientIp:  clientIp,
	}
	DB.Model(&VisitorExt{}).Create(v)
}
func CountVisitorExtByVistorId(visitorId string) uint {
	var v uint
	DB.Model(&VisitorExt{}).Where("visitor_id = ?", visitorId).Count(&v)
	return v
}
func FindVisitorExtByVistorId(visitorId string, page uint, pagesize uint) []VisitorExt {
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pagesize
	if offset <= 0 {
		offset = 0
	}
	var v []VisitorExt
	DB.Model(&VisitorExt{}).Where("visitor_id = ?", visitorId).Order("id desc").Offset(offset).Limit(pagesize).Find(&v)
	return v
}
func FindVisitorByVistorId(visitorId string) Visitor {
	var v Visitor
	DB.Where("visitor_id = ?", visitorId).First(&v)
	return v
}
func FindVisitorsByWhere(page uint, pagesize uint, query interface{}, args ...interface{}) []Visitor {
	offset := (page - 1) * pagesize
	if offset < 0 {
		offset = 0
	}
	var visitors []Visitor
	DB.Where(query, args...).Offset(offset).Limit(pagesize).Order("status desc, updated_at desc").Find(&visitors)
	return visitors
}
func FindVisitorsByKefuId(page uint, pagesize uint, kefuId string) []Visitor {
	offset := (page - 1) * pagesize
	if offset < 0 {
		offset = 0
	}
	var visitors []Visitor
	DB.Where("to_id=?", kefuId).Offset(offset).Limit(pagesize).Order("status desc, updated_at desc").Find(&visitors)
	return visitors
}
func FindVisitorsByEntId(page uint, pagesize uint, entId string, orderBy string, search string, args ...interface{}) []Visitor {
	offset := (page - 1) * pagesize
	if offset < 0 {
		offset = 0
	}
	var visitors []Visitor
	query := "ent_id=? "
	if search != "" {
		query += search
	}
	entArgs := []interface{}{entId}
	entArgs = append(entArgs, args...)
	DB.Where(query, entArgs...).Offset(offset).Limit(pagesize).Order(orderBy).Find(&visitors)
	return visitors
}
func FindVisitorsOnline() []Visitor {
	var visitors []Visitor
	DB.Where("status = ?", 1).Find(&visitors)
	return visitors
}
func UpdateVisitorVisitorId(visitorId, newId string) int64 {
	visitor := Visitor{
		VisitorId: newId,
	}
	db := DB.Model(&visitor).Where("visitor_id = ?", visitorId).Update(visitor)
	return db.RowsAffected
}
func UpdateVisitorStatus(visitorId string, status uint) int64 {
	visitor := Visitor{
		Status: status,
	}
	db := DB.Model(&visitor).Where("visitor_id = ?", visitorId).Update(visitor)
	return db.RowsAffected
}
func UpdateVisitorRealName(name, EntId, visitorId string) {
	visitor := Visitor{
		RealName: name,
	}
	DB.Model(&visitor).Where("visitor_id = ? and ent_id = ?", visitorId, EntId).Update(visitor)
}
func UpdateVisitor(entId, name, avator, visitorId string, toId string, status uint, clientIp string, sourceIp string, refer, extra string) {
	visitor := &Visitor{
		ToId:     toId,
		Status:   status,
		ClientIp: clientIp,
		SourceIp: sourceIp,
		Refer:    refer,
		Extra:    extra,
		Name:     name,
		Avator:   avator,
		EntId:    entId,
	}
	visitor.UpdatedAt = time.Now()
	DB.Model(visitor).Where("visitor_id = ?", visitorId).Update(visitor)
}

func UpdateVisitorKefu(visitorId string, kefuId string) {
	visitor := Visitor{}
	DB.Model(&visitor).Where("visitor_id = ?", visitorId).Update("to_id", kefuId)
}

//根据条件查visitors
func FindVisitorsByCondition(query interface{}, args ...interface{}) []Visitor {
	var visitors []Visitor
	DB.Where(query, args...).Find(&visitors)
	return visitors
}

//查询条数
func CountVisitors(query interface{}, args ...interface{}) uint {
	var count uint
	DB.Model(&Visitor{}).Where(query, args...).Count(&count)
	return count
}

//查询每天条数
type EveryDayNum struct {
	Day string `json:"day"`
	Num int64  `json:"num"`
}

func CountVisitorsEveryDay(toId string) []EveryDayNum {
	var results []EveryDayNum
	DB.Raw("select DATE_FORMAT(created_at,'%Y-%m-%d') as day ,"+
		"count(*) as num from visitor where to_id=? group by day order by day desc limit 30",
		toId).Scan(&results)
	return results
}

//查询条数
func CountVisitorsByKefuId(kefuId string) uint {
	var count uint
	DB.Model(&Visitor{}).Where("to_id=?", kefuId).Count(&count)
	return count
}

//查询条数
func CountVisitorsByEntid(entId string, search string, args ...interface{}) uint {
	var count uint
	query := "ent_id=? "
	if search != "" {
		query += search
	}
	entArgs := []interface{}{entId}
	entArgs = append(entArgs, args...)
	DB.Model(&Visitor{}).Where(query, entArgs...).Count(&count)
	return count
}

//删除访客
func DelVisitor(query interface{}, args ...interface{}) {
	DB.Where(query, args...).Delete(&Visitor{})
}
