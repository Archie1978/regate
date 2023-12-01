package main

import (
	"github.com/Archie1978/regate/authentification"
	"github.com/Archie1978/regate/authentification/authentificationFlat"

	"github.com/Archie1978/regate/drivers"
	"github.com/Archie1978/regate/drivers/rdpDriver"
	"github.com/Archie1978/regate/drivers/sshDriver"
)

func init() {
	// Init Driver connection
	drivers.AddDriver(&rdpDriver.ProcessRdp{})
	drivers.AddDriver(&sshDriver.ProcessSsh{})

	// Init Driver authentification
	authentification.AddDriver(&authentificationFlat.AuthentificationFlat{})
}
