package net

import (
	"net"
	"time"

	"github.com/metopa/distributed_variable/common"
	"golang.org/x/net/ipv4"

	"github.com/metopa/distributed_variable/logger"
)

const MULTICAST_ADDR = "224.0.0.64"
const LOCAL_ADDR = "localhost"
const PORT=":7788"

var MULTICAST_GROUP = net.IPv4(224, 0, 0, 64)

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

	s.packetConnTransport, err = net.ListenPacket("udp4", "0.0.0.0" + PORT)

	s.packetConn = ipv4.NewPacketConn(s.packetConnTransport)
	logger.Info("Multicast if: %v", s.iface)
	err = s.packetConn.JoinGroup(s.iface, &net.UDPAddr{IP: MULTICAST_GROUP})
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
	//udpAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:8668")
	dst := &net.UDPAddr{IP: MULTICAST_GROUP, Port: 7788}

	conn, err := net.ListenPacket("udp",  "0.0.0.0:6886")
	if err != nil {
		logger.Fatal("%v", err)
	}
	p := ipv4.NewPacketConn(conn)
	p.SetMulticastInterface(s.iface)
	p.SetMulticastTTL(10)
	for i := 0; i <10; i++ {
		n, err := p.WriteTo([]byte(s.ownDiscoverResponse), nil, dst)
		logger.Info("Write to %v: %v, %v", dst, n, err)
		time.Sleep(time.Second)
	}
}

func (s *DiscoveryService) SendDiscoveryRequest2() {
	dst := &net.UDPAddr{IP: MULTICAST_GROUP, Port: 7788}
	for i := 0; i <10; i++ {
		n, err := s.packetConn.WriteTo([]byte(s.ownDiscoverResponse), nil, dst)
		logger.Info("Write to %v: %v, %v", dst, n, err)
	}
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
			n, cm, src, err := s.packetConn.ReadFrom(buf)
			if err != nil {
				if common.IsTimeoutError(err) {
					logger.Info("UDP Discovery timeout")
				} else {
					logger.Fatal("%v", err)
				}
			} else if cm.Dst.IsMulticast() && cm.Dst.Equal(MULTICAST_GROUP) {
				logger.Info("Datagram from %v[%v]: %v", src, cm, string(buf[:n]))
				go s.discoveryEventHandler(string(buf[:n]))
			}
		}
	}
}

func Listen() {
	const IF_NAME = "enp2s0"
	iface, err := net.InterfaceByName(IF_NAME)
	if err != nil {
		logger.Fatal("%v", err)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", MULTICAST_ADDR+":7655")
	if err != nil {
		logger.Fatal("ResolveUDPAddr: %v", err)
	}

	conn, err := net.ListenMulticastUDP("udp", iface, udpAddr)
	if err != nil {
		logger.Fatal("%v", err)
	}

	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	n, src, err := conn.ReadFrom(buf)
	logger.Info("Datagram from %v: %v", src, string(buf[:n]))
	if err != nil {
		logger.Fatal("ReadFrom: %v", err)
	}
}

func Send() {
	const IF_NAME = "enp2s0"
	iface, err := net.InterfaceByName(IF_NAME)
	if err != nil {
		logger.Fatal("%v", err)
	}
	ma, err := iface.Addrs()
	logger.Info("Local IP: %v -> %v", iface, ma)

	dstAddr, err := net.ResolveUDPAddr("udp", MULTICAST_ADDR+":7655")
	//udpAddr, err := net.ResolveUDPAddr("udp", ":7656")
	if err != nil {
		logger.Fatal("ResolveUDPAddr: %v", err)
	}
	conn, err := net.ListenUDP("udp", nil)

	if err != nil {
		logger.Fatal("%v", err)
	}

	_, err = conn.WriteTo([]byte("Hello"), dstAddr)
	logger.Info("Datagram sent")
	if err != nil {
		logger.Fatal("WriteTo: %v", err)
	}
}
