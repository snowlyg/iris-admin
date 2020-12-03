package seeder

import (
	"fmt"
	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/blog/app"
	"github.com/snowlyg/blog/libs/logging"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/azumads/faker"
	"github.com/jinzhu/configor"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/models"
	"github.com/snowlyg/easygorm"
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

	fpaths, _ := filepath.Glob(filepath.Join(libs.CWD(), "seeder", "data", "*.yml"))
	logging.Dbug.Debugf("数据填充YML文件路径：%+v\n", fpaths)
	if err := configor.Load(&Seeds, fpaths...); err != nil {
		logging.Err.Errorf("load config file err：%+v\n", err)
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
	app.Ser.App.Logger().Debugf("系统设置填充：%+v\n", configs)
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
		if err != nil && models.IsNotFound(err) {
			if perm.ID == 0 {
				perm = &models.Config{
					Model: gorm.Model{CreatedAt: time.Now()},
					Name:  m.Name,
					Value: m.Value,
				}
				if err := perm.CreateConfig(); err != nil {
					logging.Err.Errorf("seeder data create config err：%+v\n", err)
					return
				}
			}
		}
	}
}

// CreatePerms 新建权限
func CreatePerms() {
	ss := &easygorm.Search{
		Limit:  -1,
		Offset: -1,
	}
	var err error
	var perms []*models.Permission
	perms, _, err = models.GetAllPermissions(ss)
	if err != nil {
		logging.Err.Errorf("seeder data get perms err：%+v\n", err)
		return
	}
	var insertPerms []models.Permission
	for _, m := range Seeds.Perms {
		isInDB := func() bool {
			for _, perm := range perms {
				if perm.Name == m.Name && perm.Act == m.Act {
					return true
				}
			}
			return false
		}()

		if !isInDB {
			perm := models.Permission{
				Name:        m.Name,
				DisplayName: m.DisplayName,
				Description: m.Description,
				Act:         m.Act,
			}
			insertPerms = append(insertPerms, perm)
		}
	}

	logging.Dbug.Debugf("seeder data insert perms ：%+d\n", len(insertPerms))

	if len(insertPerms) == 0 {
		return
	}

	if err := models.CreatePermission(&insertPerms); err != nil {
		logging.Err.Errorf("seeder data create perms err：%+v\n", err)
		return
	}
}

// CreateAdminRole 新建管理角色
func CreateAdminRole() {
	var permIds []uint
	ss := &easygorm.Search{
		Limit:  -1,
		Offset: -1,
	}

	perms, _, err := models.GetAllPermissions(ss)
	if err != nil {
		logging.Err.Errorf("seeder data  get perms err：%+v\n", err)
		return
	}

	for _, perm := range perms {
		permIds = append(permIds, perm.ID)
	}

	s := &easygorm.Search{
		Fields: []*easygorm.Field{
			{
				Key:       "name",
				Condition: "=",
				Value:     libs.Config.Admin.RoleName,
			},
		},
	}
	role := &models.Role{}
	role, err = models.GetRole(s)
	if err != nil && models.IsNotFound(err) {
		if role.ID == 0 {
			role = &models.Role{
				Name:        libs.Config.Admin.RoleName,
				DisplayName: libs.Config.Admin.RoleDisplayName,
				Description: libs.Config.Admin.RoleDisplayName,
				Model:       gorm.Model{CreatedAt: time.Now()},
			}
			role.PermIds = permIds
			if err := role.CreateRole(); err != nil {
				logging.Err.Errorf("seeder data create role err：%+v\n", err)
				return
			}
		} else {
			role.PermIds = permIds
			if err := models.UpdateRole(role.ID, role); err != nil {
				logging.Err.Errorf("seeder data  update role err：%+v\n", err)
				return
			}
		}
	}

	logging.Dbug.Debugf("填充角色数据：%+v\n", role)
	logging.Dbug.Debugf("填充角色权限：%+v\n", role.PermIds)
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
		logging.Err.Errorf("seeder data  get user err：%+v\n", err)
		return
	}

	var roleIds []uint
	ss := &easygorm.Search{
		Limit:  -1,
		Offset: -1,
	}
	var roles []*models.Role
	roles, _, err = models.GetAllRoles(ss)
	if err != nil {
		logging.Err.Errorf("seeder data  get roles err：%+v\n", err)
		return
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
			logging.Err.Errorf("seeder data  create admin  err：%+v\n", err)
			return
		}
	} else {
		admin.Password = libs.Config.Admin.Pwd
		if err := models.UpdateUserById(admin.ID, admin); err != nil {
			app.Ser.App.Logger().Errorf("seeder data  update admin  err：%+v\n", err)
			return
		}
	}

	logging.Dbug.Debugf("管理员密码：%s\n", libs.Config.Admin.Pwd)
	logging.Dbug.Debugf("填充管理员数据：%+v", admin)

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
		logging.Err.Errorf("seeder data  auto migrate  err：%+v\n", err)
		return
	}

}

// PathName
type PathName struct {
	Name   string
	Path   string
	Method string
}

// 获取路由信息
func GetRoutes() []*models.Permission {
	var rrs []*models.Permission
	names := getPathNames(app.Ser.App.GetRoutesReadOnly())
	logging.Dbug.Debugf("路由权限集合：%v", names)
	logging.Dbug.Debugf("Iris App ：%v", app.Ser.App)
	for _, pathName := range names {
		if !isPermRoute(pathName.Name) {
			rr := &models.Permission{Name: pathName.Path, DisplayName: pathName.Name, Description: pathName.Name, Act: pathName.Method}
			rrs = append(rrs, rr)
		}
	}
	return rrs
}

// getPathNames
func getPathNames(routeReadOnly []context.RouteReadOnly) []*PathName {
	var pns []*PathName
	app.Ser.App.Logger().Debugf("routeReadOnly：%v", routeReadOnly)
	for _, s := range routeReadOnly {
		pn := &PathName{
			Name:   s.Name(),
			Path:   s.Path(),
			Method: s.Method(),
		}
		pns = append(pns, pn)
	}
	return pns
}

// 过滤非必要权限
func isPermRoute(name string) bool {
	exceptRouteName := []string{"OPTIONS", "GET", "POST", "HEAD", "PUT", "PATCH", "payload"}
	for _, er := range exceptRouteName {
		if strings.Contains(name, er) {
			return true
		}
	}
	return false
}
