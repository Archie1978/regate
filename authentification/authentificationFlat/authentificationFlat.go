package authentificationFlat

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"time"

	"bitbucket.org/avd/go-ipc/mmf"
	"bitbucket.org/avd/go-ipc/shm"
	"gorm.io/gorm"

	"github.com/Archie1978/regate/authentification"
	"github.com/Archie1978/regate/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-password/password"
)

/*
 * authentificationFlat is authentication by code in the get parameters.
 * It is used in the standalone program.
 *
 * The get request has code from shared memory. This allows
 *   * To know if the program is still started thanks to a timestamp in shared memory.
 *   * If the program is started, it generally retrieves the code and creates a request http://....?code=acode
 *   * Memory is shared only by the user of a system
 *   * The code is generated every TimeRefreshCode second.
 */

var TimeRefreshCode int = 7

var ErrAppNotStarted = errors.New("App not started")

type AuthentificationFlat struct {

	// Code is code into share memory
	Code string

	// Penultimate code where the code is refresh
	CodeOld string

	// CanalGenerate regerate Code manuel
	CanalGenerate chan bool
}

// Generate New authentificationFlat
func (authentificationFlat *AuthentificationFlat) New(configuration *url.URL) authentification.DriverAuthentfication {
	return &AuthentificationFlat{CanalGenerate: make(chan bool, 10)}
}

// Start Authentification
func (authentificationFlat *AuthentificationFlat) Start() {

	if _, err := authentificationFlat.GetCode(); err == ErrAppNotStarted {
		go authentificationFlat.engineGenerateCode()
	}

}

// Start SetRouteurGin
func (authentificationFlat *AuthentificationFlat) SetRouteurGin(router *gin.Engine) {

	var secret, err = password.Generate(64, 10, 10, true, true)
	if err != nil {
		log.Fatal(err)
	}
	router.Use(sessions.Sessions("mysession", cookie.NewStore([]byte(secret))))
	router.Use(authentificationFlat.Authgin())

}

// authentificationFlat: Allow authentification before get the page
func (authentificationFlat *AuthentificationFlat) Authgin() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get session and check parametre ( private is changed when the appalication reloaded )
		session := sessions.Default(c)
		authsession := session.Get("auth")
		authValid := false

		if authsession != nil {
			if authsession == "true" {
				// Authentification good: get page or webservices
				authValid = true
				c.Next()
				return
			}
		}

		// Check the last code into  share memory
		if c.Query("code") == authentificationFlat.Code && authentificationFlat.Code != "" {

			//Reset all code
			authentificationFlat.Code = ""
			authentificationFlat.CanalGenerate <- true
			authentificationFlat.CanalGenerate <- true
			authValid = true
		} else {

			// Check the seconde to last code into  share memory
			if c.Query("code") == authentificationFlat.CodeOld && authentificationFlat.CodeOld != "" {

				// Not Reset code
				authentificationFlat.CodeOld = ""
				authValid = true

			} else {

				// Anything code are good, return errror
				c.JSON(http.StatusOK, gin.H{
					"Error": "Authentication failed, restart",
				})

				c.String(http.StatusUnauthorized, "Authentication failed, restart")
				c.Abort()
				return
			}
		}

		// Save the auth in the session for others requests
		session.Set("auth", "true")
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}

		// Clean URL for the first connection and use session cookie
		if authValid && c.Query("code") != "" {
			http.Redirect(c.Writer, c.Request, "/", http.StatusSeeOther)
			c.Abort()
			return
		}

		c.Next()
	}
}

// getMemoryZone: Create memory share
func (authentificationFlat *AuthentificationFlat) getMemoryZone() (rwRegion *mmf.MemoryRegion, err error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	nameShare := "regate_" + u.Name + "_ontime"
	obj, err := shm.NewMemoryObject(nameShare, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
	if err != nil {
		if err != nil {
			obj, err = shm.NewMemoryObject(nameShare, os.O_RDWR, 0600)
		}
		if err != nil {
			return nil, err
		}
	}

	if err := obj.Truncate(2048); err != nil {
		panic("truncate")
	}

	// create  region for reading and writing.
	return mmf.NewMemoryRegion(obj, mmf.MEM_READWRITE, 0, 2018)
}

// engineGenerateCode: Generate Code into memorie share ( use goroutine )
func (authentificationFlat *AuthentificationFlat) engineGenerateCode() {
	rwRegion, err := authentificationFlat.getMemoryZone()

	if err != nil {
		log.Fatal(err)
	}
	writer := mmf.NewMemoryRegionWriter(rwRegion)

	authentificationFlat.CanalGenerate <- true
	for {

		ticker := time.NewTicker(time.Duration(TimeRefreshCode) * time.Second)

		select {
		case <-authentificationFlat.CanalGenerate:
		case <-ticker.C:
			// TimeOut refresh Code
		}

		buffer := bytes.NewBufferString("")
		authentificationFlat.CodeOld = authentificationFlat.Code
		authentificationFlat.Code, err = password.Generate(64, 10, 10, false, false)

		if err != nil {
			authentificationFlat.Code = ""
		}
		ts := time.Now().Unix()
		err = binary.Write(buffer, binary.LittleEndian, &ts)
		if err != nil {
			log.Fatal("Note found")
		}

		written, err := buffer.Write([]byte(authentificationFlat.Code))
		if written != len(authentificationFlat.Code) {
			log.Fatal(err)
		}

		writer.WriteAt(buffer.Bytes(), 0)

	}

}

// GetCode: GetCode get code into share memory or error with application not stated ( var ErrAppNotStarted )
func (authentificationFlat *AuthentificationFlat) GetCode() (code string, err error) {
	rwRegion, err := authentificationFlat.getMemoryZone()

	if err != nil {
		return "", err
	}

	reader := mmf.NewMemoryRegionReader(rwRegion)
	var ts int64
	err = binary.Read(reader, binary.LittleEndian, &ts)
	tm := time.Unix(ts, 0)
	if err != nil {
		return "", err
	}
	// Compare timeRefreshBefore
	fmt.Println(tm)
	if time.Now().Compare(tm.Add(time.Duration(TimeRefreshCode)*time.Second)) >= 0 && time.Now().Compare(time.Unix(4000000000, 0)) < 0 {
		return "", ErrAppNotStarted
	}

	var codeByte [64]byte
	err = binary.Read(reader, binary.LittleEndian, &codeByte)

	return string(string(codeByte[:])), nil
}

// GetProfile: get user into database
func (authentificationFlat *AuthentificationFlat) GetProfile(c *gin.Context) (*database.UserProfile, error) {
	profilUser, err := database.LoadUser("authentificationFlat")
	if err == gorm.ErrRecordNotFound || profilUser.Name == "" {
		profilUser, err = database.NewUser()
		if err != nil {
			log.Fatal(err)
		}
		profilUser.Name = "Unknown"
		profilUser.Referance = "authentificationFlat"
		database.SaveUser(profilUser)
	}
	return profilUser, nil
}
