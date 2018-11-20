package models

import (
	"time"
)

type Migrations struct {
	IdMigration        int       `xorm:"not null pk autoincr comment('surrogate key') INT(10)"`
	Name               string    `xorm:"comment('migration name, unique') VARCHAR(255)"`
	CreatedAt          time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('date migrated or rolled back') TIMESTAMP"`
	Statements         string    `xorm:"comment('SQL statements for this migration') LONGTEXT"`
	RollbackStatements string    `xorm:"comment('SQL statment for rolling back migration') LONGTEXT"`
	Status             string    `xorm:"comment('update indicates it is a normal migration while rollback means this migration is rolled back') ENUM('rollback','update')"`
}
