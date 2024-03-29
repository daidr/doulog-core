package daos

import (
	"github.com/daidr/doulog-core/lib/models"
	"github.com/pkg/errors"
)

type User struct {
	db *models.DB
}

func NewUser(db *models.DB) *User {
	return &User{db}
}

func (d *User) Get(uid uint64) (*models.TUser, error) {
	upper := models.TUser{}
	if err := d.db.PgSQL.First(&upper, uid).Error; err != nil {
		return nil, err
	}
	return &upper, nil
}

func (d *User) GetB(uid uint64) (*models.BUser, error) {
	var u models.BUser

	err := d.db.PgSQL.Model(&models.TUser{}).
		Select("t_users.id, t_users.name, t_users.attr, t_users.email, t_users.email_hash, t_users.is_admin, t_users.homepage, t_users.created_at").
		Where("t_users.id = ?", uid).First(&u).Error
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get user")
	}

	return &u, nil
}

func (d *User) IsAdmin(uid uint64) (bool, error) {
	var u models.BUser
	err := d.db.PgSQL.Model(&models.TUser{}).
		Select("t_users.is_admin").
		Where("t_users.id = ?", uid).First(&u).Error
	if err != nil {
		return false, errors.WithMessage(err, "failed to get user")
	}
	return u.IsAdmin, nil
}

func (d *User) Add(user *models.TUser) error {
	if err := d.db.PgSQL.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (d *User) ChangeUserNameByUID(uid uint64, newName string) error {
	err := d.db.PgSQL.Model(&models.TUser{}).
		Where("id = ?", uid).
		Update("name", newName).Error

	if err != nil {
		return errors.WithMessage(err, "failed to change username")
	}
	return nil
}

func (d *User) ChangeEmailByUID(uid uint64, newEmail string) error {
	err := d.db.PgSQL.Model(&models.TUser{}).
		Where("id = ?", uid).
		Update("email", newEmail).Error

	if err != nil {
		return errors.WithMessage(err, "failed to change email")
	}
	return nil
}

func (d *User) ChangeHomepageByUID(uid uint64, newHomepage string) error {
	err := d.db.PgSQL.Model(&models.TUser{}).
		Where("id = ?", uid).
		Update("homepage", newHomepage).Error

	if err != nil {
		return errors.WithMessage(err, "failed to change homepage")
	}
	return nil
}
