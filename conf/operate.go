package conf

import (
	"strings"
)

var CONFIG = Operate{
	Except: Route{
		Uri:    "api/v1/upload;api/v1/upload",
		Method: "post;put",
	},
	Include: Route{
		Uri:    "api/v1/menus",
		Method: "get",
	},
}

// Operate
// Except set which routers don't generate system log, use ';' to separate.
// Include set which routers need to generate system log, use ';' to separate.
type Operate struct {
	Except  Route `mapstructure:"except" json:"except" yaml:"except"`
	Include Route `mapstructure:"include" json:"include" yaml:"include"`
}

// type Route struct {
// 	Uri    string `mapstructure:"uri" json:"uri" yaml:"uri"`
// 	Method string `mapstructure:"method" json:"method" yaml:"method"`
// }

// GetExcept return routers which need to excepted
func (op Operate) GetExcept() ([]string, []string) {
	uri := strings.Split(op.Except.Uri, ";")
	method := strings.Split(op.Except.Method, ";")
	return uri, method
}

// GetInclude return routers which need to included
func (op Operate) GetInclude() ([]string, []string) {
	uri := strings.Split(op.Include.Uri, ";")
	method := strings.Split(op.Include.Method, ";")
	return uri, method
}

// IsInclude check whether the current route needs to belong to the included data
func (op Operate) IsInclude(uri, method string) bool {
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
func (op Operate) IsExcept(uri, method string) bool {
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
