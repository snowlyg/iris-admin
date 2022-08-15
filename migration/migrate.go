package migration

import (
	"errors"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/snowlyg/iris-admin/seed"
	"github.com/snowlyg/iris-admin/server/database"
)

// MigrationCmd 迁移 cmd
// MigrationCollection 迁移集合,数据表迁移方法
// SeedCollection 数据填充集合
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

// AddMigration 添加 *gormigrate.Migration
func (mc *MigrationCmd) AddMigration(m ...*gormigrate.Migration) {
	mc.MigrationCollection = append(mc.MigrationCollection, m...)
}

// MigrationLen MigrationCollection 的长度
func (mc *MigrationCmd) MigrationLen() int {
	return len(mc.MigrationCollection)
}

// AddSeed 添加 seed
func (mc *MigrationCmd) AddSeed(sf ...seed.SeedFunc) {
	mc.SeedCollection = append(mc.SeedCollection, sf...)
}

// SeedlLen SeedCollection 的长度
func (mc *MigrationCmd) SeedlLen() int {
	return len(mc.SeedCollection)
}

// Refresh 重置项目迁移
func (mc *MigrationCmd) Refresh() error {
	if mc.getFirstMigrateion() == "" {
		return nil
	}
	err := mc.rollbackTo(mc.getFirstMigrateion())
	if !errors.Is(gormigrate.ErrMigrationIDDoesNotExist, err) && err != nil {
		return err
	}
	return mc.Migrate()
}

// rollbackTo 回滚迁移到 migrationId 位置
func (mc *MigrationCmd) rollbackTo(migrationId string) error {
	return mc.gormigrate().RollbackTo(migrationId)
}

// Rollback 回滚迁移到
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

// rollbackLast 回滚最后一次迁移
func (mc *MigrationCmd) rollbackLast() error {
	return mc.gormigrate().RollbackLast()
}

// Migrate 执行迁移
func (mc *MigrationCmd) Migrate() error {
	m := mc.gormigrate()
	err := m.Migrate()
	if err != nil {
		return err
	}
	return nil
}

// Seed 填充数据
func (mc *MigrationCmd) Seed() error {
	if mc.SeedCollection == nil {
		return nil
	}
	return seed.Seed(mc.SeedCollection...)
}

// getFirstMigrateion 第一次迁移ID
func (mc *MigrationCmd) getFirstMigrateion() string {
	if mc.MigrationLen() == 0 {
		return ""
	}
	return mc.MigrationCollection[0].ID
}

// gormigrate 新建 *gormigrate.Gormigrate
func (mc *MigrationCmd) gormigrate() *gormigrate.Gormigrate {
	return gormigrate.New(database.Instance(), gormigrate.DefaultOptions, mc.MigrationCollection)
}
