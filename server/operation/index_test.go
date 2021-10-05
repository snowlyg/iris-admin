package operation

import (
	"testing"
)

func TestCreateOplog(t *testing.T) {
	t.Run("test create op log", func(t *testing.T) {
		record := &Oplog{}
		if err := CreateOplog(record); err != nil {
			t.Errorf("test create op log get %s", err.Error())
		}
	})
}
