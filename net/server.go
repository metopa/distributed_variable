package net

import (
	"encoding/json"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/logger"
)

var localSessionIdCounter uint64

type TcpServer struct {
	listener *net.TCPListener
	stop     bool
	ctx      *common.Context
	sync     sync.Mutex
}

func NewTcpServer(context *common.Context) *TcpServer {
	return &TcpServer{ctx: context}
}

func (s *TcpServer) Listen() {
	if s.listener != nil {
		logger.Fatal("Main server is already running")
	}

	var err error
	s.listener, err = net.ListenTCP("tcp4", &net.TCPAddr{})
	if err != nil {
		logger.Fatal("%v", err)
	}
	logger.Info("Started main server on %v", s.Port())
	go s.accept()
}


func (s *TcpServer) Port() int {
	if s.listener == nil {
		logger.Fatal("Main server is not running")
	}
	addr, err := net.ResolveTCPAddr(s.listener.Addr().Network(), s.listener.Addr().String())
	if err != nil {
		logger.Fatal("%v", err)
	}
	return addr.Port
}

func (s *TcpServer) accept() {
	for !s.ctx.StopFlag {
		s.listener.SetDeadline(time.Now().Add(time.Second * 3))
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			if common.IsTimeoutError(err) {
				continue
			} else {
				logger.Fatal("%v", err)
			}
		}
		go s.handleConnection(conn)
	}
	s.listener.Close()
	logger.Info("Stopped main server")
}

func (s *TcpServer) handleConnection(conn *net.TCPConn) {
	sessionId := atomic.AddUint64(&localSessionIdCounter, 1)
	cmd, err := decodeCommand(conn)
	conn.Close()
	if err != nil {
		logger.Warn("Session #%v(%v): Error: %v",
			sessionId, conn.RemoteAddr().String(), err)
		return
	}

	s.ctx.Clock.SyncAfter(cmd.Clock, 1)

	senderAddr := cmd.From
	if senderAddr == "" {
		logger.Warn("Session #%v(%v): field FROM is empty",
			sessionId, conn.RemoteAddr().String())
	}
	sender := s.ctx.ResolvePeerName(senderAddr)

	//logger.Info("Session #%v(%v): Received %v[%v]", sessionId, sender, cmd, cmd.Clock.Value)

	if cmd.Destination != s.ctx.ServerAddr {
		if cmd.Destination == "BROADCAST" {
			if cmd.Source == s.ctx.ServerAddr {
				return
			} else {
				go SendToHi(s.ctx, cmd, false)
				common.DispatchCommand(s.ctx.GetState(), senderAddr, cmd)
				return
			}
		} else {
			cmd.Ttl--
			if cmd.Ttl > 0 {
				ForwardInRing(s.ctx, senderAddr, cmd)
			} else {
				logger.Warn("Session #%v(%v): TTL expired", sessionId, sender)
			}
			return
		}
	} else {
		common.DispatchCommand(s.ctx.GetState(), senderAddr, cmd)
		return
	}
}

func decodeCommand(conn *net.TCPConn) (common.Command, error) {
	conn.SetDeadline(time.Now().Add(time.Second * 15))
	d := json.NewDecoder(conn)
	var cmd common.Command
	err := d.Decode(&cmd)
	return cmd, err
}
