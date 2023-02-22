package g

const (
	ConfigType     = "json"            // config's type
	ConfigDir      = "config"          // config's dir
	CasbinFileName = "rbac_model.conf" // casbin rule file's name
)

// Status 0:unknown,1:true,2:false
const (
	StatusUnknown int = iota
	StatusTrue
	StatusFalse
)
