// +build test

package easygorm

import (
	"github.com/snowlyg/blog/libs"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"path/filepath"
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
			Adapter:         "mysql", // 类型
			Name:            "blog_test",
			Username:        "root",                                       // 用户名
			Pwd:             "123456",                                     // 密码
			Host:            "127.0.0.1",                                  // 地址
			Port:            3305,                                         // 端口
			CasbinModelPath: filepath.Join(libs.CWD(), "rbac_model.conf"), // casbin 模型规则路径
		})
		if Egm.Db == nil {
			t.Errorf("TestInitMysql error")
		}
		if Egm.Enforcer == nil {
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
			Adapter:         "sqlite3", // 类型
			Name:            "blog_test",
			CasbinModelPath: filepath.Join(libs.CWD(), "rbac_model.conf"), // casbin 模型规则路径
		})
		if Egm.Db == nil {
			t.Errorf("TestInitSqlite error")
		}
		if Egm.Enforcer == nil {
			t.Errorf("TestInitSqlite error")
		}
	})
}
func TestCwd(t *testing.T) {
	t.Run("TestCwdWithEnv", func(t *testing.T) {
		getenv := os.Getenv("TRAVIS_BUILD_DIR")
		if cwd() != getenv {
			t.Errorf("TestCwd %s != %s", cwd(), getenv)
		}

	})
	t.Run("TestCwdNoEnv", func(t *testing.T) {
		os.Setenv("TRAVIS_BUILD_DIR", "")
		path, _ := os.Executable()

		if cwd() != filepath.Dir(path) {
			t.Errorf("TestCwd %s != %s", cwd(), filepath.Dir(path))
		}

	})
}
