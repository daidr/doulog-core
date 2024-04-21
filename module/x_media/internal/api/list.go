package api

import (
	"github.com/daidr/doulog-core/lib/ecode"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/daidr/doulog-core/module/x_media/internal/dto"
	"github.com/daidr/doulog-core/module/x_media/internal/e"
	"github.com/daidr/doulog-core/module/x_media/internal/service"
	"github.com/gin-gonic/gin"
)

func ListAllMedia(c *gin.Context) {
	sp := utils.GetScope(c)
	var req dto.MediaListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespLogger(c).Debugw("failed to bind request", "error", err)
		format.HTTP(c, ecode.InvalidParams, nil)
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 10 {
		req.PageSize = 10
	}

	if req.PageSize > 100 {
		req.PageSize = 10
	}

	resp, err := service.ListAllMedia(sp.DB, &req)

	if err != nil {
		sp.Log.Debugw("failed to list all media",
			"error", err,
			"req", req)
		format.HTTP(c, e.ErrGetMediaInfo, nil)
		return
	}

	format.HTTP(c, ecode.Success, resp)
}
