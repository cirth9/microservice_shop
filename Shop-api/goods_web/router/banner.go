package router

import (
	"Shop-api/goods_web/api/banner"
	"Shop-api/goods_web/middleware"
	"github.com/gin-gonic/gin"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	BannerRouter := Router.Group("banners")
	//.Use(middleware.Trace())
	{
		BannerRouter.GET("", banner.List)                                                        // 轮播图列表页
		BannerRouter.DELETE("/:id", middleware.JWTAuth(), middleware.AdminAuth(), banner.Delete) // 删除轮播图
		BannerRouter.POST("", middleware.JWTAuth(), middleware.AdminAuth(), banner.New)          //新建轮播图
		BannerRouter.PUT("/:id", middleware.JWTAuth(), middleware.AdminAuth(), banner.Update)    //修改轮播图信息
	}
}
