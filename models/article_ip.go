package models

import (
	"github.com/snowlyg/easygorm"
	"gorm.io/gorm"
	"sync"
)

type ArticleIp struct {
	sync.Mutex
	gorm.Model

	Mun  int64  `gorm:"not null;default(0)" json:"mun" comment:"访问次数"`
	Type int8   `gorm:"not null;default(0)" json:"type" comment:"访问类型：0 无操作, 1 阅读，2 点赞，3 阅读点赞"`
	Addr string `gorm:"not null;default(0);type:varchar(20)" json:"addr" comment:"ip 地址"`

	ArticleID uint
	Article   *Article
}

// GetArticleIps get all ips
func GetArticleIps(s *easygorm.Search) ([]*ArticleIp, error) {
	var ips []*ArticleIp
	err := easygorm.All(&ArticleIp{}, &ips, s)
	if err != nil {
		return ips, err
	}

	return ips, nil
}

// CreateArticleIp add article ip
func (p *ArticleIp) CreateArticleIp() error {
	if err := easygorm.Create(p); err != nil {
		return err
	}
	return nil
}

// UpdateType add article ip type
func (p *ArticleIp) UpdateType() error {
	p.Lock()
	defer p.Unlock()
	p.Type++
	if err := easygorm.UpdateWithFilde(&ArticleIp{}, map[string]interface{}{"Type": p.Type}, p.ID); err != nil {
		return err
	}
	return nil
}

// AddArticleIpMun add article mun
func (p *ArticleIp) AddArticleIpMun() error {
	p.Lock()
	defer p.Unlock()

	p.Mun++
	if err := easygorm.UpdateWithFilde(&ArticleIp{}, map[string]interface{}{"Mun": p.Mun}, p.ID); err != nil {
		return err
	}
	return nil
}
