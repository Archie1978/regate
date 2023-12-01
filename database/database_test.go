package database

import (
	"fmt"
	"os"
	"testing"
)

func TestDatabase(t *testing.T) {
	fmt.Println("==")
	os.Remove("/tmp/test.sqlite")
	OpenDatabase("/tmp/test.sqlite")

	// Add record
	group, err := SavePathGroup("Group1/Group1.1/Group1.1.1")
	if err != nil {
		t.Fatal(err)
	}
	if group == nil {
		t.Fatal(err)
	}
	fmt.Println("==")
	SavePathGroup("Group1/Group1.1/Group1.1.2")
	SavePathGroup("Group1/Group1.1/Group1.1.3")
	SavePathGroup("Group1/Group1.2/Group1.2.1")
	SavePathGroup("Group2/Group2.1/Group2.1.1")

	// load Serveur
	getGroup1, err := GetPathGroup("Group1/Group1.1/Group1.1.1")
	if err != nil {
		t.Fatal(err)
	}
	if getGroup1.ID != group.ID {
		t.Fatal("Id Not identique")
	}

	server := Server{Name: "polux", URL: "ssh://ronron@192.168.1.1"}
	err = SaveServer("Group1/Group1.2/Group1.2.1", &server)
	if err != nil {
		t.Fatal(err)
	}

	if server.ID != 1 {
		t.Fatal("Server can't have int !=1 ")
	}

	if server.ServerGroupID != 6 {
		t.Fatal("ServerGroupID can't have int !=6 ")
	}

	server2 := Server{Name: "polux2", URL: "ssh://ronron@192.168.1.1"}
	SaveServer("Group1/Group1.2/Group1.2.1", &server2)

	server3 := Server{Name: "polux3", URL: "ssh://ronron@192.168.1.1"}
	SaveServer("Group1/Group1.1/Group1.1.2", &server3)
	fmt.Println("===2=============")
	composite, err := GetServerGroupComposit()
	fmt.Println(composite, err)
}
