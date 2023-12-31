//go:generate go-winres make --product-version=git-tag

// Generate client standalone
package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Archie1978/regate/authentification"
	"github.com/Archie1978/regate/authentification/authentificationFlat"
	"github.com/Archie1978/regate/authentification/authentificationNone"

	"github.com/Archie1978/regate/configuration"
	"github.com/Archie1978/regate/database"
	"github.com/Archie1978/regate/version"
	"github.com/Archie1978/regate/webservice"

	"github.com/pkg/browser"
	"github.com/takama/daemon"
	"github.com/tomatome/grdp/glog"
)

const (

	// name of the service
	name        = "regate-standalone"
	description = "Regate: Webervice remote desktop via eb browser"

	// port which daemon should be listen
	port = "localhost:8354"
)

// dependencies that are NOT required by the service, but might be used
var dependencies = []string{"dummy.service"}

var stdlog, errlog *log.Logger

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

// AddParametreURL: add code into URL
func AddParametreURL(urlString string, cle string, valeur string) (string, error) {
	// Analyser l'URL
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	// Obtenir les valeurs de la requête actuelles
	valeurs, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", err
	}

	// Ajouter ou remplacer le paramètre
	valeurs.Set(cle, valeur)

	// Mettre à jour la chaîne de requête
	u.RawQuery = valeurs.Encode()

	// Renvoyer la nouvelle URL
	return u.String(), nil
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {

	usage := "Usage: " + os.Args[0] + " [version]"

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "unsafe":
			authentification.AddDriver(&authentificationNone.AuthentificationNone{})
		case "version":
			return service.Version()
		default:

		}
	}

	fmt.Println("Load configuration: configuration.json")
	err := configuration.LoadConfiguration("configuration.json")
	if err != nil {
		log.Fatal(err)
	}

	// Get Authenfication
	authweb, err := configuration.ConfigurationGlobal.GetAuthentification()
	if err != nil {
		log.Fatal(err)
	}
	switch authweb.(type) {
	case *authentificationFlat.AuthentificationFlat:
		// Check start server web and start programme
		authFlat := authweb.(*authentificationFlat.AuthentificationFlat)
		code, err := authFlat.GetCode()
		if code != "" && err == nil {
			uCode, err := AddParametreURL(configuration.ConfigurationGlobal.GetConnectURL(), "code", code)
			if err != nil {
				log.Fatal(err)
			}
			browser.OpenURL(uCode)
			return usage, fmt.Errorf("Start browser to " + uCode)
		}
		if err != authentificationFlat.ErrAppNotStarted {
			return usage, err
		}

		// Start horodator
		authFlat.Start()

		// Wait the application is started and start browser
		go func() {
			for i := 0; i < 5; i++ {
				authFlat := authweb.(*authentificationFlat.AuthentificationFlat)
				code, err := authFlat.GetCode()

				if code != "" && err == nil {
					uCode, err := AddParametreURL(configuration.ConfigurationGlobal.GetConnectURL(), "code", code)
					if err != nil {
						log.Fatal(err)
					}
					browser.OpenURL(uCode)
					return
				}
				<-time.After(5 * time.Second)
			}
		}()
	case *authentificationNone.AuthentificationNone:
		browser.OpenURL(configuration.ConfigurationGlobal.GetConnectURL())
	default:
	}

	// Start Application

	// OpenDatabase
	err = database.OpenDatabase("database.sqlite")
	if err != nil {
		fmt.Println("Opendatabase", err)
	}

	// Check User Exist

	// Do something, call your goroutines, etc

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-interrupt
		service.Stop()
	}()

	webservice.StartWebservice()

	// never happen, but need to complete code
	return usage, nil
}

func (service *Service) Stop() (string, error) {
	os.Exit(1)
	return "Exit service", nil
}

func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)

	//glog.SetLevel(glog.LEVEL(0))
	glog.SetLevel(glog.DEBUG)
	logger := log.New(os.Stdout, "", 0)
	glog.SetLogger(logger)
}

func (service *Service) Version() (string, error) {
	versionBinary := version.Version()
	date := version.Date()
	return fmt.Sprintf("Version:%v\nDate:%v\n", versionBinary, date.Format(time.RFC3339)), nil
}

func main() {

	// Init service like daemon
	srv, err := daemon.New(name, description, daemon.SystemDaemon, dependencies...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}

	// Start Service
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)

}
