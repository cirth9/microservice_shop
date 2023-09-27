package router

import (
	"Shop-api/user_web/api"
	"github.com/gin-gonic/gin"
)

func InitBaseRouter(r *gin.RouterGroup) {
	routerGroup := r.Group("base")
	{
		routerGroup.GET("/captcha", api.GetCaptcha)
		routerGroup.POST("/sendSms", api.SendSms)
	}
}
