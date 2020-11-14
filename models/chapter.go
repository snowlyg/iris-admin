package models

import (
	"fmt"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/libs/easygorm"
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
	Sort         int64     `gorm:"not null;default(0)" json:"sort" comment:"排序"`

	DocID uint
	Doc   *Doc
}

type MiniChapter struct {
	Id   uint  `json:"id" validate:"required"`
	Sort int64 `json:"sort" validate:"required"`
}

type SortChapter struct {
	OldId   uint  `json:"old_id" validate:"required"`
	OldSort int64 `json:"old_sort" validate:"required"`
	NewId   uint  `json:"new_id" validate:"required"`
	NewSort int64 `json:"new_sort" validate:"required"`
}

func NewChapter() *Chapter {
	return &Chapter{
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// GetChapterTableName
func GetChapterTableName() string {
	return fmt.Sprintf("%s%s", libs.Config.DB.Prefix, "chapters")
}

// GetDocReads 获取文章阅读量
func GetDocReads() (*easygorm.SumRes, error) {
	var sumRes easygorm.SumRes
	err := easygorm.Egm.Db.Model(&Chapter{}).Select("sum(`read`) as total").Scan(&sumRes).Error
	if err != nil {
		return &sumRes, err
	}
	return &sumRes, nil
}

// GetChapter 获取
func GetChapter(search *easygorm.Search) (*Chapter, error) {
	t := NewChapter()
	err := easygorm.First(t, search)
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
		err := easygorm.Save(p)
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
	err := easygorm.Save(p)
	if err != nil {
		return err
	}
	return nil
}

// DeleteChapterById 删除
func DeleteChapterById(id uint) error {
	t := NewChapter()
	if err := easygorm.DeleteById(t, id); err != nil {
		color.Red(fmt.Sprintf("DeleteChapterByIdError:%s \n", err))
		return err
	}
	return nil
}

// GetAllChapters
func GetAllChapters(search *easygorm.Search) ([]*Chapter, int64, error) {
	var chapters []*Chapter
	db, count, err := easygorm.Paginate(&Chapter{}, search)
	if err != nil {
		return nil, count, err
	}

	if err := db.Find(&chapters).Error; err != nil {
		return chapters, count, err
	}

	return chapters, count, nil
}

// getDoc get doc
func (p *Chapter) getDoc() {
	if p.Doc != nil {
		if p.Doc.ID > 0 {
			s := &easygorm.Search{
				Fields: []*easygorm.Field{
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
	if err := easygorm.Create(p); err != nil {
		return err
	}
	return nil
}

// UpdateChapterById update chapter by id
func UpdateChapterById(id uint, np *Chapter) error {
	np.getDoc()
	if err := easygorm.Update(&Chapter{}, np, id); err != nil {
		return err
	}
	return nil
}

func Sort(sc *SortChapter) error {
	err := easygorm.Egm.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Chapter{}).Where("id = ?", sc.NewId).Update("sort", sc.NewSort).Error; err != nil {
			return err
		}

		if err := tx.Model(&Chapter{}).Where("id = ?", sc.OldId).Update("sort", sc.OldSort).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
