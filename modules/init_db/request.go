package init_db

type Request struct {
	Sql       Sql    `json:"sql"`
	SqlType   string `json:"sqlType"`
	Cache     Cache  `json:"cache"`
	CacheType string `json:"cacheType"`
	Level     string `json:"level"` // debug,release,test
	Addr      string `json:"addr"`
}

type Sql struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password"`
	DBName   string `json:"dbName" binding:"required"`
}
type Cache struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	PoolSize int    `json:"poolSize"`
}
