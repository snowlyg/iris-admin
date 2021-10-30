package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/snowlyg/iris-admin/modules/seed"
	"github.com/snowlyg/iris-admin/server/database"
	"gorm.io/gorm"
)

type MigrationCmd struct {
	MigrationCollection []*gormigrate.Migration
	ModelCollection     []interface{}
	SeedCollection      []seed.SeedFunc
}

// New 构建 MigrationCmd
func New() *MigrationCmd {
	mc := &MigrationCmd{
		MigrationCollection: nil,
		ModelCollection:     nil,
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

// AddModel 添加 model
func (mc *MigrationCmd) AddModel(dst ...interface{}) {
	mc.ModelCollection = append(mc.ModelCollection, dst...)
}

// ModelLen ModelCollection 的长度
func (mc *MigrationCmd) ModelLen() int {
	return len(mc.ModelCollection)
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
	if err != nil {
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
	if migrationId == "" {
		mc.rollbackLast()
	}
	return mc.rollbackTo(migrationId)
}

// rollbackLast 回滚最后一次迁移
func (mc *MigrationCmd) rollbackLast() error {
	return mc.gormigrate().RollbackLast()
}

// Migrate 执行迁移
func (mc *MigrationCmd) Migrate() error {
	m := mc.gormigrate()
	m.InitSchema(func(tx *gorm.DB) error {
		err := tx.AutoMigrate(mc.ModelCollection...)
		if err != nil {
			return err
		}
		return nil
	})
	return m.Migrate()
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

// gormigrate get gormigrate
func (mc *MigrationCmd) gormigrate() *gormigrate.Gormigrate {
	return gormigrate.New(database.Instance(), gormigrate.DefaultOptions, mc.MigrationCollection)
}
