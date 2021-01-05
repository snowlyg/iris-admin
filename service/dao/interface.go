package dao

import (
	"encoding/json"
	"fmt"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/snowlyg/blog/application/libs/easygorm"
	"github.com/snowlyg/blog/application/models"
	"github.com/snowlyg/blog/service/auth"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/application/libs/logging"
)

type Dao interface {
	ModelName() string
	All(name, sort, orderBy string, page, pageSize int) (map[string]interface{}, error)
	Create(object map[string]interface{}) error
	Update(id uint, object map[string]interface{}) error
	Find(id uint) error
	Delete(id uint) error
}

const (
	ActionList   = "列表查询"
	ActionOne    = "单个查询"
	ActionAdd    = "添加"
	ActionUpdate = "更新"
	ActionDel    = "删除"
	ActionLogin  = "登录"
	ActionLogout = "登出"
)

// GetAuthId
func GetAuthId(ctx iris.Context) (uint, error) {
	authDriver := auth.NewAuthDriver()
	defer authDriver.Close()
	token := ctx.Values().Get("jwt").(*jwt.Token).Raw
	id, err := authDriver.GetAuthId(token)
	if err != nil {
		logging.ErrorLogger.Errorf("get user id err", err)
		return 0, err
	}
	return id, err
}

// Add
func Add(ctx iris.Context, model, action, content string) error {
	uId, err := GetAuthId(ctx)
	if err != nil {
		logging.ErrorLogger.Errorf("dao all get auth id get err ", err)
		return err
	}
	err = CreateOplog(model, action, content, uId)
	if err != nil {
		return err
	}
	return nil
}

// CreateOplog
func CreateOplog(model string, action string, content string, uId uint) error {
	oplog := models.Oplog{
		ModelName:  model,
		ActionName: action,
		Content:    content,
		UserID:     uId,
	}
	err := easygorm.GetEasyGormDb().Model(&models.Oplog{}).Create(&oplog).Error
	if err != nil {
		logging.ErrorLogger.Errorf("add oplog  get err ", err)
		return err
	}
	return nil
}

func All(d Dao, ctx iris.Context, name, sort, orderBy string, page, pageSize int) (map[string]interface{}, error) {
	all, err := d.All(name, sort, orderBy, page, pageSize)
	if err != nil {
		logging.ErrorLogger.Errorf("dao all get err ", err)
		return nil, err
	}
	var content []byte
	content, err = json.Marshal(all)
	err = Add(ctx, d.ModelName(), ActionList, string(content))
	if err != nil {
		logging.ErrorLogger.Errorf("dao all add oplog get err ", err)
		return nil, err
	}

	return all, err
}

func Create(d Dao, ctx iris.Context, object map[string]interface{}) error {
	err := d.Create(object)
	if err != nil {
		logging.ErrorLogger.Errorf("dao create get err ", err)
		return err
	}
	var content []byte
	content, err = json.Marshal(object)
	err = Add(ctx, d.ModelName(), ActionAdd, string(content))
	if err != nil {
		logging.ErrorLogger.Errorf("dao create add oplog get err ", err)
		return err
	}

	return nil
}

func Update(d Dao, ctx iris.Context, id uint, object map[string]interface{}) error {
	err := d.Update(id, object)
	if err != nil {
		logging.ErrorLogger.Errorf("dao update get err ", err)
		return err
	}
	var content []byte
	content, err = json.Marshal(object)
	err = Add(ctx, d.ModelName(), ActionUpdate, string(content))
	if err != nil {
		logging.ErrorLogger.Errorf("dao update add oplog get err ", err)
		return err
	}

	return nil
}

func Find(d Dao, ctx iris.Context, id uint) error {
	err := d.Find(id)
	if err != nil {
		logging.ErrorLogger.Errorf("dao find by id  get err ", err)
		return err
	}
	err = Add(ctx, d.ModelName(), ActionOne, fmt.Sprintf("%d", id))
	if err != nil {
		logging.ErrorLogger.Errorf("dao find by id add oplog get err ", err)
		return err
	}

	return nil
}

func Delete(d Dao, ctx iris.Context, id uint) error {
	err := d.Delete(id)
	if err != nil {
		logging.ErrorLogger.Errorf("dao delete  get err ", err)
		return err
	}
	err = Add(ctx, d.ModelName(), ActionDel, fmt.Sprintf("%d", id))
	if err != nil {
		logging.ErrorLogger.Errorf("dao delete add oplog get err ", err)
		return err
	}
	return nil
}
