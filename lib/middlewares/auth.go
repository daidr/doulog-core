package middlewares

import (
	"context"
	"github.com/daidr/doulog-core/lib/conf"
	"github.com/daidr/doulog-core/lib/daos"
	"github.com/daidr/doulog-core/lib/ecode"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strconv"
)

// Auth force为true为强制要求登录
// false为游客也可以访问，游客uid=0，若有登录则uid!=0
//
// 即在必须登陆操作的接口传入true,在游客、登录者均可访问的接口(某些字段因访问者而不同)传入false，完全不需要登录的接口不使用该中间件
func Auth(force bool, admin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := getUID(utils.GetScope(c).DB, c.GetHeader("Authorization"))
		if err != nil && force {
			c.Abort()
			format.HTTP(c, ecode.Unauthorized, nil)
			return
		}

		if admin {
			isAdmin, err := daos.NewUser(utils.GetScope(c).DB).IsAdmin(uid)
			if err != nil {
				c.Abort()
				format.HTTP(c, ecode.Unauthorized, nil)
				return
			}
			if !isAdmin {
				c.Abort()
				format.HTTP(c, ecode.PermissionDenied, nil)
				return
			}
		}

		c.Set("UID", uid)
	}
}

func getUID(db *models.DB, token string) (uint64, error) {
	// start with "Bearer "
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	if token == "" {
		return 0, errors.New("token is empty")
	}

	u, err := db.Redis.Get(context.Background(), format.Key.AuthToken(token)).Result()
	if err != nil {
		return 0, err
	}

	uid, err := strconv.ParseUint(u, 10, 64)
	if err != nil {
		return 0, err
	}

	if err = db.Redis.Expire(
		context.Background(),
		format.Key.AuthToken(token),
		conf.TokenExpire).Err(); err != nil {
		return 0, err
	}

	return uid, nil
}
