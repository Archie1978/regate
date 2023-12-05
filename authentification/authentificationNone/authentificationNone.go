package authentificationNone

import (
	"net/url"

	"github.com/Archie1978/regate/authentification"
	"github.com/Archie1978/regate/database"
	"github.com/gin-gonic/gin"
)

type AuthentificationNone struct {
}

func (authentificationNone *AuthentificationNone) New(configuration *url.URL) authentification.DriverAuthentfication {
	return &AuthentificationNone{}
}

func (authentificationNone *AuthentificationNone) SetRouteurGin(router *gin.Engine) {

}
func (authentificationNone *AuthentificationNone) Authgin() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
func (authentificationNone *AuthentificationNone) GetProfile(c *gin.Context) (*database.UserProfile, error) {
	profilUser, err := database.NewUser()
	if err != nil {
		return nil, err
	}
	profilUser.Name = "Unknow User"

	return profilUser, nil
}
