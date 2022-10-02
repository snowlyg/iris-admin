package operation

import (
	"encoding/json"
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

// IsExist config file is exist
func IsExist() bool {
	return getViperConfig().IsFileExist()
}

// Remove remove config file
func Remove() error {
	return getViperConfig().Remove()
}

// Recover
func Recover() error {
	b, err := json.Marshal(CONFIG)
	if err != nil {
		return err
	}
	return getViperConfig().Recover(b)
}

// Operation
// Except set which routers don't generate system log, use ';' to separate.
// Include set which routers need to generate system log, use ';' to separate.
type Operation struct {
	Except  Route `mapstructure:"except" json:"except" yaml:"except"`
	Include Route `mapstructure:"include" json:"include" yaml:"include"`
}

type Route struct {
	Uri    string `mapstructure:"uri" json:"uri" yaml:"uri"`
	Method string `mapstructure:"method" json:"method" yaml:"method"`
}

// GetExcept return routers which need to excepted
func (op Operation) GetExcept() ([]string, []string) {
	uri := strings.Split(op.Except.Uri, ";")
	method := strings.Split(op.Except.Method, ";")
	return uri, method
}

// GetInclude return routers which need to included
func (op Operation) GetInclude() ([]string, []string) {
	uri := strings.Split(op.Include.Uri, ";")
	method := strings.Split(op.Include.Method, ";")
	return uri, method
}

// IsInclude check whether the current route needs to belong to the included data
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

// IsExcept check whether the current route needs to belong to the excepted data
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

// getViperConfig get viper config
func getViperConfig() viper_server.ViperConfig {
	configName := "operation"
	return viper_server.ViperConfig{
		Debug:     true,
		Directory: g.ConfigDir,
		Name:      configName,
		Type:      g.ConfigType,
		Watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(&CONFIG); err != nil {
				return fmt.Errorf("get Unarshal error: %v", err)
			}
			// watch config file change
			vi.SetConfigName(configName)
			return nil
		},
		//
		Default: []byte(`
{
	"except":{ 
		"uri": "` + CONFIG.Except.Uri + `",
		"method": "` + CONFIG.Except.Method + `"
	},	
  "include":
	{
		"uri": "` + CONFIG.Include.Uri + `",
		"method": "` + CONFIG.Include.Method + `"
	} 
 }`),
	}
}
