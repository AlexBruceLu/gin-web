package product

import (
	"gin-web/app/controller/param_bind"
	"gin-web/app/controller/param_verify"
	"gin-web/app/utils/bind"
	"gin-web/app/utils/response"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

// 添加
func Add(c *gin.Context) {
	utilGin := response.Gin{Ctx: c}
	// 绑定参数
	s, e := bind.Bind(&param_bind.ProductAdd{}, c)
	if e != nil {
		utilGin.Response(-1, e.Error(), nil)
		return
	}

	// 参数验证
	validate := validator.New()
	if err := validate.RegisterValidation("NameValid", param_verify.NameValid); err != nil {
		utilGin.Response(-1, err.Error(), nil)
		return
	}

	if err := validate.Struct(s); err != nil {
		utilGin.Response(-1, err.Error(), nil)
		return
	}

	utilGin.Response(1, "success", nil)

}

// 修改
func Edit(c *gin.Context) {
	log.Info(c.Request.RequestURI)
}

// 查询
func Detail(c *gin.Context) {
	log.Info(c.Request.RequestURI)
}

// 删除
func Delete(c *gin.Context) {
	log.Info(c.Request.RequestURI)
}
