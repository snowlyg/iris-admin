package models

import (
	"fmt"
	"github.com/snowlyg/blog/libs"
	"time"

	"github.com/fatih/color"
	"gorm.io/gorm"
)

type Type struct {
	gorm.Model
	Name     string `gorm:"not null ;type:varchar(256)" json:"name" validate:"required,gte=0,lte=256" comment:"分类名称"`
	Articles []*Article
}

func NewType() *Type {
	return &Type{
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

/**
 * 通过 id 获取 type 记录
 * @method GetTypeById
 * @param  {[type]}       type  *Type [description]
 */
func GetTypeById(id uint) (*Type, error) {
	t := NewType()
	err := IsNotFound(libs.Db.Where("id = ?", id).First(t).Error)
	if err != nil {
		return nil, err
	}

	return t, nil
}

/**
 * 通过 name 获取 type 记录
 * @method GetTypeByName
 * @param  {[type]}       type  *Type [description]
 */
func GetTypeByName(name string) (*Type, error) {
	t := NewType()
	err := IsNotFound(libs.Db.Where("name = ?", name).First(t).Error)
	if err != nil {
		return nil, err
	}

	return t, nil
}

/**
 * 通过 id 删除权限
 * @method DeleteTypeById
 */
func DeleteTypeById(id uint) error {
	t := NewType()
	t.ID = id
	if err := libs.Db.Delete(t).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteTypeByIdError:%s \n", err))
		return err
	}
	return nil
}

/**
 * 获取所有的权限
 * @method GetAllTypes
 * @param  {[type]} name string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllTypes(name, orderBy string, offset, limit int) ([]*Type, error) {
	var types []*Type
	all, _ := GetAll(&Type{}, name, orderBy, offset, limit)
	if err := all.Find(&types).Error; err != nil {
		return nil, err
	}

	return types, nil
}

/**
 * 创建
 * @method CreateType
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (p *Type) CreateType() error {
	if err := libs.Db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

/**
 * 更新
 * @method UpdateType
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func UpdateTypeById(id uint, np *Type) error {
	if err := Update(&Type{}, np, id); err != nil {
		return err
	}
	return nil
}
