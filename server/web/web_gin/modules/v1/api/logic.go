package api

import (
	"github.com/snowlyg/iris-admin/server/casbin"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/database/scope"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreatenInBatches 批量加入
func CreatenInBatches(db *gorm.DB, perms ApiCollection) error {
	err := db.Model(&Api{}).CreateInBatches(&perms, 500).Error
	if err != nil {
		zap_server.ZAPLOG.Error("批量导入权限", zap.String("错误:", err.Error()))
		return err
	}
	return nil
}

func Delete(req DeleteApiReq) error {
	if err := orm.Delete(database.Instance(), &Api{}, scope.IdScope(req.Id)); err != nil {
		return err
	}
	if err := casbin.ClearCasbin(1, req.Path, req.Method); err != nil {
		return err
	}
	return nil
}

func BatcheDelete(ids []uint) error {
	apis := &PageResponse{}
	err := orm.Find(database.Instance(), apis, scope.InIdsScope(ids))
	if err != nil {
		return err
	}

	err = database.Instance().Transaction(func(tx *gorm.DB) error {
		if err := orm.Delete(tx, &Api{}, scope.InIdsScope(ids)); err != nil {
			return err
		}
		for _, api := range apis.Item {
			if err := casbin.ClearCasbin(1, api.Path, api.Method); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
