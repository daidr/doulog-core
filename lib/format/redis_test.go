package format

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestRedis_Key(t *testing.T) {
	if r := Key.Gen("basic", "user", "username", strconv.FormatUint(7829, 10)); r != "test.user:basic:user:username:7829" {
		t.Log(r)
		t.FailNow()
	}
}

func BenchmarkRedis_Key(b *testing.B) {
	for i := 0; i < b.N; i++ {
		a := Key.Gen("op", "111", "222", "333", "444")
		_ = a
	}
}
func BenchmarkRedis_Key1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		a := fmt.Sprintf("%s:%s:%s", "test.mod", "op", strings.Join([]string{"111", "222", "333", "444"}, ":"))
		_ = a
	}
}

func BenchmarkRedis_Key2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := &bytes.Buffer{}
		buf.WriteString("test.mod")
		buf.WriteString(":")
		buf.WriteString("op")
		buf.WriteString(":")
		buf.WriteString(strings.Join([]string{"111", "222", "333", "444"}, ":"))
		a := buf.String()
		_ = a
	}
}

func BenchmarkRedis_Key3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		builder := strings.Builder{}
		builder.WriteString("test.mod")
		builder.WriteString(":")
		builder.WriteString("op")
		builder.WriteString(":")
		builder.WriteString(strings.Join([]string{"111", "222", "333", "444"}, ":"))
		a := builder.String()
		_ = a
	}
}

func BenchmarkRedis_Key4(b *testing.B) {
	pool := sync.Pool{
		New: func() interface{} {
			buf := &bytes.Buffer{}
			buf.Grow(32)
			return buf
		},
	}
	for i := 0; i < b.N; i++ {
		buf := pool.Get().(*bytes.Buffer)
		buf.WriteString("test.mod")
		buf.WriteString(":")
		buf.WriteString("op")
		buf.WriteString(":")
		buf.WriteString(strings.Join([]string{"111", "222", "333", "444"}, ":"))
		a := buf.String()
		_ = a
		buf.Reset()
		pool.Put(buf)
	}
}
