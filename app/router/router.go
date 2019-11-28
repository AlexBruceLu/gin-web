package router

import (
	"gin-web/app/controller/product"

	"github.com/gin-gonic/gin"
)

func SetupRouter(engine *gin.Engine) {
	ProductRoute := engine.Group("")
	{
		// 新增产品
		ProductRoute.POST("/product", product.Add)
		// 更新产品
		ProductRoute.PUT("/product/:id", product.Edit)
		// 删除产品
		ProductRoute.DELETE("/product/:id", product.Delete)
		// 获取产品详情
		ProductRoute.GET("/product/:id", product.Detail)
	}
}
