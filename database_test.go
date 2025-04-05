package admin

import (
	"testing"

	"github.com/snowlyg/iris-admin/conf"
)

func TestGormDb(t *testing.T) {
	c := conf.NewConf()
	if err := c.Recover(); err != nil {
		t.Fatal(err.Error())
	}
	if c.Mysql.Password != "123456" {
		t.Errorf("mysql password want '123456' but get '%s'", c.Mysql.Password)
	}
	gormDb, err := gormDb(&c.Mysql)
	if err != nil {
		t.Fatal(err.Error())
	}
	if gormDb.Exec("show databases;").Error != nil {
		t.Fatal(err.Error())
	}
}
