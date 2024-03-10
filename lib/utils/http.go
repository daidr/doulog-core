package utils

import (
	"github.com/daidr/doulog-core/lib/conf"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func GetScope(c *gin.Context) *models.Scope {
	return c.MustGet("scope").(*models.Scope)
}

func RespLogger(c *gin.Context) *zap.SugaredLogger {
	return GetScope(c).Log.Named("resp")
}

func UrlAppend(rawUrl string, queries map[string][]string) string {
	uu, err := url.Parse(rawUrl)
	if err != nil {
		return rawUrl
	}
	q1 := uu.Query()
	for k, v := range queries {
		q1[k] = v
	}
	uu.RawQuery = q1.Encode()
	return uu.String()
}

func GetMPFDContentType(src io.Reader) (string, error) {
	buffer := make([]byte, 512)
	if _, err := src.Read(buffer); err != nil {
		return "", err
	}

	return http.DetectContentType(buffer), nil
}

func IsAllowedFrontendCallback(callback string) bool {
	frontendCallbackPrefix := conf.C.Auth.FrontendCallbackPrefix
	if len(frontendCallbackPrefix) == 0 {
		// if no prefix is set, allow all
		return true
	}

	for _, prefix := range frontendCallbackPrefix {
		if strings.HasPrefix(callback, prefix) {
			return true
		}
	}
	
	return false
}
