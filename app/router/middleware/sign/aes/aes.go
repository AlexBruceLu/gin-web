package sign_aes

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

var AppSecret string

//aes 对称加密
func SetUp() gin.HandlerFunc {

	return func(c *gin.Context) {

	}

}

// 验证签名
func verifySign(c *gin.Context) (map[string]string, error) {

}

// 创建签名
func createSign(params url.Values) string {

}
func createEncryptStr(params url.Values) string {

}
