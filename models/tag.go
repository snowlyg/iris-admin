package models

import (
	"github.com/snowlyg/blog/libs/logging"
	"github.com/snowlyg/easygorm"
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name     string     `gorm:"not null ;type:varchar(256)" json:"name" validate:"required,gte=0,lte=256" comment:"标签名称"`
	Articles []*Article `gorm:"many2many:article_tags;"`
}

// GetTag get tag
func GetTag(s *easygorm.Search) (*Tag, error) {
	t := &Tag{}
	err := easygorm.First(t, s)
	if err != nil {
		logging.Err.Errorf("get tag err: %+v", err)
		return t, err
	}

	return t, nil
}

/**
 * 通过 id 删除权限
 * @method DeleteTagById
 */
func DeleteTagById(id uint) error {
	t := &Tag{}
	if err := easygorm.DeleteById(t, id); err != nil {
		logging.Err.Errorf("del tag by id err: %+v", err)
		return err
	}
	return nil
}

// GetAllTags get all tags
func GetAllTags(s *easygorm.Search) ([]*Tag, int64, error) {
	var tags []*Tag
	count, err := easygorm.Paginate(&Tag{}, &tags, s)
	if err != nil {
		logging.Err.Errorf("get all tags err: %+v", err)
		return tags, count, err
	}

	return tags, count, nil
}

// CreateTag create tag
func (p *Tag) CreateTag() error {
	if err := easygorm.Create(p); err != nil {
		logging.Err.Errorf("create tag err: %+v", err)
		return err
	}
	return nil
}

// UpdateTagById update tag by id
func UpdateTagById(id uint, pj *Tag) error {
	if err := easygorm.Update(&Tag{}, pj, nil, id); err != nil {
		logging.Err.Errorf("update tag by id err: %+v", err)
		return err
	}
	return nil
}
