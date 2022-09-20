package migration

import (
	"os"
	"strings"
	"testing"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"gorm.io/gorm"
)

type Test struct {
	gorm.Model
}

var id = "20211214120700_create_tests_table"
var m = &gormigrate.Migration{
	ID: id,
	Migrate: func(tx *gorm.DB) error {
		return tx.AutoMigrate(&Test{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("tests")
	},
}

func TestAddMigration(t *testing.T) {
	migrate := New()
	t.Run("migrate add migration", func(t *testing.T) {
		migrate.AddMigration(m)
		l := migrate.MigrationLen()
		if l != 1 {
			t.Errorf("MigrationLen want %d but get %d", 1, l)
		}
	})
}

type testSeed struct{}

func (ts *testSeed) Init() error {
	return nil
}
func TestAddSeed(t *testing.T) {
	migrate := New()
	t.Run("migrate add seed", func(t *testing.T) {
		migrate.AddSeed(&testSeed{})
		l := migrate.SeedlLen()
		if l != 1 {
			t.Errorf("SeedlLen want %d but get %d", 1, l)
		}
	})
}
func TestMigrate(t *testing.T) {
	defer zap_server.Remove()
	defer database.Remove()
	database.CONFIG.Path = strings.TrimSpace(os.Getenv("mysqlAddr"))
	database.CONFIG.Password = os.Getenv("mysqlPwd")
	migrate := New()
	migrate.AddMigration(m)
	t.Run("migrate migrate", func(t *testing.T) {
		err := migrate.Migrate()
		if err != nil {
			t.Error(err.Error())
		}
	})
}
func TestRollback(t *testing.T) {
	defer zap_server.Remove()
	defer database.Remove()
	database.CONFIG.Path = strings.TrimSpace(os.Getenv("mysqlAddr"))
	database.CONFIG.Password = os.Getenv("mysqlPwd")
	migrate := New()
	t.Run("migrate rollback no migrate with id", func(t *testing.T) {
		err := migrate.Rollback(id)
		if err != nil {
			t.Error(err.Error())
		}
	})
	t.Run("migrate rollback no migrate without id", func(t *testing.T) {
		err := migrate.Rollback("")
		if err != nil {
			t.Error(err.Error())
		}
	})
	migrate.AddMigration(m)
	err := migrate.Migrate()
	if err != nil {
		t.Error(err.Error())
	}

	t.Run("migrate rollback after migrate with id", func(t *testing.T) {
		err := migrate.Rollback(id)
		if err != nil {
			t.Error(err.Error())
		}
	})
	t.Run("migrate rollback after migrate without id", func(t *testing.T) {
		err := migrate.Rollback("")
		if err != nil {
			t.Error(err.Error())
		}
	})
}

func TestRefresh(t *testing.T) {
	defer zap_server.Remove()
	defer database.Remove()
	database.CONFIG.Path = strings.TrimSpace(os.Getenv("mysqlAddr"))
	database.CONFIG.Password = os.Getenv("mysqlPwd")
	migrate := New()

	t.Run("migrate refresh no migrate", func(t *testing.T) {
		err := migrate.Refresh()
		if err != nil {
			t.Error(err.Error())
		}
	})

	migrate.AddMigration(m)
	err := migrate.Migrate()
	if err != nil {
		t.Error(err.Error())
	}

	t.Run("migrate refresh after migrate", func(t *testing.T) {
		err := migrate.Refresh()
		if err != nil {
			t.Error(err.Error())
		}
	})

}
