package initialize

import (
	"encoding/json"
	"fmt"
	"mx_shop/goods_srv/global"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv() // 读取环境变量
	return viper.GetBool(env)
}

func InitConfig() {
	// 因为 配置文件已经交给nacos管理了，所以只需要从nacos中获取配置即可
	debug := GetEnvInfo("MXSHOP_DEV")
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
	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息:%v", global.NacosConfig)

	// ==========================> nacos配置中心 <==============================
	// 创建clientConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}
	// 创建动态配置客户端
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	// 创建动态配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}
	// 获取配置
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatal("读取配置失败:%s", err.Error())
	}
	// fmt.Println(serverConfig)

	// 监听配置文件的变化
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: "user-web.yaml",
		Group:  "dev",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("nacos配置文件发生了变化")
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})
}

func InitConfig2() {
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
