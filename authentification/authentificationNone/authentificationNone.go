package authentificationNone

import (
	"net/url"

	"github.com/Archie1978/regate/authentification"
	"github.com/gin-gonic/gin"
)

type AuthentificationNone struct {
}

func (authentficationFlat *AuthentificationNone) New(configuration *url.URL) authentification.DriverAuthentfication {
	return &AuthentificationNone{}
}

func (authentficationFlat *AuthentificationNone) SetRouteurGin(router *gin.Engine) {

}
func (authentficationFlat *AuthentificationNone) Authgin() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
