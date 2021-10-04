package perm

import (
	"errors"
	"fmt"

	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreatenInBatches 批量加入
func CreatenInBatches(db *gorm.DB, perms PermCollection) error {
	err := db.Model(&Permission{}).CreateInBatches(&perms, 500).Error
	if err != nil {
		zap_server.ZAPLOG.Error("添加权限失败", zap.String("错误:", err.Error()))
		return err
	}
	return nil
}

// CheckNameAndAct 检测权限是否存在
func CheckNameAndAct(scopes ...func(db *gorm.DB) *gorm.DB) bool {
	perm := &Response{}
	err := perm.First(database.Instance(), scopes...)
	return errors.Is(err, gorm.ErrRecordNotFound)
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
