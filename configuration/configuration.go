package configuration

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
	"time"

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

func (configuration *Configuration) GetConnectURL() (url string) {
	if strings.HasPrefix(configuration.Listen, ":") {
		return "http://127.0.0.1" + configuration.Listen
	}
	if strings.HasPrefix(configuration.Listen, "0.0.0.0:") {
		return "http://127.0.0.1" + configuration.Listen[10:]
	}
	return "http://" + configuration.Listen
}

// Load Configuration By Path
func LoadConfigurationSystem() error {

	listFolderPath := []string{"configuration.json", "./configuration.json"}

	// User configuration
	user, err := user.Current()
	if err == nil {
		listFolderPath = append(listFolderPath, path.Join(user.HomeDir, "/.config/regate.json"))
	}

	// Create folder configuration
	for _, folderPath := range listFolderPath {
		err = LoadConfiguration(folderPath)
		if err == nil {
			log.Println("Le configuration selected: ", folderPath)
			return nil
		}
	}

	return fmt.Errorf("Configuration not found in: %s", strings.Join(listFolderPath, ","))
}

// Load Configuration
func LoadConfiguration(pathConf string) error {

	// Check existe
	_, err := os.Stat(pathConf)
	if err != nil {
		return err
	}

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
		err = watcher.Add(pathConf)
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
					watcher.Remove(pathConf)
				}
				if event.Op&fsnotify.Write == fsnotify.Write {

					// Ne lis pas dessuite le fichier car il peut y avoir plusieurs notification d'ecriture à la suite en chaine
					if !readConfigurationTimer {
						readConfigurationTimer = true
						go func() {
							<-time.After(time.Duration(5) * time.Second)
							log.Println("Configuration modified: reload")
							err := loadConfiguration(pathConf)
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

	err = loadConfiguration(pathConf)
	return err
}

func loadConfiguration(pathConfiguration string) error {

	// Open Json
	jsonFile, err := os.Open(pathConfiguration)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	// Get path configuration
	pathParent := path.Dir(pathConfiguration)

	// Décode flux
	dec := json.NewDecoder(jsonFile)
	err = dec.Decode(&ConfigurationGlobal)
	if err != nil {
		return fmt.Errorf("error Load:%v", err)
	}

	// Authentification flat
	if ConfigurationGlobal.Authentification == "" {
		ConfigurationGlobal.Authentification = "none:///"
	}

	// Init default
	if ConfigurationGlobal.DatabasePath == "" {
		ConfigurationGlobal.DatabasePath = path.Join(pathParent, "database.sqlite")
	}

	if ConfigurationGlobal.Listen == "" {
		ConfigurationGlobal.Listen = ":5537"
	}

	if ConfigurationGlobal.KeyCrypt == nil {
		return fmt.Errorf("KeyCrypt not configure into configuration.json")
	}

	return nil
}

// Load PasswordCrypte
func SaveConfiguration(path string) error {

	// Créez un nouveau fichier avec os.Create
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Erreur lors de la création du fichier : %v", err)
	}
	defer file.Close()

	// Encodage des données en JSON
	jsonData, err := json.MarshalIndent(ConfigurationGlobal, "", "  ")
	if err != nil {
		return err
	}

	// Écriture du JSON dans le fichier
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	// Fermeture du fichier
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}
