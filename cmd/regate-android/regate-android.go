package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/Archie1978/regate/authentification"
	"github.com/Archie1978/regate/authentification/authentificationFlat"

	"github.com/Archie1978/regate/configuration"
	"github.com/Archie1978/regate/database"
	"github.com/Archie1978/regate/webservice"
)

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

// splitEnv divise une chaîne "clé=valeur" en un tableau [clé, valeur].
func splitEnv(env string) []string {
	for i := 0; i < len(env); i++ {
		if env[i] == '=' {
			return []string{env[:i], env[i+1:]}
		}
	}
	return []string{env, ""}
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Start regate")
	widget1 := widget.NewLabel("Click the button to open a browser")
	widgeterror := widget.NewLabel("")

	/// visionage
	r := ""

	// Récupérer toutes les variables d'environnement
	envVars := os.Environ()
	// Afficher toutes les variables d'environnement
	for _, env := range envVars {
		pair := splitEnv(env)
		r += fmt.Sprintf("%s=%s\n", pair[0], pair[1])
	}

	ex, err := os.Executable()
	if err != nil {
		fmt.Println("Erreur lors de la récupération du répertoire courant :", err)
		r += err.Error() + "\n"
	}
	r += "Exec:" + ex + "\n"
	/*	dir, err := os.Getwd()
		if err != nil {
			fmt.Println("Erreur lors de la récupération du répertoire courant :", err)
			r += "EEXEC:" + err.Error() + "\n"
		}
	*/
	files, err := os.ReadDir(os.Getenv("FILESDIR"))
	if err != nil {
		fmt.Println("Erreur lors de la lecture du répertoire :", err)
		r += "Elist:" + err.Error() + "\n"
	}

	fmt.Println("Fichiers et répertoires dans le répertoire courant :")
	for _, file := range files {
		r += "Entree:" + file.Name() + "\n"
	}

	// Load configuration and ignore error
	//fmt.Println("Load configuration: configuration.json")
	confPath := filepath.Join(os.Getenv("FILESDIR"), "configuration.json")
	_, err = os.Stat(confPath)

	// If the configuration file don t exist, the program generate one default conf
	if os.IsNotExist(err) {
		configuration.ConfigurationGlobal.Listen = "127.0.0.1:5543"
		// Generate random key
		keyCrypt := make([]byte, 16)
		source := rand.NewSource(time.Now().Unix())
		for i := 0; i < len(keyCrypt); i++ {
			keyCrypt[i] = byte(source.Int63())
		}
		configuration.ConfigurationGlobal.KeyCrypt = keyCrypt
		configuration.ConfigurationGlobal.Authentification = "flat:///"
		configuration.ConfigurationGlobal.DatabasePath = filepath.Join(os.Getenv("FILESDIR"), "database.sqlite")
		err = configuration.SaveConfiguration(confPath)
	}
	if err != nil {
		widgeterror.SetText(err.Error())
	}

	// reload configuration
	err = configuration.LoadConfiguration(confPath)
	if err != nil {
		widgeterror.SetText(err.Error())
	}

	// Get Authenfication
	authweb, err := authentification.GetAuthentification()
	if err != nil {
		widgeterror.SetText(err.Error())
	}
	fmt.Println(authweb)

	// OpenDatabase
	err = database.OpenDatabase(configuration.ConfigurationGlobal.DatabasePath)
	if err != nil {
		widgeterror.SetText(fmt.Sprintf("Error opening database: %v", err))
	}

	// start webservice
	go func() {
		webservice.StartWebservice()
	}()

	// Start authenfication
	authFlat := authweb.(*authentificationFlat.AuthentificationFlat)
	authFlat.Start()

	// Créer un bouton pour ouvrir une URL
	openBrowserButton := widget.NewButton("Open Browser", func() {

		code, err := authFlat.GetCode()
		if err != nil {
			widget1.SetText(fmt.Sprintf("code authentification:\n%v", err.Error()[15:]))
		}
		widget1.SetText(fmt.Sprintf("code authentification:%v\n", code))

		if code != "" && err == nil {
			uCode, err := AddParametreURL(configuration.ConfigurationGlobal.GetConnectURL(), "code", code)
			if err != nil {
				log.Fatal(err)
			}

			// URL à ouvrir
			urlToOpen, _ := url.Parse(uCode)
			err = fyne.CurrentApp().OpenURL(urlToOpen)
			if err != nil {
				// Gérer les erreurs si le navigateur ne s'ouvre pas
				println("Failed to open browser:", err.Error())
			}
		}
	})

	// Ajouter le bouton à la fenêtre
	myWindow.SetContent(container.NewVBox(
		widget1,
		widgeterror,
		openBrowserButton,
	))

	myWindow.ShowAndRun()
}
