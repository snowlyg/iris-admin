package models

import (
	"fmt"
	"time"

	"IrisAdminApi/database"
	"IrisAdminApi/validates"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
)

type App struct {
	gorm.Model

	Name string `gorm:"not null VARCHAR(191)"`
}

func NewApp(id uint, name string) *App {
	return &App{
		Model: gorm.Model{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: name,
	}
}

func NewAppByStruct(ru *validates.AppRequest) *App {
	return &App{
		Model: gorm.Model{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: ru.Name,
	}
}

func (u *App) GetAppByName() {
	IsNotFound(database.GetGdb().Where("name = ?", u.Name).First(u).Error)
}

func (u *App) GetAppById() {
	IsNotFound(database.GetGdb().Where("id = ?", u.ID).First(u).Error)
}

/**
 * 通过 id 删除用户
 * @method DeleteAppById
 */
func (u *App) DeleteApp() {
	if err := database.GetGdb().Delete(u).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteAppByIdErr:%s \n ", err))
	}
}

/**
 * 获取所有的账号
 * @method GetAllApp
 * @param  {[type]} name string [description]
 * @param  {[type]} username string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllApps(name, orderBy string, offset, limit int) []*App {
	var users []*App
	q := GetAll(name, orderBy, offset, limit)
	if err := q.Find(&users).Error; err != nil {
		color.Red(fmt.Sprintf("GetAllAppErr:%s \n ", err))
		return nil
	}
	return users
}

/**
 * 创建
 * @method CreateApp
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (u *App) CreateApp() {
	if err := database.GetGdb().Create(u).Error; err != nil {
		color.Red(fmt.Sprintf("CreateAppErr:%s \n ", err))
	}

	return
}

/**
 * 更新
 * @method UpdateApp
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (u *App) UpdateApp(uj *validates.AppRequest) {
	if err := database.Update(u, uj); err != nil {
		color.Red(fmt.Sprintf("UpdateAppErr:%s \n ", err))
	}
}
