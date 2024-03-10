package conf

import "time"

type RouterNS string

const (
	RouterNSMain RouterNS = "x"
	RouterNSAuth RouterNS = "auth"
	RouterNSTest RouterNS = "test"
)

func (r RouterNS) NS() string {
	return string(r)
}

const (
	TokenExpire = 7 * 24 * time.Hour
)

const (
	AuthGitHub = "github"
)

const (
	ImageSrcLocal = iota + 1
	ImageSrcAliYun
)
