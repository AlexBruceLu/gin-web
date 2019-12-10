package sign_aes

import (
	"errors"
	"fmt"
	"gin-web/app/config"
	"gin-web/app/tools/sign/aes"
	"gin-web/app/utils/response"
	"gin-web/example/jaeger/speak/app/util"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var AppSecret string

//aes 对称加密
func SetUp() gin.HandlerFunc {

	return func(c *gin.Context) {
		utilGin := response.Gin{Ctx: c}

		sign, err := verifySign(c)
		if sign != nil {
			utilGin.Response(-1, "debug sign", sign)
			c.Abort()
			return
		}
		if err != nil {
			utilGin.Response(-1, err.Error(), sign)
			c.Abort()
			return
		}

		c.Next()

	}

}

// 验证签名
func verifySign(c *gin.Context) (map[string]string, error) {
	_ = c.Request.ParseForm()

	req := c.Request.Form
	debug := strings.Join(c.Request.Form["debug"], "")
	ak := strings.Join(c.Request.Form["ak"], "")
	sn := strings.Join(c.Request.Form["sn"], "")
	ts := strings.Join(c.Request.Form["ts"], "")

	value, ok := config.APIAuthConfig[ak]
	if ok {
		AppSecret = value["aes"]
	} else {
		return nil, errors.New("ak error")
	}

	if debug == "1" {
		currentUnix := util.GetCurrentUnix()
		req.Set("ts", strconv.FormatInt(currentUnix, 10))
		sn, err := createSign(req)
		if err != nil {
			return nil, errors.New("sn exception")
		}

		res := map[string]string{
			"ts": strconv.FormatInt(currentUnix, 10),
			"sn": sn,
		}
		return res, nil
	}

	timeStamp := time.Now().Unix()
	exp, _ := strconv.ParseInt(config.AppSignExpiry, 10, 64)
	tsInt, _ := strconv.ParseInt(ts, 10, 64)
	if tsInt > timeStamp || timeStamp-tsInt >= exp {
		return nil, errors.New("ts error")
	}

	if sn == "" {
		return nil, errors.New("sn error")
	}

	decryptStr, err := aes.Decrypt(sn, []byte(AppSecret), AppSecret)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	if decryptStr != createEncryptStr(req) {
		return nil, errors.New("sn error")
	}

	return nil, nil

}

// 创建签名
func createSign(params url.Values) (string, error) {
	return aes.Encrypt(createEncryptStr(params), []byte(AppSecret), AppSecret)
}
func createEncryptStr(params url.Values) string {
	var key []string
	var str = ""
	for k := range params {
		if k != "sn" && k != "debug" {
			key = append(key, k)
		}
	}
	sort.Strings(key)
	for i := 0; i < len(key); i++ {
		if i == 0 {
			str = fmt.Sprintf("%v=%v", key[i], params.Get(key[i]))
		} else {
			str = str + fmt.Sprintf("&%v=%v", key[i], params.Get(key[i]))
		}
	}
	return str
}
