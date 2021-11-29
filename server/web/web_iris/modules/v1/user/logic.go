package user

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/snowlyg/helper/arr"
	"github.com/snowlyg/iris-admin/server/casbin"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/database/scope"
	"github.com/snowlyg/iris-admin/server/web/web_iris/modules/v1/role"
	"github.com/snowlyg/iris-admin/server/zap_server"
	multi "github.com/snowlyg/multi/iris"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrUserNameInvalide = errors.New("用户名名称已经被使用")

// getRoles
func getRoles(db *gorm.DB, users ...*Response) {
	var roleIds []uint
	userRoleIds := make(map[uint][]string, 10)
	if len(users) == 0 {
		return
	}
	for _, user := range users {
		user.ToString()
		userRoleId := casbin.GetRolesForUser(user.Id)
		userRoleIds[user.Id] = userRoleId
		for _, roleId := range userRoleId {
			id, err := strconv.ParseUint(roleId, 10, 64)
			if err != nil {
				continue
			}
			roleIds = append(roleIds, uint(id))
		}
	}

	roles, err := role.FindInId(db, roleIds)
	if err != nil {
		zap_server.ZAPLOG.Error("get role get err ", zap.String("错误:", err.Error()))
	}

	for _, user := range users {
		for _, role := range roles {
			sRoleId := strconv.FormatInt(int64(role.Id), 10)
			if arr.InArrayS(userRoleIds[user.Id], sRoleId) {
				user.Roles = append(user.Roles, role.DisplayName)
			}
		}
	}
}

// FindByName
func FindByUserName(scopes ...func(db *gorm.DB) *gorm.DB) (*Response, error) {
	user := &Response{}
	err := user.First(database.Instance(), scopes...)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func FindPasswordByUserName(db *gorm.DB, username string, ids ...uint) (LoginResponse, error) {
	user := LoginResponse{}
	db = db.Model(&User{}).Select("id,password").
		Where("username = ?", username)
	if len(ids) == 1 {
		db.Where("id != ?", ids[0])
	}
	err := db.First(&user).Error
	if err != nil {
		zap_server.ZAPLOG.Error("根据用户名查询用户错误", zap.String("用户名:", username), zap.Uints("ids:", ids), zap.String("错误:", err.Error()))
		return user, err
	}
	return user, nil
}

func Create(req *Request) (uint, error) {
	if _, err := FindByUserName(UserNameScope(req.Username)); !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, ErrUserNameInvalide
	}
	user := &User{BaseUser: req.BaseUser, RoleIds: req.RoleIds}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		zap_server.ZAPLOG.Error("密码加密错误", zap.String("错误:", err.Error()))
		return 0, err
	}

	zap_server.ZAPLOG.Info("添加用户", zap.String("hash:", req.Password), zap.ByteString("hash:", hash))

	user.Password = string(hash)
	id, err := orm.Create(database.Instance(), user)
	if err != nil {
		return 0, err
	}

	if err := AddRoleForUser(user); err != nil {
		zap_server.ZAPLOG.Error("添加用户角色错误", zap.String("错误:", err.Error()))
		return 0, err
	}

	return id, nil
}

func IsAdminUser(id uint) error {
	user := &Response{}
	err := user.First(database.Instance(), scope.IdScope(id))
	if err != nil {
		return err
	}
	if arr.InArrayS(user.Roles, role.GetAdminRoleName()) {
		return errors.New("不能操作超级管理员")
	}
	return nil
}

func FindById(db *gorm.DB, id uint) (Response, error) {
	user := Response{}
	err := db.Model(&User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		zap_server.ZAPLOG.Error("find user err ", zap.String("错误:", err.Error()))
		return user, err
	}

	getRoles(db, &user)

	return user, nil
}

// AddRoleForUser add roles for user
func AddRoleForUser(user *User) error {
	userId := strconv.FormatUint(uint64(user.ID), 10)
	oldRoleIds, err := casbin.Instance().GetRolesForUser(userId)
	if err != nil {
		zap_server.ZAPLOG.Error("获取用户角色错误", zap.String("错误:", err.Error()))
		return err
	}

	if len(oldRoleIds) > 0 {
		if _, err := casbin.Instance().DeleteRolesForUser(userId); err != nil {
			zap_server.ZAPLOG.Error("添加角色到用户错误", zap.String("错误:", err.Error()))
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
		zap_server.ZAPLOG.Error("添加角色到用户错误", zap.String("错误:", err.Error()))
		return err
	}

	return nil
}

// DelToken 删除token
func DelToken(token string) error {
	err := multi.AuthDriver.DelUserTokenCache(token)
	if err != nil {
		zap_server.ZAPLOG.Error("del token", zap.Any("err", err))
		return fmt.Errorf("del token %w", err)
	}
	return nil
}

// CleanToken 清空 token
func CleanToken(authorityType int, userId string) error {
	err := multi.AuthDriver.CleanUserTokenCache(authorityType, userId)
	if err != nil {
		zap_server.ZAPLOG.Error("clean token", zap.Any("err", err))
		return fmt.Errorf("clean token %w", err)
	}
	return nil
}

func UpdateAvatar(db *gorm.DB, id uint, avatar string) error {
	return nil
}
