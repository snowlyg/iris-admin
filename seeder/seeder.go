package seeder

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/snowlyg/easygorm"
	"math/rand"
	"path/filepath"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/snowlyg/blog/libs"

	"github.com/azumads/faker"
	"github.com/jinzhu/configor"
	"github.com/snowlyg/blog/models"
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
	if libs.Config.Debug {
		fmt.Println(fmt.Sprintf("数据填充YML文件路径：%+v\n", filepaths))
	}
	if err := configor.Load(&Seeds, filepaths...); err != nil {
		logger.Println(err)
	}
}

func Run() {
	AutoMigrates()

	CreateConfigs()
	fmt.Println(fmt.Sprintf("系统设置填充完成！！"))
	CreatePerms()
	fmt.Println(fmt.Sprintf("权限填充完成！！"))
	CreateAdminRole()
	fmt.Println(fmt.Sprintf("管理角色填充完成！！"))
	CreateAdminUser()
	fmt.Println(fmt.Sprintf("管理员填充完成！！"))
}

func AddPerm() {
	fmt.Println(fmt.Sprintf("开始同步权限！！"))
	CreatePerms()
	CreateAdminRole()
	CreateAdminUser()
	fmt.Println(fmt.Sprintf("权限同步完成！！"))
}

// CreateConfigs 新建权限
func CreateConfigs() {
	configs := []*models.Config{
		{
			Name:  "imageHost",
			Value: "https://www.snowlyg.com",
		},
		{
			Name:  "beianhao",
			Value: "",
		},
	}

	if libs.Config.Debug {
		color.Yellow(fmt.Sprintf("系统设置填充：%+v\n", configs))
	}
	for _, m := range configs {
		s := &easygorm.Search{
			Fields: []*easygorm.Field{
				{
					Key:       "name",
					Condition: "=",
					Value:     m.Name,
				},
			},
		}
		perm, err := models.GetConfig(s)
		if !models.IsNotFound(err) {
			if perm.ID == 0 {
				perm = &models.Config{
					Model: gorm.Model{CreatedAt: time.Now()},
					Name:  m.Name,
					Value: m.Value,
				}
				if err := perm.CreateConfig(); err != nil {
					color.Red("系统设置填充错误：%+v\n", err)
				}
			}
		}
	}
}

// CreatePerms 新建权限
func CreatePerms() {
	if libs.Config.Debug {
		fmt.Println(fmt.Sprintf("填充权限：%+v\n", Seeds))
	}
	for _, m := range Seeds.Perms {
		s := &easygorm.Search{
			Fields: []*easygorm.Field{
				{
					Key:       "name",
					Condition: "=",
					Value:     m.Name,
				}, {
					Key:       "act",
					Condition: "=",
					Value:     m.Act,
				},
			},
		}
		perm, err := models.GetPermission(s)
		if err != nil && !models.IsNotFound(err) {
			if perm.ID == 0 {
				perm = &models.Permission{
					Model:       gorm.Model{CreatedAt: time.Now()},
					Name:        m.Name,
					DisplayName: m.DisplayName,
					Description: m.Description,
					Act:         m.Act,
				}
				if err := perm.CreatePermission(); err != nil {
					logger.Println(fmt.Sprintf("权限填充错误：%+v\n", err))
				}
				logger.Println(fmt.Sprintf("权限填充：%+v\n", perm))
			}
		}
	}
}

// CreateAdminRole 新建管理角色
func CreateAdminRole() {
	s := &easygorm.Search{
		Fields: []*easygorm.Field{
			{
				Key:       "name",
				Condition: "=",
				Value:     libs.Config.Admin.RoleName,
			},
		},
	}
	role, err := models.GetRole(s)
	if err != nil && !models.IsNotFound(err) {
		fmt.Println(fmt.Sprintf("角色获取失败：%+v\n", err))
	}
	var permIds []uint
	ss := &easygorm.Search{
		Limit:  -1,
		Offset: -1,
	}
	var perms []*models.Permission
	perms, _, err = models.GetAllPermissions(ss)
	if err != nil {
		fmt.Println(fmt.Sprintf("权限获取失败：%+v\n", err))
	}

	for _, perm := range perms {
		permIds = append(permIds, perm.ID)
	}

	if err == nil {
		if role.ID == 0 {
			role = &models.Role{
				Name:        libs.Config.Admin.RoleName,
				DisplayName: libs.Config.Admin.RoleDisplayName,
				Description: libs.Config.Admin.RoleDisplayName,
				Model:       gorm.Model{CreatedAt: time.Now()},
			}
			role.PermIds = permIds
			if err := role.CreateRole(); err != nil {
				logger.Println(fmt.Sprintf("create 管理角色填充错误：%+v\n", err))
			}
		} else {
			role.PermIds = permIds
			if err := models.UpdateRole(role.ID, role); err != nil {
				logger.Println(fmt.Sprintf("update 管理角色填充错误：%+v\n", err))
			}
		}
	}
	if libs.Config.Debug {
		fmt.Println(fmt.Sprintf("填充角色数据：%+v\n", role))
		fmt.Println(fmt.Sprintf("填充角色权限：%+v\n", role.PermIds))
	}

}

// CreateAdminUser 新建管理员
func CreateAdminUser() {
	s := &easygorm.Search{
		Fields: []*easygorm.Field{
			{
				Key:       "username",
				Condition: "=",
				Value:     libs.Config.Admin.UserName,
			},
		},
	}
	admin, err := models.GetUser(s)
	if err != nil && !models.IsNotFound(err) {
		fmt.Println(fmt.Sprintf("Get admin error：%+v\n", err))
	}

	var roleIds []uint
	ss := &easygorm.Search{
		Limit:  -1,
		Offset: -1,
	}
	roles, _, err := models.GetAllRoles(ss)
	if libs.Config.Debug {
		if err != nil {
			fmt.Println(fmt.Sprintf("角色获取失败：%+v\n", err))
		}
	}

	for _, role := range roles {
		roleIds = append(roleIds, role.ID)
	}
	admin.RoleIds = roleIds

	if admin.ID == 0 {
		admin = &models.User{
			Username: libs.Config.Admin.UserName,
			Name:     libs.Config.Admin.Name,
			Password: libs.Config.Admin.Pwd,
			Avatar:   "https://wx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIPbZRufW9zPiaGpfdXgU7icRL1licKEicYyOiace8QQsYVKvAgCrsJx1vggLAD2zJMeSXYcvMSkw9f4pw/132",
			Intro:    "超级弱鸡程序猿一枚！！！！",
			Model:    gorm.Model{CreatedAt: time.Now()},
		}
		admin.RoleIds = roleIds
		if err := admin.CreateUser(); err != nil {
			logger.Println(fmt.Sprintf("管理员填充错误：%+v\n", err))
		}
	} else {
		admin.Password = libs.Config.Admin.Pwd
		if err := models.UpdateUserById(admin.ID, admin); err != nil {
			logger.Println(fmt.Sprintf("管理员填充错误：%+v\n", err))
		}
	}

	if libs.Config.Debug {
		fmt.Println(fmt.Sprintf("管理员密码：%s\n", libs.Config.Admin.Pwd))
		fmt.Println(fmt.Sprintf("填充管理员数据：%+v", admin))
	}
}

/*
	AutoMigrates 重置数据表
	easygorm.Egm.Db.DropTableIfExists 删除存在数据表
*/
func AutoMigrates() {
	models.DropTables("")
	if err := easygorm.Migrate([]interface{}{
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.Article{},
		&models.Config{},
		&models.Tag{},
		&models.Type{},
		&models.Doc{},
		&models.Chapter{},
	}); err != nil {
		logger.Println(fmt.Sprintf("AutoMigrates 重置数据表错误：%+v\n", err))
	}

}
