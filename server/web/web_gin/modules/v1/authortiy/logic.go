package authority

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/snowlyg/iris-admin/server/casbin"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/database/scope"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ErrRoleNameInvalide = errors.New("角色名称已经被使用")

// GetAdminRoleName 获管理员角色名称
func GetAdminRoleName() string {
	return "admin"
}

func Copy(req *AuthorityCopyResponse) (uint, error) {
	return 0, nil
}

// Create 添加
func Create(req *Request) (uint, error) {
	if _, err := FindByName(AuthorityNameScope(req.AuthorityName)); !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, ErrRoleNameInvalide
	}
	authority := &Authority{BaseAuthority: req.BaseAuthority}
	id, err := orm.Create(database.Instance(), authority)
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
	authority := &Response{}
	err := authority.First(database.Instance(), scope.IdScope(id))
	if err != nil {
		return err
	}
	if authority.AuthorityName == GetAdminRoleName() {
		return errors.New("不能操作超级管理员")
	}
	return nil
}

func FindById(db *gorm.DB, id uint) (Response, error) {
	authority := Response{}
	err := db.Model(&Authority{}).Where("id = ?", id).First(&authority).Error
	if err != nil {
		zap_server.ZAPLOG.Error("根据id查询角色错误", zap.String("错误:", err.Error()))
		return authority, err
	}
	return authority, nil
}

func FindInId(db *gorm.DB, ids []uint) ([]*Response, error) {
	authorities := PageResponse{}
	err := authorities.Find(database.Instance(), scope.InIdsScope(ids))
	if err != nil {
		zap_server.ZAPLOG.Error("通过ids查询角色错误", zap.String("错误:", err.Error()))
		return nil, err
	}
	return authorities.Item, nil
}

// AddPermForRole
func AddPermForRole(id uint, perms [][]string) error {
	roleId := strconv.FormatUint(uint64(id), 10)
	oldPerms := casbin.Instance().GetPermissionsForUser(roleId)
	_, err := casbin.Instance().RemovePolicies(oldPerms)
	if err != nil {
		zap_server.ZAPLOG.Error("add policy err: %+v", zap.String("错误:", err.Error()))
		return err
	}

	if len(perms) == 0 {
		zap_server.ZAPLOG.Debug("没有权限")
		return nil
	}
	var newPerms [][]string
	for _, perm := range perms {
		newPerms = append(newPerms, append([]string{roleId}, perm...))
	}
	zap_server.ZAPLOG.Info("添加权限到角色", zap_server.Strings("新权限", newPerms))
	_, err = casbin.Instance().AddPolicies(newPerms)
	if err != nil {
		zap_server.ZAPLOG.Error("add policy err: %+v", zap.String("错误:", err.Error()))
		return err
	}

	return nil
}

func GetRoleIds() ([]uint, error) {
	var roleIds []uint
	err := database.Instance().Model(&Authority{}).Select("id").Find(&roleIds).Error
	if err != nil {
		return roleIds, fmt.Errorf("获取角色ids错误 %w", err)
	}
	return roleIds, nil
}
