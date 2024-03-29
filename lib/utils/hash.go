package utils

import (
	"crypto/md5"
	"fmt"
)

func GetMD5(s string) string {
	data := []byte(s)
	return fmt.Sprintf("%x", md5.Sum(data))
}
