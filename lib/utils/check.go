package utils

import (
	"errors"
	"strconv"
)

type check struct {
	kv map[string]string
}

var (
	ErrKeyNotFound = errors.New("key not found")
)

func newCheck(kv map[string]string) *check {
	return &check{kv: kv}
}

func (c *check) Uint64(key string) error {
	v, ok := c.kv[key]
	if !ok {
		return ErrKeyNotFound
	}

	if _, err := strconv.ParseUint(v, 10, 64); err != nil {
		return err
	}
	return nil
}
