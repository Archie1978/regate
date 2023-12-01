package main

import (
	"fmt"

	"github.com/Archie1978/regate/authentification"
	"github.com/Archie1978/regate/authentification/authentificationBasic"
	"github.com/Archie1978/regate/authentification/authentificationNone"

	"github.com/Archie1978/regate/drivers"
	"github.com/Archie1978/regate/drivers/rdpDriver"
	"github.com/Archie1978/regate/drivers/sshDriver"
)

func init() {
	// Init Driver connection
	drivers.AddDriver(&rdpDriver.ProcessRdp{})
	drivers.AddDriver(&sshDriver.ProcessSsh{})

	// Init Driver authentification
	authentification.AddDriver(&authentificationBasic.AuthentificationBasic{})
	authentification.AddDriver(&authentificationNone.AuthentificationNone{})
	fmt.Println("List drivers", drivers.ListDriver())
}
