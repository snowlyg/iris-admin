package models

import (
	"fmt"
	"github.com/snowlyg/blog/libs/database"
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

// GetType get type
func GetType(search *Search) (*Type, error) {
	t := NewType()
	err := Found(search).First(t).Error
	if !IsNotFound(err) {
		return t, err
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
	if err := database.Singleton().Db.Delete(t).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteTypeByIdError:%s \n", err))
		return err
	}
	return nil
}

// GetAllTypes get all types
func GetAllTypes(s *Search) ([]*Type, int64, error) {
	var types []*Type
	var count int64
	all := GetAll(&Type{}, s)
	if err := all.Count(&count).Error; err != nil {
		return nil, count, err
	}
	all = all.Scopes(Paginate(s.Offset, s.Limit), Relation(s.Relations))
	if err := all.Find(&types).Error; err != nil {
		return nil, count, err
	}

	return types, count, nil
}

// CreateType create type
func (p *Type) CreateType() error {
	if err := database.Singleton().Db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTypeById update type by id
func UpdateTypeById(id uint, np *Type) error {
	if err := Update(&Type{}, np, id); err != nil {
		return err
	}
	return nil
}
