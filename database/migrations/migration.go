package migrations

import (
	"sahamrakyat_test/database"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&database.Histories{}, &database.Orders{}, &database.Users{})
}