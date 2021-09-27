package role

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/snowlyg/iris-admin/server/casbin"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/database/scope"
	myzap "github.com/snowlyg/iris-admin/server/zap"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ErrRoleNameInvalide = errors.New("角色名称已经被使用")

// GetAdminRoleName 获管理员角色名称
func GetAdminRoleName() string {
	return "admin"
}

// Create 添加
func Create(req *Request) (uint, error) {
	if _, err := FindByName(NameScope(req.Name)); !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, ErrRoleNameInvalide
	}
	role := &Role{BaseRole: req.BaseRole}
	id, err := orm.Create(database.Instance(), role)
	if err != nil {
		return 0, err
	}
	err = AddPermForRole(id, req.Perms)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// FindByName
func FindByName(scopes ...func(db *gorm.DB) *gorm.DB) (*Response, error) {
	role := &Response{}
	err := role.First(database.Instance(), scopes...)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func IsAdminRole(id uint) error {
	role := &Response{}
	err := role.First(database.Instance(), scope.IdScope(id))
	if err != nil {
		return err
	}
	if role.Name == GetAdminRoleName() {
		return errors.New("不能操作超级管理员")
	}
	return nil
}

func FindById(db *gorm.DB, id uint) (Response, error) {
	role := Response{}
	err := db.Model(&Role{}).Where("id = ?", id).First(&role).Error
	if err != nil {
		myzap.ZAPLOG.Error("根据id查询角色错误", zap.String("错误:", err.Error()))
		return role, err
	}
	return role, nil
}

func FindInId(db *gorm.DB, ids []string) ([]*Response, error) {
	roles := PageResponse{}
	err := roles.Find(database.Instance(), scope.InIdsScope(ids))
	if err != nil {
		myzap.ZAPLOG.Error("通过ids查询角色错误", zap.String("错误:", err.Error()))
		return nil, err
	}
	return roles, nil
}

// AddPermForRole
func AddPermForRole(id uint, perms [][]string) error {
	roleId := strconv.FormatUint(uint64(id), 10)
	oldPerms := casbin.GetPermissionsForUser(roleId)
	_, err := casbin.Instance().RemovePolicies(oldPerms)
	if err != nil {
		myzap.ZAPLOG.Error("add policy err: %+v", zap.String("错误:", err.Error()))
		return err
	}

	if len(perms) == 0 {
		myzap.ZAPLOG.Debug("没有权限")
		return nil
	}
	var newPerms [][]string
	for _, perm := range perms {
		newPerms = append(newPerms, append([]string{roleId}, perm...))
	}
	myzap.ZAPLOG.Info("添加权限到角色", myzap.Strings("新权限", newPerms))
	_, err = casbin.Instance().AddPolicies(newPerms)
	if err != nil {
		myzap.ZAPLOG.Error("add policy err: %+v", zap.String("错误:", err.Error()))
		return err
	}

	return nil
}

func GetRoleIds() ([]uint, error) {
	var roleIds []uint
	err := database.Instance().Model(&Role{}).Select("id").Find(&roleIds).Error
	if err != nil {
		return roleIds, fmt.Errorf("获取角色ids错误 %w", err)
	}
	return roleIds, nil
}
