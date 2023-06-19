package initialize

import "go.uber.org/zap"

func InitLogger() {
	logger, _ := zap.NewDevelopment() // 初始化一个logger
	zap.ReplaceGlobals(logger)        // 替换全局的logger
}
