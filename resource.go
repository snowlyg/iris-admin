package admin

import (
	"github.com/gin-gonic/gin"
)

type Model interface {
	TableName() string
	List() []map[string]any
	// Detail() map[string]any
	// Create() error
	// Update() error
	// Delete() error
}

func (ws *WebServe) Resource(group *gin.RouterGroup, model Model) {
	r := group.Group(model.TableName())
	{
		// list,create,update,delete,detail
		// should add database operation
		r.GET("/list", func(ctx *gin.Context) {
			if ws.db == nil {
				FailWithMessage("database not found", ctx)
				return
			}
			list := model.List()

			if err := ws.db.Table(model.TableName()).Scopes(SoftDeleteScope()).Find(&list).Error; err != nil {
				FailWithMessage(err.Error(), ctx)
				return
			}
			OkWithData(list, ctx)
		})
	}
}
