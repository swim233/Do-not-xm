package config

import (
	"encoding/json"
	"os"

	"github.com/spf13/viper"
	"github.com/swim233/do_not_xm/logger"
)

var Config *viper.Viper

// 初始化配置文件
func InitConfig() {
	defer logger.Logger.Sync()
	Config = viper.New()
	Config.SetConfigType("json")
	var GlobalFile *os.File
	if file, err := os.Open("../config/config.json"); err != nil {
		logger.Suger.Warnf("Fail to open config file : %s , trying to create...", err.Error())
		if newConfigFile, err := os.Create("../config/config.json"); err != nil {
			logger.Suger.Panicf("Fail to create config file reason is : %s /n Please check ./config/config.json manually", err.Error())
		} else {
			if _, err := newConfigFile.Write(createNewConfigFile()); err != nil {
				logger.Suger.Panicf("Fail to create config file reason is : %s /n Please check ./config/config.json manually", err.Error())
			}
			logger.Suger.Panicf("Successful create config file,please entry it and reboot program")
		}

	} else {
		if file != nil {
			GlobalFile = file
		}
		//读取配置文件
		if err := Config.ReadConfig(GlobalFile); err != nil {
			logger.Suger.Errorf("Fail to read config file : %s", err.Error())
		}
		defer GlobalFile.Close()
	}
}

// 创建默认配置内容
func createNewConfigFile() []byte {
	config := struct {
		BotToken  string `json:"bot_token"`
		DebugMode bool   `json:"debug_mode"`
	}{
		BotToken:  "",
		DebugMode: false,
	}
	if configJson, err := json.Marshal(config); err != nil {
		logger.Suger.Panicf("Fail to create default config content : %s", err.Error())
		return nil
	} else {
		return configJson
	}
}
