// Package e 错误码定义
package e

import "github.com/daidr/doulog-core/lib/ecode"

var (
	ErrMediaUpload  = ecode.New(110001)
	ErrGetFormFile  = ecode.New(110002)
	ErrGetMediaInfo = ecode.New(110003)
)

var ECode = map[ecode.Code]string{
	ErrMediaUpload:  "图片上传失败",
	ErrGetFormFile:  "获取上传文件失败",
	ErrGetMediaInfo: "获取图片信息失败",
}
