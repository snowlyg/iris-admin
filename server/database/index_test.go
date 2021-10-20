package database

import (
	"reflect"
	"testing"

	"gorm.io/gorm/logger"
)

func TestInstanceMysql(t *testing.T) {
	t.Run("test instance mysql", func(t *testing.T) {
		mysql := Instance()
		if mysql == nil {
			t.Error("mysql instance is nil")
		}
	})
}

func TestGormMysql(t *testing.T) {
	t.Run("test gorm mysql", func(t *testing.T) {
		gormDb := gormMysql()
		if gormDb == nil {
			t.Error("gorm db is nil")
		}
	})
}

func TestGormConfig(t *testing.T) {
	t.Run("test gorm config", func(t *testing.T) {
		gormConfig := gormConfig(false)
		if !reflect.DeepEqual(gormConfig.Logger, Default.LogMode(logger.Silent)) {
			t.Errorf("gorm config logger want %+v but get %+v", Default.LogMode(logger.Silent), gormConfig.Logger)
		}
	})
}
