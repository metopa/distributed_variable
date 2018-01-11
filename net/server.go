package net

import (
	"encoding/json"
	"net"
	"sync/atomic"
	"time"

	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/logger"
)

var localSessionIdCounter uint64

type TcpServer struct {
	listener     *net.TCPListener
	stop         bool
	eventHandler CommandHandler
	ctx          *common.Context
}

func NewTcpServer(handler CommandHandler, context *common.Context) *TcpServer {
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

func (s *TcpServer) Stop() {
	s.stop = true
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

func (s *TcpServer) SetEventHandler(handler CommandHandler) {
	s.eventHandler = handler
}

func (s *TcpServer) accept() {
	for !s.stop {
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
	logger.Info("Stopped main server")
}

func (s *TcpServer) handleConnection(conn *net.TCPConn) {
	//TODO Resolve alias
	sessionId := atomic.AddUint64(&localSessionIdCounter, 1)
	cmd, err := decodeCommand(conn)
	conn.Close()
	if err != nil {
		logger.Warn("Session #%v(%v): Error: %v", sessionId, conn.RemoteAddr(), err)
		return
	}
	logger.Info("Session #%v(%v): Received %v", sessionId, conn.RemoteAddr(), cmd)

	cmd.Ttl--
	if cmd.Ttl <= 0 {
		logger.Warn("Session #%v(%v): TTL expired", sessionId, conn.RemoteAddr())
		return
	}

	if cmd.Destination != s.ctx.ServerAddr {
		//t.Forward(cmd)
		//TODO
		return
	}

	dispatchCommand(s.eventHandler, common.PeerAddr(conn.RemoteAddr().String()), cmd)
}

func decodeCommand(conn *net.TCPConn) (TcpCommand, error) {
	conn.SetDeadline(time.Now().Add(time.Second * 15))
	d := json.NewDecoder(conn)
	var cmd TcpCommand
	err := d.Decode(&cmd)
	return cmd, err
}
