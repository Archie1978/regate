package authentificationBasic

import (
	"log"
	"net/url"

	"github.com/Archie1978/regate/authentification"
	"github.com/Archie1978/regate/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	//"github.com/gin-contrib/sessions"
	//"github.com/gin-contrib/sessions/cookie"
)

type AuthentificationBasic struct {
	accounts gin.Accounts
}

func (authentificationBasic *AuthentificationBasic) New(configuration *url.URL) authentification.DriverAuthentfication {
	return &AuthentificationBasic{accounts: gin.Accounts{"toto": "to"}}
}

func (authentificationFlat *AuthentificationBasic) SetRouteurGin(router *gin.Engine) {
}

func (authentificationBasic *AuthentificationBasic) Authgin() gin.HandlerFunc {

	// Utilisez la méthode `authorized` de votre struct `AuthentificationBasic`
	return gin.BasicAuth(authentificationBasic.accounts)

}

// GetProfile: get user into database
func (authentificationBasic *AuthentificationBasic) GetProfile(c *gin.Context) (*database.UserProfile, error) {
	profilUser, err := database.LoadUser("authentificationFlat")

	if err == gorm.ErrRecordNotFound {
		profilUser, err = database.NewUser()
		if err != nil {
			log.Fatal(err)
		}
		profilUser.Name = "Unkown"
	}
	return profilUser, nil
}
