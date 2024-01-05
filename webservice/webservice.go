package webservice

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/Archie1978/regate/authentification"
	"github.com/Archie1978/regate/configuration"
	"github.com/Archie1978/regate/database"
	"github.com/Archie1978/regate/drivers"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"

	_ "embed"
)

var AuthentificationWeb authentification.DriverAuthentfication

// Start all service
func StartWebservice() {
	router := gin.Default()

	// End point webservice
	initService(router)

	// Listen and serve on 0.0.0.0:5543
	router.Run(configuration.ConfigurationGlobal.Listen)

}

// init interne service
func initService(router *gin.Engine) {

	// Get authentification
	AuthentificationWeb, err := authentification.GetAuthentification()
	if err != nil {
		log.Fatal(err)
	}

	// Init authentification
	AuthentificationWeb.SetRouteurGin(router)

	// Init map
	listWS = make(map[int]*websocket.Conn)
	listSession = make(map[int]drivers.DriverRP)

	// Service page static interface
	filememory := NewServeFileSystem(systemFileInternal, "/www/regate/dist")
	router.Use(static.Serve("/", filememory))

	// Load ws
	router.GET("/ws", func(c *gin.Context) {
		wshandler(c)
	})
	router.GET("/downloadFile", func(c *gin.Context) {
		downloadFile(c)
	})
	router.POST("/uploadFile", func(c *gin.Context) {
		uploadFile(c)
	})

	// Page dynamique for addon reote connexion
	router.GET("/addon-local.js", funcRouterJavascriptRM)

}

// upload javascript from drivers
func funcRouterJavascriptRM(c *gin.Context) {
	content := `
		var listPlugin=new Object();
	`
	content += drivers.GetAllCodeJavascript()
	c.String(http.StatusOK, content)
}

// Gestion WebSocket connexion
var listWS map[int]*websocket.Conn

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Add websocket into map
func AddSocket(ws *websocket.Conn) (slotNumber int) {
	for {
		i := rand.Int()
		if _, ok := listWS[i]; !ok {
			listWS[i] = ws
			return i
		}
	}
	return -1
}

// Delete socket
func DelSocket(slotNumber int) {
	delete(listWS, slotNumber)
}

/*
 * MessageRTM
 *   Session: session nmber
 *   Command: command of session, if number = 0  it's general
 *   TypeProtocol: Name Procotol of driver
 *   Type: ?
 *   Msg: Message for command session (Object or string or int )
 */
type MessageRTM struct {
	Session int

	Command      string
	TypeProtocol string //Use for command "start"

	Msg interface{}
}

var listSession map[int]drivers.DriverRP

// Get wsHandler for webservice
func wshandler(c *gin.Context) {
	w := c.Writer
	r := c.Request

	fmt.Println("Web service start")
	connexionWebsocket, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	slotNumber := AddSocket(connexionWebsocket)
	defer DelSocket(slotNumber)
	defer connexionWebsocket.Close()

	// Get authentification
	AuthentificationWeb, err := authentification.GetAuthentification()
	if err != nil {
		log.Fatal(err)
	}
	profile, err := AuthentificationWeb.GetProfile(c)
	if err != nil {
		log.Fatal("GetProfile:")
	}

	// Channel socket
	chanelWebSocket := make(chan interface{}, 100)
	go func() {

		for {
			interfaceObject := <-chanelWebSocket
			glog.Info("Send client:", interfaceObject)

			v := reflect.ValueOf(interfaceObject)
			switch v.Kind() {
			case reflect.Ptr:
				chanelWebSocket <- v.Elem().Interface()

			case reflect.Slice:
				if v.Type().Elem().Kind() == reflect.Uint8 {
					connexionWebsocket.WriteMessage(websocket.TextMessage, interfaceObject.([]byte))
				} else {
					glog.Info("Error websocket: type send not found")
				}

			default:
				jsonText, _ := json.Marshal(interfaceObject)
				connexionWebsocket.WriteMessage(websocket.TextMessage, jsonText)
			}
		}
	}()

	//Send profile
	chanelWebSocket <- MessageRTM{Command: "USER_PROFILE", Msg: profile}

	// Send directly profile

	// Management request from interface
	for {
		_, message, err := connexionWebsocket.ReadMessage()
		if err != nil {
			glog.Info("Error websocket:", err)
			return
		}

		var messageRTM MessageRTM
		err = json.Unmarshal(message, &messageRTM)
		glog.Info("message get: %v", messageRTM.Command)
		if err != nil {
			fmt.Println("Error decode message:", err)
		} else {
			switch messageRTM.Command {
			case "VERSION":
				// VERSION
				//chanelWebSocket <- MessageRTM{Command: "VERSION", Msg: string(VERSION)}

			case "LISTSERVER":
				// List of serveur into group
				serverGroup, err := database.GetServerGroupComposit()
				if err != nil {
					glog.Error("Error ServerGroup:", err)
				}
				chanelWebSocket <- MessageRTM{Command: "LISTSERVER", Msg: serverGroup}

			case "SaveConnection":
				// Save connection
				fmt.Println("SaveConnexion", messageRTM)
				if options, ok := messageRTM.Msg.(map[string]interface{}); ok {

					var server database.Server
					server.Name = fmt.Sprintf("%v", options["Name"])
					server.URL = fmt.Sprintf("%v", options["URL"])

					// Set Server Group for furtur dev
					if sgi, ok := options["ServerGroupId"]; ok {
						v, _ := strconv.Atoi(fmt.Sprintf("%v", sgi))
						server.ServerGroupID = uint(v)
					}
					// Update ID
					if _, ok := options["Id"]; ok {
						i, _ := strconv.Atoi(fmt.Sprintf("%v", options["Id"]))
						server.ID = uint(i)
					}

					// Add server GroupRoot
					if server.ServerGroupID == 0 {
						server.ServerGroupID = 1
					}

					// Password not change get old server
					passwordclear, _ := server.GetPassword()
					if passwordclear == "Regate%3A N%2FA" && server.ID != 0 {

						// Get the old Password if not change
						serverIntoDB, err := database.GetServerById(int(server.ID))
						if err == nil {
							password, _ := serverIntoDB.GetPassword()
							server.UpdatePassword(password)
						}

					} else {

						// Password change so crypt this
						err := server.UpdatePassword(passwordclear)
						fmt.Println(err)
					}

					// Save database
					ret := database.DB.Save(&server)
					if ret.Error != nil {
						chanelWebSocket <- MessageRTM{Command: "ERROR", Msg: ret.Error}
						return
					}

					serverGroup, err := database.GetServerGroupComposit()
					if err != nil {
						glog.Error("Error ServerGroup:", err)
					}
					chanelWebSocket <- MessageRTM{Command: "LISTSERVER", Msg: serverGroup}

				} else {
					chanelWebSocket <- MessageRTM{Command: "ERROR", Msg: "SaveConnexion error type"}
				}
			case "DeleteConnection":
				fmt.Println("DeleteConnexion", messageRTM)
				if messageRTM.Msg == nil {
					chanelWebSocket <- MessageRTM{Command: "ERROR", Msg: "Delete Connexion: Add ID"}
					return
				}

				idString := fmt.Sprintf("%v", messageRTM.Msg)
				id, _ := strconv.Atoi(idString)

				var server database.Server
				server.ID = uint(id)
				ret := database.DB.Delete(&server)
				if ret.Error != nil {
					chanelWebSocket <- MessageRTM{Command: "ERROR", Msg: "Delete Connexion: " + ret.Error.Error()}
					return
				}

				// Reload Server
				serverGroup, err := database.GetServerGroupComposit()
				if err != nil {
					glog.Error("Error ServerGroup:", err)
				}
				chanelWebSocket <- MessageRTM{Command: "LISTSERVER", Msg: serverGroup}

			// Start session init Session
			case "START":
				// Create session
				var numeroSession int
				for {
					numeroSession = rand.Int() % 99999999
					if _, ok := listSession[numeroSession]; !ok {
						break
					}
				}
				chanelWebSocket <- MessageRTM{Command: "STARTED", Session: numeroSession}
				if messageRTM.Msg != nil {

					switch reflect.TypeOf(messageRTM.Msg).Kind().String() {
					case "string":

						if strings.HasPrefix(messageRTM.Msg.(string), "Id:") {
							idServer, _ := strconv.Atoi(messageRTM.Msg.(string)[3:])
							server, errSql := database.GetServerById(idServer)
							if errSql == nil {
								messageRTM.Msg = server.URL
								/*
									u, err := url.Parse(server.URL)
									if err != nil {
										glog.Error("Config driver parse URL (", server.URL, "):", err)
									} else {
										msg := make(map[string]string)
										msg["Port"] = fmt.Sprintf("%v", u.Port())
										hsplit := strings.Split(u.Host, "|")
										msg["Ip"] = hsplit[0]
										if len(hsplit) > 1 {
											msg["Domain"] = hsplit[1]
										}

										msg["Username"] = u.User.Username()
										msg["Password"], _ = u.User.Password()
										messageRTM.Msg = msg
									}
								*/
							} else {
								glog.Error("Server not found")
							}
						} else {
							glog.Error("ID not found")
						}
					default:
						glog.Error("Type of Message not found: %v", reflect.TypeOf(messageRTM.Msg).Kind().String())
					}
				}

				if err != nil {
					glog.Error("Error:", err)
				} else {
					driver, err := drivers.GetDriver("Process" + messageRTM.TypeProtocol)
					if err != nil {
						glog.Error("Error get driver:", err)
					} else {

						listSession[numeroSession] = driver.New()

						// Start can too long
						go func() {
							listSession[numeroSession].Start(chanelWebSocket, numeroSession, fmt.Sprintf("%v", messageRTM.Msg))
						}()
					}
				}
			case "Close":
				if s, ok := listSession[messageRTM.Session]; ok {
					s.Close()
				}
				delete(listSession, messageRTM.Session)
			default:
				// Default CMD
				if s, ok := listSession[messageRTM.Session]; ok {
					s.Process(messageRTM.Msg)
				} else {
					glog.Error("Error get session:", messageRTM.Session, " ", messageRTM.Command)
				}
			}
		}
	}
}

// Download file into session
func downloadFile(c *gin.Context) {
	sessionNumberString := c.Query("sessionNumber")
	sessionNumber, _ := strconv.Atoi(sessionNumberString)

	// get Session
	if session, ok := listSession[sessionNumber]; !ok {
		fmt.Println("session: failed")
		c.String(500, "session not found or unauthorized.")
		c.Abort()
		return
	} else {
		err := session.DownloadFile(c)
		if err != nil {

		}
	}
}

// Upload file into session
func uploadFile(c *gin.Context) {
	sessionNumberString := c.Query("sessionNumber")
	sessionNumber, _ := strconv.Atoi(sessionNumberString)

	// get Session
	if session, ok := listSession[sessionNumber]; !ok {
		fmt.Println("session: failed")
		c.String(500, "session not found or unauthorized.")
		c.Abort()
		return
	} else {
		err := session.UploadFile(c)
		if err != nil {
			fmt.Println("UploadFile: error:", err)
		}
	}
}
