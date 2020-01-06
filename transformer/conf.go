package transformer

type Conf struct {
	App      App
	Database Database
	Mongodb  Mongodb
	Sqlite   Sqlite
	Redis    Redis
	TestData TestData
}

type App struct {
	Name string
	URl  string
	Port string
}

type Database struct {
	DirverName string
	Name       string
	UserName   string
	Password   string
}

type Mongodb struct {
	Connect string
}

type Sqlite struct {
	DirverName string
	Connect    string
}

type Redis struct {
	Addr     string
	Password string
	DB       string
}

type TestData struct {
	UserName        string
	Name            string
	Pwd             string
	DataBaseDriver  string
	DataBaseConnect string
}
