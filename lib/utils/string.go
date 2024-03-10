package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"unicode/utf8"
)

// parseKV k1=v1,k2=v2,k3=v3
func parseKV(content string) map[string]string {
	kvMap := make(map[string]string)
	kvs := strings.Split(content, ",")
	for _, kv := range kvs {
		kvArr := strings.Split(kv, "=")
		if len(kvArr) != 2 {
			continue
		}
		kvMap[kvArr[0]] = kvArr[1]
	}
	return kvMap
}

func RandString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func InStringSlice(s string, slice []string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func Strings2Uint64s(strs []string) ([]uint64, error) {
	ids := make([]uint64, 0, len(strs))
	for _, str := range strs {
		id, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func SubStrDecodeRuneInString(s string, length int) string {
	var size, n int
	for i := 0; i < length && n < len(s); i++ {
		_, size = utf8.DecodeRuneInString(s[n:])
		n += size
	}

	return s[:n]
}
