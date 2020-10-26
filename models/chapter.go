package models

import (
	"fmt"
	"github.com/snowlyg/blog/libs"
	"time"

	"github.com/fatih/color"
	"gorm.io/gorm"
)

type Chapter struct {
	gorm.Model
	Name string `gorm:"not null;type:varchar(256)" json:"name" validate:"required,gte=0,lte=256" comment:"章节名称"`

	DocID    uint
	Doc      *Doc
	Articles []*Article
}

func NewChapter() *Chapter {
	return &Chapter{
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

/**
 * 通过 id 获取 chapter 记录
 * @method GetChapterById
 * @param  {[chapter]}       chapter  *Chapter [description]
 */
func GetChapterById(id uint) (*Chapter, error) {
	t := NewChapter()
	err := IsNotFound(libs.Db.Where("id = ?", id).First(t).Error)
	if err != nil {
		return nil, err
	}

	return t, nil
}

/**
 * 通过 name 获取 chapter 记录
 * @method GetChapterByName
 * @param  {[chapter]}       chapter  *Chapter [description]
 */
func GetChapterByName(name string) (*Chapter, error) {
	t := NewChapter()
	err := IsNotFound(libs.Db.Where("name = ?", name).Preload("Doc").First(t).Error)
	if err != nil {
		return nil, err
	}

	return t, nil
}

/**
 * 通过 id 删除权限
 * @method DeleteChapterById
 */
func DeleteChapterById(id uint) error {
	t := NewChapter()
	t.ID = id
	if err := libs.Db.Delete(t).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteChapterByIdError:%s \n", err))
		return err
	}
	return nil
}

/**
 * 获取所有的权限
 * @method GetAllChapters
 * @param  {[chapter]} name string [description]
 * @param  {[chapter]} orderBy string [description]
 * @param  {[chapter]} offset int    [description]
 * @param  {[chapter]} limit int    [description]
 */
func GetAllChapters(docId uint, name, orderBy string, offset, limit int) ([]*Chapter, error) {
	var chapters []*Chapter
	all := GetAll(&Chapter{}, name, orderBy, offset, limit)
	if docId > 0 {
		all = all.Where("doc_id = ? ", docId)
	}
	if err := all.Preload("Doc").Find(&chapters).Error; err != nil {
		return nil, err
	}

	return chapters, nil
}

func (p *Chapter) getDoc() {
	if p.Doc != nil {
		if p.Doc.ID > 0 {
			doc, err := GetDocById(p.Doc.ID)
			if err == nil && doc.ID > 0 {
				p.DocID = doc.ID
				p.Doc = doc
			}
		}
	}
}

/**
 * 创建
 * @method CreateChapter
 * @param  {[chapter]} kw string [description]
 * @param  {[chapter]} cp int    [description]
 * @param  {[chapter]} mp int    [description]
 */
func (p *Chapter) CreateChapter() error {
	p.getDoc()
	if err := libs.Db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

/**
 * 更新
 * @method UpdateChapter
 * @param  {[chapter]} kw string [description]
 * @param  {[chapter]} cp int    [description]
 * @param  {[chapter]} mp int    [description]
 */
func UpdateChapterById(id uint, np *Chapter) error {
	np.getDoc()
	if err := Update(&Chapter{}, np, id); err != nil {
		return err
	}
	return nil
}
