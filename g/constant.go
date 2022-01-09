package g

const (
	ConfigType     = "json"            // config's type
	ConfigDir      = "config"          // config's dir
	CasbinFileName = "rbac_model.conf" // casbin's name
)

// status 0:unkown,1:true,2:false
const (
	StatusUnknown int = iota
	StatusTrue
	StatusFalse
)
