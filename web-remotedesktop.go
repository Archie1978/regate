//	"getInfoInternetWebservice/git.private.idesi.fr/lisa/extractor"

// Example of a daemon with echo service
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/takama/daemon"
	"github.com/golang/glog"
)

const (

	// name of the service
	name        = "web-remotedektop"
	description = "Web remoteServeur"

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

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {

	usage := "Usage: myservice install | remove | start | stop | status"

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	fmt.Println("Load configuration: configuration.json")
	err := LoadConfiguration("configuration.json")
	if err != nil {
		log.Fatal(err)
	}

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

	// OpenDatabase
	OpenDatabase("database.sqlite")

	// /admin/secrets endpoint
	r := gin.Default()
	initService(r)

	// Listen and serve on 0.0.0.0:5543
	r.Run(ConfigurationGlobal.Listen)

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

func main() {
	srv, err := daemon.New(name, description, daemon.SystemDaemon, dependencies...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
