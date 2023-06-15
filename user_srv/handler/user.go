package handle

import (
	"context"
	"crypto/sha512"
	"fmt"
	"mx_shop/user_srv/global"
	"mx_shop/user_srv/model"
	"mx_shop/user_srv/proto"
	"strings"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/jinzhu/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserServer struct {
}

func Paginate(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}

func ModelToRsp(u model.User) *proto.UserInfoResponse {
	return &proto.UserInfoResponse{
		Id:       u.ID,
		Mobile:   u.Mobile,
		NickName: u.NickName,
		Password: u.Password,
		Gender:   u.Gender,
		Role:     int32(u.Role),
	}
}

func (u *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	userDb := global.DB.Where(&users)

	// if result.Error != nil {
	// 	return nil, result.Error
	// }

	var rsp = &proto.UserListResponse{}
	// rsp.Total = int32(result.RowsAffected)
	var count int64
	userDb.Count(&count)
	rsp.Total = int32(count)
	global.DB.Scopes(Paginate(int(req.Page), int(req.Size))).Find(&users)

	for _, v := range users {
		data := ModelToRsp(v)
		rsp.Data = append(rsp.Data, data)
	}
	return rsp, nil
}

// 获取用户信息
func (u *UserServer) GetUserInfo(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	// fmt.Println(req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	var userInfo proto.UserInfoResponse

	copier.Copy(&userInfo, &user)
	return &userInfo, nil
}

func (u *UserServer) GetUserMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).Find(&user)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	var userInfo proto.UserInfoResponse
	copier.Copy(&userInfo, &user)
	return &userInfo, nil
}

func (u *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserReq) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).Find(&user)
	if result.RowsAffected == 1 {
		return nil, status.Error(codes.AlreadyExists, "手机号码已绑定其他用户")
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
	var userInfo proto.UserInfoResponse
	copier.Copy(&userInfo, &user)
	return &userInfo, nil
}

func (u *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserReq) (*proto.Response, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)

	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	copier.Copy(&user, &req)

	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}

	return &proto.Response{Code: 0, Msg: "修改成功"}, nil
}

// func (u *UserServer) CheckPassword(ctx context.Context, req *proto.CheckPasswordReq) (*proto.CheckPasswordResponse, error) {
// 	return nil, status.Errorf(codes.Unimplemented, "method CheckPassword not implemented")
// }

// CheckPassword 校验密码
func (s *UserServer) CheckPassword(ctx context.Context, req *proto.CheckPasswordReq) (*proto.CheckPasswordResponse, error) {
	options := &password.Options{16, 100, 32, sha512.New}
	// 对加密后的密码进行解析  ？？？  为什么这里不是password
	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	check := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], options)
	return &proto.CheckPasswordResponse{
		IsValid: check,
	}, nil
}
