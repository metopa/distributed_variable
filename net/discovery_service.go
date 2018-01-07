package net

import (
	"fmt"
	"net"
	"time"

	"github.com/metopa/distributed_variable/common"
	"golang.org/x/net/ipv4"

	"github.com/metopa/distributed_variable/logger"
)

const MULTICAST_PORT = 7788

var MULTICAST_IP = net.IPv4(224, 0, 0, 64)
var MULTICAST_ADDR = &net.UDPAddr{IP: MULTICAST_IP, Port: MULTICAST_PORT}

type DiscoveryService struct {
	started               bool
	ownDiscoverResponse   string
	packetConnTransport   net.PacketConn
	packetConn            *ipv4.PacketConn
	iface                 *net.Interface
	stopListenChan        chan struct{}
	discoveryEventHandler func(discoverResponse string)
}

func NewDiscoveryService(discoverResponse string, iface *net.Interface,
	discoveryEventHandler func(discoverResponse string)) *DiscoveryService {
	return &DiscoveryService{
		ownDiscoverResponse:   discoverResponse,
		iface:                 iface,
		stopListenChan:        make(chan struct{}, 1),
		discoveryEventHandler: discoveryEventHandler}
}

func (s *DiscoveryService) Start() {
	var err error
	if s.started {
		logger.Warn("Tried to start UDP Discovery service twice")
		return
	}
	logger.Info("Multicast interface: %v", s.iface)

	s.packetConnTransport, err = net.ListenPacket("udp4", fmt.Sprintf(":%d", MULTICAST_PORT))
	if err != nil {
		logger.Fatal("%v", err)
	}

	s.packetConn = ipv4.NewPacketConn(s.packetConnTransport)

	err = s.packetConn.JoinGroup(s.iface, &net.UDPAddr{IP: MULTICAST_IP})
	if err != nil {
		logger.Fatal("%v", err.Error())
	}
	err = s.packetConn.SetControlMessage(ipv4.FlagDst, true)
	if err != nil {
		logger.Fatal("%v", err.Error())
	}
	err = s.packetConn.SetMulticastLoopback(true)
	if err != nil {
		logger.Fatal("%v", err.Error())
	}
	s.started = true

	go s.listen()
	logger.Info("Started UDP Discovery service")
}

func (s *DiscoveryService) Stop() {
	if !s.started {
		logger.Warn("Tried to stop UDP Discovery service twice")
		return
	}
	s.stopListenChan <- struct{}{}
}

func (s *DiscoveryService) SendDiscoveryRequest() {

	conn, err := net.ListenPacket("udp", ":0")
	defer conn.Close()
	if err != nil {
		logger.Fatal("%v", err)
	}
	p := ipv4.NewPacketConn(conn)
	p.SetMulticastInterface(s.iface)
	p.SetMulticastTTL(2)

	for {
		n, err := p.WriteTo([]byte(s.ownDiscoverResponse), nil, MULTICAST_ADDR)

		if err != nil {
			logger.Fatal("%v", err)
		}
		if n != len(s.ownDiscoverResponse) {
			logger.Warn("Discovery response was not transmitted as whole, repeating")
			continue
		}

		break
	}

	logger.Info("Transmitted discovery response to %v", MULTICAST_ADDR)
}

func (s *DiscoveryService) listen() {
	defer s.packetConn.Close()
	defer func() { s.started = false }()

	buf := make([]byte, 1024)
	for {
		select {
		case _ = <-s.stopListenChan:
			logger.Info("Shut down UDP Discovery service")
			return
		default:
			s.packetConnTransport.SetReadDeadline(time.Now().Add(time.Second * 5))
			n, cm, _, err := s.packetConn.ReadFrom(buf)
			if err != nil {
				if common.IsTimeoutError(err) {
					logger.Info("UDP Discovery timeout")
				} else {
					logger.Fatal("%v", err)
				}
			} else if cm.Dst.IsMulticast() && cm.Dst.Equal(MULTICAST_IP) {
				response := string(buf[:n])
				if response != s.ownDiscoverResponse {
					go s.discoveryEventHandler(response)
				} else {
					logger.Info("Got loopback response")
				}
			}
		}
	}
}
