package models

import (
	"gorm.io/gorm"
)

type Config struct {
	gorm.Model
	Name  string `gorm:"unique;not null;type:varchar(256)" json:"name" validate:"required" comment:"name"`
	Value string `gorm:"not null;type:varchar(1024)" json:"value" validate:"required"  comment:"value"`
}

//// GetConfig get config
//func GetConfig(search *easygorm.Search) (*Config, error) {
//	r := &Config{}
//	err := easygorm.Found(search).First(r).Error
//	if err != nil {
//		logging.ErrorLogger.Errorf("get config err: %+v", err)
//		return r, err
//	}
//	if r.ID == 0 {
//		return r, errors.New("数据不存在")
//	}
//	return r, nil
//}
//
//// DeleteConfig del config
//func DeleteConfig(id uint) error {
//	u := &Config{}
//	if err := easygorm.DeleteById(u, id); err != nil {
//		logging.ErrorLogger.Errorf("del config err: %+v", err)
//		return err
//	}
//	return nil
//}
//
//// GetAllConfigs get all configs
//func GetAllConfigs(s *easygorm.Search) ([]*Config, error) {
//	var configs []*Config
//	if err := easygorm.All(&Config{}, &configs, s); err != nil {
//		logging.ErrorLogger.Errorf("get all configs err: %+v", err)
//		return configs, err
//	}
//	return configs, nil
//}
//
//// CreateConfig create config
//func (u *Config) CreateConfig() error {
//	if err := easygorm.Create(u); err != nil {
//		logging.ErrorLogger.Errorf("create config err: %+v", err)
//		return err
//	}
//	return nil
//}
//
//// UpdateConfig update config
//func UpdateConfig(id uint, nu *Config) error {
//	if err := easygorm.Update(&Config{}, nu, nil, id); err != nil {
//		logging.ErrorLogger.Errorf("update config err: %+v", err)
//		return err
//	}
//	return nil
//}
