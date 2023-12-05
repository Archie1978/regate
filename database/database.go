package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func OpenDatabase(pathdatabase string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(pathdatabase), &gorm.Config{})
	if err != nil {
		return err
	}

	DB.AutoMigrate(&ServerGroup{}, &Server{}, &UserProfile{})
	return nil
}
