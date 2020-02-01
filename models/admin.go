package models

import (
	"fmt"
	"time"

	"IrisAdminApi/database"
	"IrisAdminApi/validates"
	"github.com/fatih/color"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/gorm"
)

type Admin struct {
	gorm.Model

	Name      string `gorm:"not null VARCHAR(191)"`
	Adminname string `gorm:"unique;VARCHAR(191)"`
	Password  string `gorm:"not null VARCHAR(191)"`
}

func NewAdmin(id uint, username string) *Admin {
	return &Admin{
		Model: gorm.Model{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Adminname: username,
	}
}

func NewAdminByStruct(ru *validates.CreateUpdateAdminRequest) *Admin {
	return &Admin{
		Model: gorm.Model{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Adminname: ru.Adminname,
		Name:      ru.Name,
		Password:  HashPassword(ru.Password),
	}
}

func (u *Admin) GetAdminByAdminname() {
	IsNotFound(database.GetGdb().Where("username = ?", u.Adminname).First(u).Error)
}

func (u *Admin) GetAdminById() {
	IsNotFound(database.GetGdb().Where("id = ?", u.ID).First(u).Error)
}

/**
 * 通过 id 删除用户
 * @method DeleteAdminById
 */
func (u *Admin) DeleteAdmin() {
	if err := database.GetGdb().Delete(u).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteAdminByIdErr:%s \n ", err))
	}
}

/**
 * 获取所有的账号
 * @method GetAllAdmin
 * @param  {[type]} name string [description]
 * @param  {[type]} username string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllAdmins(name, orderBy string, offset, limit int) []*Admin {
	var users []*Admin
	q := GetAll(name, orderBy, offset, limit)
	if err := q.Find(&users).Error; err != nil {
		color.Red(fmt.Sprintf("GetAllAdminErr:%s \n ", err))
		return nil
	}
	return users
}

/**
 * 创建
 * @method CreateAdmin
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (u *Admin) CreateAdmin(aul *validates.CreateUpdateAdminRequest) {
	u.Password = HashPassword(aul.Password)
	if err := database.GetGdb().Create(u).Error; err != nil {
		color.Red(fmt.Sprintf("CreateAdminErr:%s \n ", err))
	}

	return
}

/**
 * 更新
 * @method UpdateAdmin
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (u *Admin) UpdateAdmin(uj *validates.CreateUpdateAdminRequest) {
	uj.Password = HashPassword(uj.Password)
	if err := database.Update(u, uj); err != nil {
		color.Red(fmt.Sprintf("UpdateAdminErr:%s \n ", err))
	}
}

/**
 * 判断用户是否登录
 * @method CheckLogin
 * @param  {[type]}  id       int    [description]
 * @param  {[type]}  password string [description]
 */
func (u *Admin) CheckLogin(password string) (*Token, bool, string) {
	if u.ID == 0 {
		return nil, false, "用户不存在"
	} else {
		if ok := bcrypt.Match(password, u.Password); ok {
			token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"exp": time.Now().Add(time.Hour * time.Duration(1)).Unix(),
				"iat": time.Now().Unix(),
			})
			tokenString, _ := token.SignedString([]byte("HS2JDFKhu7Y1av7b"))

			oauthToken := new(OauthToken)
			oauthToken.Token = tokenString
			oauthToken.UserId = u.ID
			oauthToken.Secret = "secret"
			oauthToken.Revoked = false
			oauthToken.ExpressIn = time.Now().Add(time.Hour * time.Duration(1)).Unix()
			oauthToken.CreatedAt = time.Now()

			response := oauthToken.OauthTokenCreate()

			return response, true, "登陆成功"
		} else {
			return nil, false, "用户名或密码错误"
		}
	}
}

/**
* 用户退出登陆
* @method AdminAdminLogout
* @param  {[type]} ids string [description]
 */
func AdminAdminLogout(userId uint) bool {
	ot := OauthToken{}
	ot.UpdateOauthTokenByUserId(userId)
	return ot.Revoked
}
