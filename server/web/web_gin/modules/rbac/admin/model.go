package admin

import (
	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Password string `json:"-"  gorm:"not null;type:varchar(128);comment:用户登录密码"`
	BaseAdmin
	AuthorityIds []uint `json:"authorityIds" gorm:"-"`
}

type BaseAdmin struct {
	Username string `json:"userName" gorm:"not null;type:varchar(32);comment:用户登录名"`
	Status   int    `gorm:"column:status;type:tinyint(1);not null;default:1" json:"status"`   // 账号冻结 1为正常，2为禁止
	IsShow   int    `gorm:"column:is_show;type:tinyint(1);not null;default:1" json:"is_show"` // 是否显示 1为正常，2为禁止
	Email    string `json:"email" gorm:"default:'';comment:员工邮箱"`
	Phone    string `json:"phone" gorm:"type:char(15);default:'';comment:员工手机号" `
	NickName string `json:"nickName" gorm:"type:varchar(16);default:'员工姓名';comment:员工姓名" `
	Avatar
}

type Avatar struct {
	HeaderImg string `json:"headerImg" gorm:"default:http://qmplusimg.henrongyi.top/head.png;comment:用户头像"`
}

// Create 添加
func (item *Admin) Create(db *gorm.DB) (uint, error) {
	err := db.Model(item).Create(item).Error
	if err != nil {
		zap_server.ZAPLOG.Error("添加失败", zap.String("(item *Admin) Create()", err.Error()))
		return item.ID, err
	}
	return item.ID, nil
}

// Update 更新
func (item *Admin) Update(db *gorm.DB, id uint, scopes ...func(db *gorm.DB) *gorm.DB) error {
	err := db.Model(item).Scopes(scopes...).Updates(item).Error
	if err != nil {
		zap_server.ZAPLOG.Error("更新失败", zap.String("(item *Admin) Update() ", err.Error()))
		return err
	}
	return nil
}

// Delete 删除
func (item *Admin) Delete(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) error {
	err := db.Model(item).Unscoped().Scopes(scopes...).Delete(item).Error
	if err != nil {
		zap_server.ZAPLOG.Error("删除失败", zap.String("(item *Admin) Delete()", err.Error()))
		return err
	}
	return nil
}
