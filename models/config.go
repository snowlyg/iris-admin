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

	Name  string `gorm:"not null; type:varchar(256)" json:"name" validate:"required" comment:"name"`
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

func GetConfigByName(name string) (*Config, error) {
	config := NewConfig()
	if err := libs.Db.Where("name = ?", name).First(config).Error; err != nil {
		return nil, err
	}
	return config, nil
}

func GetConfigById(id uint) (*Config, error) {
	config := NewConfig()
	if err := libs.Db.Where("id = ?", id).First(config).Error; err != nil {
		return nil, err
	}
	return config, nil
}

/**
 * 通过 id 删除
 * @method DeleteConfig
 */
func DeleteConfig(id uint) error {
	u, err := GetConfigById(id)
	if err != nil {
		return err
	}
	if err := libs.Db.Delete(u).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteConfigByIdErr:%s \n ", err))
		return err
	}
	return nil
}

/**
 * 获取所有的账号
 * @method GetAllConfig
 * @param  {[type]} name string [description]
 * @param  {[type]} configname string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllConfigs(name, orderBy string, offset, limit int) ([]*Config, error) {
	var configs []*Config
	q, _ := GetAll(&Config{}, name, orderBy, offset, limit)
	if err := q.Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

/**
 * 创建
 * @method CreateConfig
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (u *Config) CreateConfig() error {
	if err := libs.Db.Create(u).Error; err != nil {
		return err
	}

	return nil
}

/**
 * 更新
 * @method UpdateConfig
 * @param  {[type]} kw string [description]
 */
func UpdateConfig(id uint, nu *Config) error {

	if err := Update(&Config{}, nu, id); err != nil {
		return err
	}

	return nil
}
