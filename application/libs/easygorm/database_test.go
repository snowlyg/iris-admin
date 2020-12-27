package easygorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"testing"
)

func TestInitMysql(t *testing.T) {
	t.Run("TestInitMysql", func(t *testing.T) {
		Init(&Config{
			GormConfig: &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					TablePrefix:   "iris_", // 表名前缀，`User` 的表名应该是 `t_users`
					SingularTable: false,   // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
				},
			},
			Adapter: "mysql", // 类型
		})
		if EasyGorm.DB == nil {
			t.Errorf("TestInitMysql error")
		}
		if EasyGorm.Enforcer == nil {
			t.Errorf("TestInitMysql error")
		}
	})
}

func TestInitSqlite(t *testing.T) {
	t.Run("TestInitSqlite", func(t *testing.T) {
		Init(&Config{
			GormConfig: &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					TablePrefix:   "iris_", // 表名前缀，`User` 的表名应该是 `t_users`
					SingularTable: false,   // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
				},
			},
			Adapter: "sqlite3", // 类型
		})
		if EasyGorm.DB == nil {
			t.Errorf("TestInitSqlite error")
		}
		if EasyGorm.Enforcer == nil {
			t.Errorf("TestInitSqlite error")
		}
	})
}
