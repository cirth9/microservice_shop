package router

import (
	"Shop-api/user_web/api"
	"Shop-api/user_web/middleware"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.RouterGroup) {
	userGroup := r.Group("user")
	{
		//middleware.AdminAuth(),
		userGroup.GET("/list", middleware.JWTAuth(), middleware.AdminAuth(), api.GetUserList)
		userGroup.POST("/login", api.PasswordLogin)
		userGroup.POST("/register", api.UserRegister)
	}
}
