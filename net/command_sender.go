package net

import (
	"encoding/json"
	"net"
	"time"

	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/logger"
)

const DEFAULT_TTL = 20

func SendToDirectly(ctx *common.Context, destination common.PeerAddr, cmd common.TcpCommand) {
	if cmd.Source == "" {
		cmd.Source = ctx.ServerAddr
	}
	if cmd.Destination == "" {
		cmd.Destination = destination
	}
	if cmd.Ttl == 0 {
		cmd.Ttl = DEFAULT_TTL
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
		logger.Info("New transmission: %v",
			common.GetTransmissionInfoString(
				cmd.Source, ctx.ServerAddr, destination, cmd.Destination))
		return
	}

	logger.Warn("Failed to send data to %v after %v retries", destination, ctx.SendNumRetries)
}

func SendToRingLeader(ctx *common.Context, cmd common.TcpCommand) {
	cmd.Destination = ctx.ServerAddr
	//TODO Choose direction basing on leader distance
	SendToHi(ctx, cmd)
}

func SendToHi(ctx *common.Context, cmd common.TcpCommand) {
	addr := ctx.LinkedPeers[1]
	if len(addr) == 0 {
		logger.Warn("Hi peer is unknown, send canceled")
	} else {
		SendToDirectly(ctx, addr, cmd)
	}
}

func SendToLo(ctx *common.Context, cmd common.TcpCommand) {
	//TODO Send in different direction
	addr := ctx.LinkedPeers[0]
	if len(addr) == 0 {
		logger.Warn("Lo peer is unknown, send canceled")
	} else {
		SendToDirectly(ctx, addr, cmd)
	}
}

func TransmitInRing(ctx *common.Context, from common.PeerAddr, cmd common.TcpCommand) {
	if from == ctx.LinkedPeers[0] {
		SendToHi(ctx, cmd)
	} else if from == ctx.LinkedPeers[1] {
		SendToLo(ctx, cmd)
	} else {
		logger.Warn("Tried to transmit %v from %v, but linked peers are %v", cmd, from, ctx.LinkedPeers)
		SendToHi(ctx, cmd)
	}
}

func BroadcastInRing(ctx *common.Context, cmd common.TcpCommand) {
	cmd.Destination = "BROADCAST"
	SendToHi(ctx, cmd)
}
