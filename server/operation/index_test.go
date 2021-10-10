package operation

import (
	"testing"

	"github.com/snowlyg/iris-admin/server/database"
)

func TestCreateOplog(t *testing.T) {
	database.Instance().AutoMigrate(&Oplog{})
	t.Run("test create op log", func(t *testing.T) {
		record := &Oplog{}
		if err := CreateOplog(record); err != nil {
			t.Errorf("test create op log get %s", err.Error())
		}
	})
}
