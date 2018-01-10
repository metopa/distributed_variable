package net

import (
	"encoding/json"
	"net"
	"time"

	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/logger"
)

func SendToDirectly(ctx *common.Context, destination common.PeerAddr, cmd TcpCommand) {
	if cmd.Source == "" {
		cmd.Source = ctx.ServerAddr
	}
	if cmd.Destination == "" {
		cmd.Destination = destination
	}
	destAddr, err := net.ResolveTCPAddr("tcp", string(destination))
	if err != nil {
		logger.Warn("Can't resolve %v: %v", destination, destAddr)
		return
	}
	for i := 0; i < ctx.SendNumRetries; i++ {
		conn, err := net.DialTCP("tcp", nil, destAddr)
		if err != nil {
			logger.Warn("Can't dial %v: %v", destination, err)
			time.Sleep(ctx.SendRetryPause)
			continue
		}
		e := json.NewEncoder(conn)
		err = e.Encode(&cmd)
		if err != nil {
			logger.Warn("Failed to send data to %v: %v", destination, err)
			time.Sleep(ctx.SendRetryPause)
			continue
		}
		return
	}

	logger.Warn("Failed to send data to %v after %v retries", destination, ctx.SendNumRetries)
}

func SendToRing(ctx *common.Context, requestSender common.PeerAddr, destination common.PeerAddr, cmd TcpCommand) {
}
func SendToRingLeader(ctx *common.Context, cmd TcpCommand) {}
