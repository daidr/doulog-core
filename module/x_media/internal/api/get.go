package api

import (
	"github.com/daidr/doulog-core/lib/ecode"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/localstorage"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/daidr/doulog-core/module/x_media/internal/e"
	"github.com/daidr/doulog-core/module/x_media/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetMediaMetadata(c *gin.Context) {
	sp := utils.GetScope(c)

	req := struct {
		TargetPid uint64 `uri:"id"`
	}{}

	if err := c.ShouldBindUri(&req); err != nil {
		format.HTTPInvalidParams(c)
		return
	}

	resp, err := service.GetMediaMetadata(sp.DB, req.TargetPid)
	if err != nil {
		sp.Log.Debugw("failed to get media metadata",
			"error", err,
			"TargetPid", req.TargetPid)
		format.HTTP(c, e.ErrGetMediaInfo, nil)
		return
	}
	format.HTTP(c, ecode.Success, resp)
}

func GetMediaById(c *gin.Context) {
	scope := utils.GetScope(c)

	req := struct {
		TargetPid uint64 `uri:"id"`
	}{}

	if err := c.ShouldBindUri(&req); err != nil {
		format.HTTPInvalidParams(c)
		return
	}

	thumbnail := c.Query("thumbnail")

	resp, err := service.GetMediaMetadata(scope.DB, req.TargetPid)
	if err != nil {
		scope.Log.Debugw("failed to get media metadata",
			"error", err,
			"TargetPid", req.TargetPid)
		format.HTTP(c, e.ErrGetMediaInfo, nil)
		return
	}

	etag := resp.Etag

	var hexEtag string

	if hexEtag, err = localstorage.EtagToHex(etag); err != nil {
		scope.Log.Debugw("failed to convert etag to hex",
			"error", err,
			"etag", etag)
		format.HTTP(c, e.ErrGetMediaInfo, nil)
		return
	}

	path := "/store/" + localstorage.MediaStore.Biz + "/" + hexEtag[:2] + "/" + hexEtag[2:4] + "/"

	// TODO: CDN support
	if thumbnail != "" {
		c.Redirect(http.StatusFound, path+hexEtag+"_thumbnail.webp")
	} else {
		c.Redirect(http.StatusFound, path+hexEtag+".webp")
	}
}
