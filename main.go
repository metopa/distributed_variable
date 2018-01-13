package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
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
	stdInChan := make(chan string)
	go stdInStream(stdInChan)

MAIN_LOOP:
	for {
		runDistributedApp(iface, stdInChan)

		for {
			input, ok := <-stdInChan
			if !ok || input == "exit" {
				break MAIN_LOOP
			} else if input == "restart" {
				break
			} else {
				fmt.Print("Valid commands:\n\texit\n\trestart\n")
			}
		}
	}
	fmt.Println("Terminating...")
	time.Sleep(time.Second * 2)
	fmt.Println("Terminated")
}

func runDistributedApp(iface *net.Interface, stdInChan chan string) {
	ctx := common.NewContext(common.PickRandomName(), 3, time.Second, time.Second*2)
	logger.SetContext(ctx)
	ctx.SetState(&state.DiscoveryState{Ctx: ctx})

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

	discoveryServer := dv_net.NewDiscoveryServer(ctx, string(ctx.ServerAddr),
		func(response string) {
			logger.Info("New discovery request from %v", response)
			dv_net.SendToDirectly(ctx, common.PeerAddr(response),
				common.NewPeerInfoRequestCommand(ctx.Name, ctx.Leader))
		})
	discoveryServer.StartOn(iface)
	discoveryServer.SendDiscoveryRequestOn(iface)

	console.ListenConsole(ctx, stdInChan)
	ctx.StopFlag = true
	return
}

func stdInStream(ch chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			close(ch)
			return
		}
		s = strings.TrimRight(s, "\n \t")
		ch <- s
	}
}
