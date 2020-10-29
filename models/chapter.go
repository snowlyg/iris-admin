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

// GetChapter 获取
func GetChapter(search *Search) (*Chapter, error) {
	t := NewChapter()
	err := Found(search).First(t).Error
	if !IsNotFound(err) {
		return t, err
	}
	return t, nil
}

// ReadChapter 增加阅读量
func (p *Chapter) ReadChapter(rh *http.Request) error {
	ip := p.Ips
	ips := strings.Split(ip, ",")
	publicIp := libs.ClientPublicIp(rh)
	if !libs.InArrayS(ips, publicIp) {
		p.Lock()
		defer p.Unlock()

		p.Read++
		ips = append(ips, publicIp)
		p.Ips = strings.Join(ips, ",")
		err := libs.Db.Save(p).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// LikeChapter 点赞
func (p *Chapter) LikeChapter() error {
	p.Lock()
	defer p.Unlock()

	p.Like++
	err := libs.Db.Save(p).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteChapterById 删除
func DeleteChapterById(id uint) error {
	t := NewChapter()
	t.ID = id
	if err := libs.Db.Delete(t).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteChapterByIdError:%s \n", err))
		return err
	}
	return nil
}

// GetAllChapters
func GetAllChapters(search *Search) ([]*Chapter, int64, error) {
	var chapters []*Chapter
	var count int64
	all := GetAll(&Chapter{}, search)

	if err := all.Count(&count).Error; err != nil {
		return nil, count, err
	}

	all = all.Scopes(Paginate(search.Offset, search.Limit))

	if err := all.Find(&chapters).Error; err != nil {
		return nil, count, err
	}

	return chapters, count, nil
}

// getDoc get doc
func (p *Chapter) getDoc() {
	if p.Doc != nil {
		if p.Doc.ID > 0 {
			s := &Search{
				Fields: []*Filed{
					{
						Key:       "id",
						Condition: "=",
						Value:     p.Doc.ID,
					},
				},
			}
			doc, err := GetDoc(s)
			if err == nil && doc.ID > 0 {
				p.DocID = doc.ID
				p.Doc = doc
			}
		}
	}
}

// CreateChapter create chapter
func (p *Chapter) CreateChapter() error {
	p.getDoc()
	if err := libs.Db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

// UpdateChapterById update chapter by id
func UpdateChapterById(id uint, np *Chapter) error {
	np.getDoc()
	if err := Update(&Chapter{}, np, id); err != nil {
		return err
	}
	return nil
}
