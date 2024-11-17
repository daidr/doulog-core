package daos

import (
	"github.com/daidr/doulog-core/lib/conf"
	"github.com/daidr/doulog-core/lib/haikunator_zh"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/lib/search"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/pkg/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type User struct {
	db *models.DB
}

func NewUser(db *models.DB) *User {
	return &User{db}
}

func (d *User) CheckIsInitialized() (bool, error) {
	var total int64
	if err := d.db.PgSQL.Model(&models.TUser{}).Count(&total).Error; err != nil {
		return false, err
	}
	return total > 0, nil
}

func (d *User) CreateUser(Name string, Email string, Password string, IsAdmin bool) error {
	encodedPassword, err := utils.EncodeBcrypt(Password, conf.C.Auth.BcryptRounds)
	if err != nil {
		return errors.WithMessage(err, "failed to encode password")
	}
	user := &models.TUser{
		Name:      Name,
		Email:     Email,
		EmailHash: utils.GetMD5(Email),
		Password:  encodedPassword,
		IsAdmin:   IsAdmin,
		Attr:      0,
	}
	err = d.db.PgSQL.Transaction(func(tx *gorm.DB) error {
		if err := d.db.PgSQL.Create(user).Error; err != nil {
			return err
		}
		if err := search.IndexUser(search.UserSearch{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (d *User) Get(uid uint64) (*models.TUser, error) {
	upper := models.TUser{}
	if err := d.db.PgSQL.First(&upper, uid).Error; err != nil {
		return nil, err
	}
	return &upper, nil
}

func (d *User) GetWithCredentials(uid uint64) (*models.TUser, error) {
	upper := models.TUser{}
	if err := d.db.PgSQL.Model(models.TUser{}).Preload("Credentials").First(&upper, uid).Error; err != nil {
		return nil, err
	}
	return &upper, nil
}

func (d *User) GetByEmail(email string) (*models.TUser, error) {
	upper := models.TUser{}
	if err := d.db.PgSQL.Where("email = ?", email).First(&upper).Error; err != nil {
		return nil, err
	}
	return &upper, nil
}

func (d *User) CheckEmailExists(email string) (bool, error) {
	var total int64
	if err := d.db.PgSQL.Model(&models.TUser{}).Where("email = ?", email).Count(&total).Error; err != nil {
		return false, err
	}
	return total > 0, nil
}

func (d *User) GetB(uid uint64) (*models.BUser, error) {
	var u models.BUser

	err := d.db.PgSQL.Model(&models.TUser{}).
		Select("t_users.id, t_users.name, t_users.attr, t_users.email, t_users.motto, t_users.email_hash, t_users.is_admin, t_users.homepage, t_users.created_at").
		Where("t_users.id = ?", uid).First(&u).Error
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get user")
	}

	return &u, nil
}

func (d *User) IsAdmin(uid uint64) (bool, int, error) {
	var u models.BUser
	err := d.db.PgSQL.Model(&models.TUser{}).
		Select("t_users.is_admin, t_users.attr").
		Where("t_users.id = ?", uid).First(&u).Error
	if err != nil {
		return false, 0, errors.WithMessage(err, "failed to get user")
	}
	return u.IsAdmin, u.Attr, nil
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

var hai = haikunator_zh.New(time.Now().UTC().UnixNano())

func getRandomTag() string {
	return hai.Haikunate()
}

func (d *User) AddCredential(uid uint64, credential *webauthn.Credential) error {
	var user = &models.TUser{
		ID: uid,
	}
	err := d.db.PgSQL.Model(&user).
		Where("id = ?", uid).Association("Credentials").Append(&models.TWebAuthnCredential{
		UserID:     uid,
		Label:      getRandomTag(),
		LastUsedAt: -1,
		Credential: datatypes.NewJSONType(*credential),
	})

	if err != nil {
		return errors.WithMessage(err, "failed to add credential")
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

func (d *User) ChangeMottoByUID(uid uint64, newMotto string) error {
	err := d.db.PgSQL.Model(&models.TUser{}).
		Where("id = ?", uid).
		Update("motto", newMotto).Error

	if err != nil {
		return errors.WithMessage(err, "failed to change motto")
	}
	return nil
}
