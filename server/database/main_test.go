package database

import (
	_ "embed"
	"os"
	"strings"
	"testing"

	"github.com/bwmarrin/snowflake"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/server/zap_server"
)

func TestMain(m *testing.M) {

	node, _ := snowflake.NewNode(1)
	uuid := str.Join("database", "_", node.Generate().String())

	CONFIG.Dbname = uuid
	CONFIG.Path = strings.TrimSpace(os.Getenv("mysqlAddr"))
	CONFIG.Password = strings.TrimSpace(os.Getenv("mysqlPwd"))

	Instance()

	code := m.Run()

	err := DorpDB(CONFIG.BaseDsn(), "mysql", uuid)
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
	}

	db, _ := Instance().DB()
	if db != nil {
		db.Close()
	}
	Remove()
	zap_server.Remove()
	os.Exit(code)
}
