package authentification

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/Archie1978/regate/configuration"
	"github.com/Archie1978/regate/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Liste of driver remote protocol
var listDriverAuthentfication map[string]DriverAuthentfication

var authConf DriverAuthentfication
var authLocker sync.Mutex

// Interface Authentification
type DriverAuthentfication interface {
	New(*url.URL) DriverAuthentfication
	Authgin() gin.HandlerFunc
	SetRouteurGin(router *gin.Engine)
	GetProfile(c *gin.Context) (*database.UserProfile, error)
}

// Init list drivers authentification
func init() {
	listDriverAuthentfication = make(map[string]DriverAuthentfication)

}

// Add driver
func AddDriver(d DriverAuthentfication) {
	if reflect.ValueOf(d).Kind() != reflect.Ptr {
		return
	}
	var nameStruct = reflect.TypeOf(d).Elem().Name()
	listDriverAuthentfication[nameStruct] = d
}

// Get Driver authentification client WEB
func GetDriverURL(nameDriverURL string) (DriverAuthentfication, error) {
	u, _ := url.Parse(nameDriverURL)
	if !strings.Contains(nameDriverURL, "://") {
		u.Scheme = "flat"
	}
	caser := cases.Title(language.Und)
	key := "Authentification" + caser.String(u.Scheme)
	if d, ok := listDriverAuthentfication[key]; ok {
		return d.New(u), nil
	}
	return nil, fmt.Errorf("Driver authentification ["+key+"] not found into %v", ListDriver())
}

// Lis Driver authentification availaible
func ListDriver() (list []string) {
	list = make([]string, 0)
	for key, _ := range listDriverAuthentfication {
		list = append(list, key)
	}
	return
}

// Get driver authentification, create Object If not existe
func GetAuthentification() (auth DriverAuthentfication, err error) {
	authLocker.Lock()
	defer authLocker.Unlock()

	if authConf == nil {
		authConf, err = GetDriverURL(configuration.ConfigurationGlobal.Authentification)
	}
	return authConf, err
}
