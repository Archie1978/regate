package main

import (
	"fmt"

	"github.com/Archie1978/regate/drivers"
	"github.com/Archie1978/regate/drivers/rdpDriver"
	"github.com/Archie1978/regate/drivers/sshDriver"
)

func init() {
	drivers.AddDriver(&rdpDriver.ProcessRdp{})
	drivers.AddDriver(&sshDriver.ProcessSsh{})
	fmt.Println(drivers.ListDriver())
}
