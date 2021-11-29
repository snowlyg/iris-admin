package authority

import (
	"time"

	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthorityCollection []Authority

type Authority struct {
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-" sql:"index"`

	AuthorityId uint `json:"authorityId" gorm:"not null;primary_key;type:varchar(90)" binding:"required"`

	BaseAuthority
	Menus    []BaseMenu  `json:"menus" gorm:"many2many:authority_menus;"`
	Children []Authority `json:"children" gorm:"-"`
	Perms    [][]string  `json:"perms" gorm:"-"`
}

type BaseAuthority struct {
	AuthorityName string `json:"authorityName" gorm:"comment:角色名" binding:"required"`
	AuthorityType int    `json:"authorityType" gorm:"comment:角色类型"`
	ParentId      uint   `json:"parentId" gorm:"comment:父角色ID" binding:"required"`
	DefaultRouter string `json:"defaultRouter" gorm:"comment:默认菜单;default:dashboard"`
}

type BaseMenu struct {
	gorm.Model
	Pid        uint        `gorm:"index:pid;column:pid;type:int unsigned;not null;default:0" json:"pid"`         // 父级id
	Path       string      `gorm:"column:path;type:varchar(512);not null" json:"path"`                           // 路径
	Icon       string      `gorm:"column:icon;type:varchar(32);default:''" json:"icon"`                          // 图标
	MenuName   string      `gorm:"column:menu_name;type:varchar(128);not null;default:''" json:"menu_name"`      // 按钮名
	Route      string      `gorm:"column:route;type:varchar(64);not null" json:"route"`                          // 路由名称
	Params     string      `gorm:"column:params;type:varchar(128);not null;default:''" json:"params"`            // 参数
	Sort       int8        `gorm:"column:sort;type:tinyint;not null;default:0" json:"sort"`                      // 排序
	Hidden     int         `gorm:"column:hidden;type:tinyint unsigned;not null;default:1" json:"hidden"`         // 是否显示
	IsTenancy  int         `gorm:"column:is_tenancy;type:tinyint unsigned;not null;default:1" json:"is_tenancy"` // 模块，1 平台， 2商户
	IsMenu     int         `gorm:"column:is_menu;type:tinyint unsigned;not null;default:1" json:"is_menu"`       // 类型，1菜单 2 权限
	Authoritys []Authority `json:"authoritys" gorm:"many2many:authority_menus;"`
	Children   []BaseMenu  `json:"children" gorm:"-"`
}

// Create 添加
func (item *Authority) Create(db *gorm.DB) (uint, error) {
	err := db.Model(item).Create(item).Error
	if err != nil {
		zap_server.ZAPLOG.Error("添加失败", zap.String("(item *Authority) Create()", err.Error()))
		return item.AuthorityId, err
	}
	return item.AuthorityId, nil
}

// Update 更新
func (item *Authority) Update(db *gorm.DB, id uint, scopes ...func(db *gorm.DB) *gorm.DB) error {
	err := db.Model(item).Scopes(scopes...).Updates(item).Error
	if err != nil {
		zap_server.ZAPLOG.Error("更新失败", zap.String("(item *Authority) Update() ", err.Error()))
		return err
	}
	return nil
}

// Delete 删除
func (item *Authority) Delete(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) error {
	err := db.Model(item).Unscoped().Scopes(scopes...).Delete(item).Error
	if err != nil {
		zap_server.ZAPLOG.Error("删除失败", zap.String("(item *Authority) Delete()", err.Error()))
		return err
	}
	return nil
}
