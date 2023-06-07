/*
需要实现的 grpc 中的方法
type UserServer interface {
	GetUserList(context.Context, *PageInfo) (*UserListResponse, error)
	GetUserByMobile(context.Context, *MobileRequest) (*UserInfoResponse, error)
	GetUserById(context.Context, *IdRequest) (*UserInfoResponse, error)
	CreateUser(context.Context, *CreateUserInfo) (*UserInfoResponse, error)
	UpdateUser(context.Context, *UpdateUserInfo) (*emptypb.Empty, error)
	CheckPassword(context.Context, *PasswordCheckInfo) (*CheckResponse, error)
	mustEmbedUnimplementedUserServer()
}
*/

package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"mx_shop/user_srv/global"
	"mx_shop/user_srv/model"
	"mx_shop/user_srv/proto"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type UserServer struct {
}

// 有问题？？？ 为什么会多这么一个方法
func (s *UserServer) mustEmbedUnimplementedUserServer() {
	//TODO implement me
	panic("implement me")
}

// 将从数据库中查询出来的数据转换为RPC需要返回的数据格式

func ModelToResponse(user model.User) proto.UserInfoResponse {
	//在grpc的message中字段有默认值，你不能随便赋值nil进去，容易出错
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		NickName: user.NickName,
		Role:     int32(user.Role),
	}
	// 因为birthday可能是没有的
	if user.Birthday != nil {
		//Unix() 函数用于将一个 time.Time 类型的时间转换为 Unix 时间戳
		userInfoRsp.Birthday = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

// 分页功能

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 10:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// 获取用户列表
// (s *UserServer) 表示这个函数是在 UserServer 结构体上定义的方法。

func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.UserListResponse{} // 赋值地址
	fmt.Println(rsp)
	rsp.Total = int32(result.RowsAffected) // 挂载total属性

	// 官方翻页案例
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	for _, user := range users {
		userInfoRsp := ModelToResponse(user)
		// 将RPC中定义好的data数据格式挂载
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}
	return rsp, nil
}

// 通过号码查询用户

func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	//	没有查询到数据
	if result.RowsAffected == 0 {
		// status  是由grpc库包提供
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	// 查询数据库出现错误
	if result.Error != nil {
		return nil, result.Error
	}
	// 将数据转换成对应的格式
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

// GetUserById 通过ID查询用户
func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	//result := global.DB.Where(&model.User{ID: req.Id}).First(&user)
	result := global.DB.First(&user, req.Id) // 主键便捷查询
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

// CreateUser 创建用户
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	// 先查询用户是否已注册
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	// 用户已存在的逻辑
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	user.Mobile = req.Mobile
	user.NickName = req.NickName

	// 自定义 密码 加盐md5加密
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(req.Password, options)
	// 加盐操作
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	//fmt.Println(user.Password, "-------------------------")
	result = global.DB.Create(&user)
	if result.Error != nil {
		// codes.Internal 表示内部错误
		return nil, status.Errorf(codes.Internal, result.Error.Error()) // Error() 自动转换错误类型
	}
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

// UpdateUser 更新用户
func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	// 先查询用户是否存在
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	// birthday  需要将 uint 类型转换为 time 类型
	birthDay := time.Unix(int64(req.Birthday), 0)
	user.NickName = req.NickName
	user.Birthday = &birthDay
	user.Gender = req.Gender
	result = global.DB.Save(user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil
}

// CheckPassword 校验密码
func (s *UserServer) CheckPassword(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	options := &password.Options{16, 100, 32, sha512.New}
	// 对加密后的密码进行解析  ？？？  为什么这里不是password
	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	check := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], options)
	return &proto.CheckResponse{
		Success: check,
	}, nil
}
