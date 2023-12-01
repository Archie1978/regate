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

	"github.com/Archie1978/regate/authentification"
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

// Generate New authentficationFlat
func (authentficationFlat *AuthentificationFlat) New(configuration *url.URL) authentification.DriverAuthentfication {
	return &AuthentificationFlat{CanalGenerate: make(chan bool, 10)}
}

// Start Authentification
func (authentficationFlat *AuthentificationFlat) Start() {

	if _, err := authentficationFlat.GetCode(); err == ErrAppNotStarted {
		go authentficationFlat.engineGenerateCode()
	}

}

// Start SetRouteurGin
func (authentficationFlat *AuthentificationFlat) SetRouteurGin(router *gin.Engine) {

	var secret, err = password.Generate(64, 10, 10, true, true)
	if err != nil {
		log.Fatal(err)
	}
	router.Use(sessions.Sessions("mysession", cookie.NewStore([]byte(secret))))
	router.Use(authentficationFlat.Authgin())

}

// authentficationFlat: Allow authentification before get the page
func (authentficationFlat *AuthentificationFlat) Authgin() gin.HandlerFunc {
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
		if c.Query("code") == authentficationFlat.Code && authentficationFlat.Code != "" {

			//Reset all code
			authentficationFlat.Code = ""
			authentficationFlat.CanalGenerate <- true
			authentficationFlat.CanalGenerate <- true
			authValid = true
		} else {

			// Check the seconde to last code into  share memory
			if c.Query("code") == authentficationFlat.CodeOld && authentficationFlat.CodeOld != "" {

				// Not Reset code
				authentficationFlat.CodeOld = ""
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
func (authentficationFlat *AuthentificationFlat) getMemoryZone() (rwRegion *mmf.MemoryRegion, err error) {
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
func (authentficationFlat *AuthentificationFlat) engineGenerateCode() {
	rwRegion, err := authentficationFlat.getMemoryZone()

	if err != nil {
		log.Fatal(err)
	}
	writer := mmf.NewMemoryRegionWriter(rwRegion)

	authentficationFlat.CanalGenerate <- true
	for {

		ticker := time.NewTicker(time.Duration(TimeRefreshCode) * time.Second)

		select {
		case <-authentficationFlat.CanalGenerate:
		case <-ticker.C:
			// TimeOut refresh Code
		}

		buffer := bytes.NewBufferString("")
		authentficationFlat.CodeOld = authentficationFlat.Code
		authentficationFlat.Code, err = password.Generate(64, 10, 10, false, false)

		if err != nil {
			authentficationFlat.Code = ""
		}
		ts := time.Now().Unix()
		err = binary.Write(buffer, binary.LittleEndian, &ts)
		if err != nil {
			log.Fatal("Note found")
		}

		written, err := buffer.Write([]byte(authentficationFlat.Code))
		if written != len(authentficationFlat.Code) {
			log.Fatal(err)
		}

		writer.WriteAt(buffer.Bytes(), 0)

	}

}

// GetCode: GetCode get code into share memory or error with application not stated ( var ErrAppNotStarted )
func (authentficationFlat *AuthentificationFlat) GetCode() (code string, err error) {
	rwRegion, err := authentficationFlat.getMemoryZone()

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
