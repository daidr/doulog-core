package service

import (
	"context"
	"github.com/daidr/doulog-core/lib/auth"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/lib/utils"
	authUtils "github.com/daidr/doulog-core/module/auth_login/internal/utils"
	"github.com/pkg/errors"
	"strings"
)

// Callback 返回设置在缓存中的token
func Callback(db *models.DB, state string, code string) (string, string, error) {
	states := strings.Split(state, "_")
	// 校验state格式 platform_mark
	if len(states) != 2 || !auth.PlatExists(states[0]) {
		return "", "", errors.New("format not matches")
	}

	if ex, err := db.Redis.
		Exists(context.Background(), format.Key.AuthCallback(states[1])).Result(); ex != 1 || err != nil {
		return "", "", errors.New("mark not exists")
	}

	// 获取callback url
	callback, err := db.Redis.Get(context.Background(), format.Key.AuthCallback(states[1])).Result()
	if err != nil {
		return "", "", err
	}

	// 校验 callback 是否合法
	if !utils.IsAllowedFrontendCallback(callback) {
		return "", callback, errors.New("invalid callback")
	}

	// 删除mark缓存
	_ = db.Redis.Del(context.Background(), format.Key.AuthCallback(states[1]))

	// 校验
	open, err := auth.Validate(states[0], code, state)
	if err != nil {
		return "", callback, err
	}

	o := models.TOauth{}
	// 先找是否绑定过，绑定过则设置完token直接重定向回去
	err = db.PgSQL.Where("platform = ? AND open_id = ?", states[0], open.Id).
		First(&o).Error
	if err == nil {
		token := authUtils.SetToken(db, o.User)
		return token, callback, nil
	}

	//// 没绑定过则新建用户
	//token := ""
	//err = db.PgSQL.Transaction(func(tx *gorm.DB) error {
	//	name := open.Name
	//	if name == "" {
	//		name = open.Login
	//	}
	//	u := models.TUser{
	//		Name:      name,
	//		Email:     open.Email,
	//		EmailHash: utils.GetMD5(open.Email),
	//		Homepage:  open.Homepage,
	//		IsAdmin:   false,
	//		Attr:      0,
	//	}
	//	if err := tx.Create(&u).Error; err != nil {
	//		return err
	//	}
	//
	//	// 如果 id 为 1 则设置为管理员
	//	if u.ID == 1 {
	//		u.IsAdmin = true
	//		if err := tx.Save(&u).Error; err != nil {
	//			return err
	//		}
	//	}
	//
	//	// 新建oauth绑定关系
	//	o = models.TOauth{
	//		User:     u.ID,
	//		Platform: states[0],
	//		OpenID:   open.ID,
	//	}
	//
	//	if err := tx.Create(&o).Error; err != nil {
	//		return err
	//	}
	//
	//	if err := search.IndexUser(search.UserSearch{
	//		ID:    u.ID,
	//		Name:  u.Name,
	//		Email: u.Email,
	//	}); err != nil {
	//		return err
	//	}
	//
	//	token = setToken(db, u.ID)
	//	return nil
	//})
	//if err != nil {
	//	return "", callback, err
	//}
	// 没绑定过则返回OAuth未绑定错误
	return "", callback, errors.New("The OAuth account has not been bound to a specific user yet.")
}
