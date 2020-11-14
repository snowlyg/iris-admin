package models

import (
	"fmt"
	"github.com/snowlyg/blog/libs/easygorm"
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
func GetTag(s *easygorm.Search) (*Tag, error) {
	t := NewTag()
	err := easygorm.Found(s).First(t).Error
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
	if err := easygorm.DeleteById(t, id); err != nil {
		color.Red(fmt.Sprintf("DeleteTagByIdError:%s \n", err))
		return err
	}
	return nil
}

// GetAllTags get all tags
func GetAllTags(s *easygorm.Search) ([]*Tag, int64, error) {
	var tags []*Tag
	db, count, err := easygorm.Paginate(&Tag{}, s)
	if err != nil {
		return tags, count, err
	}
	if err := db.Find(&tags).Error; err != nil {
		return tags, count, err
	}

	return tags, count, nil
}

// CreateTag create tag
func (p *Tag) CreateTag() error {
	if err := easygorm.Create(p); err != nil {
		return err
	}
	return nil
}

// UpdateTagById update tag by id
func UpdateTagById(id uint, pj *Tag) error {
	if err := easygorm.Update(&Tag{}, pj, id); err != nil {
		return err
	}
	return nil
}
