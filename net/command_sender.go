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
	initServiceFields(&cmd, ctx, destination)
	err := sendImpl(&cmd, ctx, destination)
	infoStr := common.GetTransmissionInfoString(
		cmd.Source, ctx.ServerAddr, destination, cmd.Destination)
	if err != nil {
		logger.Warn("Send %v(%v): Error: %v", cmd, infoStr, err)
	} else {
		logger.Info("Send %v(%v): OK", cmd, infoStr)
	}
}

func SendToReliable(ctx *common.Context, mainDest common.PeerAddr, altDest common.PeerAddr, cmd common.Command) {
	initServiceFields(&cmd, ctx, mainDest)
	err := sendImpl(&cmd, ctx, mainDest)
	infoStr := common.GetTransmissionInfoString(
		cmd.Source, ctx.ServerAddr, mainDest, cmd.Destination)
	if err == nil {
		logger.Info("Send %v(%v): OK", cmd, infoStr)
		return
	}

	logger.Warn("Send %v(%v): Trying alt link; Error: %v", cmd, infoStr, err)

	if altDest == mainDest || cmd.Destination == mainDest || len(altDest) == 0 {
		reason := "altDest == mainDest"
		if cmd.Destination == mainDest {
			reason = "cmd.Destination == mainDest"
		} else if len(altDest) == 0 {
			reason = "len(altDest) == 0"
		}
		logger.Warn("Send %v(%v): Alt link unavailable: %v", cmd, infoStr, reason)
		return
	}
	err = sendImpl(&cmd, ctx, altDest)
	if err != nil {
		//TODO report peer
		logger.Warn("Send %v(%v): Alt: Error: %v", cmd, infoStr, err)
	} else {
		logger.Info("Send %v(%v): Alt: OK", cmd, infoStr)
	}
}

func SendToHi(ctx *common.Context, cmd common.Command) {
	addr := ctx.LinkedPeers[1]
	if len(addr) == 0 {
		logger.Warn("Hi peer is unknown, send canceled")
	} else {
		SendToDirectly(ctx, addr, cmd)
	}
}

func SendToLo(ctx *common.Context, cmd common.Command) {
	addr := ctx.LinkedPeers[0]
	if len(addr) == 0 {
		logger.Warn("Lo peer is unknown, send canceled")
	} else {
		SendToDirectly(ctx, addr, cmd)
	}
}

func SendToLoReliable(ctx *common.Context, cmd common.Command) {
	SendToReliable(ctx, ctx.LinkedPeers[0], ctx.LinkedPeers[1], cmd)
}

func SendToHiReliable(ctx *common.Context, cmd common.Command) {
	SendToReliable(ctx, ctx.LinkedPeers[1], ctx.LinkedPeers[0], cmd)
}

func SendToRingLeader(ctx *common.Context, cmd common.Command) {
	cmd.Destination = ctx.Leader
	if ctx.LeaderDistance[1] == -1 ||
		(ctx.LeaderDistance[0] != -1 && ctx.LeaderDistance[0] < ctx.LeaderDistance[1]) {
		SendToLoReliable(ctx, cmd)
	} else {
		SendToHiReliable(ctx, cmd)
	}
}

func ForwardInRing(ctx *common.Context, from common.PeerAddr, cmd common.Command) {
	if from == ctx.LinkedPeers[0] {
		SendToHiReliable(ctx, cmd)
	} else if from == ctx.LinkedPeers[1] {
		SendToLoReliable(ctx, cmd)
	} else {
		logger.Warn("Tried to forward %v from %v, but linked peers are %v", cmd, from, ctx.LinkedPeers)
		SendToHiReliable(ctx, cmd)
	}
}

func ReplyInRing(ctx *common.Context, from common.PeerAddr, cmd common.Command) {
	if from == ctx.LinkedPeers[0] {
		SendToLoReliable(ctx, cmd)
	} else if from == ctx.LinkedPeers[1] {
		SendToHiReliable(ctx, cmd)
	} else {
		logger.Warn("Tried to reply %v to %v, but linked peers are %v", cmd, from, ctx.LinkedPeers)
		SendToHiReliable(ctx, cmd)
	}
}

func BroadcastInRing(ctx *common.Context, cmd common.Command) {
	cmd.Destination = "BROADCAST"
	SendToHi(ctx, cmd)
}

func initServiceFields(cmd *common.Command, ctx *common.Context, destination common.PeerAddr) {
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
func sendImpl(cmd *common.Command, ctx *common.Context, destination common.PeerAddr) error {
	destAddr, err := net.ResolveTCPAddr("tcp", string(destination))
	if err != nil {
		return err
	}

	for i := 0; i < ctx.SendNumRetries; i++ {
		conn, xerr := net.DialTCP("tcp", nil, destAddr)
		err = xerr
		if err != nil {
			logger.Warn("%v: Error: %v", destination, err)
			time.Sleep(ctx.SendRetryPause)
			continue
		}
		e := json.NewEncoder(conn)
		err = e.Encode(&cmd)
		if err != nil {
			logger.Warn("%v: Error: %v", destination, err)
			time.Sleep(ctx.SendRetryPause)
			continue
		}
		logger.Info("Sent %v to %v", cmd, destination)
		return nil
	}
	logger.Warn("Exit: Error: %v", err)

	return err
}
