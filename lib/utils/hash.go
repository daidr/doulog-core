package utils

import (
	"crypto/md5"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func GetMD5(s string) string {
	data := []byte(s)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func EncodeBcrypt(s string, round int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), round)
	return string(hash), err
}

func CompareBcrypt(s, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
	return err == nil
}
