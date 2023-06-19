package initialize

import (
	"fmt"
	"mx_shop/user_srv/global"

	"github.com/spf13/viper"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv() // 读取环境变量
	return viper.GetBool(env)
}

func InitConfig() {
	// 从配置文件中读取配置
	debug := GetEnvInfo("MXSHOP_DEV") // 获取系统环境变量
	configFilePrefix := "user_srv/config"
	configFileName := fmt.Sprintf("%s_prod.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("%s_dev.yaml", configFilePrefix)
	}
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
}
