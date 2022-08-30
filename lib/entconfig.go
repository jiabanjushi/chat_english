package lib

import (
	"go-fly-muti/models"
)

type EntConfig struct {
	LandHost string
}

//配置项
func NewEntConfig(entId string) (*EntConfig, error) {
	configs := models.FindEntConfigs(entId)
	landHost := ""
	for _, config := range configs {

		if config.ConfKey == "LandHost" {
			landHost = config.ConfValue
		}
	}

	entConfig := &EntConfig{
		LandHost: landHost,
	}
	return entConfig, nil
}
