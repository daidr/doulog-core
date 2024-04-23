package service

import (
	"github.com/daidr/doulog-core/lib/models"
)

// DemoLogin 测试用户登录
func DemoLogin(db *models.DB) string {
	token := setToken(db, 2)
	return token
}
