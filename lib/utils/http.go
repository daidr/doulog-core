package utils

import (
	"github.com/daidr/doulog-core/lib/conf"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
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

	if len(buffer) >= 12 && string(buffer[4:12]) == "ftypheic" {
		return "image/heic", nil
	}

	return http.DetectContentType(buffer), nil
}

func IsAllowedFrontendCallback(callback string) bool {
	frontendCallbackPrefix := conf.C.Auth.FrontendCallbackPrefix
	if len(frontendCallbackPrefix) == 0 {
		// if no prefix is set, allow all
		return true
	}
	u, err := url.Parse(callback)
	if err != nil {
		return false
	}
	scheme := u.Scheme
	host := u.Hostname()
	port := u.Port()
	for _, prefix := range frontendCallbackPrefix {
		prefixUrl, err := url.Parse(prefix)
		if err != nil {
			continue
		}

		if host == prefixUrl.Hostname() && port == prefixUrl.Port() && scheme == prefixUrl.Scheme {
			return true
		}
	}

	return false
}
