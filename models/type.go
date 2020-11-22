package models

import (
	"fmt"
	"github.com/snowlyg/easygorm"
	"time"

	"github.com/fatih/color"
	"gorm.io/gorm"
)

type Type struct {
	gorm.Model
	Name     string `gorm:"not null ;type:varchar(256)" json:"name" validate:"required,gte=0,lte=256" comment:"分类名称"`
	Articles []*Article
}

type ListType struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
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
func GetType(search *easygorm.Search) (*Type, error) {
	t := NewType()
	err := easygorm.First(t, search)
	if err != nil {
		return t, err
	}
	return t, nil
}

// GetTypeById get type
func GetTypeById(id uint) (*Type, error) {
	t := NewType()
	err := easygorm.FindById(t, id)
	if err != nil {
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
	if err := easygorm.DeleteById(t, id); err != nil {
		color.Red(fmt.Sprintf("DeleteTypeByIdError:%s \n", err))
		return err
	}
	return nil
}

// GetAllTypes get all types
func GetAllTypes(s *easygorm.Search) ([]*ListType, int64, error) {
	var types []*ListType

	count, err := easygorm.Paginate(&Type{}, &types, s)
	if err != nil {
		return nil, count, err
	}

	return types, count, nil
}

// CreateType create type
func (p *Type) CreateType() error {
	if err := easygorm.Create(p); err != nil {
		return err
	}
	return nil
}

// UpdateTypeById update type by id
func UpdateTypeById(id uint, np *Type) error {
	if err := easygorm.Update(&Type{}, np, nil, id); err != nil {
		return err
	}
	return nil
}
