package seeder

import (
	"fmt"
	"github.com/snowlyg/IrisAdminApi/server/libs"
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
		DisplayName string `json:"displayname"`
		Description string `json:"description"`
		Act         string `json:"act"`
	}
}{}

func init() {
	Fake, _ = faker.New("en")
	Fake.Rand = rand.New(rand.NewSource(42))
	rand.Seed(time.Now().UnixNano())

	filepaths, _ := filepath.Glob(filepath.Join(libs.CWD(), "../server/seeder/data", "*.yml"))
	if config.Config.Debug {
		fmt.Println(fmt.Sprintf("数据填充YML文件路径：%v", filepaths))
	}
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

// CreatePerms 新建权限
func CreatePerms() {
	if config.Config.Debug {
		fmt.Println(fmt.Sprintf("填充权限：%v", Seeds))
	}
	for _, m := range Seeds.Perms {
		perm := &models.Permission{
			Model:       gorm.Model{CreatedAt: time.Now()},
			Name:        m.Name,
			DisplayName: m.DisplayName,
			Description: m.Description,
			Act:         m.Act,
		}
		perm.GetPermissionByNameAct()
		if perm.ID == 0 {
			if err := perm.CreatePermission(); err != nil {
				panic(fmt.Sprintf("权限填充错误：%v", err))
			}
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

	role.GetRoleByName()
	if role.ID == 0 {
		var permIds []uint
		perms, _ := models.GetAllPermissions("", "", 0, 0)
		for _, perm := range perms {
			permIds = append(permIds, perm.ID)
		}

		role.PermIds = permIds
		if config.Config.Debug {
			fmt.Println(fmt.Sprintf("填充角色数据：%v", role))
		}
		if err := role.CreateRole(); err != nil {
			panic(fmt.Sprintf("管理角色填充错误：%v", err))
		}
	}

}

// CreateAdminUser 新建管理员
func CreateAdminUser() {
	admin := &models.User{
		Username: "username",
		Name:     "超级管理员",
		Password: "123456",
		Avatar:   "https://wx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIPbZRufW9zPiaGpfdXgU7icRL1licKEicYyOiace8QQsYVKvAgCrsJx1vggLAD2zJMeSXYcvMSkw9f4pw/132",
		Intro:    "超级弱鸡程序猿一枚！！！！",
		Model:    gorm.Model{CreatedAt: time.Now()},
	}
	var roleIds []uint
	roles, _ := models.GetAllRoles("", "", 0, 0)
	for _, role := range roles {
		roleIds = append(roleIds, role.ID)
	}
	admin.RoleIds = roleIds
	if config.Config.Debug {
		fmt.Println(fmt.Sprintf("填充管理员数据：%v", admin))
	}
	admin.GetUserByUsername()
	if admin.ID == 0 {
		if err := admin.CreateUser(); err != nil {
			panic(fmt.Sprintf("管理员填充错误：%v", err))
		}
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
