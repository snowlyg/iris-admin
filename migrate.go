package admin

import (
	"errors"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// AddMigration add *gormigrate.Migration
func (ws *WebServe) AddMigration(m ...*gormigrate.Migration) {
	ws.items = append(ws.items, m...)
}

// MigrationLen length of MigrationCollection
func (ws *WebServe) MigrationLen() int {
	return len(ws.items)
}

// Refresh refresh migration
func (mc *WebServe) Refresh() error {
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
func (ws *WebServe) rollbackTo(migrationId string) error {
	return ws.m.RollbackTo(migrationId)
}

// Rollback roolback migrations
func (ws *WebServe) Rollback(migrationId string) error {
	if ws.MigrationLen() == 0 {
		return nil
	}
	if migrationId == "" {
		err := ws.rollbackLast()
		if !errors.Is(err, gormigrate.ErrMigrationIDDoesNotExist) && err != nil {
			return err
		}
		return nil
	}
	err := ws.rollbackTo(migrationId)
	if !errors.Is(err, gormigrate.ErrMigrationIDDoesNotExist) && err != nil {
		return err
	}
	return nil
}

// rollbackLast roolback the lasted migration
func (ws *WebServe) rollbackLast() error {
	return ws.m.RollbackLast()
}

// Migrate exec migration cmd
func (ws *WebServe) Migrate() error {
	// add migrations
	ws.AddMigration(
		&gormigrate.Migration{
			ID: "init_system",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&Router{}, &Menu{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(new(Router).TableName(), new(Menu).TableName())
			},
		},
		// add more migrations
	)
	if ws.m == nil {
		ws.m = gormigrate.New(ws.db, gormigrate.DefaultOptions, ws.items)
	}
	if err := ws.m.Migrate(); err != nil {
		return err
	}
	return nil
}

// getFirstMigration get first migration's id
func (ws *WebServe) getFirstMigration() string {
	if ws.MigrationLen() == 0 {
		return ""
	}
	return ws.items[0].ID
}
