package me

import (
	"github.com/daidr/doulog-core/lib/daos"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/lib/search"
)

/*
ChangeNickName 修改昵称
false err -> dao出错
false nil -> 存在相同昵称
*/
func ChangeNickName(db *models.DB, uid uint64, newName string) (bool, error) {
	if err := db.PgSQL.Where("uname = ?", newName).First(&models.TUser{}).Error; err == nil {
		return false, nil
	}

	// 更新昵称
	if err := daos.NewUser(db).ChangeUserNameByUID(uid, newName); err != nil {
		return false, err
	}

	// 更新搜索 TODO: 逻辑需要优化，现在有可能MySQL写入，搜索炸了
	now, err := search.UserGet(uid)
	if err != nil {
		return false, err
	}
	if err = search.UserEdit(search.UserSearch{ID: uid, Name: newName, Email: now.Email}); err != nil {
		return false, err
	}

	return true, nil
}
