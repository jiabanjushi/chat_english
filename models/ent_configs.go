package models

import "fmt"

type EntConfig struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	ConfName  string `json:"conf_name"`
	ConfKey   string `json:"conf_key"`
	ConfValue string `json:"conf_value"`
	EntId     string `json:"ent_id"`
}

func CreateEntConfig(kefuId interface{}, name, key, value string) {
	c := &EntConfig{
		ConfName:  name,
		ConfKey:   key,
		ConfValue: value,
		EntId:     fmt.Sprintf("%v", kefuId),
	}
	DB.Create(c)
}
func FindEntConfigs(kefuId interface{}) []EntConfig {
	var configs []EntConfig
	DB.Where("ent_id = ?", kefuId).Find(&configs)
	return configs
}
func FindEntConfigByEntid(entId interface{}) []EntConfig {
	var configs []EntConfig
	DB.Where("ent_id = ?", entId).Find(&configs)
	return configs
}
func FindEntConfig(kefuId interface{}, key string) EntConfig {
	var config EntConfig
	DB.Where("ent_id = ? and conf_key=?", kefuId, key).Find(&config)
	return config
}
func UpdateEntConfig(kefuId interface{}, name, key, value string) {
	c := map[string]string{
		"conf_name":  name,
		"conf_key":   key,
		"conf_value": value,
		"ent_id":     fmt.Sprintf("%v", kefuId),
	}
	DB.Model(&EntConfig{}).Where("ent_id = ? and conf_key = ?", fmt.Sprintf("%v", kefuId), key).Update(c)
}
