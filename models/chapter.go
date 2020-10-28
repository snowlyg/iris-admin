package models

import (
	"fmt"
	"github.com/snowlyg/blog/libs"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"gorm.io/gorm"
)

type Chapter struct {
	sync.Mutex
	gorm.Model

	Title        string    `gorm:"not null;default:'';type:varchar(256)" json:"title" validate:"required,gte=2,lte=256" comment:"标题"`
	ContentShort string    `gorm:"not null;default:'';type:varchar(512)" json:"content_short" validate:"required,gte=2,lte=512" comment:"简介"`
	Author       string    `gorm:"not null;default:'';type:varchar(30)" json:"author" comment:"作者" validate:"required,gte=4,lte=30"`
	ImageUri     string    `gorm:"type:longText" json:"image_uri" comment:"封面" `
	SourceUri    string    `gorm:"not null;default:'';type:varchar(512)" json:"source_uri" comment:"来源"`
	IsOriginal   bool      `gorm:"not null;default:true;type:tinyint(1)" json:"is_original" comment:"是否原创" validate:""`
	Content      string    `gorm:"type:longText" json:"content" comment:"内容" validate:"required,gte=6"`
	Status       string    `gorm:"not null;default:'';type:varchar(10)" json:"status" comment:"文章状态" validate:"required,gte=1,lte=10"`
	DisplayTime  time.Time `json:"display_time" comment:"发布时间" validate:"required"`
	Like         int64     `gorm:"not null;default(0)" json:"like" comment:"点赞"`
	Read         int64     `gorm:"not null;default(0)" json:"read" comment:"阅读量"`
	Ips          string    `gorm:"not null;default(0);type:varchar(1024)" json:"ips" comment:"ip 地址"`

	DocID uint
	Doc   *Doc
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
	err := IsNotFound(libs.Db.Where("id = ?", id).Preload("Doc").First(t).Error)
	if err != nil {
		return nil, err
	}

	return t, nil
}

/**
 * 通过 id 获取 chapter 记录
 * @method GetPublishedChapterById
 * @param  {[chapter]}       chapter  *Chapter [description]
 */
func GetPublishedChapterById(id uint) (*Chapter, error) {
	t := NewChapter()
	err := IsNotFound(libs.Db.Where("id = ?", id).Preload("Doc").Where("status = ?", "published").First(t).Error)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (c *Chapter) ReadChapter(rh *http.Request) error {
	ip := c.Ips
	ips := strings.Split(ip, ",")
	publicIp := libs.ClientPublicIp(rh)
	if !libs.InArrayS(ips, publicIp) {
		c.Lock()
		defer c.Unlock()

		c.Read++
		ips = append(ips, publicIp)
		c.Ips = strings.Join(ips, ",")
		err := libs.Db.Save(c).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Chapter) LikeChapter() error {
	c.Lock()
	defer c.Unlock()

	c.Like++
	err := libs.Db.Save(c).Error
	if err != nil {
		return err
	}
	return nil
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
func GetAllChapters(docId int, published, name, orderBy string, offset, limit int) ([]*Chapter, int64, error) {
	var chapters []*Chapter
	var count int64
	all := GetAll(&Chapter{}, name, orderBy, offset, limit)
	if len(published) > 0 {
		all = all.Where("status = ?", "published")
	}

	if docId > 0 {
		all = all.Where("doc_id = ? ", docId)
	}

	if err := all.Count(&count).Error; err != nil {
		return nil, count, err
	}

	if err := all.Preload("Doc").Find(&chapters).Error; err != nil {
		return nil, count, err
	}

	return chapters, count, nil
}

func (p *Chapter) getDoc() {
	if p.Doc != nil {
		if p.Doc.ID > 0 {
			doc, err := GetDocById(p.Doc.ID, "")
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
