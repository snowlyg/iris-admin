package seeder

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/azumads/faker"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	"github.com/snowlyg/IrisAdminApi/server/config"
	"github.com/snowlyg/IrisAdminApi/server/models"
	"github.com/snowlyg/IrisAdminApi/server/sysinit"
)

var Fake *faker.Faker

var Seeds = struct {
	Perms []struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
		Description string `json:"description"`
		Act         string `json:"act"`
	}
}{}

func init() {
	Fake, _ = faker.New("en")
	Fake.Rand = rand.New(rand.NewSource(42))
	rand.Seed(time.Now().UnixNano())

	filepaths, _ := filepath.Glob(filepath.Join("seeder", "data", "*.yml"))
	if err := configor.Load(&Seeds, filepaths...); err != nil {
		panic(err)
	}
}

func Run() {

	AutoMigrates()

	fmt.Println(fmt.Sprintf("权限填充开始！！"))
	CreatePerms()
	fmt.Println(fmt.Sprintf("权限填充完成！！"))

	fmt.Println(fmt.Sprintf("管理角色填充开始！！"))
	CreateAdminRole()
	fmt.Println(fmt.Sprintf("管理角色填充完成！！"))

	fmt.Println(fmt.Sprintf("管理员填充开始！！"))
	CreateAdminUser()
	fmt.Println(fmt.Sprintf("管理员填充完成！！"))

}

// CreatePerms 新建菜单
func CreatePerms() {
	for _, m := range Seeds.Perms {
		menu := &models.Permission{
			Model:       gorm.Model{CreatedAt: time.Now()},
			Name:        m.Name,
			DisplayName: m.DisplayName,
			Description: m.Description,
			Act:         m.Act,
		}

		if err := menu.CreatePermission(); err != nil {
			panic(fmt.Sprintf("菜单填充错误：%v", err))
		}
	}

}

// CreateAdminRole 新建管理角色
func CreateAdminRole() {
	role := &models.Role{
		Name:        config.Config.Admin.RoleName,
		DisplayName: config.Config.Admin.RoleDisplayName,
		Description: config.Config.Admin.RoleDisplayName,
		Model:       gorm.Model{CreatedAt: time.Now()},
	}

	var permIds []uint
	perms, _ := models.GetAllPermissions("", "", 1, 20)
	for _, perm := range perms {
		permIds = append(permIds, perm.ID)
	}

	role.PermissionsIds = permIds
	if err := role.CreateRole(); err != nil {
		panic(fmt.Sprintf("管理角色填充错误：%v", err))
	}
}

// CreateAdminUser 新建管理员
func CreateAdminUser() {
	admin := &models.User{
		Username: "username",
		Name:     "超级管理员",
		Password: "123456",
		Model:    gorm.Model{CreatedAt: time.Now()},
	}

	if err := admin.CreateUser(); err != nil {
		panic(fmt.Sprintf("管理员填充错误：%v", err))
	}
}

/*
	AutoMigrates 重置数据表
	sysinit.Db.DropTableIfExists 删除存在数据表
	sysinit.Db.AutoMigrate 重建数据表
*/
func AutoMigrates() {
	sysinit.Db.DropTableIfExists("users", "permissions", "roles", "casbin_rule")

	sysinit.Db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&gormadapter.CasbinRule{},
	)
}
