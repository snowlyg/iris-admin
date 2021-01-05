package easygorm

import (
//"gorm.io/gorm"
//"gorm.io/gorm/schema"
//"path/filepath"
//"testing"
)

//func TestOrm(t *testing.T) {
//	t.Run("TestInitMysql", func(t *testing.T) {
//		Init(&Config{
//			GormConfig: &gorm.Config{
//				NamingStrategy: schema.NamingStrategy{
//					TablePrefix:   "iris_", // 表名前缀，`User` 的表名应该是 `t_users`
//					SingularTable: false,   // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
//				},
//			},
//			Adapter:         "mysql", // 类型
//			Name:            "blog_test",
//			Username:        "root",                                  // 用户名
//			Pwd:             "",                                      // 密码
//			Host:            "127.0.0.1",                             // 地址
//			Port:            3306,                                    // 端口
//			CasbinModelPath: filepath.Join(cwd(), "rbac_model.conf"), // casbin 模型规则路径
//		})
//		if Egm.Db == nil {
//			t.Errorf("TestInitMysql error")
//		}
//		if Egm.Enforcer == nil {
//			t.Errorf("TestInitMysql error")
//		}
//	})
//}
