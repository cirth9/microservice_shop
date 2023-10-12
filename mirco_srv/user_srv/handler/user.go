package handler

import (
	"MircoServer/user_srv/global"
	"MircoServer/user_srv/model"
	"MircoServer/user_srv/proto"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

var Options = &password.Options{SaltLen: 16, Iterations: 10000, KeyLen: 32, HashFunction: sha256.New}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

type UserServer struct {
	*proto.UnimplementedUserServer
}

func ModelToResponse(user model.User) proto.UserInfoResponse {
	userinfoResp := proto.UserInfoResponse{
		Id:                user.ID,
		Password:          user.Password,
		MobilePhoneNumber: user.MobilePhoneNumber,
		Nickname:          user.NickName,
		Gender:            user.Gender,
		Role:              user.Role,
	}
	if user.BirthDay != nil {
		userinfoResp.Birthday = uint64(user.BirthDay.Unix())
	}
	return userinfoResp
}

func (userServer *UserServer) GetUserList(ctx context.Context, in *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []*model.User
	result := global.MysqlDB.Find(&users)
	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}
	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	global.MysqlDB.Scopes(Paginate(int(in.Page), int(in.PSize))).Find(&users)
	for _, v := range users {
		userinfoResp := ModelToResponse(*v)
		rsp.Data = append(rsp.Data, &userinfoResp)
	}
	zap.S().Infof(">>>>>>[GetUserList] data %+v", rsp.Data)
	return rsp, nil
}

func (userServer *UserServer) GetUserByPhoneNumber(ctx context.Context, in *proto.MobilePhoneRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.MysqlDB.Where("mobile_phone_number = ?", in.Mobile).First(&user)
	log.Printf("[GetUserByPhoneNumber] result = %+v", result)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	rsp := ModelToResponse(user)
	return &rsp, nil
}

func (userServer *UserServer) GetUserByUid(ctx context.Context, in *proto.UserIdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.MysqlDB.Where("id = ?", in.Id).First(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	rsp := ModelToResponse(user)
	return &rsp, nil
}

func (userServer *UserServer) CreateUser(ctx context.Context, in *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.MysqlDB.Where("mobile_phone_number = ?", in.MobilePhone).First(&user)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "该用户已存在！")
	}
	user.MobilePhoneNumber = in.MobilePhone
	user.NickName = in.Nickname

	// Using custom options /   encode
	salt, encodedPwd := password.Encode(in.Password, Options)
	user.Password = fmt.Sprintf("pbkdf2-sha256$%s$%s", salt, encodedPwd)
	user.CreatedTime = time.Now()
	user.UpdatedTime = time.Now()
	user.IsDeleted = false

	result1 := global.MysqlDB.Create(&user)
	if result1.Error != nil {
		return nil, status.Errorf(codes.Internal, result1.Error.Error())
	}

	userInfoResp := ModelToResponse(user)
	return &userInfoResp, nil
}

func (userServer *UserServer) UpdateUser(ctx context.Context, in *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	var user model.User
	result := global.MysqlDB.Where("id = ?", in.Id).First(&user)
	if result.Error != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return &emptypb.Empty{}, status.Errorf(codes.NotFound, "用户不存在")
	}

	updateBirth := time.Unix(int64(in.Birthday), 0)
	user.NickName = in.Nickname
	user.BirthDay = &updateBirth
	user.Gender = in.Gender

	result1 := global.MysqlDB.Save(&user)
	if result1.Error != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, result.Error.Error())
	}

	return &emptypb.Empty{}, nil
}

func (userServer *UserServer) CheckPassword(ctx context.Context, in *proto.PasswordInfo) (*proto.CheckPasswordResponse, error) {
	passwordInfo := strings.Split(in.EncryptPassword, "$")
	for _, v := range passwordInfo {
		log.Println(v)
	}
	verify := password.Verify(in.Password, passwordInfo[1], passwordInfo[2], Options)
	return &proto.CheckPasswordResponse{
		IsOK: verify,
	}, nil
}
