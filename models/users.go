package models

import (
	"fmt"
	"strconv"
	"time"

	"IrisAdminApi/transformer"
	"github.com/fatih/color"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model

	Name     string `gorm:"not null VARCHAR(191)"`
	Username string `gorm:"unique;VARCHAR(191)"`
	Password string `gorm:"not null VARCHAR(191)"`
}

type UserRequest struct {
	Username string `json:"username" validate:"required,gte=2,lte=50"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required,gte=2,lte=50"`
	RoleIds  []uint `json:"role_ids" validate:"required"`
}

/**
 * 校验用户登录
 * @method UserAdminCheckLogin
 * @param  {[type]}       username string [description]
 */
func UserAdminCheckLogin(username string) *User {
	user := new(User)
	IsNotFound(Db.Where("username = ?", username).First(user).Error)

	return user
}

/**
 * 通过 id 获取 user 记录
 * @method GetUserById
 * @param  {[type]}       user  *User [description]
 */
func GetUserById(id uint) *User {
	user := new(User)
	IsNotFound(Db.Where("id= ?", id).First(user).Error)
	return user
}

/**
 * 通过 username 获取 user 记录
 * @method GetUserByUserName
 * @param  {[type]}       user  *User [description]
 */
func GetUserByUserName(username string) *User {
	user := new(User)
	IsNotFound(Db.Where("username= ?", username).First(user).Error)
	return user
}

/**
 * 通过 id 删除用户
 * @method DeleteUserById
 */
func DeleteUserById(id uint) {
	u := &User{
		Model: gorm.Model{
			ID: id,
		},
	}
	if err := Db.Delete(u).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteUserByIdErr:%s \n ", err))
	}
}

/**
 * 获取所有的账号
 * @method GetAllUser
 * @param  {[type]} name string [description]
 * @param  {[type]} username string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllUsers(name, orderBy string, offset, limit int) []*User {
	var users []*User
	q := GetAll(name, orderBy, offset, limit)
	if err := q.Find(&users).Error; err != nil {
		color.Red(fmt.Sprintf("GetAllUserErr:%s \n ", err))
		return nil
	}
	return users
}

/**
 * 创建
 * @method CreateUser
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func CreateUser(aul *UserRequest) (user *User) {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(aul.Password, salt)

	user = &User{
		Username: aul.Username,
		Password: hash,
		Name:     aul.Name,
	}

	if err := Db.Create(user).Error; err != nil {
		color.Red(fmt.Sprintf("CreateUserErr:%s \n ", err))
	}

	if len(aul.RoleIds) > 0 {
		userId := strconv.FormatUint(uint64(user.ID), 10)
		if _, err = Enforcer.DeleteRolesForUser(userId); err != nil {
			color.Red(fmt.Sprintf("CreateUserErr:%s \n ", err))
		}

		for _, roleId := range aul.RoleIds {
			roleId := strconv.FormatUint(uint64(roleId), 10)
			if _, err = Enforcer.AddRoleForUser(userId, roleId); err != nil {
				color.Red(fmt.Sprintf("CreateUserErr:%s \n ", err))
			}
		}
	}

	return
}

/**
 * 更新
 * @method UpdateUser
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func UpdateUser(uj *UserRequest, id uint) *User {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(uj.Password, salt)

	user := &User{
		Model: gorm.Model{
			ID: id,
		},
		Password: hash,
	}

	if err := Db.Model(user).Updates(uj).Error; err != nil {
		color.Red(fmt.Sprintf("UpdateUserErr:%s \n ", err))
	}

	return user
}

/**
 * 判断用户是否登录
 * @method CheckLogin
 * @param  {[type]}  id       int    [description]
 * @param  {[type]}  password string [description]
 */
func CheckLogin(username, password string) (response Token, status bool, msg string) {
	user := UserAdminCheckLogin(username)
	if user.ID == 0 {
		msg = "用户不存在"
		return
	} else {
		if ok := bcrypt.Match(password, user.Password); ok {

			token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"exp": time.Now().Add(time.Hour * time.Duration(1)).Unix(),
				"iat": time.Now().Unix(),
			})
			tokenString, _ := token.SignedString([]byte("HS2JDFKhu7Y1av7b"))

			oauthToken := new(OauthToken)
			oauthToken.Token = tokenString
			oauthToken.UserId = user.ID
			oauthToken.Secret = "secret"
			oauthToken.Revoked = false
			oauthToken.ExpressIn = time.Now().Add(time.Hour * time.Duration(1)).Unix()
			oauthToken.CreatedAt = time.Now()

			response = oauthToken.OauthTokenCreate()
			status = true
			msg = "登陆成功"

			return

		} else {
			msg = "用户名或密码错误"
			return
		}
	}
}

/**
* 用户退出登陆
* @method UserAdminLogout
* @param  {[type]} ids string [description]
 */
func UserAdminLogout(userId uint) bool {
	ot := UpdateOauthTokenByUserId(userId)
	return ot.Revoked
}

/**
*创建系统管理员
*@param role_id uint
*@return   *models.AdminUserTranform api格式化后的数据格式
 */
func CreateSystemAdmin(roleId uint, rc *transformer.Conf) *User {
	aul := &UserRequest{
		Username: rc.TestData.UserName,
		Password: rc.TestData.Pwd,
		Name:     rc.TestData.Name,
		RoleIds:  []uint{roleId},
	}

	user := GetUserByUserName(aul.Username)
	if user.ID == 0 {
		return CreateUser(aul)
	} else {
		return user
	}
}

//func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
//	fmt.Println( fmt.Sprintf("AfterUpdate: %s", u.Name))
//	return
//}
