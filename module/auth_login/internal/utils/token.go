package utils

import (
	"context"
	"github.com/daidr/doulog-core/lib/conf"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/lib/utils"
)

func SetToken(db *models.DB, uid uint64) string {
	for {
		token := utils.RandString(40)

		ok, err := db.Redis.
			SetNX(context.Background(),
				format.Key.AuthToken(token), uid, conf.TokenExpire).
			Result()
		if ok && err == nil {
			return token
		}
	}
}
