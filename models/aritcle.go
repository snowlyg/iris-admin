package models

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/snowlyg/IrisAdminApi/sysinit"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sync"
	"time"
)

type Article struct {
	sync.Mutex
	gorm.Model

	Title        string    `gorm:"not null;default:'';type:varchar(256)" json:"title" validate:"required,gte=4,lte=256" comment:"标题"`
	ContentShort string    `gorm:"not null;default:'';type:varchar(512)" json:"content_short" validate:"required,gte=6,lte=512" comment:"简介"`
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
}

func NewArticle() *Article {
	return &Article{}
}

/**
 * 通过 id 获取 role 记录
 * @method GetArticleById
 * @param  {[type]}       role  *Article [description]
 */
func GetPublishedArticleById(id uint) (*Article, error) {
	r := NewArticle()
	err := IsNotFound(sysinit.Db.Where("id = ?", id).Where("status = ?", "published").Preload(clause.Associations).First(r).Error)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Article) ReadArticle() error {
	r.Lock()
	defer r.Unlock()

	r.Read++
	err := sysinit.Db.Save(r).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Article) LikeArticle() error {
	r.Lock()
	defer r.Unlock()

	r.Like++
	err := sysinit.Db.Save(r).Error
	if err != nil {
		return err
	}
	return nil
}

/**
 * 通过 id 获取 role 记录
 * @method GetArticleById
 * @param  {[type]}       role  *Article [description]
 */
func GetArticleById(id uint) (*Article, error) {
	r := NewArticle()
	err := IsNotFound(sysinit.Db.Where("id = ?", id).Preload(clause.Associations).First(r).Error)
	if err != nil {
		return nil, err
	}
	return r, nil
}

/**
 * 通过 id 删除角色
 * @method DeleteArticleById
 */
func DeleteArticleById(id uint) error {
	r := NewArticle()
	r.ID = id
	if err := sysinit.Db.Delete(r).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteArticleErr:%s \n", err))
		return err
	}
	return nil
}

/**
 * 获取所有的角色
 * @method GetAllArticle
 * @param  {[type]} searchStr string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllArticles(searchStr, orderBy, published string, offset, limit, tagId int) ([]*Article, int64, error) {
	var articles []*Article
	var count int64

	getAll := GetAll(&Article{}, searchStr, orderBy, offset, limit)
	if err := getAll.Count(&count).Error; err != nil {
		return nil, count, err
	}
	if len(published) > 0 {
		getAll = getAll.Where("status = ?", "published")
	}

	if tagId > 0 {
		var tagArticleIds []int64
		tag, err := GetTagById(uint(tagId), true)
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

	if err := getAll.Preload(clause.Associations).Find(&articles).Error; err != nil {
		return nil, count, err
	}

	return articles, count, nil
}

/**
 * 创建
 * @method CreateArticle
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (r *Article) CreateArticle() error {
	r.getTagTypes()

	if err := sysinit.Db.Create(r).Error; err != nil {
		return err
	}

	return nil
}

func (r *Article) getTagTypes() {
	if err := sysinit.Db.Model(r).Association("Tags").Clear(); err != nil {
		fmt.Println(fmt.Sprintf("Tags 清空关系错误:%+v\n", err))
	}

	if r.Type != nil && len(r.Type.Name) > 0 {
		tt, err := GetTypeByName(r.Type.Name)
		if err == nil && tt.ID > 0 {
			r.TypeID = tt.ID
			r.Type = tt
		}
	}

	if len(r.TagNames) > 0 {
		var tags []*Tag
		for _, tagName := range r.TagNames {
			tag, err := GetTagByName(tagName)
			if err == nil {
				if tag.ID == 0 {
					tag.Name = tagName
					err := tag.CreateTag()
					if err != nil {
						fmt.Println(fmt.Sprintf("标签新建错误:%+v\n", err))
					}
				}
				tags = append(tags, tag)
			}
		}

		r.Tags = tags
	}
}

/**
 * 更新
 * @method UpdateArticle
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func UpdateArticle(id uint, nr *Article) error {
	r, err := GetArticleById(id)
	if err != nil {
		return err
	}
	nr.getTagTypes()

	if err := Update(r, nr); err != nil {
		return err
	}

	return nil
}
