package main

import (
	"flag"
	"fmt"
	"mx_shop/user_srv/global"
	handle "mx_shop/user_srv/handler"
	"mx_shop/user_srv/initialize"
	"mx_shop/user_srv/proto"
	"net"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50051, "端口号")
	flag.Parse() // 解析用户输入

	initialize.InitLogger() // 初始化日志
	initialize.InitConfig() // 初始化端口等配置
	initialize.InitDB()     // 初始化数据库

	// 这里为什么是指针 ？？？
	zap.S().Info("ip---", *IP)
	zap.S().Info("port--", *Port)

	// 开启 grpc 服务
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handle.UserServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("监听端口失败" + err.Error())
	}

	// 注册服务健康检测
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 服务注册
	cfg := api.DefaultConfig() // 获取默认配置
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulConfig.Host,
		global.ServerConfig.ConsulConfig.Port)

	client, err := api.NewClient(cfg) // 创建客户端
	if err != nil {
		panic(err)
	}
	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      global.ServerConfig.Name,
		Name:    global.ServerConfig.Name,
		Address: "172.18.81.229", // consul 地址
		Port:    *Port,           // 监听的服务的端口
		Tags:    []string{"coder", "imooc"},
		Check: &api.AgentServiceCheck{ // 健康检查
			GRPC:                           "172.25.16.1:50051", // 监听的GPRC服务
			Timeout:                        "3s",                // 超时时间
			Interval:                       "5s",                // 健康检查间隔
			DeregisterCriticalServiceAfter: "30s",               // 注销时间，相当于过期时间
		},
	})
	if err != nil {
		panic(err)
	}

	err = server.Serve(lis)
	if err != nil {
		panic("启动GRPC失败" + err.Error())
	}
}
