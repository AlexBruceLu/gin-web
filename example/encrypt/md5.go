package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
)

func main() {
	params := map[string]interface{}{
		"name": "Tomato",
		"pwd":  "12345",
		"age":  22,
	}
	fmt.Printf("sign %s\n", CreateSign(params))
	// sign 3a4c4ae9984fd286f839111f3a17758b
}
func MD5(str string) string {
	s := md5.New()
	s.Write([]byte(str))
	return hex.EncodeToString(s.Sum(nil))
}

// 生成签名
func CreateSign(params map[string]interface{}) string {
	var key []string
	var str = ""
	for k := range params {
		key = append(key, k)
	}
	sort.Strings(key)
	for i := 0; i < len(key); i++ {
		if i == 0 {
			str = fmt.Sprintf("%v=%v", key[i], params[key[i]])
		}
	}

	scret := "1234567890"

	sign := MD5(MD5(str) + MD5(scret))
	return sign
}
