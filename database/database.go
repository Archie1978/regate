package database

import (
	"log"

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

	DB.AutoMigrate(&ServerGroup{}, &Server{}, &UserProfile{}, &SettingSecurity{})

	// Get Setting Security
	_, err = GetSettingSecurity()
	if err != nil {
		log.Fatal(err)
	}

	// Create first group if never group exist
	err = CreateFirstGroup()
	if err != nil {
		log.Fatal("Error root group server not found")
	}
	return nil
}
