package init_db

type Request struct {
	Sql       Sql    `json:"sql"`
	SqlType   string `json:"sqlType" validate:"required"`
	Cache     Cache  `json:"cache"`
	CacheType string `json:"cacheType"  validate:"required"`
	Level     string `json:"level"  validate:"required"` // debug,release,test
	Addr      string `json:"addr"  validate:"required"`
}

type Sql struct {
	Host     string `json:"host"  validate:"required"`
	Port     string `json:"port"  validate:"required"`
	UserName string `json:"userName" validate:"required"`
	Password string `json:"password"  validate:"required"`
	DBName   string `json:"dbName" validate:"required"`
}
type Cache struct {
	Host     string `json:"host"  validate:"required"`
	Port     string `json:"port"  validate:"required"`
	Password string `json:"password"`
	PoolSize int    `json:"poolSize"`
}
