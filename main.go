package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	dv_common "github.com/metopa/distributed_variable/common"
	dv_net "github.com/metopa/distributed_variable/net"
)

func main() {
	iface, err := net.InterfaceByName("wlp3s0")
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())
	name := dv_common.PickRandomName()
	fmt.Printf("Peer name: %s\n", name)
	service := dv_net.NewDiscoveryService(name, iface,
		func(response string) { fmt.Printf("Hello from %s\n", response) })
	service.Start()
	defer service.Stop()
	time.Sleep(time.Second * 2)
	service.SendDiscoveryRequest2()
	time.Sleep(time.Second * 600)

	//go dv_net.Listen()
	//time.Sleep(time.Second)
	//dv_net.Send()
	//time.Sleep(time.Second)
}
