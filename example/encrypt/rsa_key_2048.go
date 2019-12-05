package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	str := `
			登  高
			 杜甫
	风急天高猿啸哀，渚清沙白鸟飞回。
	无边落木萧萧下，不尽长江滚滚来。
	万里悲秋常作客，百年多病独登台。
	艰难苦恨繁霜鬓，潦倒新停浊酒杯。
	`
	fmt.Println("加密前：", str)

	encrypted, err := RsaEncrypt(str, "public_key.pem")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("加密后：", encrypted)
	}

	decrypted, err := RsaDecrypt(encrypted, "private_key.pem")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("解密后：", decrypted)
	}

}

// 公钥加密
func RsaEncrypt(encryptStr, path string) (string, error) {
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
func RsaDecrypt(decryptStr, path string) (string, error) {
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

/*
加密前：
                        登  高
                         杜甫
        风急天高猿啸哀，渚清沙白鸟飞回。
        无边落木萧萧下，不尽长江滚滚来。
        万里悲秋常作客，百年多病独登台。
        艰难苦恨繁霜鬓，潦倒新停浊酒杯。

加密后： DrZzpWwyIO/76eA0uiu4OPb8X94KDJMX3hBkrwYaudE4qclftiLFtcHcdUeg3G3wnU23X+Q/aBbxM8hf76CbSyjcSgw2RetVA9uJy8w5STrxTuq8MF3ScdEhpKj8BWBIOG1cCleBm2wXuIH/f4vzoEF+l1ikK5TdK0c8E5+4r03OOUK
QqJ/ekpujMAaEKvKG5eAN4ddNIEttU+nwSURH+RUIBlOa8ZppKK3IAPH6KGNH/6Ui1+LXTgLdzC1NAzTOrxFgBy4cAYAjStmrpyrAoy5chRu2r4qvLdVigiGYGwCsYqtgXC/bRUXpgEWhWhkOzmUumGfz3oSnrxmnfOjPRQ==
解密后：
                        登  高
                         杜甫
        风急天高猿啸哀，渚清沙白鸟飞回。
        无边落木萧萧下，不尽长江滚滚来。
        万里悲秋常作客，百年多病独登台。
        艰难苦恨繁霜鬓，潦倒新停浊酒杯。

*/
