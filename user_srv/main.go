package main

import (
	"flag"
	"fmt"
	handle "mx_shop/user_srv/handler"
	"mx_shop/user_srv/proto"
	"net"

	"google.golang.org/grpc"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50051, "端口号")
	flag.Parse() // 解析用户输入

	// 这里为什么是指针 ？？？
	fmt.Println("ip---", *IP)
	fmt.Println("port--", *Port)

	// 开启 grpc 服务
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handle.UserServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("监听端口失败" + err.Error())
	}
	err = server.Serve(lis)
	if err != nil {
		panic("启动GRPC失败" + err.Error())
	}
}
