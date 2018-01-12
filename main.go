package main

import (
	"flag"
	"math/rand"
	"net"
	"time"

	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/console"
	"github.com/metopa/distributed_variable/logger"
	dv_net "github.com/metopa/distributed_variable/net"
	"github.com/metopa/distributed_variable/state"
)

func main() {
	ifaceNames := common.GetInterfaceNames()

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

	ctx := common.NewContext(common.PickRandomName(), 3, time.Second, time.Second*5)
	ctx.State = &state.DiscoveryState{Ctx: ctx}
	server := dv_net.NewTcpServer(ctx)

	server.Listen()
	ifAddr, err := common.GetInterfaceIPv4Addr(iface)
	if err != nil {
		logger.Fatal("%v", err)
	}
	ctx.ServerAddr = common.PeerAddr((&net.TCPAddr{IP: ifAddr, Port: server.Port()}).String())
	logger.Info("Peer name:    %s", ctx.Name)
	logger.Info("Peer id:      %d", ctx.PeerId)
	logger.Info("Peer address: %s", string(ctx.ServerAddr))

	discoveryServer := dv_net.NewDiscoveryServer(string(ctx.ServerAddr),
		func(response string) {
			logger.Info("New discovery response from %v", response)
			dv_net.SendToDirectly(ctx, common.PeerAddr(response),
				common.NewPeerInfoRequestCommand(ctx.Name))
		})
	discoveryServer.StartOn(iface)
	discoveryServer.SendDiscoveryRequestOn(iface)

	stop := make(chan struct{}, 2)
	console.ListenConsole(ctx, &stop)
	time.Sleep(time.Hour)
}
