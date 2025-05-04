package admin

import (
	"testing"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/snowlyg/iris-admin/conf"
	"gorm.io/gorm"
)

func testInitMigrate() (*Migrate, error) {
	c := conf.NewConf()
	if err := c.Recover(); err != nil {
		return nil, err
	}
	db, err := gormDb(c.Mysql)
	if err != nil {
		return nil, err
	}
	return &Migrate{
		db:    db,
		items: nil,
		seeds: nil,
	}, nil
}

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
	migrate, err := testInitMigrate()
	if err != nil {
		t.Fatal(err.Error())
	}
	migrate.AddMigration(m)
	l := migrate.MigrationLen()
	if l != 1 {
		t.Errorf("MigrationLen want %d but get %d", 1, l)
	}
}

type testSeed struct{}

func (ts *testSeed) Init() error {
	return nil
}
func TestAddSeed(t *testing.T) {
	migrate, err := testInitMigrate()
	if err != nil {
		t.Fatal(err.Error())
	}
	migrate.AddSeed(&testSeed{})
	l := migrate.SeedlLen()
	if l != 1 {
		t.Errorf("SeedlLen want %d but get %d", 1, l)
	}
}
func TestMigrate(t *testing.T) {
	migrate, err := testInitMigrate()
	if err != nil {
		t.Fatal(err.Error())
	}
	migrate.AddMigration(m)
	if err := migrate.Migrate(); err != nil {
		t.Error(err.Error())
	}
}
func TestRollback(t *testing.T) {
	migrate, err := testInitMigrate()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err := migrate.Rollback(id); err != nil {
		t.Error(err.Error())
	}
	if err := migrate.Rollback(""); err != nil {
		t.Error(err.Error())
	}
	migrate.AddMigration(m)
	if err := migrate.Migrate(); err != nil {
		t.Error(err.Error())
	}

	if err := migrate.Rollback(id); err != nil {
		t.Error(err.Error())
	}
	if err := migrate.Rollback(""); err != nil {
		t.Error(err.Error())
	}
}

func TestRefresh(t *testing.T) {
	migrate, err := testInitMigrate()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err := migrate.Refresh(); err != nil {
		t.Error(err.Error())
	}

	migrate.AddMigration(m)
	if err := migrate.Migrate(); err != nil {
		t.Error(err.Error())
	}

	if err := migrate.Refresh(); err != nil {
		t.Error(err.Error())
	}

}
