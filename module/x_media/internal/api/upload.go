package api

import (
	"github.com/daidr/doulog-core/lib/ecode"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/daidr/doulog-core/module/x_media/internal/e"
	"github.com/daidr/doulog-core/module/x_media/internal/service"
	"github.com/gin-gonic/gin"
)

func MediaUpload(c *gin.Context) {
	file, err := c.FormFile("media")
	if err != nil {
		format.HTTP(c, e.ErrGetFormFile, nil)
		return
	}

	sp := utils.GetScope(c)
	uid := c.GetUint64("UID")

	id, err := service.Upload(sp, uid, file)
	if err != nil {
		sp.Log.Debugw("failed to upload media",
			"error", err,
			"uid", uid)
		format.HTTP(c, e.ErrMediaUpload, nil)
		return
	}

	format.HTTP(c, ecode.Success, gin.H{
		"pid": id,
	})
}
