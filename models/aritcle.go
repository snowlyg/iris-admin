package models

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"github.com/snowlyg/IrisAdminApi/sysinit"
	"time"
)

type Article struct {
	gorm.Model

	Title        string    `gorm:"not null;default:'';VARCHAR(256)" json:"title" validate:"required,gte=4,lte=256" comment:"标题"`
	ContentShort string    `gorm:"not null;default:'';VARCHAR(512)" json:"content_short" validate:"required,gte=6,lte=512" comment:"简介"`
	Author       string    `gorm:"not null;default:'';VARCHAR(30)" json:"author" comment:"作者" validate:"required,gte=4,lte=30"`
	ImageUri     string    `gorm:"not null;default:'';VARCHAR(2048)" json:"image_uri" comment:"封面" validate:"required"`
	SourceUri    string    `gorm:"not null;default:'';VARCHAR(512)" json:"source_uri" comment:"来源" `
	IsOriginal   bool      `gorm:"not null;default:true" json:"is_original" comment:"是否原创" validate:"required"`
	Content      string    `gorm:"not null;default:'';VARCHAR(2048)" json:"content" comment:"内容" validate:"required,gte=6,lte=2048"`
	Status       string    `gorm:"not null;default:'';VARCHAR(10)" json:"status" comment:"文章状态" validate:"required,gte=6,lte=10"`
	DisplayTime  time.Time `json:"display_time" validate:"required,datetime"`
}

func NewArticle() *Article {
	return &Article{}
}

/**
 * 通过 id 获取 role 记录
 * @method GetArticleById
 * @param  {[type]}       role  *Article [description]
 */
func (r *Article) GetArticleById(id uint) *Article {
	IsNotFound(sysinit.Db.Where("id = ?", id).First(r).Error)
	return r
}

/**
 * 通过 id 删除角色
 * @method DeleteArticleById
 */
func (r *Article) DeleteArticleById() *Article {
	if err := sysinit.Db.Delete(r).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteArticleErr:%s \n", err))
	}
	return r
}

/**
 * 获取所有的角色
 * @method GetAllArticle
 * @param  {[type]} name string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllArticles(name, orderBy string, offset, limit int) ([]*Article, int64, error) {
	var roles []*Article
	var count int64

	getAll := GetAll(&Article{}, name, orderBy, offset, limit)
	if err := getAll.Count(&count).Error; err != nil {
		return nil, count, err
	}

	if err := getAll.Find(&roles).Error; err != nil {
		return nil, count, err
	}

	return roles, count, nil
}

/**
 * 创建
 * @method CreateArticle
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (r *Article) CreateArticle() error {
	if err := sysinit.Db.Create(r).Error; err != nil {
		return err
	}

	return nil
}

func (r *Article) PaseTime(t string) error {
	displayTime, err := time.Parse(time.RFC3339Nano, t)
	if err != nil {
		return err
	}
	r.DisplayTime = displayTime
	return nil
}

//func addPerms(permIds []uint, role *Article) {
//	if len(permIds) > 0 {
//		roleId := strconv.FormatUint(uint64(role.ID), 10)
//		if _, err := sysinit.Enforcer.DeletePermissionsForUser(roleId); err != nil {
//			color.Red(fmt.Sprintf("AppendPermsErr:%s \n", err))
//		}
//		var perms []Permission
//		sysinit.Db.Where("id in (?)", permIds).Find(&perms)
//		for _, perm := range perms {
//			if _, err := sysinit.Enforcer.AddPolicy(roleId, perm.Name, perm.Act); err != nil {
//				color.Red(fmt.Sprintf("AddPolicy:%s \n", err))
//			}
//		}
//	}
//}

/**
 * 更新
 * @method UpdateArticle
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (r *Article) UpdateArticle(rj *Article) {

	if err := Update(r, rj); err != nil {
		color.Red(fmt.Sprintf("UpdatArticleErr:%s \n", err))
	}

	return
}
