package api

import (
	"Shop-api/user_web/forms"
	"Shop-api/user_web/global"
	"Shop-api/user_web/middleware"
	"Shop-api/user_web/proto"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"time"
)

func GrpcErrorToHTTP(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
		}
	}
}

func GetUserList(c *gin.Context) {

	zap.S().Debug("grpcAddress", global.GrpcAddress)
	UserListResp, err := global.UserClient.GetUserList(context.Background(), &proto.PageInfo{
		Page:  1,
		PSize: 2,
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接成功但是调用失败",
			"error", err.Error())
		return
	}

	zap.S().Debug("Get UserList")

	userList := make([]interface{}, 0)
	for _, val := range UserListResp.Data {
		userInfo := make(map[string]interface{}, 0)
		userInfo["birthday"] = val.Birthday
		userInfo["nickname"] = val.Nickname
		userInfo["gender"] = val.Gender
		userInfo["role"] = val.Role
		userInfo["mobilePhoneNumber"] = val.MobilePhoneNumber

		userList = append(userList, userInfo)
	}
	c.JSON(http.StatusOK, userList)
}

func PasswordLogin(c *gin.Context) {
	passwordLoginForm := forms.PasswordLoginForm{}
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		zap.S().Errorw("获取数据失败", "error", err.Error())
		return
	}
	zap.S().Debug("passwordLoginForm:", passwordLoginForm)

	//验证验证码
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.CaptchaNumber, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	}

	if UserInfo, err1 := global.UserClient.GetUserByPhoneNumber(context.Background(), &proto.MobilePhoneRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err1 != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "用户不存在",
		})
		return

	} else {

		if checkPassword, err2 := global.UserClient.CheckPassword(context.Background(), &proto.PasswordInfo{
			Password:        passwordLoginForm.Password,
			EncryptPassword: UserInfo.Password,
		}); err2 != nil {

			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "登录失败",
			})
			return

		} else {

			if checkPassword.IsOK {
				token, err3 := middleware.GetToken(UserInfo.Nickname, uint(UserInfo.Id), uint(UserInfo.Role))
				if err3 != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "登陆失败，token生成失败",
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"id":         UserInfo.Id,
					"nick_name":  UserInfo.Nickname,
					"token":      token,
					"expired_at": time.Now().Add(global.TokenExpireDuration).Unix(),
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "登陆失败",
				})
			}

		}

	}
}

func UserRegister(c *gin.Context) {
	userRegister := forms.UserRegister{}
	if err := c.ShouldBindJSON(&userRegister); err != nil {
		zap.S().Errorw("[UserRegister] 表单读取错误 Error ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	result, err := global.Rdb.Get(context.Background(), userRegister.Mobile).Result()
	code, _ := strconv.ParseInt(result, 10, 64)
	if err != nil {
		zap.S().Errorw("[UserRegister] Redis Get data failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	} else {
		if code != int64(userRegister.CheckCode) {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "验证码错误",
			})
			return
		}
	}

	userResp, err := global.UserClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Nickname:    userRegister.Mobile,
		Password:    userRegister.Password,
		MobilePhone: userRegister.Mobile,
	})
	if err != nil {
		zap.S().Errorw("[UserRegister] CreateUser failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	token, err3 := middleware.GetToken(userResp.Nickname, uint(userResp.Id), uint(userResp.Role))
	if err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "登陆失败，token生成失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":                  "添加成功",
		"token":                token,
		"expired_at":           time.Now().Add(global.TokenExpireDuration).Unix(),
		"create_user_id":       userResp.Id,
		"create_user_nickname": userResp.Nickname,
	})
}
