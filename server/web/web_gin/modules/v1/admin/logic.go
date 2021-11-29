package admin

import (
	"errors"
	"strconv"

	"github.com/snowlyg/helper/arr"
	"github.com/snowlyg/iris-admin/server/casbin"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/database/scope"
	"github.com/snowlyg/iris-admin/server/web/web_iris/modules/v1/role"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrUserNameInvalide = errors.New("用户名名称已经被使用")

// getRoles
func getRoles(db *gorm.DB, admins ...*Response) {
	var roleIds []uint
	userRoleIds := make(map[uint][]string, 10)
	if len(admins) == 0 {
		return
	}
	for _, admin := range admins {
		admin.ToString()
		userRoleId := casbin.GetRolesForUser(admin.Id)
		userRoleIds[admin.Id] = userRoleId
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

	for _, admin := range admins {
		for _, role := range roles {
			sRoleId := strconv.FormatInt(int64(role.Id), 10)
			if arr.InArrayS(userRoleIds[admin.Id], sRoleId) {
				admin.AuthorityIds = append(admin.AuthorityIds, role.DisplayName)
			}
		}
	}
}

// FindByName
func FindByUserName(scopes ...func(db *gorm.DB) *gorm.DB) (*Response, error) {
	admin := &Response{}
	err := admin.First(database.Instance(), scopes...)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func FindPasswordByUserName(db *gorm.DB, username string, ids ...uint) (LoginResponse, error) {
	admin := LoginResponse{}
	db = db.Model(&Admin{}).Select("id,password").
		Where("username = ?", username)
	if len(ids) == 1 {
		db.Where("id != ?", ids[0])
	}
	err := db.First(&admin).Error
	if err != nil {
		zap_server.ZAPLOG.Error("根据用户名查询用户错误", zap.String("用户名:", username), zap.Uints("ids:", ids), zap.String("错误:", err.Error()))
		return admin, err
	}
	return admin, nil
}

func Create(req *Request) (uint, error) {
	if _, err := FindByUserName(UserNameScope(req.Username)); !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, ErrUserNameInvalide
	}
	admin := &Admin{BaseAdmin: req.BaseAdmin, AuthorityIds: req.AuthorityIds}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		zap_server.ZAPLOG.Error("密码加密错误", zap.String("错误:", err.Error()))
		return 0, err
	}

	zap_server.ZAPLOG.Info("添加用户", zap.String("hash:", req.Password), zap.ByteString("hash:", hash))

	admin.Password = string(hash)
	id, err := orm.Create(database.Instance(), admin)
	if err != nil {
		return 0, err
	}

	if err := AddRoleForUser(admin); err != nil {
		zap_server.ZAPLOG.Error("添加用户角色错误", zap.String("错误:", err.Error()))
		return 0, err
	}

	return id, nil
}

func IsAdminUser(id uint) error {
	admin := &Response{}
	err := admin.First(database.Instance(), scope.IdScope(id))
	if err != nil {
		return err
	}
	if arr.InArrayS(admin.AuthorityIds, role.GetAdminRoleName()) {
		return errors.New("不能操作超级管理员")
	}
	return nil
}

func FindById(db *gorm.DB, id uint) (Response, error) {
	admin := Response{}
	err := db.Model(&Admin{}).Where("id = ?", id).First(&admin).Error
	if err != nil {
		zap_server.ZAPLOG.Error("find admin err ", zap.String("错误:", err.Error()))
		return admin, err
	}

	getRoles(db, &admin)

	return admin, nil
}

// AddRoleForUser add roles for user
func AddRoleForUser(admin *Admin) error {
	userId := strconv.FormatUint(uint64(admin.ID), 10)
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
	if len(admin.AuthorityIds) == 0 {
		return nil
	}

	var roleIds []string
	for _, userRoleId := range admin.AuthorityIds {
		authId := strconv.FormatUint(uint64(userRoleId), 10)
		roleIds = append(roleIds, authId)
	}

	if _, err := casbin.Instance().AddRolesForUser(userId, roleIds); err != nil {
		zap_server.ZAPLOG.Error("添加角色到用户错误", zap.String("错误:", err.Error()))
		return err
	}

	return nil
}

func UpdateAvatar(db *gorm.DB, id uint, avatar string) error {
	return nil
}
