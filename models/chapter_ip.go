package models

import (
	"github.com/snowlyg/blog/libs/logging"
	"github.com/snowlyg/easygorm"
	"gorm.io/gorm"
	"sync"
)

type ChapterIp struct {
	sync.Mutex
	gorm.Model

	Mun  int64  `gorm:"not null;default(0)" json:"mun" comment:"访问次数"`
	Type int8   `gorm:"not null;default(0)" json:"type" comment:"访问类型：0 无操作, 1 阅读，2 点赞，3 阅读点赞"`
	Addr string `gorm:"not null;default(0);type:varchar(20)" json:"addr" comment:"ip 地址"`

	ChapterID uint
	Chapter   *Chapter
}

// GetChapterIps get all ips
func GetChapterIps(s *easygorm.Search) ([]*ChapterIp, error) {
	var ips []*ChapterIp
	err := easygorm.All(&ChapterIp{}, &ips, s)
	if err != nil {
		logging.Err.Errorf("get chapter ips err :%+v\n", err)
		return ips, err
	}
	return ips, nil
}

// CreateChapterIp add chapter ip
func (p *ChapterIp) CreateChapterIp() error {
	if err := easygorm.Create(p); err != nil {
		logging.Err.Errorf("create chapter ip err :%+v\n", err)
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
		logging.Err.Errorf("add chapter ip err :%+v\n", err)
		return err
	}
	return nil
}

// UpdateType add article ip type
func (p *ChapterIp) UpdateType() error {
	p.Lock()
	defer p.Unlock()
	p.Type++
	if err := easygorm.UpdateWithFilde(&ChapterIp{}, map[string]interface{}{"Type": p.Type}, p.ID); err != nil {
		logging.Err.Errorf("update chapter type err :%+v\n", err)
		return err
	}
	return nil
}
