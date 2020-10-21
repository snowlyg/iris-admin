package seeder

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/snowlyg/IrisAdminApi/libs"

	"github.com/azumads/faker"
	"github.com/jinzhu/configor"
	"github.com/snowlyg/IrisAdminApi/config"
	"github.com/snowlyg/IrisAdminApi/models"
	"gorm.io/gorm"
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

	filepaths, _ := filepath.Glob(filepath.Join(libs.CWD(), "seeder", "data", "*.yml"))
	if config.Config.Debug {
		fmt.Println(fmt.Sprintf("数据填充YML文件路径：%v", filepaths))
	}
	if err := configor.Load(&Seeds, filepaths...); err != nil {
		logger.Println(err)
	}
}

func Run() {
	AutoMigrates()
	CreatePerms()
	fmt.Println(fmt.Sprintf("权限填充完成！！"))
	CreateAdminRole()
	fmt.Println(fmt.Sprintf("管理角色填充完成！！"))
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
				logger.Println(fmt.Sprintf("权限填充错误：%v", err))
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

	var permIds []uint
	perms, err := models.GetAllPermissions("", "", -1, -1)
	if config.Config.Debug {
		if err != nil {
			fmt.Println(fmt.Sprintf("权限获取失败：%v", err))
		}
	}

	for _, perm := range perms {
		permIds = append(permIds, perm.ID)
	}

	role.PermIds = permIds

	role.GetRoleByName()
	if role.ID == 0 {
		if err := role.CreateRole(); err != nil {
			logger.Println(fmt.Sprintf("管理角色填充错误：%v", err))
		}
	} else {
		if err := role.UpdateRole(); err != nil {
			logger.Println(fmt.Sprintf("管理角色填充错误：%v", err))
		}
	}

	if config.Config.Debug {
		fmt.Println(fmt.Sprintf("填充角色数据：%v", role))
		fmt.Println(fmt.Sprintf("填充角色权限：%v", permIds))
	}

}

// CreateAdminUser 新建管理员
func CreateAdminUser() {

	admin, err := models.GetUserByUsername(config.Config.Admin.UserName)
	if err != nil {
		fmt.Println(fmt.Sprintf("Get admin error：%v", err))
	}
	password := config.Config.Admin.Pwd
	if admin.ID == 0 {
		admin = &models.User{
			Username: config.Config.Admin.UserName,
			Name:     config.Config.Admin.Name,
			Password: password,
			Avatar:   "https://wx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIPbZRufW9zPiaGpfdXgU7icRL1licKEicYyOiace8QQsYVKvAgCrsJx1vggLAD2zJMeSXYcvMSkw9f4pw/132",
			Intro:    "超级弱鸡程序猿一枚！！！！",
			Model:    gorm.Model{CreatedAt: time.Now()},
		}
		var roleIds []uint
		roles, err := models.GetAllRoles("", "", -1, -1)
		if config.Config.Debug {
			if err != nil {
				fmt.Println(fmt.Sprintf("角色获取失败：%v", err))
			}
		}

		for _, role := range roles {
			roleIds = append(roleIds, role.ID)
		}
		admin.RoleIds = roleIds
		if err := admin.CreateUser(); err != nil {
			logger.Println(fmt.Sprintf("管理员填充错误：%v", err))
		}
	} else {
		if err := admin.UpdateUser(password); err != nil {
			logger.Println(fmt.Sprintf("管理员填充错误：%v", err))
		}
	}

	if config.Config.Debug {
		fmt.Println(fmt.Sprintf("填充管理员数据：%v", admin))
	}
}

/*
	AutoMigrates 重置数据表
	sysinit.Db.DropTableIfExists 删除存在数据表
	sysinit.Db.AutoMigrate 重建数据表
*/
func AutoMigrates() {
	models.DropTables()
	models.Migrate()
}
