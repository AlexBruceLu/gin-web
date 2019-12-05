package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

const (
	key = "HFu8Z5SjAT7CudQc"
	iv  = "HFu8Z5SjAT7CudQc" // 加密的iv向量
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

	encrypted, err := encrypt(str, []byte(key))
	if err != nil {
		fmt.Println("ERROR: encrypt failed:", err)
	} else {
		fmt.Println("加密后", encrypted)
	}

	decrypted, err := decrypt(encrypted, []byte(key))
	if err != nil {
		fmt.Println("ERROR: decrypt failed:", err)
	} else {
		fmt.Println("解密后", decrypted)
	}
}

// 加密
func encrypt(encryptStr string, key []byte) (string, error) {
	encryptBytes := []byte(encryptStr)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	encryptBytes = PKCS5Padding(encryptBytes, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	encrypted := make([]byte, len(encryptBytes))
	blockMode.CryptBlocks(encrypted, encryptBytes)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// 解密
func decrypt(decryptStr string, key []byte) (string, error) {
	decryptBytes, err := base64.StdEncoding.DecodeString(decryptStr)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	decrypted := make([]byte, len(decryptBytes))

	blockMode.CryptBlocks(decrypted, decryptBytes)
	decrypted = PKCS5UnPadding(decrypted)

	return string(decrypted), nil
}

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}
func PKCS5UnPadding(decrypted []byte) []byte {
	length := len(decrypted)
	unPadding := int(decrypted[length-1])
	return decrypted[:(length - unPadding)]
}

/*加密前：
                        登  高
                         杜甫
        风急天高猿啸哀，渚清沙白鸟飞回。
        无边落木萧萧下，不尽长江滚滚来。
        万里悲秋常作客，百年多病独登台。
        艰难苦恨繁霜鬓，潦倒新停浊酒杯。

加密后 45MPXZ4Pep90RBZXwFeqo1Qcu0m+taA6hHY3+r8Qjd+bcuoBkLv31vxlP+ZA7MJgQj4SRM7N7WZA/NJhP0r5qRJSzIU4TzgKN+d5cyrxnDJh/f4aGx5LczIOgxqvso/3KsdlZeOlXKvgAKTn7iyieeY63gLx+rCwQXUGE2e2QmbgqfbDT
VFP6W0tsOXi+paiqEloMJQIFfvuWWrvzbjZFit60ryoIqTGwpa8CHsJtX/grl1R/XzT2QHAu/XaHCaqCqsmhmbawDJvVPgcF2bBaEuTICcUymOBqmO5a5RtEcuxwCx5kTWK4gGT3+bXXfWw
解密后
                        登  高
                         杜甫
        风急天高猿啸哀，渚清沙白鸟飞回。
        无边落木萧萧下，不尽长江滚滚来。
        万里悲秋常作客，百年多病独登台。
        艰难苦恨繁霜鬓，潦倒新停浊酒杯。
*/
