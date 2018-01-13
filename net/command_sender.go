package net

import (
	"encoding/json"
	"net"
	"time"

	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/logger"
)

const DEFAULT_TTL = 20

func SendToDirectly(ctx *common.Context, destination common.PeerAddr, cmd common.Command) {
	doSend(ctx, destination, cmd)
}

func SendToRingLeader(ctx *common.Context, cmd common.Command) {
	cmd.Destination = ctx.Leader
	if ctx.LeaderDistance[1] == -1 ||
		(ctx.LeaderDistance[0] != -1 && ctx.LeaderDistance[0] < ctx.LeaderDistance[1]) {
		doSendReliable(ctx, ctx.LinkedPeers[0], ctx.LinkedPeers[1], cmd)
	} else {
		doSendReliable(ctx, ctx.LinkedPeers[1], ctx.LinkedPeers[0], cmd)
	}
}

func SendToHi(ctx *common.Context, cmd common.Command) {
	addr := ctx.LinkedPeers[1]
	if len(addr) == 0 {
		logger.Warn("Hi peer is unknown, send canceled")
	} else {
		doSend(ctx, addr, cmd)
	}
}

func SendToLo(ctx *common.Context, cmd common.Command) {
	addr := ctx.LinkedPeers[0]
	if len(addr) == 0 {
		logger.Warn("Lo peer is unknown, send canceled")
	} else {
		doSend(ctx, addr, cmd)
	}
}

func ForwardInRing(ctx *common.Context, from common.PeerAddr, cmd common.Command) {
	if from == ctx.LinkedPeers[0] {
		doSendReliable(ctx, ctx.LinkedPeers[1], ctx.LinkedPeers[0], cmd)
		SendToHi(ctx, cmd)
	} else if from == ctx.LinkedPeers[1] {
		doSendReliable(ctx, ctx.LinkedPeers[1], ctx.LinkedPeers[0], cmd)
	} else {
		logger.Warn("Tried to forward %v from %v, but linked peers are %v", cmd, from, ctx.LinkedPeers)
		doSendReliable(ctx, ctx.LinkedPeers[1], ctx.LinkedPeers[0], cmd)
	}
}

func ReplyInRing(ctx *common.Context, from common.PeerAddr, cmd common.Command) {
	if from == ctx.LinkedPeers[0] {
		doSendReliable(ctx, ctx.LinkedPeers[0], ctx.LinkedPeers[1], cmd)
	} else if from == ctx.LinkedPeers[1] {
		doSendReliable(ctx, ctx.LinkedPeers[1], ctx.LinkedPeers[0], cmd)
	} else {
		logger.Warn("Tried to reply %v to %v, but linked peers are %v", cmd, from, ctx.LinkedPeers)
		SendToHi(ctx, cmd)
	}
}

func BroadcastInRing(ctx *common.Context, cmd common.Command) {
	cmd.Destination = "BROADCAST"
	doSendReliable(ctx, ctx.LinkedPeers[0], ctx.LinkedPeers[1], cmd)
}

func initServiceFields(ctx *common.Context, destination common.PeerAddr, cmd *common.Command) {
	cmd.From = ctx.ServerAddr
	cmd.Clock = ctx.Clock
	if cmd.Source == "" {
		cmd.Source = ctx.ServerAddr
	}
	if cmd.Destination == "" {
		cmd.Destination = destination
	}
	if cmd.Ttl == 0 {
		cmd.Ttl = DEFAULT_TTL
	}
}

func sendImpl(ctx *common.Context, destination *net.TCPAddr, cmd *common.Command) error {
	var err error
	for i := 0; i < ctx.SendNumRetries; i++ {
		err = nil
		conn, err := net.DialTCP("tcp", nil, destination)
		if err != nil {
			time.Sleep(ctx.SendRetryPause)
			continue
		}
		e := json.NewEncoder(conn)
		err = e.Encode(&cmd)
		if err != nil {
			time.Sleep(ctx.SendRetryPause)
			continue
		}

		return nil
	}
	return err
}

func doSend(ctx *common.Context, destination common.PeerAddr, cmd common.Command) error {
	initServiceFields(ctx, destination, &cmd)
	destAddr, err := net.ResolveTCPAddr("tcp", string(destination))
	if err != nil {
		logger.Warn("Can't resolve %v: %v", destination, destAddr)
		return err
	}

	err = sendImpl(ctx, destAddr, &cmd)

	if err != nil {
		logger.Warn("Transmission: %v(%v): Error after %v retries: %v",
			cmd, common.GetTransmissionInfoString(
				cmd.Source, ctx.ServerAddr, destination, cmd.Destination),
			ctx.SendNumRetries, err)

	} else {
		logger.Info("Transmission: %v(%v): OK", cmd,
			common.GetTransmissionInfoString(
				cmd.Source, ctx.ServerAddr, destination, cmd.Destination))
	}

	return err
}

func doSendReliable(ctx *common.Context, mainAddr common.PeerAddr,
	backupAddr common.PeerAddr, cmd common.Command) {
	err := doSend(ctx, mainAddr, cmd)
	if err == nil || mainAddr == backupAddr{
		return
	}
	go doSend(ctx, backupAddr, cmd)
	logger.Warn("Reporting peer %v", mainAddr)
	report := common.NewReportPeerCommand(mainAddr)
	report.Destination = ctx.Leader
	initServiceFields(ctx, backupAddr, &report)
	doSend(ctx, backupAddr, report)
}