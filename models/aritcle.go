package models

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/easygorm"
	"gorm.io/gorm"
	"net/http"
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

	TypeID   uint
	Type     *Type
	Tags     []*Tag   `gorm:"many2many:article_tags;"`
	TagNames []string `gorm:"-" json:"tag_names"`
	Ips      []*ArticleIp
}

func NewArticle() *Article {
	return &Article{}
}

func (r *Article) ReadArticle(rh *http.Request) error {
	search := &easygorm.Search{
		Fields: []*easygorm.Field{
			{
				Key:       "article_id",
				Condition: "=",
				Value:     r.ID,
			},
		},
	}
	articleIps, err := GetArticleIps(search)
	if err != nil {
		return err
	}

	publicIp := libs.ClientPublicIp(rh)
	if publicIp == "" {
		return nil
	}

	for _, articleIp := range articleIps {
		// 原来ip增加访问次数
		if articleIp.Addr == publicIp {
			err := articleIp.AddArticleIpMun()
			if err != nil {
				return err
			}
			return nil
		}
	}

	r.Lock()
	defer r.Unlock()
	r.Read++
	if err := easygorm.UpdateWithFilde(&Article{}, map[string]interface{}{"Read": r.Read}, r.ID); err != nil {
		return err
	}

	// 没有的话就创建新的 ip
	articleIp := ArticleIp{
		Mun:       1,
		Addr:      publicIp,
		ArticleID: r.ID,
		Article:   r,
	}
	err = articleIp.CreateArticleIp()
	if err != nil {
		return err
	}

	return nil

}

func (r *Article) LikeArticle() error {
	r.Lock()
	defer r.Unlock()

	r.Like++
	if err := easygorm.UpdateWithFilde(&Article{}, map[string]interface{}{"Like": r.Like}, r.ID); err != nil {
		return err
	}
	return nil
}

// GetArticle 获取文章
func GetArticle(s *easygorm.Search) (*Article, error) {
	r := NewArticle()
	err := easygorm.First(r, s)
	if !IsNotFound(err) {
		return r, err
	}
	return r, nil
}

// GetArticleById 通过 id 获取文章
func GetArticleById(id uint) (*Article, error) {
	r := NewArticle()
	err := easygorm.FindById(r, id)
	if !IsNotFound(err) {
		return r, err
	}
	return r, nil
}

// GetArticleCount 获取文章数量
func GetArticleCount() (int64, error) {
	var count int64
	r := NewArticle()
	err := easygorm.Egm.Db.Model(r).Count(&count).Error
	if err != nil {
		return count, err
	}
	return count, nil
}

// GetArticleReads 获取文章阅读量
func GetArticleReads() (int64, error) {
	sr, err := easygorm.Count(&Article{}, "read")
	if err != nil {
		return sr, err
	}
	return sr, nil
}

// DeleteArticleById 删除
func DeleteArticleById(id uint) error {
	r := NewArticle()
	if err := easygorm.DeleteById(r, id); err != nil {
		color.Red(fmt.Sprintf("DeleteArticleErr:%s \n", err))
		return err
	}
	return nil
}

// GetAllArticles 获取集合
func GetAllArticles(search *easygorm.Search, tagId int) ([]*Article, int64, error) {
	var articles []*Article

	// 多对多标签搜索
	if tagId > 0 {
		var tagArticleIds []int
		s := &easygorm.Search{
			Fields: []*easygorm.Field{
				{
					Key:       "id",
					Condition: "=",
					Value:     tagId,
				},
			},
			Relations: []*easygorm.Relate{
				{
					Value: "Articles",
				},
			},
		}
		tag, err := GetTag(s)
		if err != nil {
			return nil, 0, err
		}

		for _, tagArticle := range tag.Articles {
			if tagArticle.ID > 0 {
				tagArticleIds = append(tagArticleIds, int(tagArticle.ID))
			}
		}

		field := &easygorm.Field{
			Condition: "IN",
			Key:       "id",
			Value:     tagArticleIds,
		}
		search.Fields = append(search.Fields, field)
	}

	count, err := easygorm.Paginate(&Article{}, &articles, search)
	if err != nil {
		return nil, count, err
	}

	return articles, count, nil
}

// CreateArticle 创建
func (r *Article) CreateArticle() error {
	r.getTypes()
	if err := easygorm.Create(r); err != nil {
		return err
	}
	r.addTags()
	return nil
}

// addTags 添加标签关联关系
func (r *Article) addTags() {
	if len(r.Tags) > 0 {
		if err := easygorm.Egm.Db.Model(r).Association("Tags").Clear(); err != nil {
			color.Red(fmt.Sprintf("Tags 清空关系错误:%+v\n", err))
		}
	}
	if len(r.TagNames) > 0 {
		var tags []*Tag
		for _, tagName := range r.TagNames {
			s := &easygorm.Search{
				Fields: []*easygorm.Field{
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

		err := easygorm.Egm.Db.Model(r).Association("Tags").Append(tags)
		if err != nil {
			color.Red(fmt.Sprintf("标签添加错误:%+v\n tags:%+v", err, tags))
		}
	}
}

// getTypes 获取关联数据
func (r *Article) getTypes() {
	if r.Type != nil {
		if r.Type.ID > 0 {
			s := &easygorm.Search{
				Fields: []*easygorm.Field{
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
	if err := easygorm.Update(&Article{}, nr, []interface{}{"Title", "ContentShort", "Author", "ImageUri", "SourceUri", "IsOriginal", "Content", "Status"}, id); err != nil {
		return err
	}
	nr.addTags()
	return nil
}
