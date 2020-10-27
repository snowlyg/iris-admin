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

/**
 * 通过 id 获取 type 记录
 * @method GetDocById
 * @param  {[type]}       type  *Doc [description]
 */
func GetDocById(id uint, relation string) (*Doc, error) {
	t := NewDoc()
	err := IsNotFound(libs.Db.Scopes(Relation(relation)).Where("id = ?", id).First(t).Error)
	if err != nil {
		return nil, err
	}

	return t, nil
}

/**
 * 通过 name 获取 doc 记录
 * @method GetDocByName
 * @param  {[doc]}       doc  *Doc [description]
 */
func GetDocByName(name string) (*Doc, error) {
	t := NewDoc()
	err := IsNotFound(libs.Db.Where("name = ?", name).First(t).Error)
	if err != nil {
		return nil, err
	}

	return t, nil
}

/**
 * 通过 id 删除权限
 * @method DeleteDocById
 */
func DeleteDocById(id uint) error {
	t := NewDoc()
	t.ID = id
	if err := libs.Db.Delete(t).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteDocByIdError:%s \n", err))
		return err
	}
	return nil
}

/**
 * 获取所有的权限
 * @method GetAllDocs
 * @param  {[doc]} name string [description]
 * @param  {[doc]} orderBy string [description]
 * @param  {[doc]} offset int    [description]
 * @param  {[doc]} limit int    [description]
 */
func GetAllDocs(name, orderBy string, offset, limit int) ([]*Doc, error) {
	var docs []*Doc
	all := GetAll(&Doc{}, name, orderBy, offset, limit)
	if err := all.Find(&docs).Error; err != nil {
		return nil, err
	}

	return docs, nil
}

/**
 * 创建
 * @method CreateDoc
 * @param  {[doc]} kw string [description]
 * @param  {[doc]} cp int    [description]
 * @param  {[doc]} mp int    [description]
 */
func (p *Doc) CreateDoc() error {
	if err := libs.Db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

/**
 * 更新
 * @method UpdateDoc
 * @param  {[doc]} kw string [description]
 * @param  {[doc]} cp int    [description]
 * @param  {[doc]} mp int    [description]
 */
func UpdateDocById(id uint, np *Doc) error {
	if err := Update(&Doc{}, np, id); err != nil {
		return err
	}
	return nil
}
