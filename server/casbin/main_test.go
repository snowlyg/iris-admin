package casbin

import (
	_ "embed"
	"os"
	"testing"

	"github.com/bwmarrin/snowflake"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/zap_server"
)

func TestMain(m *testing.M) {

	node, _ := snowflake.NewNode(1)
	uuid := str.Join("casbin", "_", node.Generate().String())

	database.CONFIG.DbName = uuid
	database.CONFIG.Path = g.TestMysqlAddr
	database.CONFIG.Password = g.TestMysqlPwd

	Instance()

	code := m.Run()

	err := database.DorpDB(database.CONFIG.BaseDsn(), "mysql", uuid)
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
	}

	db, _ := database.Instance().DB()
	if db != nil {
		db.Close()
	}
	Remove()
	zap_server.Remove()
	database.Remove()
	os.Exit(code)
}
