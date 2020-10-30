package models

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/snowlyg/blog/libs"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Article struct {
	sync.Mutex
	gorm.Model

	Title        string    `gorm:"not null;default:'';type:varchar(256)" json:"title" validate:"required,gte=2,lte=256" comment:"标题"`
	ContentShort string    `gorm:"not null;default:'';type:varchar(512)" json:"content_short" validate:"required,gte=2,lte=512" comment:"简介"`
	Author       string    `gorm:"not null;default:'';type:varchar(30)" json:"author" comment:"作者" validate:"required,gte=4,lte=30"`
	ImageUri     string    `gorm:"type:longText" json:"image_uri" comment:"封面" validate:"required"`
	SourceUri    string    `gorm:"not null;default:'';type:varchar(512)" json:"source_uri" comment:"来源"`
	IsOriginal   bool      `gorm:"not null;default:true;type:tinyint(1)" json:"is_original" comment:"是否原创" validate:""`
	Content      string    `gorm:"type:longText" json:"content" comment:"内容" validate:"required,gte=6"`
	Status       string    `gorm:"not null;default:'';type:varchar(10)" json:"status" comment:"文章状态" validate:"required,gte=1,lte=10"`
	DisplayTime  time.Time `json:"display_time" comment:"发布时间" validate:"required"`
	Like         int64     `gorm:"not null;default(0)" json:"like" comment:"点赞"`
	Read         int64     `gorm:"not null;default(0)" json:"read" comment:"阅读量"`
	Ips          string    `gorm:"not null;default(0);type:varchar(1024)" json:"ips" comment:"ip 地址"`

	TypeID   uint
	Type     *Type
	Tags     []*Tag   `gorm:"many2many:article_tags;"`
	TagNames []string `gorm:"-" json:"tag_names"`
}

func NewArticle() *Article {
	return &Article{}
}

func (r *Article) ReadArticle(rh *http.Request) error {
	ip := r.Ips
	ips := strings.Split(ip, ",")
	publicIp := libs.ClientPublicIp(rh)
	if !libs.InArrayS(ips, publicIp) {
		r.Lock()
		defer r.Unlock()

		r.Read++
		ips = append(ips, publicIp)
		r.Ips = strings.Join(ips, ",")
		err := libs.Db.Save(r).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Article) Addip() error {
	r.Lock()
	defer r.Unlock()

	r.Read++
	err := libs.Db.Save(r).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Article) LikeArticle() error {
	r.Lock()
	defer r.Unlock()

	r.Like++
	err := libs.Db.Save(r).Error
	if err != nil {
		return err
	}
	return nil
}

// GetArticle 获取文章
func GetArticle(search *Search) (*Article, error) {
	r := NewArticle()
	err := Found(search).First(r).Error
	if !IsNotFound(err) {
		return r, err
	}
	return r, nil
}

// DeleteArticleById 删除
func DeleteArticleById(id uint) error {
	r := NewArticle()
	r.ID = id
	if err := libs.Db.Delete(r).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteArticleErr:%s \n", err))
		return err
	}
	return nil
}

// GetAllArticles 获取集合
func GetAllArticles(search *Search, tagId int) ([]*Article, int64, error) {
	var articles []*Article
	var count int64

	getAll := GetAll(&Article{}, search)

	// 多对多标签搜索
	if tagId > 0 {
		var tagArticleIds []int64
		s := &Search{
			Fields: []*Filed{
				{
					Key:       "id",
					Condition: "=",
					Value:     tagId,
				},
			},
			Relations: []*Relate{
				{
					Value: "Articles",
				},
			},
		}
		tag, err := GetTag(s)
		if err != nil {
			return nil, count, err
		}

		for _, tagArticle := range tag.Articles {
			if tagArticle.ID > 0 {
				tagArticleIds = append(tagArticleIds, int64(tagArticle.ID))
			}
		}

		getAll = getAll.Where("id IN ?", tagArticleIds)
	}

	if err := getAll.Count(&count).Error; err != nil {
		return nil, count, err
	}

	getAll = getAll.Scopes(Paginate(search.Offset, search.Limit))

	if err := getAll.Find(&articles).Error; err != nil {
		return nil, count, err
	}

	color.Yellow(fmt.Sprintf("Searach :%+v", search))

	return articles, count, nil
}

// CreateArticle 创建
func (r *Article) CreateArticle() error {
	r.getTypes()

	if err := libs.Db.Create(r).Error; err != nil {
		return err
	}

	r.addTags()

	return nil
}

func (r *Article) addTags() {
	if err := libs.Db.Model(r).Association("Tags").Clear(); err != nil {
		color.Red(fmt.Sprintf("Tags 清空关系错误:%+v\n", err))
	}
	if len(r.TagNames) > 0 {
		var tags []*Tag
		for _, tagName := range r.TagNames {
			s := &Search{
				Fields: []*Filed{
					{
						Key:       "name",
						Condition: "=",
						Value:     tagName,
					},
				},
			}
			tag, err := GetTag(s)
			if err != nil || tag.ID == 0 {
				tag = NewTag()
				tag.Name = tagName
				err := tag.CreateTag()
				if err != nil {
					color.Red(fmt.Sprintf("标签新建错误:%+v\n", err))
				}
			}
			tags = append(tags, tag)
		}

		err := libs.Db.Model(r).Association("Tags").Append(tags)
		if err != nil {
			color.Red(fmt.Sprintf("标签添加错误:%+v\n tags:%+v", err, tags))
		}
	}
}

// getTypes 获取关联数据
func (r *Article) getTypes() {
	if r.Type != nil {
		if r.Type.ID > 0 {
			s := &Search{
				Fields: []*Filed{
					{
						Key:       "id",
						Condition: "=",
						Value:     r.Type.ID,
					},
				},
			}
			tt, err := GetType(s)
			if err == nil && tt.ID > 0 {
				r.TypeID = tt.ID
				r.Type = tt
			}
		}
	}

}

// UpdateArticle 更新
func UpdateArticle(id uint, nr *Article) error {
	nr.getTypes()

	if err := Update(&Article{}, nr, id); err != nil {
		return err
	}

	nr.addTags()
	return nil
}
