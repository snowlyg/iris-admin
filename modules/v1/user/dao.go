package user

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/helper/arr"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/modules/v1/role"
	"github.com/snowlyg/iris-admin/server/casbin"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/multi"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Paginate(db *gorm.DB, req ReqPaginate) (map[string]interface{}, error) {
	var count int64
	users := []*Response{}
	db = db.Model(&User{})
	if len(req.Name) > 0 {
		db = db.Where("name LIKE ?", fmt.Sprintf("%s%%", req.Name))
	}
	err := db.Count(&count).Error
	if err != nil {
		g.ZAPLOG.Error("获取用户总数错误", zap.String("错误:", err.Error()))
		return nil, err
	}

	err = db.Scopes(database.PaginateScope(req.Page, req.PageSize, req.Sort, req.OrderBy)).
		Find(&users).Error
	if err != nil {
		g.ZAPLOG.Error("获取用户分页数据错误", zap.String("错误:", err.Error()))
		return nil, err
	}

	// 查询用户角色
	getRoles(db, users...)

	list := iris.Map{"items": users, "total": count, "pageSize": req.PageSize, "page": req.Page}
	return list, nil
}

// getRoles
func getRoles(db *gorm.DB, users ...*Response) {
	var roleIds []string
	userRoleIds := make(map[uint][]string, 10)
	if len(users) == 0 {
		return
	}
	for _, user := range users {
		userRoleId := casbin.GetRolesForUser(user.Id)
		userRoleIds[user.Id] = userRoleId
		roleIds = append(roleIds, userRoleId...)
	}

	roles, err := role.FindInId(db, roleIds)
	if err != nil {
		g.ZAPLOG.Error("get role get err ", zap.String("错误:", err.Error()))
	}

	for _, user := range users {
		for _, role := range roles {
			sRoleId := strconv.FormatInt(int64(role.Id), 10)
			if arr.InArrayS(userRoleIds[user.Id], sRoleId) {
				user.Roles = append(user.Roles, role.Name)
			}
		}
	}
}

func FindByUserName(db *gorm.DB, username string, ids ...uint) (Response, error) {
	user := Response{}
	db = db.Model(&User{}).Where("username = ?", username)
	if len(ids) == 1 {
		db.Where("id != ?", ids[0])
	}
	err := db.First(&user).Error
	if err != nil {
		g.ZAPLOG.Error("根据用户名查询用户错误", zap.String("用户名:", username), zap.Uints("ids:", ids), zap.String("错误:", err.Error()))
		return user, err
	}
	return user, nil
}

func FindPasswordByUserName(db *gorm.DB, username string, ids ...uint) (LoginResponse, error) {
	user := LoginResponse{}
	db = db.Model(&User{}).Select("id,password").Where("username = ?", username)
	if len(ids) == 1 {
		db.Where("id != ?", ids[0])
	}
	err := db.First(&user).Error
	if err != nil {
		g.ZAPLOG.Error("根据用户名查询用户错误", zap.String("用户名:", username), zap.Uints("ids:", ids), zap.String("错误:", err.Error()))
		return user, err
	}
	return user, nil
}

func Create(db *gorm.DB, req Request) (uint, error) {
	if _, err := FindByUserName(db, req.Username); !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, fmt.Errorf("用户名 %s 已经被使用", req.Username)
	}
	user := User{BaseUser: req.BaseUser, RoleIds: req.RoleIds}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		g.ZAPLOG.Error("密码加密错误", zap.String("错误:", err.Error()))
		return 0, err
	}

	g.ZAPLOG.Info("添加用户", zap.String("hash:", req.Password), zap.ByteString("hash:", hash))

	user.Password = string(hash)
	err = db.Model(&User{}).Create(&user).Error
	if err != nil {
		g.ZAPLOG.Error("添加用户错误", zap.String("错误:", err.Error()))
		return 0, err
	}

	if err := AddRoleForUser(&user); err != nil {
		g.ZAPLOG.Error("添加用户角色错误", zap.String("错误:", err.Error()))
		return 0, err
	}

	return user.ID, nil
}

func Update(db *gorm.DB, id uint, req Request) error {
	if b, err := IsAdminUser(db, id); err != nil {
		return err
	} else if b {
		return errors.New("不能编辑超级管理员")
	}
	if _, err := FindByUserName(db, req.Username, id); !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	user := User{BaseUser: req.BaseUser}
	err := db.Model(&User{}).Where("id = ?", id).Updates(&user).Error
	if err != nil {
		g.ZAPLOG.Error("更新用户错误", zap.String("错误:", err.Error()))
		return err
	}

	if err := AddRoleForUser(&user); err != nil {
		g.ZAPLOG.Error("添加用户角色错误", zap.String("错误:", err.Error()))
		return err
	}

	return nil
}

func IsAdminUser(db *gorm.DB, id uint) (bool, error) {
	user, err := FindById(db, id)
	if err != nil {
		return false, err
	}
	return arr.InArrayS(user.Roles, role.GetAdminRoleName()), nil
}

func FindById(db *gorm.DB, id uint) (Response, error) {
	user := Response{}
	err := db.Model(&User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		g.ZAPLOG.Error("find user err ", zap.String("错误:", err.Error()))
		return user, err
	}
	getRoles(db, &user)
	return user, nil
}

func DeleteById(db *gorm.DB, id uint) error {
	err := db.Unscoped().Delete(&User{}, id).Error
	if err != nil {
		g.ZAPLOG.Error("delete user by id get  err ", zap.String("错误:", err.Error()))
		return err
	}
	return nil
}

// AddRoleForUser add roles for user
func AddRoleForUser(user *User) error {
	userId := strconv.FormatUint(uint64(user.ID), 10)
	oldRoleIds, err := casbin.Instance().GetRolesForUser(userId)
	if err != nil {
		g.ZAPLOG.Error("获取用户角色错误", zap.String("错误:", err.Error()))
		return err
	}

	if len(oldRoleIds) > 0 {
		if _, err := casbin.Instance().DeleteRolesForUser(userId); err != nil {
			g.ZAPLOG.Error("添加角色到用户错误", zap.String("错误:", err.Error()))
			return err
		}
	}
	if len(user.RoleIds) == 0 {
		return nil
	}

	var roleIds []string
	for _, userRoleId := range user.RoleIds {
		roleIds = append(roleIds, strconv.FormatUint(uint64(userRoleId), 10))
	}

	if _, err := casbin.Instance().AddRolesForUser(userId, roleIds); err != nil {
		g.ZAPLOG.Error("添加角色到用户错误", zap.String("错误:", err.Error()))
		return err
	}

	return nil
}

// DelToken 删除token
func DelToken(token string) error {
	err := multi.AuthDriver.DelUserTokenCache(token)
	if err != nil {
		g.ZAPLOG.Error("del token", zap.Any("err", err))
		return fmt.Errorf("del token %w", err)
	}
	return nil
}

// CleanToken 清空 token
func CleanToken(authorityType int, userId string) error {
	err := multi.AuthDriver.CleanUserTokenCache(authorityType, userId)
	if err != nil {
		g.ZAPLOG.Error("clean token", zap.Any("err", err))
		return fmt.Errorf("clean token %w", err)
	}
	return nil
}
