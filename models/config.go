package models

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/snowlyg/blog/libs"
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
func GetConfig(search *Search) (*Config, error) {
	r := NewConfig()
	err := Found(search).First(r).Error
	if !IsNotFound(err) {
		return r, err
	}
	return r, nil
}

// DeleteConfig del config
func DeleteConfig(id uint) error {
	u := NewConfig()
	u.ID = id

	if err := libs.Db.Delete(u).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteConfigByIdErr:%s \n ", err))
		return err
	}
	return nil
}

// GetAllConfigs get all configs
func GetAllConfigs(s *Search) ([]*Config, error) {
	var configs []*Config
	q := GetAll(&Config{}, s)
	q = q.Scopes(Relation(s.Relations))
	if err := q.Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// CreateConfig create config
func (u *Config) CreateConfig() error {
	if err := libs.Db.Create(u).Error; err != nil {
		return err
	}

	return nil
}

// UpdateConfig update config
func UpdateConfig(id uint, nu *Config) error {

	if err := Update(&Config{}, nu, id); err != nil {
		return err
	}

	return nil
}
