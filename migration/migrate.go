package migration

import (
	"errors"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/snowlyg/iris-admin/seed"
	"github.com/snowlyg/iris-admin/server/database"
)

// MigrationCmd migration cmd
// MigrationCollection migration collections
// SeedCollection data seed collection
type MigrationCmd struct {
	MigrationCollection []*gormigrate.Migration
	SeedCollection      []seed.SeedFunc
}

// New MigrationCmd
func New() *MigrationCmd {
	mc := &MigrationCmd{
		MigrationCollection: nil,
		SeedCollection:      nil,
	}

	return mc
}

// AddMigration add *gormigrate.Migration
func (mc *MigrationCmd) AddMigration(m ...*gormigrate.Migration) {
	mc.MigrationCollection = append(mc.MigrationCollection, m...)
}

// MigrationLen length of MigrationCollection
func (mc *MigrationCmd) MigrationLen() int {
	return len(mc.MigrationCollection)
}

// AddSeed add SeedFunc
func (mc *MigrationCmd) AddSeed(sf ...seed.SeedFunc) {
	mc.SeedCollection = append(mc.SeedCollection, sf...)
}

// SeedlLen length of  SeedCollection
func (mc *MigrationCmd) SeedlLen() int {
	return len(mc.SeedCollection)
}

// Refresh refresh migration
func (mc *MigrationCmd) Refresh() error {
	if mc.getFirstMigration() == "" {
		return nil
	}
	err := mc.rollbackTo(mc.getFirstMigration())
	if !errors.Is(gormigrate.ErrMigrationIDDoesNotExist, err) && err != nil {
		return err
	}
	return mc.Migrate()
}

// rollbackTo roolback migration to migrationId
func (mc *MigrationCmd) rollbackTo(migrationId string) error {
	return mc.gormigrate().RollbackTo(migrationId)
}

// Rollback roolback migrations
func (mc *MigrationCmd) Rollback(migrationId string) error {
	if mc.MigrationLen() == 0 {
		return nil
	}
	if migrationId == "" {
		err := mc.rollbackLast()
		if !errors.Is(gormigrate.ErrMigrationIDDoesNotExist, err) && err != nil {
			return err
		}
		return nil
	}
	err := mc.rollbackTo(migrationId)
	if !errors.Is(gormigrate.ErrMigrationIDDoesNotExist, err) && err != nil {
		return err
	}
	return nil
}

// rollbackLast roolback the lasted migration
func (mc *MigrationCmd) rollbackLast() error {
	return mc.gormigrate().RollbackLast()
}

// Migrate exec migration cmd
func (mc *MigrationCmd) Migrate() error {
	m := mc.gormigrate()
	err := m.Migrate()
	if err != nil {
		return err
	}
	return nil
}

// Seed seed data into database
func (mc *MigrationCmd) Seed() error {
	if mc.SeedCollection == nil {
		return nil
	}
	return seed.Seed(mc.SeedCollection...)
}

// getFirstMigration get first migration's id
func (mc *MigrationCmd) getFirstMigration() string {
	if mc.MigrationLen() == 0 {
		return ""
	}
	return mc.MigrationCollection[0].ID
}

// gormigrate create *gormigrate.Gormigrate
func (mc *MigrationCmd) gormigrate() *gormigrate.Gormigrate {
	return gormigrate.New(database.Instance(), gormigrate.DefaultOptions, mc.MigrationCollection)
}
