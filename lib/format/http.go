package format

import (
	"github.com/daidr/doulog-core/lib/ecode"
	"github.com/gin-gonic/gin"
	"net/http"
)

type resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// HTTP code 为错误码，可以传入公共错误码，也可以传入模块内部自定义错误码
func HTTP(c *gin.Context, code ecode.Code, data interface{}) {
	c.JSON(http.StatusOK, &resp{
		Code: code.Code(),
		Msg:  code.Msg(),
		Data: data,
	})
}

func HTTPInvalidParams(c *gin.Context) {
	c.JSON(http.StatusOK, &resp{
		Code: ecode.InvalidParams.Code(),
		Msg:  ecode.InvalidParams.Msg(),
		Data: nil,
	})
}
