package authentificationBasic

import (
	"net/url"

	"github.com/Archie1978/regate/authentification"
	"github.com/gin-gonic/gin"
	//"github.com/gin-contrib/sessions"
	//"github.com/gin-contrib/sessions/cookie"
)

type AuthentificationBasic struct {
	accounts gin.Accounts
}

func (authentificationBasic *AuthentificationBasic) New(configuration *url.URL) authentification.DriverAuthentfication {
	return &AuthentificationBasic{accounts: gin.Accounts{"toto": "to"}}
}

func (authentficationFlat *AuthentificationBasic) SetRouteurGin(router *gin.Engine) {
}

func (authentificationBasic *AuthentificationBasic) Authgin() gin.HandlerFunc {

	// Utilisez la méthode `authorized` de votre struct `AuthentificationBasic`
	return gin.BasicAuth(authentificationBasic.accounts)

}
