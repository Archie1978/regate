package database

import (
	"errors"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func OpenDatabase(pathdatabase string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(pathdatabase), &gorm.Config{})
	if err != nil {
		return err
	}

	DB.AutoMigrate(&ServerGroup{}, &Server{})
	return nil
}

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

// Save Servers
func SaveServer(pathGroup string, server *Server) error {

	groupFinal, err := SavePathGroup(pathGroup)
	if err != nil {
		return err
	}

	server.ServerGroupID = groupFinal.ID
	ret := DB.Save(&server)

	return ret.Error

}

// Get All servers by composit
type ServerGroupComposit struct {
	ServerGroup
	ServerGroupChildren []ServerGroupComposit
	Servers             Servers
}

func GetServerGroupComposit() (ServerGroupComposit, error) {
	// get servers
	var servers Servers
	DB.Debug().Model(&Server{}).Find(&servers)

	// get group
	var serverGroups ServerGroups
	DB.Debug().Model(&ServerGroup{}).Find(&serverGroups)

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
