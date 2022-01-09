package casbin

import (
	_ "embed"
	"os"
	"strings"
	"testing"

	"github.com/bwmarrin/snowflake"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/zap_server"
)

func TestMain(m *testing.M) {
	defer Remove()
	defer zap_server.Remove()
	defer database.Remove()

	node, _ := snowflake.NewNode(1)
	uuid := str.Join("database", "_", node.Generate().String())

	database.CONFIG.Dbname = uuid
	database.CONFIG.Password = strings.TrimSpace(os.Getenv("mysqlPwd"))

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

	os.Exit(code)
}
