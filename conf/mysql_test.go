package conf

import "testing"

func TestMysqlBaseDsn(t *testing.T) {
	m := &Mysql{
		Path:         "127.0.0.1:3306",
		Config:       "charset=utf8mb4",
		DbName:       "db_name",
		Username:     "name",
		Password:     "pwd",
		MaxIdleConns: 0,
		MaxOpenConns: 0,
		LogMode:      false,
		LogZap:       "",
	}
	b := m.BaseDsn()
	want := "name:pwd@tcp(127.0.0.1:3306)/"
	if b != want {
		t.Errorf("mysql config base dsn want '%s' but get '%s'", want, b)
	}
	dsn := m.Dsn()
	wantDsn := "name:pwd@tcp(127.0.0.1:3306)/charset=utf8mb4"
	if b != wantDsn {
		t.Errorf("mysql config base dsn want '%s' but get '%s'", wantDsn, dsn)
	}
}
