package user

import (
	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	BaseUser
	Password string `gorm:"type:varchar(250)" json:"password" validate:"required"`
	RoleIds  []uint `gorm:"-" json:"role_ids"`
}

type BaseUser struct {
	Name     string `gorm:"index;not null; type:varchar(60)" json:"name"`
	Username string `gorm:"uniqueIndex;not null;type:varchar(60)" json:"username" validate:"required"`
	Intro    string `gorm:"not null; type:varchar(512)" json:"intro"`
	Avatar   string `gorm:"type:varchar(1024)" json:"avatar"`
}

type Avatar struct {
	Avatar string `json:"avatar"`
}

// Create 添加
func (item *User) Create(db *gorm.DB) (uint, error) {
	err := db.Model(item).Create(item).Error
	if err != nil {
		zap_server.ZAPLOG.Error("添加失败", zap.String("(item *User) Create()", err.Error()))
		return item.ID, err
	}
	return item.ID, nil
}

// Update 更新
func (item *User) Update(db *gorm.DB, id uint, scopes ...func(db *gorm.DB) *gorm.DB) error {
	err := db.Model(item).Scopes(scopes...).Updates(item).Error
	if err != nil {
		zap_server.ZAPLOG.Error("更新失败", zap.String("(item *User) Update() ", err.Error()))
		return err
	}
	return nil
}

// Delete 删除
func (item *User) Delete(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) error {
	err := db.Model(item).Unscoped().Scopes(scopes...).Delete(item).Error
	if err != nil {
		zap_server.ZAPLOG.Error("删除失败", zap.String("(item *User) Delete()", err.Error()))
		return err
	}
	return nil
}
