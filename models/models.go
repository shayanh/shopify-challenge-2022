package models

import "gorm.io/gorm"

var instances = []interface{}{
	&Inventory{},
	&Item{},
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(instances...)
}
