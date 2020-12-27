package easygorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"testing"
)

func TestInitMysql(t *testing.T) {
	t.Run("TestInitMysql", func(t *testing.T) {
		err := Init(&Config{
			GormConfig: &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					TablePrefix:   "iris_", // 表名前缀，`User` 的表名应该是 `t_users`
					SingularTable: false,   // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
				},
			},
			Casbin: &Casbin{
				Path:   "/Users/snowlyg/go/src/github.com/snowlyg/blog/rbac_model.conf",
				Prefix: "casbin",
			},
			Adapter: "mysql", // 类型
			Conn:    "root:123456@tcp(localhost:3306)/test?parseTime=True&loc=Local",
		})
		if err != nil {
			t.Errorf("TestInitSqlite DB error %+v", err)
		}
		if EasyGorm.DB == nil {
			t.Errorf("TestInitMysql DB error")
		}
		if EasyGorm.Enforcer == nil {
			t.Errorf("TestInitMysql Enforcer error")
		}
	})
}

func TestInitSqlite(t *testing.T) {
	t.Run("TestInitSqlite", func(t *testing.T) {
		err := Init(&Config{
			GormConfig: &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					TablePrefix:   "iris_", // 表名前缀，`User` 的表名应该是 `t_users`
					SingularTable: false,   // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
				},
			},
			Casbin: &Casbin{
				Path:   "/Users/snowlyg/go/src/github.com/snowlyg/blog/rbac_model.conf",
				Prefix: "casbin",
			},
			Adapter: "sqlite3", // 类型
			Conn:    "/Users/snowlyg/go/src/github.com/snowlyg/blog/test.db",
		})

		if err != nil {
			t.Errorf("TestInitSqlite DB error %+v", err)
		}

		if EasyGorm.DB == nil {
			t.Errorf("TestInitSqlite DB error")
		}
		if EasyGorm.Enforcer == nil {
			t.Errorf("TestInitSqlite Enforcer error")
		}
	})
}
