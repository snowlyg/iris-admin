package models

import (
	"fmt"
	"github.com/snowlyg/blog/libs"
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

// GetDocTableName
func GetDocTableName() string {
	return fmt.Sprintf("%s%s", libs.Config.DB.Prefix, "docs")
}

// GetDoc get doc
func GetDoc(search *Search) (*Doc, error) {
	t := NewDoc()
	err := Found(search).First(t).Error
	if !IsNotFound(err) {
		return t, err
	}
	return t, nil
}

// GetDocCount get doc count
func GetDocCount() (int64, error) {
	var count int64
	t := NewDoc()
	err := libs.Db.Model(t).Count(&count).Error
	if err != nil {
		return count, err
	}
	return count, nil
}

// DeleteDocById del doc by id
func DeleteDocById(id uint) error {
	t := NewDoc()
	t.ID = id
	if err := libs.Db.Delete(t).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteDocByIdError:%s \n", err))
		return err
	}
	return nil
}

// GetAllDocs get all docs
func GetAllDocs(s *Search) ([]*Doc, int64, error) {
	var docs []*Doc
	var count int64
	all := GetAll(&Doc{}, s)
	if err := all.Count(&count).Error; err != nil {
		return nil, count, err
	}

	all = all.Scopes(Relation(s.Relations))

	if err := all.Find(&docs).Error; err != nil {
		return nil, count, err
	}

	return docs, count, nil
}

// CreateDoc create doc
func (p *Doc) CreateDoc() error {
	if err := libs.Db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

// UpdateDocById update doc by id
func UpdateDocById(id uint, np *Doc) error {
	if err := Update(&Doc{}, np, id); err != nil {
		return err
	}
	return nil
}
