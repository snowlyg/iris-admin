package admin

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func GetMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20211214120700_create_oplogs_table",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Oplog{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("oplogs")
		},
	}
}

// CreateOplog
func CreateOplog(ol *Oplog) error {
	// err := Instance().Model(&Oplog{}).Create(ol).Error
	// if err != nil {
	// 	return err
	// }
	return nil
}

// Oplog middleware model
type Oplog struct {
	gorm.Model
	Ip           string        `json:"ip" form:"ip" gorm:"column:ip;comment:ip"`
	Method       string        `json:"method" form:"method" gorm:"column:method;comment:method" validate:"required"`
	Path         string        `json:"path" form:"path" gorm:"column:path;comment:path" validate:"required"`
	Status       int           `json:"status" form:"status" gorm:"column:status;comment:status" validate:"required"`
	Latency      time.Duration `json:"latency" form:"latency" gorm:"column:latency;comment:latency"`
	Agent        string        `json:"agent" form:"agent" gorm:"column:agent;comment:agent"`
	ErrorMessage string        `json:"errorMessage" form:"errorMessage" gorm:"column:error_message;comment:error_message"`
	Body         string        `json:"body" form:"body" gorm:"type:longtext;column:body;comment:body"`
	Resp         string        `json:"resp" form:"resp" gorm:"type:longtext;column:resp;comment:resp"`
	UserID       uint          `json:"userId" form:"userId" gorm:"column:user_id;comment:user_id"`
	IsSuperAdmin bool          `json:"isSuperAdmin" form:"isSuperAdmin" gorm:"column:tenancy_id;comment:isSuperAdmin"`
}
