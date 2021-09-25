package role

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/casbin"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/scope"
	customZap "github.com/snowlyg/iris-admin/server/zap"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const adminRoleName = "admin"

func GetAdminRoleName() string {
	return adminRoleName
}

// Paginate
func Paginate(db *gorm.DB, req ReqPaginate) (map[string]interface{}, error) {
	var count int64
	var roles []*Response
	db = db.Model(&Role{})
	if req.Name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%s%%", req.Name))
	}
	err := db.Count(&count).Error
	if err != nil {
		g.ZAPLOG.Error("获取角色总数错误", zap.String("错误:", err.Error()))
		return nil, err
	}
	err = db.Scopes(scope.PaginateScope(req.Page, req.PageSize, req.Sort, req.OrderBy)).Find(&roles).Error
	if err != nil {
		g.ZAPLOG.Error("获取角色分页数据错误", zap.String("错误:", err.Error()))
		return nil, err
	}
	list := iris.Map{"items": roles, "total": count, "pageSize": req.PageSize, "page": req.Page}
	return list, nil
}

// FindByName
func FindByName(db *gorm.DB, name string, ids ...uint) (Response, error) {
	role := Response{}
	db = db.Model(&Role{}).Where("name = ?", name)
	if len(ids) == 1 {
		db.Where("id != ?", ids[0])
	}
	err := db.First(&role).Error
	if err != nil {
		g.ZAPLOG.Error("根据名称查询角色错误", zap.String("名称:", name), zap.String("错误:", err.Error()))
		return role, err
	}
	return role, nil
}

func Create(db *gorm.DB, req Request) (uint, error) {
	role := Role{BaseRole: req.BaseRole}
	_, err := FindByName(db, req.Name)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		g.ZAPLOG.Error("角色名称已经被使用")
		return 0, err
	}

	err = db.Create(&role).Error
	if err != nil {
		g.ZAPLOG.Error("create data err ", zap.String("错误:", err.Error()))
		return 0, err
	}

	err = AddPermForRole(role.ID, req.Perms)
	if err != nil {
		g.ZAPLOG.Error("添加权限到角色错误", zap.String("错误:", err.Error()))
		return 0, err
	}

	return role.ID, nil
}

func Update(db *gorm.DB, id uint, req Request) error {
	if b, err := IsAdminRole(db, id); err != nil {
		return err
	} else if b {
		return errors.New("不能编辑超级管理员")
	}
	_, err := FindByName(db, req.Name, id)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		g.ZAPLOG.Error("角色名称已经被使用")
		return err
	}
	role := Role{BaseRole: req.BaseRole}
	err = db.Model(&Role{}).Where("id = ?", id).Updates(&role).Error
	if err != nil {
		g.ZAPLOG.Error("更新角色错误", zap.String("错误:", err.Error()))
		return err
	}
	err = AddPermForRole(role.ID, req.Perms)
	if err != nil {
		g.ZAPLOG.Error("添加权限到角色错误", zap.String("错误:", err.Error()))
		return err
	}
	return nil
}

func IsAdminRole(db *gorm.DB, id uint) (bool, error) {
	role, err := FindById(db, id)
	if err != nil {
		return false, err
	}
	return role.Name == GetAdminRoleName(), nil
}

func FindById(db *gorm.DB, id uint) (Response, error) {
	role := Response{}
	err := db.Model(&Role{}).Where("id = ?", id).First(&role).Error
	if err != nil {
		g.ZAPLOG.Error("根据id查询角色错误", zap.String("错误:", err.Error()))
		return role, err
	}
	return role, nil
}

func DeleteById(db *gorm.DB, id uint) error {
	if b, err := IsAdminRole(db, id); err != nil {
		return err
	} else if b {
		return errors.New("不能删除超级管理员")
	}
	err := db.Unscoped().Delete(&Role{}, id).Error
	if err != nil {
		g.ZAPLOG.Error("删除角色错误", zap.String("错误:", err.Error()))
		return err
	}
	return nil
}

func FindInId(db *gorm.DB, ids []string) ([]*Response, error) {
	roles := []*Response{}
	err := db.Model(&Role{}).Where("id in ?", ids).Find(&roles).Error
	if err != nil {
		g.ZAPLOG.Error("通过ids查询角色错误", zap.String("错误:", err.Error()))
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
		g.ZAPLOG.Error("add policy err: %+v", zap.String("错误:", err.Error()))
		return err
	}

	if len(perms) == 0 {
		g.ZAPLOG.Debug("没有权限")
		return nil
	}
	var newPerms [][]string
	for _, perm := range perms {
		newPerms = append(newPerms, append([]string{roleId}, perm...))
	}
	g.ZAPLOG.Debug("添加权限到角色", customZap.Strings("新权限", newPerms))
	_, err = casbin.Instance().AddPolicies(newPerms)
	if err != nil {
		g.ZAPLOG.Error("add policy err: %+v", zap.String("错误:", err.Error()))
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
