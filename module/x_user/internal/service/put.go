package service

import (
	"github.com/daidr/doulog-core/lib/daos"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/lib/search"
)

/*
UpdateName 修改昵称
false err -> dao出错
false nil -> 存在相同昵称
*/
func UpdateName(db *models.DB, uid uint64, newName string) (bool, error) {
	if err := db.PgSQL.Where("name = ?", newName).First(&models.TUser{}).Error; err == nil {
		return false, nil
	}

	// 更新昵称
	if err := daos.NewUser(db).ChangeUserNameByUID(uid, newName); err != nil {
		return false, err
	}

	// 更新搜索
	now, err := search.UserGet(uid)
	if err != nil {
		return false, err
	}
	if err = search.UserEdit(search.UserSearch{ID: uid, Name: newName, Email: now.Email}); err != nil {
		return false, err
	}

	return true, nil
}

/*
UpdateEmail 修改邮箱
false err -> dao出错
false nil -> 存在相同邮箱
*/
func UpdateEmail(db *models.DB, uid uint64, newEmail string) (bool, error) {
	if err := db.PgSQL.Where("email = ?", newEmail).First(&models.TUser{}).Error; err == nil {
		return false, nil
	}

	// 更新邮箱
	if err := daos.NewUser(db).ChangeEmailByUID(uid, newEmail); err != nil {
		return false, err
	}

	// 更新搜索
	now, err := search.UserGet(uid)
	if err != nil {
		return false, err
	}
	if err = search.UserEdit(search.UserSearch{ID: uid, Name: newEmail, Email: now.Email}); err != nil {
		return false, err
	}

	return true, nil
}

/*
UpdateHomepage 修改主页
false err -> dao出错
*/
func UpdateHomepage(db *models.DB, uid uint64, newHomepage string) (bool, error) {
	// 更新主页
	if err := daos.NewUser(db).ChangeHomepageByUID(uid, newHomepage); err != nil {
		return false, err
	}

	return true, nil
}
