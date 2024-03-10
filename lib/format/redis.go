package format

import (
	"bytes"
	"strconv"
	"strings"
	"sync"
)

var keyPool = sync.Pool{
	New: func() interface{} {
		b := &bytes.Buffer{}
		b.Grow(32)
		return b
	},
}

type key struct{}

var Key key

func (key) Gen(business string, indexes ...string) string {
	buf := keyPool.Get().(*bytes.Buffer)

	buf.WriteString(business)
	buf.WriteString(":")
	buf.WriteString(strings.Join(indexes, ":"))

	key := buf.String()
	buf.Reset()
	keyPool.Put(buf)
	return key
}

func (key) Views(tp int, oid uint64) string {
	return Key.Gen("vw", strconv.Itoa(tp), strconv.FormatUint(oid, 10))
}

func (key) Likes(tp int, oid uint64) string {
	return Key.Gen("lk", strconv.Itoa(tp), strconv.FormatUint(oid, 10))
}

func (key) AuthToken(token string) string {
	return Key.Gen("au", "t", token)
}

func (key) AuthCallback(mark string) string {
	return Key.Gen("au", "c", mark)
}

func (key) LikeHistory(tp int, uid uint64) string {
	return Key.Gen("ul", strconv.Itoa(tp), strconv.FormatUint(uid, 10))
}
