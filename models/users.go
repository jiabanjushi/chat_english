package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-fly-muti/common"
	"go-fly-muti/types"
	"strconv"
	"time"
)

type User struct {
	ID           uint       `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ExpiredAt    types.Time `json:"expired_at"`
	Name         string     `json:"name"`
	Password     string     `json:"password"`
	Nickname     string     `json:"nickname"`
	Avator       string     `json:"avator"`
	Pid          uint       `json:"pid"`
	RecNum       uint       `json:"rec_num"`
	AgentNum     uint       `json:"agent_num"`
	Status       uint       `json:"status"`
	OnlineStatus uint       `json:"online_status"`
	EntId        int        `json:"ent_id" sql:"-"`
	RoleName     string     `json:"role_name" sql:"-"`
	RoleId       string     `json:"role_id" sql:"-"`
	orderBy      string
}

func CreateUser(name string, password string, avator string, nickname string, pid, agentNum uint) uint {
	user := &User{
		Name:         name,
		Password:     password,
		Avator:       avator,
		Nickname:     nickname,
		Pid:          pid,
		RecNum:       0,
		Status:       1,
		OnlineStatus: 1,
		AgentNum:     agentNum,
	}
	user.UpdatedAt = time.Now()
	DB.Create(user)
	return user.ID
}
func UpdateUser(id string, name string, password string, avator string, nickname string, agentNum uint) {
	user := &User{
		Name:     name,
		Avator:   avator,
		Nickname: nickname,
		AgentNum: agentNum,
	}
	user.UpdatedAt = time.Now()
	if password != "" {
		user.Password = password
	}
	DB.Model(&User{}).Where("id = ?", id).Update(user)
}
func UpdateUserRecNum(name string, num interface{}) {
	user := &User{}
	DB.Model(user).Where("name = ?", name).Update("RecNum", gorm.Expr("rec_num + ?", num))
}
func UpdateUserRecNumZero(name string) {
	values := map[string]uint{
		"rec_num": 0,
	}
	DB.Model(&User{}).Where("name = ?", name).Update(values)
}
func UpdateUserStatusWhere(status uint, query interface{}, args ...interface{}) {
	values := map[string]uint{
		"Status": status,
	}
	DB.Model(&User{}).Where(query, args...).Update(values)
}
func UpdateUserPass(name string, pass string) {
	user := &User{
		Password: pass,
	}
	user.UpdatedAt = time.Now()
	DB.Model(user).Where("name = ?", name).Update("Password", pass)
}
func UpdateUserAvator(name string, avator string) {
	user := &User{
		Avator: avator,
	}
	user.UpdatedAt = time.Now()
	DB.Model(user).Where("name = ?", name).Update("Avator", avator)
}
func UpdateKefuInfoByName(name, avator, nickname string) {
	user := &User{
		Avator:   avator,
		Nickname: nickname,
	}
	user.UpdatedAt = time.Now()
	DB.Model(user).Where("name = ?", name).Update(user)
}
func FindUser(username string) User {
	var user User
	DB.Where("name = ?", username).First(&user)
	return user
}
func FindUserByUid(id interface{}) User {
	var user User
	DB.Where("id = ?", id).First(&user)
	return user
}
func FindUsersStatus(status interface{}) []User {
	var users []User
	DB.Where("status=?", status).Order("rec_num asc").Find(&users)
	return users
}
func FindUserById(id interface{}) User {
	var user User
	DB.Select("user.*,role.name role_name,role.id role_id").Joins("join user_role on user.id=user_role.user_id").Joins("join role on user_role.role_id=role.id").Where("user.id = ?", id).First(&user)
	return user
}
func FindUserByIdPid(pid interface{}, id interface{}) User {
	var user User
	DB.Select("user.*,role.name role_name,role.id role_id").Joins("join user_role on user.id=user_role.user_id").Joins("join role on user_role.role_id=role.id").Where("user.id = ? and user.pid = ?", id, pid).First(&user)
	return user
}
func DeleteUserById(id string) {
	DB.Where("id = ?", id).Delete(User{})
}
func DeleteUserByIdPid(id string, pid interface{}) {
	DB.Where("id = ? and pid=?", id, pid).Delete(User{})
}
func FindUsers() []User {
	var users []User
	DB.Select("user.*,role.name role_name").Joins("left join user_role on user.id=user_role.user_id").Joins("left join role on user_role.role_id=role.id").Order("user.id desc").Find(&users)
	return users
}
func FindUsersByEntId(entId interface{}) []User {
	var users []User
	users = FindUsersByPid(entId)
	return users
}

func FindUsersByPid(pid interface{}) []User {
	var users []User
	DB.Where("pid=? or id=?", pid, pid).Order("rec_num asc").Find(&users)
	return users
}
func FindUsersWhere(query interface{}, args ...interface{}) []User {
	var users []User
	DB.Where(query, args...).Order("rec_num asc").Find(&users)
	return users
}
func FindUserRole(query interface{}, id interface{}) User {
	var user User
	DB.Select(query).Where("user.id = ?", id).Joins("join user_role on user.id=user_role.user_id").Joins("join role on user_role.role_id=role.id").First(&user)
	return user
}

//查询条数
func CountUsers() uint {
	var count uint
	DB.Model(&User{}).Count(&count)
	return count
}

//根据where查询条数
func CountUsersWhere(query interface{}, args ...interface{}) uint {
	var count uint
	DB.Model(&User{}).Where(query, args...).Count(&count)
	return count
}

//根据where分页查询
func FindUsersOwn(page uint, pagesize uint, query interface{}, args ...interface{}) []User {
	if pagesize == 0 {
		pagesize = common.PageSize
	}
	offset := (page - 1) * pagesize
	if offset < 0 {
		offset = 0
	}
	var users []User
	DB.Select("user.*,role.name role_name,role.id role_id").Joins("left join user_role on user.id=user_role.user_id").Joins("left join role on user_role.role_id=role.id").Where(query, args...).Offset(offset).Order("user.updated_at desc").Limit(pagesize).Find(&users)
	return users
}
func FindUsersPages(page uint, pagesize uint) []User {
	offset := (page - 1) * pagesize
	if offset < 0 {
		offset = 0
	}
	var users []User
	DB.Select("user.*,role.name role_name").Joins("left join user_role on user.id=user_role.user_id").Joins("left join role on user_role.role_id=role.id").Offset(offset).Order("user.id desc").Limit(pagesize).Find(&users)
	return users
}

//获取一条用户信息
func (user *User) AddUser() uint {
	DB.Create(user)
	return user.ID
}

//获取一条用户信息
func (user *User) GetOneUser(fields string) User {
	var dUser User
	var userRole User_role
	myDB := user.buildQuery()
	myDB.Select(fields).First(&dUser)
	if dUser.ID != 0 {
		userRole = FindRoleByUserId(dUser.ID)
		dUser.RoleId = strconv.Itoa(int(userRole.RoleId))
	}
	if userRole.RoleId != 0 {
		role := FindRole(userRole.RoleId)
		dUser.RoleName = role.Name
	}
	return dUser
}

//获取多条用户信息
func (user *User) GetUsers(fields string) []User {
	var users []User
	myDB := user.buildQuery()
	myDB.Select(fields).Find(&users)
	return users
}

//更新user
func (user *User) UpdateUser() {
	user.UpdatedAt = time.Now()
	user.buildQuery().Model(&User{}).Update(user)
}

//查询构造
func (user *User) buildQuery() *gorm.DB {
	userDB := DB
	userDB.Model(user)
	if user.ID != 0 {
		userDB = userDB.Where("id = ?", user.ID)
	}
	if user.Pid != 0 {
		userDB = userDB.Where("pid = ?", user.Pid)
	}
	if user.Name != "" {
		userDB = userDB.Where("name = ?", user.Name)
	}
	if user.orderBy != "" {
		userDB = userDB.Order(user.orderBy)
	}
	return userDB
}

//设置属性
func (this *User) SetOrder(orderBy string) {
	this.orderBy = orderBy
}

//获检查用户的状态
func (this *User) CheckStatusExpired() (bool, uint, string) {
	if this.ID == 0 || this.Status == 0 {
		return false, types.ApiCode.ACCOUNT_NO_EXIST, types.ApiCode.GetMessage(types.ApiCode.ACCOUNT_NO_EXIST)
	}
	if int64(this.Status) == types.Constant.AccountForbidden {
		return false, types.ApiCode.ACCOUNT_FORBIDDEN, types.ApiCode.GetMessage(types.ApiCode.ACCOUNT_FORBIDDEN)
	}
	//查看过期时间
	nowSecond := time.Now().Unix()
	expireSecond := this.ExpiredAt.Unix()
	if expireSecond != 0 && expireSecond < nowSecond {
		return false, types.ApiCode.ACCOUNT_EXPIRED, types.ApiCode.GetMessage(types.ApiCode.ACCOUNT_EXPIRED)
	}
	return true, 0, ""
}
