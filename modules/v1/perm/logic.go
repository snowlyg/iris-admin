package perm

import (
	"errors"
	"fmt"

	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/scope"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Create 添加
func Create(req *Request) (uint, error) {
	perm := &Permission{BasePermission: req.BasePermission}
	if !checkNameAndAct(NameScope(req.Name), ActScope(req.Act)) {
		return perm.ID, fmt.Errorf("权限[%s-%s]已存在", req.Name, req.Act)
	}
	return perm.Create()
}

// CreatenInBatches 批量加入
func CreatenInBatches(db *gorm.DB, perms PermCollection) error {
	err := db.Model(&Permission{}).CreateInBatches(&perms, 500).Error
	if err != nil {
		g.ZAPLOG.Error("添加权限失败", zap.String("错误:", err.Error()))
		return err
	}
	return nil
}

// Update
func Update(id uint, req *Request) error {
	if !checkNameAndAct(NameScope(req.Name), ActScope(req.Act), NeIdScope(id)) {
		return fmt.Errorf("权限[%s-%s]已存在", req.Name, req.Act)
	}
	perm := &Permission{BasePermission: req.BasePermission}
	err := perm.Update(scope.IdScope(id))
	if err != nil {
		return err
	}
	return nil
}

// checkNameAndAct 检测权限是否存在
func checkNameAndAct(scopes ...func(db *gorm.DB) *gorm.DB) bool {
	perm := &Response{}
	err := perm.First(scopes...)
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// DeleteById
func DeleteById(id uint) error {
	perm := Permission{}
	err := perm.Update(scope.IdScope(id))
	if err != nil {
		g.ZAPLOG.Error("删除权限失败", zap.String("错误:", err.Error()))
		return err
	}
	return nil
}

// GetPermsForRole
func GetPermsForRole() ([][]string, error) {
	var permsForRoles [][]string
	perms := PermCollection{}
	err := database.Instance().Model(&Permission{}).Find(&perms).Error
	if err != nil {
		return nil, fmt.Errorf("获取权限错误 %w", err)
	}
	for _, perm := range perms {
		permsForRole := []string{perm.Name, perm.Act}
		permsForRoles = append(permsForRoles, permsForRole)
	}
	return permsForRoles, nil
}
