package router

import (
	"Shop-api/goods_web/api/goods"
	"Shop-api/goods_web/middleware"
	"github.com/gin-gonic/gin"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	//.Use(middleware.Trace())
	{
		GoodsRouter.GET("", goods.GetGoodsList)                                                //商品列表
		GoodsRouter.POST("", middleware.JWTAuth(), middleware.AdminAuth(), goods.New)          //改接口需要管理员权限
		GoodsRouter.GET("/:id", goods.Detail)                                                  //获取商品的详情
		GoodsRouter.DELETE("/:id", middleware.JWTAuth(), middleware.AdminAuth(), goods.Delete) //删除商品
		GoodsRouter.GET("/:id/stocks", goods.Stocks)                                           //获取商品的库存

		GoodsRouter.PUT("/:id", middleware.JWTAuth(), middleware.AdminAuth(), goods.Update)
		GoodsRouter.PATCH("/:id", middleware.JWTAuth(), middleware.AdminAuth(), goods.UpdateStatus)
	}
}
