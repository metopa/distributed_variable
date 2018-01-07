package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	dv_common "github.com/metopa/distributed_variable/common"
	dv_net "github.com/metopa/distributed_variable/net"
)

func main() {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	var ifaceNames []string
	for _, i := range ifaces {
		ifaceNames = append(ifaceNames, i.Name)
	}

	ifaceName := flag.String("if", "",
		"Network interface name. Available: "+strings.Join(ifaceNames, "|"))

	flag.Parse()

	if *ifaceName == "" {
		panic(fmt.Sprintf("Usage: %s -if <%s>", flag.Arg(0), strings.Join(ifaceNames, "|")))
	}

	iface, err := net.InterfaceByName(*ifaceName)
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())
	name := dv_common.PickRandomName()
	fmt.Printf("Peer name: %s\n", name)
	seen := make(map[string]bool)
	service := dv_net.NewDiscoveryService(name, iface,
		func(response string) {
			if seen[response] == false {
				fmt.Printf("Hello from %s\n", response)
				seen[response] = true
			}
		})
	service.Start()
	time.Sleep(time.Second * 2)
	service.SendDiscoveryRequest()
	time.Sleep(time.Second * 60)
	service.Stop()
	time.Sleep(time.Second * 5)
}
