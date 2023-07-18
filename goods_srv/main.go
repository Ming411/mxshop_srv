package main

import (
	"flag"
	"fmt"
	"mx_shop/goods_srv/global"
	handle "mx_shop/goods_srv/handler"
	"mx_shop/goods_srv/initialize"
	"mx_shop/goods_srv/proto"
	"mx_shop/goods_srv/utils"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// IP := flag.String("ip", "0.0.0.0", "ip地址")
	IP := flag.String("ip", "172.29.32.1", "ip地址")
	// Port := flag.Int("port", 50051, "端口号") // 本地测试使用固定一个端口
	Port := flag.Int("port", 0, "端口号")
	flag.Parse() // 解析用户输入

	initialize.InitLogger() // 初始化日志
	initialize.InitConfig() // 初始化端口等配置
	initialize.InitDB()     // 初始化数据库

	zap.S().Info("ip---", *IP)

	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	zap.S().Info("port--", *Port)

	// 开启 grpc 服务
	server := grpc.NewServer()
	proto.RegisterGoodsServer(server, &handle.GoodsServer{})

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
	serviceID := uuid.NewV4().String() // 通过uuid生成唯一名称的服务，这样就可以同时启动多个服务
	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		// ID:      global.ServerConfig.Name,
		ID:   serviceID,
		Name: global.ServerConfig.Name,
		// Address: "172.29.32.1",
		Address: global.NacosConfig.Host,
		Port:    *Port, // 监听的服务的端口
		Tags:    global.ServerConfig.Tags,
		Check: &api.AgentServiceCheck{ // 健康检查
			GRPC:                           fmt.Sprintf("%s:%d", "172.29.32.1", *Port), // 监听的GPRC服务
			Timeout:                        "3s",                                       // 超时时间
			Interval:                       "5s",                                       // 健康检查间隔
			DeregisterCriticalServiceAfter: "30s",                                      // 注销时间，相当于过期时间
		},
	})
	if err != nil {
		panic(err)
	}

	go func() {
		err = server.Serve(lis) // 会造成阻塞,所以要以goruntine方式启动
		if err != nil {
			panic("启动GRPC失败" + err.Error())
		}
	}()

	// 接收ctrl+c终止信号，用于清除之前注册的服务
	quit := make(chan os.Signal) // chan 用于接收信号 os 用于获取系统信号
	// syscall.SIGINT 表示中断信号（通常由 Ctrl+C 发出），syscall.SIGTERM 表示终止信号。
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 接收信号
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("服务注销失败")
	}
	zap.S().Info("服务注销成功")
}
