package model

import (
	"gorm.io/gorm"
)

func SetupAutoMigrate(db *gorm.DB) {
	if err := db.Migrator().AutoMigrate(
		&File{},
		&User{},
	); err != nil {
		panic(err)
	}
}
