package models

import (
	"github.com/snowlyg/blog/libs/logging"
	"github.com/snowlyg/easygorm"
	"gorm.io/gorm"
)

type Doc struct {
	gorm.Model
	Name       string `gorm:"not null ;type:varchar(256)" json:"name" validate:"required,gte=0,lte=256" comment:"文档名称"`
	ChapterMun int64  `gorm:"not null ;default(0)" json:"chapter_mun" comment:"章节数量"`
	Chapters   []*Chapter
}

// GetDoc get doc
func GetDoc(search *easygorm.Search) (*Doc, error) {
	t := &Doc{}
	err := easygorm.First(t, search)
	if err != nil {
		logging.Err.Errorf("get doc err: %+v", err)
		return t, err
	}
	return t, nil
}

// GetDocById get doc
func GetDocById(id uint) (*Doc, error) {
	t := &Doc{}
	err := easygorm.FindById(t, id)
	if err != nil {
		logging.Err.Errorf("get doc by id err: %+v", err)
		return t, err
	}

	return t, nil
}

// GetDocCount get doc count
func GetDocCount() (int64, error) {
	var count int64
	t := &Doc{}
	err := easygorm.Egm.Db.Model(t).Count(&count).Error
	if err != nil {
		logging.Err.Errorf("get doc count err: %+v", err)
		return count, err
	}
	return count, nil
}

// DeleteDocById del doc by id
func DeleteDocById(id uint) error {
	t := &Doc{}
	if err := easygorm.DeleteById(t, id); err != nil {
		logging.Err.Errorf("del doc by id err: %+v", err)
		return err
	}
	return nil
}

// GetAllDocs get all docs
func GetAllDocs(s *easygorm.Search) ([]*Doc, int64, error) {
	var docs []*Doc
	count, err := easygorm.Paginate(&Doc{}, &docs, s)
	if err != nil {
		logging.Err.Errorf("get all docs err: %+v", err)
		return nil, count, err
	}

	return docs, count, nil
}

// CreateDoc create doc
func (p *Doc) CreateDoc() error {
	if err := easygorm.Create(p); err != nil {
		logging.Err.Errorf("create doc err: %+v", err)
		return err
	}
	return nil
}

// UpdateDocById update doc by id
func UpdateDocById(id uint, np *Doc, fileds []interface{}) error {
	if err := easygorm.Update(&Doc{}, np, fileds, id); err != nil {
		logging.Err.Errorf("update doc by id err: %+v", err)
		return err
	}
	return nil
}
