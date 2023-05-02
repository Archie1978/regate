package main

import (
	"fmt"
	"webRemotedektop/drivers"
	"webRemotedektop/drivers/rdpDriver"
	"webRemotedektop/drivers/sshDriver"
)

func init() {
	drivers.AddDriver(&rdpDriver.ProcessRdp{})
	drivers.AddDriver(&sshDriver.ProcessSsh{})
	fmt.Println(drivers.ListDriver())
}
