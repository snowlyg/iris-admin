package models

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/snowlyg/blog/libs/easygorm"
	"gorm.io/gorm"
)

type Config struct {
	gorm.Model

	Name  string `gorm:"unique;not null;type:varchar(256)" json:"name" validate:"required" comment:"name"`
	Value string `gorm:"not null;type:varchar(1024)" json:"value" validate:"required"  comment:"value"`
}

func NewConfig() *Config {
	return &Config{
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// GetConfig get config
func GetConfig(search *easygorm.Search) (*Config, error) {
	r := NewConfig()
	err := easygorm.Found(search).First(r).Error
	if !IsNotFound(err) {
		return r, err
	}
	return r, nil
}

// DeleteConfig del config
func DeleteConfig(id uint) error {
	u := NewConfig()
	u.ID = id

	if err := easygorm.Egm.Db.Delete(u).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteConfigByIdErr:%s \n ", err))
		return err
	}
	return nil
}

// GetAllConfigs get all configs
func GetAllConfigs(s *easygorm.Search) ([]*Config, error) {
	var configs []*Config
	if err := easygorm.All(&Config{}, s).Find(&configs).Error; err != nil {
		return configs, err
	}
	return configs, nil
}

// CreateConfig create config
func (u *Config) CreateConfig() error {
	if err := easygorm.Egm.Db.Create(u).Error; err != nil {
		return err
	}

	return nil
}

// UpdateConfig update config
func UpdateConfig(id uint, nu *Config) error {

	if err := easygorm.Update(&Config{}, nu, id); err != nil {
		return err
	}

	return nil
}
