package perm

import (
	"errors"
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Paginate
func Paginate(db *gorm.DB, req ReqPaginate) (iris.Map, error) {
	var count int64
	perms := []*Response{}
	db = db.Model(&Permission{})
	if req.Name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%s%%", req.Name))
	}
	err := db.Count(&count).Error
	if err != nil {
		g.ZAPLOG.Error("获取权限总数失败", zap.String("错误", err.Error()))
		return nil, err
	}
	err = db.Scopes(database.PaginateScope(req.Page, req.PageSize, req.Sort, req.OrderBy)).Find(&perms).Error
	if err != nil {
		g.ZAPLOG.Error("获取权限分页数据失败", zap.String("错误", err.Error()))
		return nil, err
	}
	list := iris.Map{"items": perms, "total": count, "limit": req.PageSize}
	return list, nil
}

// FindByNameAndAct
// db *gorm.DB
// name 名称
// act 方法
// ids 当 ids 的 len = 1 ，排除次 id 数据
func FindByNameAndAct(db *gorm.DB, name, act string, ids ...uint) (Response, error) {
	perm := Response{}
	db = db.Model(&Permission{}).Where("name = ?", name).Where("act = ?", act)
	if len(ids) == 1 {
		db.Where("id != ?", ids[0])
	}
	err := db.First(&perm).Error
	if err != nil {
		g.ZAPLOG.Error("根据名称和方法获取权限数据失败", zap.String("错误", err.Error()))
		return perm, err
	}
	return perm, nil
}

// Create
func Create(db *gorm.DB, req Request) (uint, error) {
	perm := Permission{BasePermission: req.BasePermission}
	if !checkNameAndAct(req) {
		return perm.ID, fmt.Errorf("权限[%s-%s]已存在", req.Name, req.Act)
	}
	err := db.Model(&Permission{}).Create(&perm).Error
	if err != nil {
		g.ZAPLOG.Error("添加权限失败", zap.String("错误", err.Error()))
		return perm.ID, err
	}
	return perm.ID, nil
}

// CreatenInBatches
func CreatenInBatches(db *gorm.DB, perms []Permission) error {
	err := db.Model(&Permission{}).CreateInBatches(&perms, 500).Error
	if err != nil {
		g.ZAPLOG.Error("添加权限失败", zap.String("错误", err.Error()))
		return err
	}
	return nil
}

// Update
func Update(db *gorm.DB, id uint, req Request) error {
	if !checkNameAndAct(req, id) {
		return fmt.Errorf("权限[%s-%s]已存在", req.Name, req.Act)
	}
	perm := Permission{BasePermission: req.BasePermission}
	err := db.Model(&Permission{}).Where("id = ?", id).Updates(&perm).Error
	if err != nil {
		g.ZAPLOG.Error("更新权限失败", zap.String("错误", err.Error()))
		return err
	}
	return nil
}

// checkNameAndAct
func checkNameAndAct(req Request, ids ...uint) bool {
	_, err := FindByNameAndAct(database.Instance(), req.Name, req.Act, ids...)
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// FindById
func FindById(db *gorm.DB, id uint) (Response, error) {
	res := Response{}
	err := db.Model(&Permission{}).Where("id = ?", id).First(&res).Error
	if err != nil {
		g.ZAPLOG.Error("获取权限失败", zap.String("错误", err.Error()))
		return res, err
	}
	return res, nil
}

// DeleteById
func DeleteById(db *gorm.DB, id uint) error {
	err := db.Unscoped().Delete(&Permission{}, id).Error
	if err != nil {
		g.ZAPLOG.Error("删除权限失败", zap.String("错误", err.Error()))
		return err
	}
	return nil
}

// GetPermsForRole
func GetPermsForRole() ([][]string, error) {
	var permsForRoles [][]string
	perms := []Permission{}
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
