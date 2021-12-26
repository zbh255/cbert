package ioc

import (
	"github.com/zbh255/cbert/common/config"
	"io/ioutil"
)

var (
	configObj map[string]interface{}
)

func GetProjectConfig() *config.AutoGenerated {
	return configObj["project"].(*config.AutoGenerated)
}

func GetUserConfig() *config.UserConf {
	return configObj["user"].(*config.UserConf)
}

func RegisterProjectConfig(filePath string) error {
	confBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	api := config.NewProjectConfig()
	api, err = api.DecodeConfig(confBytes)
	if err != nil {
		return err
	}
	configObj["project"] = api
	return nil
}

func RegisterUserConfig(filePath string) error {
	confBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	api := config.NewUserConfig()
	api, err = api.DecodeConfig(confBytes)
	if err != nil {
		return err
	}
	configObj["project"] = api
	return nil
}

func init() {
	configObj = make(map[string]interface{}, 2)
}
