package g

const (
	ConfigType     = "yaml"            // 配置文件类型
	ConfigDir      = "config"          // 配置目录
	CasbinFileName = "rbac_model.conf" // casbin 规则文件名称
)

// 状态值 0:未知,1:true,2:false
const (
	StatusUnknown int = iota
	StatusTrue
	StatusFalse
)
