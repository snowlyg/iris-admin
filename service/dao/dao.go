package dao

import (
	"errors"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/application/libs/easygorm"
	"github.com/snowlyg/blog/application/libs/logging"
	"github.com/snowlyg/blog/application/models"
	"github.com/snowlyg/blog/service/auth"
)

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
	jwt, ok := ctx.Values().Get("jwt").(*jwt.Token)
	if !ok {
		return 0, errors.New("jwt is nil")
	}
	id, err := authDriver.GetAuthId(jwt.Raw)
	if err != nil {
		return 0, err
	}
	return id, err
}

// CreateOplog
func CreateOplog(oplog models.Oplog) error {
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

	return all, err
}

func Create(d Dao, ctx iris.Context, object map[string]interface{}) error {
	err := d.Create(object)
	if err != nil {
		logging.ErrorLogger.Errorf("dao create get err ", err)
		return err
	}

	return nil
}

func Update(d Dao, ctx iris.Context, object map[string]interface{}) error {
	id, _ := getId(ctx)
	err := d.Update(id, object)
	if err != nil {
		logging.ErrorLogger.Errorf("dao update get err ", err)
		return err
	}

	return nil
}

func Find(d Dao, ctx iris.Context) error {
	id, _ := getId(ctx)
	err := d.Find(id)
	if err != nil {
		logging.ErrorLogger.Errorf("dao find by id  get err ", err)
		return err
	}

	return nil
}

func getId(ctx iris.Context) (uint, error) {
	id, err := ctx.Params().GetUint("id")
	if err != nil {
		logging.ErrorLogger.Errorf("dao get id get err ", err)
		return 0, err
	}

	return id, nil
}

func Delete(d Dao, ctx iris.Context) error {
	id, _ := getId(ctx)
	err := d.Delete(id)
	if err != nil {
		logging.ErrorLogger.Errorf("dao delete  get err ", err)
		return err
	}
	return nil
}
