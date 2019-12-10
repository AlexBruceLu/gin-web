package sign_rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"gin-web/app/config"
	"gin-web/app/utils/response"
	"gin-web/example/jaeger/speak/app/util"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var AppSecret string

// RSA 非对称加密
func SetUp() gin.HandlerFunc {

	return func(c *gin.Context) {
		utilGin := response.Gin{Ctx: c}

		sign, err := verifySign(c)

		if sign != nil {
			utilGin.Response(-1, "Debug Sign", sign)
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

	// 验证来源
	value, ok := config.APIAuthConfig[ak]
	if ok {
		AppSecret = value["rsa"]
	} else {
		return nil, errors.New("ak Error")
	}

	if debug == "1" {
		currentUnix := util.GetCurrentUnix()
		req.Set("ts", strconv.FormatInt(currentUnix, 10))

		sn, err := createSign(req)
		if err != nil {
			return nil, errors.New("sn Exception")
		}

		res := map[string]string{
			"ts": strconv.FormatInt(currentUnix, 10),
			"sn": sn,
		}
		return res, nil
	}

	// 验证过期时间
	timestamp := time.Now().Unix()
	exp, _ := strconv.ParseInt(config.AppSignExpiry, 10, 64)
	tsInt, _ := strconv.ParseInt(ts, 10, 64)
	if tsInt > timestamp || timestamp-tsInt >= exp {
		return nil, errors.New("ts Error")
	}

	// 验证签名
	if sn == "" {
		return nil, errors.New("sn Error")
	}

	decryptStr, decryptErr := PrivateDecrypt(sn, config.AppRSAPrivateFile)
	if decryptErr != nil {
		return nil, errors.New(decryptErr.Error())
	}
	if decryptStr != createEncryptStr(req) {
		return nil, errors.New("sn Error")
	}
	return nil, nil
}

// 创建签名
func createSign(params url.Values) (string, error) {
	return PublicEncrypt(createEncryptStr(params), AppSecret)
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

// 公钥加密
func PublicEncrypt(encryptStr, path string) (string, error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 读取文件内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)

	// pem 解码
	block, _ := pem.Decode(buf)

	// x509 解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)

	//对明文进行加密
	encryptedStr, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(encryptStr))
	if err != nil {
		return "", err
	}
	//返回密文
	return base64.StdEncoding.EncodeToString(encryptedStr), nil
}

// 私钥解密
func PrivateDecrypt(decryptStr, path string) (string, error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 获取文件内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)

	// pem 解码
	block, _ := pem.Decode(buf)

	// X509 解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	decryptBytes, err := base64.StdEncoding.DecodeString(decryptStr)

	//对密文进行解密
	decrypted, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decryptBytes)

	//返回明文
	return string(decrypted), nil
}
