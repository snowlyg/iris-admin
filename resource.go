package admin

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	databaseUnavailableMessage = "database unavailable"
	databaseQueryFailedMessage = "database query failed"
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
				ErrorWithStatus(http.StatusInternalServerError, databaseUnavailableMessage, ctx)
				return
			}
			list := model.List()

			if err := ws.db.Table(model.TableName()).Scopes(SoftDeleteScope()).Find(&list).Error; err != nil {
				log.Printf("iris-admin: list resource table=%q failed: %v", model.TableName(), err)
				ErrorWithStatus(http.StatusInternalServerError, databaseQueryFailedMessage, ctx)
				return
			}
			OkWithData(list, ctx)
		})
	}
}
