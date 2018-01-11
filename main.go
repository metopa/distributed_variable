package main

import (
	"flag"
	"math/rand"
	"net"
	"time"

	dv_common "github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/logger"
	dv_net "github.com/metopa/distributed_variable/net"
)

func main() {
	ifaceNames := dv_common.GetInterfaceNames()

	ifaceName := flag.String("if", "",
		"Network interface name. Available: "+ifaceNames)

	flag.Parse()

	if *ifaceName == "" {
		logger.Fatal("Usage: %s -if <%s>", flag.Arg(0), ifaceNames)
	}

	iface, err := net.InterfaceByName(*ifaceName)
	if err != nil {
		logger.Fatal("%v", err)
	}
	rand.Seed(time.Now().UnixNano())

	ctx := dv_common.NewContext(dv_common.PickRandomName(), 3, time.Second)

	server := dv_net.NewTcpServer(&dv_net.InitialCommandHandler{Ctx: ctx}, ctx)
	server.Listen()
	ifAddr, err := dv_common.GetInterfaceIPv4Addr(iface)
	if err != nil {
		logger.Fatal("%v", err)
	}
	ctx.ServerAddr = dv_common.PeerAddr((&net.TCPAddr{IP: ifAddr, Port: server.Port()}).String())
	logger.Info("Peer name:    %s", ctx.Name)
	logger.Info("Peer address: %s", string(ctx.ServerAddr))

	discoveryServer := dv_net.NewDiscoveryServer(string(ctx.ServerAddr),
		func(response string) {
			logger.Info("New discovery response from %v", response)
			dv_net.SendToDirectly(ctx, dv_common.PeerAddr(response),
				dv_net.NewPeerInfoRequestCommand(ctx.Name))
		})
	discoveryServer.StartOn(iface)
	discoveryServer.SendDiscoveryRequestOn(iface)
	time.Sleep(time.Hour)
}
