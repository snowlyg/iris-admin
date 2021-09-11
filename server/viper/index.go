package viper

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/iris-admin/g"
	"github.com/spf13/viper"
)

func Init() {
	config := g.ConfigFileName
	fmt.Printf("您的配置文件路径为%s\n", config)

	v := viper.New()
	v.SetConfigType("yaml")

	if !dir.IsExist(config) { //没有配置文件，写入默认配置
		var yamlDefault = []byte(`
captcha:
 key-long: 6
 img-width: 240
 img-height: 80
limit:
 limit: true
 limit: 1
 burst: 5
mysql:
 path: ""
 config: charset=utf8mb4&parseTime=True&loc=Local
 db-name: ""
 username: ""
 password: ""
 max-idle-conns: 0
 max-open-conns: 0
 log-mode: false
 log-zap: ""
redis:
 db: 0
 addr: ""
 password: ""
 pool-size: 0
system:
 level: debug # debug,release,test
 addr: "127.0.0.1:8085"
 db-type: ""
 cache-type: "" 
 time-format: "2006-01-02 15:04:05"
zap:
 level: info
 format: console
 prefix: '[OP-ONLINE]'
 director: log
 link-name: latest_log
 showLine: false
 encode-level: LowercaseColorLevelEncoder
 stacktrace-key: stacktrace
 log-in-console: true`)

		if err := v.ReadConfig(bytes.NewBuffer(yamlDefault)); err != nil {
			panic(fmt.Errorf("读取默认配置文件错误: %w ", err))
		}

		if err := v.Unmarshal(&g.CONFIG); err != nil {
			panic(fmt.Errorf("同步配置文件错误: %w ", err))
		}

		if err := v.WriteConfigAs(config); err != nil {
			panic(fmt.Errorf("写入配置文件错误: %w ", err))
		}
		return
	}

	// 存在配置文件，读取配置文件内容
	v.SetConfigFile(config)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取配置错误: %w ", err))
	}

	// 监控配置文件变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置发生变化:", e.Name)
		if err := v.Unmarshal(&g.CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&g.CONFIG); err != nil {
		fmt.Println(err)
	}
	g.VIPER = v

	casbinPath := filepath.Join(dir.GetCurrentAbPath(), g.CasbinFileName)
	fmt.Printf("casbin rbac_model.conf 位于： %s\n\n", casbinPath)
	if !dir.IsExist(casbinPath) { // casbin rbac_model.conf 文件
		var rbacModelConf = []byte(`[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch3(r.obj, p.obj) && (r.act == p.act || p.act == "*")`)

		_, err = dir.WriteBytes(casbinPath, rbacModelConf)
		if err != nil {
			panic(fmt.Errorf("初始化 casbin rbac_model.conf 文件错误: %w ", err))
		}
	}
}
