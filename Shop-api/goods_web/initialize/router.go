package initialize

import (
	"Shop-api/goods_web/middleware"
	"Shop-api/goods_web/router"
	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	GoodsGroup := r.Group("/u/v1")
	router.InitGoodsRouter(GoodsGroup)
	router.InitBrandRouter(GoodsGroup)
	router.InitCategoryRouter(GoodsGroup)
	router.InitBannerRouter(GoodsGroup)
	return r
}
