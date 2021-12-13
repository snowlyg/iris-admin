package g

const (
	ConfigType     = "yaml"            // 配置文件类型
	ConfigDir      = "config"          // 配置目录
	CasbinFileName = "rbac_model.conf" // casbin 规则文件名称
)

const (
	StatusUnknown int = iota
	StatusTrue
	StatusFalse
)
