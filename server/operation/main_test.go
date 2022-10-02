package operation

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
	uuid := str.Join("operation", "_", node.Generate().String())

	database.CONFIG.Dbname = uuid
	database.CONFIG.Path = g.TestMysqlAddr
	database.CONFIG.Password = g.TestMysqlPwd
	database.Instance()

	code := m.Run()

	err := database.DorpDB(database.CONFIG.BaseDsn(), "mysql", uuid)
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
	}

	db, _ := database.Instance().DB()
	if db != nil {
		db.Close()
	}
	database.Remove()
	Remove()
	zap_server.Remove()
	os.Exit(code)
}
