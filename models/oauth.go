package models

import (
	"go-fly-muti/types"
	"time"
)

type Oauth struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	UserId    string     `json:"user_id"`
	OauthId   string     `json:"oauth_id"`
	status    uint       `json:"status"`
	CreatedAt types.Time `json:"created_at"`
}

func CreateOauth(userId, oauthId string) Oauth {

	model := Oauth{
		OauthId: oauthId,
		UserId:  userId,
		CreatedAt: types.Time{
			time.Now(),
		},
	}
	oauth := FindOauth(userId, oauthId)
	if oauth.ID != 0 {
		UpdateOauthById(model, userId)
		return model
	}
	DB.Create(&model)
	return model
}
func FindOauthById(userId string) Oauth {
	var oauth Oauth
	DB.Where("user_id = ?", userId).First(&oauth)
	return oauth
}
func FindOauth(userId, oauthId string) Oauth {
	var oauth Oauth
	DB.Where("user_id = ? and oauth_id = ?", userId, oauthId).First(&oauth)
	return oauth
}
func DelOauth(oauthId string) {
	DB.Where("oauth_id = ?", oauthId).Delete(&Oauth{})
}
func FindOauthsById(userId string) []Oauth {
	var oauths []Oauth
	DB.Where("user_id = ?", userId).Find(&oauths)
	return oauths
}
func FindOauthsQuery(query interface{}, args []interface{}) Oauth {
	var oauths Oauth
	DB.Where(query, args...).Find(&oauths)
	return oauths
}
func UpdateOauthById(oauth Oauth, userId string) Oauth {
	DB.Where("user_id = ?", userId).Update(&oauth)
	return oauth
}
func FindOauthsByVisitors(visitorIds []string) map[string]uint {
	var oauths []Oauth
	DB.Where("user_id in (? ) ", visitorIds).Find(&oauths)
	res := make(map[string]uint)
	for _, oauth := range oauths {
		res[oauth.UserId] = 2
	}

	for _, visitorId := range visitorIds {
		if res[visitorId] != 2 {
			res[visitorId] = 1
		}
	}
	return res
}
