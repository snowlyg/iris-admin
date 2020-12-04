package models

import (
	"github.com/snowlyg/blog/libs/logging"
	"github.com/snowlyg/easygorm"
	"time"

	"gorm.io/gorm"
)

type Type struct {
	gorm.Model
	Name     string `gorm:"not null ;type:varchar(256)" json:"name" validate:"required,gte=0,lte=256" comment:"分类名称"`
	Articles []*Article
}

type TypeInfo struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// GetType get type
func GetType(search *easygorm.Search) (*TypeInfo, error) {
	t := &TypeInfo{}
	err := easygorm.First(t, search)
	if err != nil {
		logging.Err.Errorf("get type err: %+v", err)
		return t, err
	}

	return t, nil
}

// GetTypeById get type
func GetTypeById(id uint) (*TypeInfo, error) {
	t := &TypeInfo{}
	err := easygorm.FindById(t, id)
	if err != nil {
		logging.Err.Errorf("get type by id err: %+v", err)
		return t, err
	}

	return t, nil
}

/**
 * 通过 id 删除权限
 * @method DeleteTypeById
 */
func DeleteTypeById(id uint) error {
	t := &TypeInfo{}
	if err := easygorm.DeleteById(t, id); err != nil {
		logging.Err.Errorf("del type by id err: %+v", err)
		return err
	}
	return nil
}

// GetAllTypes get all types
func GetAllTypes(s *easygorm.Search) ([]*TypeInfo, int64, error) {
	var types []*TypeInfo

	count, err := easygorm.Paginate(&Type{}, &types, s)
	if err != nil {
		logging.Err.Errorf("get all type err: %+v", err)
		return nil, count, err
	}

	return types, count, nil
}

// CreateType create type
func (p *Type) CreateType() error {
	if err := easygorm.Create(p); err != nil {
		logging.Err.Errorf("create type err: %+v", err)
		return err
	}
	return nil
}

// UpdateTypeById update type by id
func UpdateTypeById(id uint, np *Type) error {
	if err := easygorm.Update(&Type{}, np, nil, id); err != nil {
		logging.Err.Errorf("update type by id err: %+v", err)
		return err
	}
	return nil
}
