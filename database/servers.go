package database

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/Archie1978/regate/configuration"
	"github.com/Archie1978/regate/crypto"
	"gorm.io/gorm"
)

type ServerGroup struct {
	gorm.Model

	ServerGroupParentID uint
	Name                string
	ServerGroupParent   *ServerGroup
}
type ServerGroups []*ServerGroup

func (serverGroups ServerGroups) ConvertMap() (list map[uint]*ServerGroup) {
	list = make(map[uint]*ServerGroup)
	for _, serverGroup := range serverGroups {
		list[serverGroup.ID] = serverGroup
	}
	return
}

type Server struct {
	gorm.Model
	Name string

	ServerGroupID uint
	ServerGroup   *ServerGroup

	URL string
}

type Servers []*Server

func GetServerById(id int) (Server, error) {
	var server Server
	ret := DB.First(&server, id)
	return server, ret.Error
}

func (servers Servers) ConvertMap() (list map[uint]*Server) {
	list = make(map[uint]*Server)
	for _, server := range servers {
		list[server.ID] = server
	}
	return
}

// Add Group Root
func CreateFirstGroup() error {
	var sg ServerGroup
	ret := DB.First(&sg)

	if ret.Error == gorm.ErrRecordNotFound {

		// create racine group
		sg.ID = 1
		sg.Name = "General"

		// Add Root$
		tx := DB.Create(&sg)
		return tx.Error
	}
	return ret.Error
}

// GetPathGroup
func GetPathGroup(pathGroup string) (*ServerGroup, error) {
	// Split Name
	listNameGroup := strings.Split(pathGroup, "/")

	var sg ServerGroup
	return GetPathSubGroup(listNameGroup[1:], &sg)
}

// GetPathSubGroup
func GetPathSubGroup(listNameGroup []string, sgParent *ServerGroup) (*ServerGroup, error) {
	var sg ServerGroup
	ret := DB.First(&sg, "name = ? AND server_group_parent_id = ?", listNameGroup[0], sgParent.ID)

	if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if ret.Error != nil {
		return nil, ret.Error
	}

	if len(listNameGroup) <= 1 {
		return &sg, nil
	}

	return GetPathSubGroup(listNameGroup[1:], &sg)
}

// SavePathGroup
func SavePathGroup(pathGroup string) (*ServerGroup, error) {
	// Split Name
	listNameGroup := strings.Split(pathGroup, "/")

	var sg ServerGroup
	return SavePathSubGroup(listNameGroup[1:], &sg)
}

func SavePathSubGroup(listNameGroup []string, sgParent *ServerGroup) (*ServerGroup, error) {
	var sg ServerGroup
	ret := DB.First(&sg, "name = ? AND server_group_parent_id = ?", listNameGroup[0], sgParent.ID)

	if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
		// Create erverGroup
		sg.Name = listNameGroup[0]
		sg.ServerGroupParentID = sgParent.ID
		ret = DB.Save(&sg)

		if ret.Error != nil {
			return nil, ret.Error
		}
	}

	if ret.Error != nil {
		return nil, ret.Error
	}
	if len(listNameGroup) <= 1 {
		return &sg, nil
	}
	return SavePathSubGroup(listNameGroup[1:], &sg)
}

// Save Server
func SaveServer(pathGroup string, server *Server) error {

	groupFinal, err := SavePathGroup(pathGroup)
	if err != nil {
		return err
	}

	server.ServerGroupID = groupFinal.ID
	ret := DB.Save(&server)

	return ret.Error
}

// Update Password
func (server *Server) UpdatePassword(passwordClear string) error {
	// Split l'URL
	u, err := url.Parse(server.URL)
	if err != nil {
		return err
	}

	// Get All information from USER with password into URL
	userInfo := u.User
	if userInfo != nil {

		// Replace the password by the new crypted password
		newUserInfo := url.UserPassword(userInfo.Username(), crypto.CryptPasswordString(passwordClear, configuration.ConfigurationGlobal.KeyCrypt))
		u.User = newUserInfo

		// Rebuild URL with new password
		server.URL = u.String()

		return nil
	}

	return fmt.Errorf("L'URL ne contient pas d'informations d'utilisateur (mot de passe)")
}

// Extract Password
func (server *Server) GetPassword() (passwordClear string, err error) {

	// Split l'URL
	u, err := url.Parse(server.URL)
	if err != nil {
		return "", err
	}

	//  Get All information USER for extract pass clear.
	userInfo := u.User
	if userInfo != nil {
		pass, _ := u.User.Password()
		return crypto.DecryptPasswordString(pass, configuration.ConfigurationGlobal.KeyCrypt), nil
	}

	return "", nil
}

// Get URL Server without password useful for upping server to interface (BECAREFULL URL in change into server)
func (server *Server) GetURLwithoutPassword() (string, error) {
	// Extract l'URL
	u, err := url.Parse(server.URL)
	if err != nil {
		return "", err
	}

	// Get All informations user and remove password
	userInfo := u.User
	if userInfo != nil {
		userInfoWithoutPassword := url.User(userInfo.Username())
		u.User = userInfoWithoutPassword
	}

	return u.String(), nil
}

/*
 *   Composite stuctures
 */
// Get All servers by composit
type ServerGroupComposit struct {
	ServerGroup
	ServerGroupChildren []ServerGroupComposit
	Servers             Servers
}

// Get All composite into database, use for menu
func GetServerGroupComposit() (ServerGroupComposit, error) {
	// get servers
	var servers Servers
	DB.Model(&Server{}).Find(&servers)

	// Remove password
	for i, _ := range servers {
		u, err := url.Parse(servers[i].URL)
		if err == nil {
			if u.User != nil {
				u.User = url.UserPassword(u.User.Username(), "Regate: N/A")
			}
		}
		servers[i].URL = u.String()
	}

	// get group
	var serverGroups ServerGroups
	DB.Model(&ServerGroup{}).Find(&serverGroups)

	// Map Server by parent
	mapServerByGroup := make(map[uint]Servers)
	for _, server := range servers {
		if _, ok := mapServerByGroup[server.ServerGroupID]; !ok {
			mapServerByGroup[server.ServerGroupID] = make(Servers, 0, 1)
		}
		mapServerByGroup[server.ServerGroupID] = append(mapServerByGroup[server.ServerGroupID],
			server)
	}

	// Map group by parent + Map id
	mapGroupParentId := make(map[uint][]ServerGroupComposit)
	for _, serverGroup := range serverGroups {
		if _, ok := mapGroupParentId[serverGroup.ServerGroupParentID]; !ok {
			mapGroupParentId[serverGroup.ServerGroupParentID] = make([]ServerGroupComposit, 0, 1)
		}

		mapGroupParentId[serverGroup.ServerGroupParentID] = append(mapGroupParentId[serverGroup.ServerGroupParentID],
			ServerGroupComposit{ServerGroup: *serverGroup})
	}

	var serverGroupComposit ServerGroupComposit
	getServerGroupCompositCreate(&serverGroupComposit,
		mapGroupParentId,
		mapServerByGroup)
	fmt.Println("serverGroupComposit:", serverGroupComposit)
	return serverGroupComposit, nil
}

func getServerGroupCompositCreate(serverGroupComposit *ServerGroupComposit,
	mapGroupParent map[uint][]ServerGroupComposit,
	mapServerByGroup map[uint]Servers) {

	serverGroupComposit.ServerGroupChildren = mapGroupParent[serverGroupComposit.ID]
	serverGroupComposit.Servers = mapServerByGroup[serverGroupComposit.ID]

	for i, _ := range serverGroupComposit.ServerGroupChildren {
		getServerGroupCompositCreate(&(serverGroupComposit.ServerGroupChildren[i]), mapGroupParent, mapServerByGroup)
	}
}
