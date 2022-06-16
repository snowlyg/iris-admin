package operation

import (
	"fmt"
	"strings"

	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/viper_server"
	"github.com/spf13/viper"
)

var CONFIG = Operation{
	Except: Route{
		Uri:    "api/v1/upload;api/v1/upload",
		Method: "post;put",
	},
	Include: Route{
		Uri:    "api/v1/menus",
		Method: "get",
	},
}

// IsExist 配置文件是否存在
func IsExist() bool {
	return getViperConfig().IsFileExist()
}

// Remove 删除配置文件
func Remove() error {
	err := getViperConfig().Remove()
	if err != nil {
		return err
	}
	return nil
}

// Operation 操作日志配置
// Except 排除生成操作日志的路由,多条使用 ; 号分割
// Include 包括生成操作日志的路由,多条使用 ; 号分割
type Operation struct {
	Except  Route `mapstructure:"except" json:"except" yaml:"except"`
	Include Route `mapstructure:"include" json:"include" yaml:"include"`
}

type Route struct {
	Uri    string `mapstructure:"uri" json:"uri" yaml:"uri"`
	Method string `mapstructure:"method" json:"method" yaml:"method"`
}

//  GetExcept 返回需要排除路由数组数据
func (op Operation) GetExcept() ([]string, []string) {
	uri := strings.Split(op.Except.Uri, ";")
	method := strings.Split(op.Except.Method, ";")
	return uri, method
}

// GetInclude 返回需要包含路由数组数据
func (op Operation) GetInclude() ([]string, []string) {
	uri := strings.Split(op.Include.Uri, ";")
	method := strings.Split(op.Include.Method, ";")
	return uri, method
}

// IsInclude 判断当前路由是否需要属于包含数据中
func (op Operation) IsInclude(uri, method string) bool {
	incUri, incMethod := op.GetInclude()
	if len(incUri) != len(incMethod) {
		return false
	}

	for i := 0; i < len(incUri); i++ {
		if uri == incUri[i] && method == incMethod[i] {
			return true
		}
	}
	return false
}

// IsExcept 判断当前路由是否需要属于排除数据中
func (op Operation) IsExcept(uri, method string) bool {
	excUri, excMethod := op.GetExcept()
	if len(excUri) != len(excMethod) {
		return false
	}

	for i := 0; i < len(excUri); i++ {
		if uri == excUri[i] && method == excMethod[i] {
			return true
		}
	}
	return false
}

// getViperConfig 获取初始化配置
func getViperConfig() viper_server.ViperConfig {
	configName := "operation"
	return viper_server.ViperConfig{
		Debug:     true,
		Directory: g.ConfigDir,
		Name:      configName,
		Type:      g.ConfigType,
		Watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(&CONFIG); err != nil {
				return fmt.Errorf("反序列化错误: %v", err)
			}
			// 监控配置文件变化
			vi.SetConfigName(configName)
			return nil
		},
		// 注意:设置默认配置值的时候,前面不能有空格等其他符号.必须紧贴左侧.
		Default: []byte(`
{
	"except":{ 
		uri: "` + CONFIG.Except.Uri + `",
		method: "` + CONFIG.Except.Method + `"
	},	
"include":
	{
		"uri": "` + CONFIG.Include.Uri + `",
		"method": "` + CONFIG.Include.Method + `"
	} 
 }`),
	}
}
