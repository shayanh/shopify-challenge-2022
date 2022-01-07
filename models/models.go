package models

import "gorm.io/gorm"

// Package models contains definition of entities, their logic and data store
// interaction. We have implemented repository objects to communicate with the
// data store using an ORM.

var instances = []interface{}{
	&Inventory{},
	&Item{},
}

// Migrate automatically migrates model schemas.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(instances...)
}
