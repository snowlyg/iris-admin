package models

import (
	"fmt"
	"github.com/snowlyg/blog/libs"
	"time"

	"github.com/fatih/color"
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

// GetTag get tag
func GetTag(s *Search) (*Tag, error) {
	t := NewTag()
	err := Found(s).First(t).Error
	if !IsNotFound(err) {
		return t, err
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
	if err := libs.Db.Delete(t).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteTagByIdError:%s \n", err))
		return err
	}
	return nil
}

// GetAllTags get all tags
func GetAllTags(s *Search) ([]*Tag, int64, error) {
	var tags []*Tag
	var count int64
	all := GetAll(&Tag{}, s)
	if err := all.Count(&count).Error; err != nil {
		return nil, count, err
	}
	all = all.Scopes(Paginate(s.Offset, s.Limit), Relation(s.Relations))
	if err := all.Find(&tags).Error; err != nil {
		return nil, count, err
	}

	return tags, count, nil
}

// CreateTag create tag
func (p *Tag) CreateTag() error {
	if err := libs.Db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTagById update tag by id
func UpdateTagById(id uint, pj *Tag) error {
	if err := Update(&Tag{}, pj, id); err != nil {
		return err
	}
	return nil
}
