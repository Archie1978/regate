package configuration

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Archie1978/regate/authentification"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
)

type Configuration struct {
	Accounts gin.Accounts
	Listen   string

	Authentification string // method of authentification  type:   method://options
	KeyCrypt         []byte
	DatabasePath     string
}

var ConfigurationGlobal Configuration

// Get
func (configuration Configuration) GetAuthentification() (authentification.DriverAuthentfication, error) {
	return authentification.GetDriverURL(configuration.Authentification)
}

// Load PasswordCrypte
func LoadConfiguration(path string) error {

	// Add watcher
	go func() {
		// Refresh conf
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Println("Notification load configuration error:", err)
			log.Println("Notification add configuration disabled")
			return
		}
		defer watcher.Close()

		// Add watcher directory
		err = watcher.Add(path)
		if err != nil {
			log.Println("Notification add configuration file:", err)
			log.Println("Notification add configuration disabled")
			return
		}

		readConfigurationTimer := false
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					// cas improble du coups backup
					watcher.Remove(path)
				}
				if event.Op&fsnotify.Write == fsnotify.Write {

					// Ne lis pas dessuite le fichier car il peut y avoir plusieurs notification d'ecriture à la suite en chaine
					if !readConfigurationTimer {
						readConfigurationTimer = true
						go func() {
							<-time.After(time.Duration(5) * time.Second)
							log.Println("Configuration modified: reload")
							err := loadConfiguration(path)
							if err != nil {
								fmt.Println("Load configuration error", err)
							}
							readConfigurationTimer = false
						}()
					}

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Notification of configuration : error:", err)
			}

		}
	}()

	err := loadConfiguration(path)
	return err
}

func loadConfiguration(path string) error {
	// Open Json
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	// Décode flux
	dec := json.NewDecoder(jsonFile)
	err = dec.Decode(&ConfigurationGlobal)
	if err != nil {
		return fmt.Errorf("Error Load:%v", err)
	}

	// Authentification flat
	if ConfigurationGlobal.Authentification == "" {
		ConfigurationGlobal.Authentification = "none:///"
	}

	// Init default
	if ConfigurationGlobal.DatabasePath == "" {
		ConfigurationGlobal.DatabasePath = "database.sqlite"
	}

	if ConfigurationGlobal.Listen == "" {
		ConfigurationGlobal.Listen = ":5537"
	}
	return nil
}
