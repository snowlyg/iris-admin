package models

import (
	"fmt"
	"gorm.io/gorm/clause"
	"time"

	"github.com/fatih/color"
	"github.com/snowlyg/IrisAdminApi/sysinit"
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name string `gorm:"not null ;type:varchar(256)" json:"name" validate:"required,gte=0,lte=256" comment:"标签名称"`

	Articles []*Article `gorm:"many2many:article_tags;"`
}

func NewTag() *Tag {
	return &Tag{
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

/**
 * 通过 id 获取 tag 记录
 * @method GetTagById
 * @param  {[type]}       tag  *Tag [description]
 */
func GetTagById(id uint, withRelation bool) (*Tag, error) {
	t := NewTag()
	get := sysinit.Db.Where("id = ?", id)
	if withRelation {
		get = get.Preload(clause.Associations)
	}
	err := IsNotFound(get.First(t).Error)
	if err != nil {
		return nil, err
	}

	return t, nil
}

/**
 * 通过 name 获取 tag 记录
 * @method GetTagByName
 * @param  {[type]}       tag  *Tag [description]
 */
func GetTagByName(name string) (*Tag, error) {
	t := NewTag()
	err := IsNotFound(sysinit.Db.Where("name = ?", name).First(t).Error)
	if err != nil {
		return nil, err
	}

	return t, nil
}

/**
 * 通过 id 删除权限
 * @method DeleteTagById
 */
func DeleteTagById(id uint) error {
	t := NewTag()
	t.ID = id
	if err := sysinit.Db.Delete(t).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteTagByIdError:%s \n", err))
		return err
	}
	return nil
}

/**
 * 获取所有的权限
 * @method GetAllTags
 * @param  {[type]} name string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllTags(name, orderBy string, offset, limit int) ([]*Tag, error) {
	var tags []*Tag
	if err := GetAll(&Tag{}, name, orderBy, offset, limit).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

/**
 * 创建
 * @method CreateTag
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (p *Tag) CreateTag() error {
	if err := sysinit.Db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

/**
 * 更新
 * @method UpdateTag
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func UpdateTagById(id uint, pj *Tag) error {
	p, err := GetUserById(id)
	if err != nil {
		return err
	}
	if err := Update(p, pj); err != nil {
		return err
	}
	return nil
}
