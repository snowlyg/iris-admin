package models

import (
	"fmt"
	"time"

	"IrisAdminApi/database"
	"IrisAdminApi/validates"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
)

type Tenancy struct {
	gorm.Model

	Name string `gorm:"not null VARCHAR(191)"`
}

func NewTenancy(id uint, name string) *Tenancy {
	return &Tenancy{
		Model: gorm.Model{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: name,
	}
}

func NewTenancyByStruct(ru *validates.TenancyRequest) *Tenancy {
	return &Tenancy{
		Model: gorm.Model{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: ru.Name,
	}
}

func (u *Tenancy) GetTenancyByName() {
	IsNotFound(database.GetGdb().Where("name = ?", u.Name).First(u).Error)
}

func (u *Tenancy) GetTenancyById() {
	IsNotFound(database.GetGdb().Where("id = ?", u.ID).First(u).Error)
}

/**
 * 通过 id 删除用户
 * @method DeleteTenancyById
 */
func (u *Tenancy) DeleteTenancy() {
	if err := database.GetGdb().Delete(u).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteTenancyByIdErr:%s \n ", err))
	}
}

/**
 * 获取所有的账号
 * @method GetAllTenancy
 * @param  {[type]} name string [description]
 * @param  {[type]} username string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllTenancies(name, orderBy string, offset, limit int) []*Tenancy {
	var users []*Tenancy
	q := GetAll(name, orderBy, offset, limit)
	if err := q.Find(&users).Error; err != nil {
		color.Red(fmt.Sprintf("GetAllTenancyErr:%s \n ", err))
		return nil
	}
	return users
}

/**
 * 创建
 * @method CreateTenancy
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (u *Tenancy) CreateTenancy() {
	if err := database.GetGdb().Create(u).Error; err != nil {
		color.Red(fmt.Sprintf("CreateTenancyErr:%s \n ", err))
	}

	return
}

/**
 * 更新
 * @method UpdateTenancy
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (u *Tenancy) UpdateTenancy(uj *validates.TenancyRequest) {
	if err := Update(u, uj); err != nil {
		color.Red(fmt.Sprintf("UpdateTenancyErr:%s \n ", err))
	}
}
