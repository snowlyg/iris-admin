package models

import (
	"errors"
	"fmt"
	"github.com/snowlyg/easygorm"
	"time"

	"github.com/fatih/color"
	"gorm.io/gorm"
)

type Doc struct {
	gorm.Model
	Name     string `gorm:"not null ;type:varchar(256)" json:"name" validate:"required,gte=0,lte=256" comment:"文档名称"`
	Chapters []*Chapter
}

func NewDoc() *Doc {
	return &Doc{
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// GetDoc get doc
func GetDoc(search *easygorm.Search) (*Doc, error) {
	t := NewDoc()
	err := easygorm.First(t, search)
	if err != nil {
		return t, err
	}
	if t.ID == 0 {
		return t, errors.New("数据不存在")
	}
	return t, nil
}

// GetDocCount get doc count
func GetDocCount() (int64, error) {
	var count int64
	t := NewDoc()
	err := easygorm.Egm.Db.Model(t).Count(&count).Error
	if err != nil {
		return count, err
	}
	return count, nil
}

// DeleteDocById del doc by id
func DeleteDocById(id uint) error {
	t := NewDoc()
	if err := easygorm.DeleteById(t, id); err != nil {
		color.Red(fmt.Sprintf("DeleteDocByIdError:%s \n", err))
		return err
	}
	return nil
}

// GetAllDocs get all docs
func GetAllDocs(s *easygorm.Search) ([]*Doc, int64, error) {
	var docs []*Doc
	count, err := easygorm.Paginate(&Doc{}, &docs, s)
	if err != nil {
		return nil, count, err
	}

	return docs, count, nil
}

// CreateDoc create doc
func (p *Doc) CreateDoc() error {
	if err := easygorm.Create(p); err != nil {
		return err
	}
	return nil
}

// UpdateDocById update doc by id
func UpdateDocById(id uint, np *Doc) error {
	if err := easygorm.Update(&Doc{}, np, nil, id); err != nil {
		return err
	}
	return nil
}
