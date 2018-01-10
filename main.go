package main

import (
	"fmt"
	"math/rand"
	"time"

	dv_common "github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/logger"
	dv_net "github.com/metopa/distributed_variable/net"
)

func main() {
	//ifaceNames := dv_common.GetInterfaceNames()
	//
	//ifaceName := flag.String("if", "",
	//	"Network interface name. Available: "+ifaceNames)
	//
	//flag.Parse()
	//
	//if *ifaceName == "" {
	//	logger.Fatal("Usage: %s -if <%s>", flag.Arg(0), ifaceNames)
	//}
	//
	//iface, err := net.InterfaceByName(*ifaceName)
	//if err != nil {
	//	panic(err)
	//}
	rand.Seed(time.Now().UnixNano())

	ctx := dv_common.NewContext(dv_common.PickRandomName(), 3, time.Second)
	fmt.Printf("Peer name: %s\n", ctx.Name)

	server := dv_net.NewTcpServer(&dv_net.InitialCommandHandler{Ctx: ctx}, ctx)
	server.Listen()

	discoveryServer := dv_net.NewDiscoveryService(string(ctx.ServerAddr),
		func(response string) {
			logger.Info("New discovery response from %v", response)
			dv_net.SendToDirectly(ctx, dv_common.PeerAddr(response),
				dv_net.NewPeerInfoRequestCommand(ctx.Name))
		})
	discoveryServer.Start()
	discoveryServer.SendDiscoveryRequest()
	for {
	}
}
