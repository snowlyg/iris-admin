package dperm

import (
	"errors"
	"fmt"
	"github.com/snowlyg/blog/application/libs/easygorm"
	"github.com/snowlyg/blog/application/libs/logging"
	"github.com/snowlyg/blog/application/models"
)

const modelName = "权限管理"

type PermResponse struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}

type PermReq struct {
	Name        string `json:"name" `
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Act         string `json:"act"`
}

func (p *PermResponse) ModelName() string {
	return modelName
}

func (p *PermResponse) Model() *models.Permission {
	return &models.Permission{}
}

func (p *PermResponse) All(name, sort, orderBy string, page, pageSize int) (map[string]interface{}, error) {
	var count int64
	var perms []*PermResponse
	db := easygorm.GetEasyGormDb().Model(p.Model())
	if len(name) > 0 {
		db = db.Where("name", "like", fmt.Sprintf("%%%s%%", name))
	}
	err := db.Count(&count).Error
	if err != nil {
		logging.ErrorLogger.Errorf("get list count err ", err)
		return nil, err
	}
	err = db.Scopes(easygorm.PaginateScope(page, pageSize, sort, orderBy)).Find(&perms).Error
	if err != nil {
		logging.ErrorLogger.Errorf("get list data err ", err)
		return nil, err
	}
	list := map[string]interface{}{"items": perms, "total": count, "limit": pageSize}
	return list, nil
}

func (p *PermResponse) FindByNameAndAct(name, act string) error {
	err := easygorm.GetEasyGormDb().Model(p.Model()).Where("name = ?", name).Where("act = ?", act).Find(p).Error
	if err != nil {
		logging.ErrorLogger.Errorf("find perm by name get err ", err)
		return err
	}
	return nil
}

func (p *PermResponse) Create(object map[string]interface{}) error {
	err := p.checkNameAndAct(object)
	if err != nil {
		logging.ErrorLogger.Errorf("check perm name and act get err ", err)
		return err
	}

	err = easygorm.GetEasyGormDb().Model(p.Model()).Create(object).Error
	if err != nil {
		logging.ErrorLogger.Errorf("create perm err ", err)
		return err
	}
	return nil
}

func (p *PermResponse) Update(id uint, object map[string]interface{}) error {
	err := p.Find(id)
	if err != nil {
		logging.ErrorLogger.Errorf("find perm by id get err ", err)
		return err
	}

	err = p.checkNameAndAct(object)
	if err != nil {
		logging.ErrorLogger.Errorf("check perm name and act get err ", err)
		return err
	}

	err = easygorm.GetEasyGormDb().Model(p.Model()).Where("id = ?", id).Updates(object).Error
	if err != nil {
		logging.ErrorLogger.Errorf("update perm  get err ", err)
		return err
	}
	return nil
}

func (p *PermResponse) checkNameAndAct(object map[string]interface{}) error {
	var name string
	var act string
	if obj, ok := object["Name"].(string); ok {
		name = obj
	}
	if obj, ok := object["Act"].(string); ok {
		act = obj
	}

	if len(name) > 0 && len(act) > 0 {
		err := p.FindByNameAndAct(name, act)
		if err != nil {
			logging.ErrorLogger.Errorf("create perm find by name get err ", err)
			return err
		}
		if p.Id > 0 {
			return errors.New(fmt.Sprintf("name %s is being used", name))
		}
	}

	return nil
}

func (p *PermResponse) Find(id uint) error {
	err := easygorm.GetEasyGormDb().Model(p.Model()).Where("id = ?", id).Find(p).Error
	if err != nil {
		logging.ErrorLogger.Errorf("find perm err ", err)
		return err
	}
	return nil
}

func (p *PermResponse) Delete(id uint) error {
	err := easygorm.GetEasyGormDb().Unscoped().Delete(p.Model(), id).Error
	if err != nil {
		logging.ErrorLogger.Errorf("delete perm by id get  err ", err)
		return err
	}
	return nil
}
