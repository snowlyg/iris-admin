package g

const (
	ConfigType     = "yaml"            // 配置文件类型
	ConfigDir      = "config"          // 配置目录
	CasbinFileName = "rbac_model.conf" // casbin 规则文件名称
)

const (
	AdminAuthorityId   = 999
	TenancyAuthorityId = 998
	LiteAuthorityId    = 997 // 小程序用户
	DeviceAuthorityId  = 996 // 床旁设备用户
)

const (
	StatusUnknown int = iota
	StatusTrue
	StatusFalse
)
