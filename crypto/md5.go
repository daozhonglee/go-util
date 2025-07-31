package crypto

import (
	"crypto/md5"
	"fmt"
)

func Md5(value string) string {
	data := []byte(value)
	has := md5.Sum(data)
	strmd5 := fmt.Sprintf("%x", has)
	return strmd5
}
