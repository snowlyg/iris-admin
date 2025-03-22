package admin

import (
	"errors"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Migrate migration cmd
// MigrationCollection migration collections
// SeedCollection data seed collection
type Migrate struct {
	db    *gorm.DB
	items []*gormigrate.Migration
	seeds []SeedFunc
}

// New MigrationCmd
func New() *Migrate {
	mc := &Migrate{
		items: nil,
		seeds: nil,
	}

	return mc
}

// AddMigration add *gormigrate.Migration
func (mc *Migrate) AddMigration(m ...*gormigrate.Migration) {
	mc.items = append(mc.items, m...)
}

// MigrationLen length of MigrationCollection
func (mc *Migrate) MigrationLen() int {
	return len(mc.items)
}

// AddSeed add SeedFunc
func (mc *Migrate) AddSeed(sf ...SeedFunc) {
	mc.seeds = append(mc.seeds, sf...)
}

// SeedlLen length of  SeedCollection
func (mc *Migrate) SeedlLen() int {
	return len(mc.seeds)
}

// Refresh refresh migration
func (mc *Migrate) Refresh() error {
	if mc.getFirstMigration() == "" {
		return nil
	}
	err := mc.rollbackTo(mc.getFirstMigration())
	if !errors.Is(err, gormigrate.ErrMigrationIDDoesNotExist) && err != nil {
		return err
	}
	return mc.Migrate()
}

// rollbackTo roolback migration to migrationId
func (mc *Migrate) rollbackTo(migrationId string) error {
	return mc.gormigrate().RollbackTo(migrationId)
}

// Rollback roolback migrations
func (mc *Migrate) Rollback(migrationId string) error {
	if mc.MigrationLen() == 0 {
		return nil
	}
	if migrationId == "" {
		err := mc.rollbackLast()
		if !errors.Is(err, gormigrate.ErrMigrationIDDoesNotExist) && err != nil {
			return err
		}
		return nil
	}
	err := mc.rollbackTo(migrationId)
	if !errors.Is(err, gormigrate.ErrMigrationIDDoesNotExist) && err != nil {
		return err
	}
	return nil
}

// rollbackLast roolback the lasted migration
func (mc *Migrate) rollbackLast() error {
	return mc.gormigrate().RollbackLast()
}

// Migrate exec migration cmd
func (mc *Migrate) Migrate() error {
	m := mc.gormigrate()
	err := m.Migrate()
	if err != nil {
		return err
	}
	return nil
}

// Seed seed data into database
func (mc *Migrate) Seed() error {
	if mc.seeds == nil {
		return nil
	}
	return Seed(mc.seeds...)
}

// getFirstMigration get first migration's id
func (mc *Migrate) getFirstMigration() string {
	if mc.MigrationLen() == 0 {
		return ""
	}
	return mc.items[0].ID
}

// gormigrate create *gormigrate.Gormigrate
func (mc *Migrate) gormigrate() *gormigrate.Gormigrate {
	return gormigrate.New(mc.db, gormigrate.DefaultOptions, mc.items)
}
