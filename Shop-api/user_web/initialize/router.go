package initialize

import (
	"Shop-api/user_web/middleware"
	"Shop-api/user_web/router"
	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	userGroup := r.Group("/u/v1")
	router.InitUserRouter(userGroup)
	router.InitBaseRouter(userGroup)
	return r
}
