package models

import (
	"github.com/snowlyg/easygorm"
	"gorm.io/gorm"
	"sync"
)

type ChapterIp struct {
	sync.Mutex
	gorm.Model

	Mun  int64  `gorm:"not null;default(0)" json:"mun" comment:"访问次数"`
	Addr string `gorm:"not null;default(0);type:varchar(20)" json:"addr" comment:"ip 地址"`

	ChapterID uint
	Chapter   *Chapter
}

// GetChapterIps get all ips
func GetChapterIps(s *easygorm.Search) ([]*ChapterIp, error) {
	var ips []*ChapterIp
	err := easygorm.All(&ChapterIp{}, &ips, s)
	if err != nil {
		return ips, err
	}

	return ips, nil
}

// CreateChapterIp add chapter ip
func (p *ChapterIp) CreateChapterIp() error {
	if err := easygorm.Create(p); err != nil {
		return err
	}
	return nil
}

// AddChapterIpMun add chapter mun
func (p *ChapterIp) AddChapterIpMun() error {
	p.Lock()
	defer p.Unlock()

	p.Mun++
	if err := easygorm.UpdateWithFilde(&ChapterIp{}, map[string]interface{}{"Mun": p.Mun}, p.ID); err != nil {
		return err
	}
	return nil
}
