package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"mx_shop/user_srv/proto"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init() {
	// 建立客户端连接
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure()) // WithInsecure 不安全
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn) // 初始化客户端
	//fmt.Println("初始化客户端完成")
}

// TestGetUserList 测试获取用户列表
func TestGetUserList() {
	pageInfo := &proto.PageInfo{
		Pn:    1,
		PSize: 5,
	}
	rsp, err := userClient.GetUserList(context.Background(), pageInfo)
	if err != nil {
		panic(err)
	}
	for _, user := range rsp.Data {
		fmt.Println(user)
		checkRsp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          "admin123",
			EncryptedPassword: user.Password,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(checkRsp.Success)
	}
}

func TestCreateUser() {
	for i := 0; i < 10; i++ {
		rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			NickName: fmt.Sprintf("bobby%d", i),
			Mobile:   fmt.Sprintf("1878222222%d", i),
			Password: "admin123",
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(rsp.Id)
	}
}

func main() {
	Init()

	TestGetUserList()
	//TestCreateUser()

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn) // 最后关闭grpc
}
