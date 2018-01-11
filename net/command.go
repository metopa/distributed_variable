package net

import (
	"fmt"

	"github.com/metopa/distributed_variable/common"
)

const (
	PEER_INFO_REQUEST_CMD = iota
	PEER_INFO_RESPONSE_CMD
	SET_LEADER_CMD
	SET_LINKED_PEERS_CMD
	REPORT_PEER_CMD
	LEADER_DISTANCE_REQUEST_CMD
	LEADER_DISTANCE_RESPONSE_CMD
	JOIN_RING_CMD
	LEAVE_RING_CMD
	GET_REQUEST_CMD
	GET_RESPONSE_CMD
	SET_REQUEST_CMD
	SET_RESPONSE_CMD
	MAX_CMD
)

var cmdNames = []string{
	"PEER_INFO_REQUEST_CMD",
	"PEER_INFO_RESPONSE_CMD",
	"SET_LEADER_CMD",
	"SET_LINKED_PEERS_CMD",
	"REPORT_PEER_CMD",
	"LEADER_DISTANCE_REQUEST_CMD",
	"LEADER_DISTANCE_RESPONSE_CMD",
	"JOIN_RING_CMD",
	"LEAVE_RING_CMD",
	"GET_REQUEST_CMD",
	"GET_RESPONSE_CMD",
	"SET_REQUEST_CMD",
	"SET_RESPONSE_CMD",
}


type TcpCommand struct {
	Op          int
	Sarg        []string
	Iarg        []int
	Source      common.PeerAddr
	Destination common.PeerAddr
	Ttl         int
}

func NewPeerInfoRequestCommand(name string) TcpCommand {
	return TcpCommand{Op: PEER_INFO_REQUEST_CMD, Sarg: []string{name}}
}

func NewPeerInfoResponseCommand(name string, leader common.PeerAddr) TcpCommand {
	return TcpCommand{Op: PEER_INFO_RESPONSE_CMD, Sarg: []string{name, string(leader)}}
}

func NewSetLeaderCommand(leader string) TcpCommand {
	return TcpCommand{Op: SET_LEADER_CMD, Sarg: []string{leader}}
}

func NewSetLinkedPeersCommand(loPeer common.PeerAddr, hiPeer common.PeerAddr) TcpCommand {
	return TcpCommand{Op: SET_LINKED_PEERS_CMD, Sarg: []string{string(loPeer), string(hiPeer)}}
}

func NewReportPeerCommand(peer common.PeerAddr) TcpCommand {
	return TcpCommand{Op: REPORT_PEER_CMD, Sarg: []string{string(peer)}}
}

func NewLeaderDistanceRequestCommand() TcpCommand {
	return TcpCommand{Op: LEADER_DISTANCE_REQUEST_CMD}
}

func NewLeaderDistanceResponseCommand(distance int) TcpCommand {
	return TcpCommand{Op: LEADER_DISTANCE_RESPONSE_CMD, Iarg: []int{distance}}
}

func NewJoinRingCommand() TcpCommand {
	return TcpCommand{Op: JOIN_RING_CMD}
}

func NewLeaveRingCommand() TcpCommand {
	return TcpCommand{Op: LEAVE_RING_CMD}
}

func NewGetRequestCommand() TcpCommand {
	return TcpCommand{Op: GET_REQUEST_CMD}
}

func NewGetResponseCommand(value int) TcpCommand {
	return TcpCommand{Op: GET_RESPONSE_CMD, Iarg: []int{value}}
}

func NewSetRequestCommand(value int) TcpCommand {
	return TcpCommand{Op: SET_REQUEST_CMD, Iarg: []int{value}}
}

func NewSetResponseCommand() TcpCommand {
	return TcpCommand{Op: SET_RESPONSE_CMD}
}

func (cmd *TcpCommand) String() string {
	if cmd.Op >= 0 && cmd.Op < MAX_CMD {
		return fmt.Sprintf("{%s to %s}", cmdNames[cmd.Op], cmd.Destination)
	} else {
		return fmt.Sprintf("{???(%d) to %s}", cmd.Op, cmd.Destination)
	}
}

func dispatchCommand(handler CommandHandler, sender common.PeerAddr, cmd TcpCommand) {
	switch cmd.Op {
	case PEER_INFO_REQUEST_CMD:
		handler.NewPeer(sender, cmd.Source, cmd.Sarg[0], true)
	case PEER_INFO_RESPONSE_CMD:
		handler.NewPeer(sender, cmd.Source, cmd.Sarg[0], false)
		if cmd.Sarg[2] != "" {
			handler.LeaderChanged(sender, common.PeerAddr(cmd.Sarg[1]))
		}
	case SET_LEADER_CMD:
		handler.LeaderChanged(sender, common.PeerAddr(cmd.Sarg[0]))
	case SET_LINKED_PEERS_CMD:
		handler.LinkedPeersChanged(sender,
			common.PeerAddr(cmd.Sarg[0]), common.PeerAddr(cmd.Sarg[1]))
	case REPORT_PEER_CMD:
		handler.PeerReported(sender, common.PeerAddr(cmd.Sarg[0]))
	case LEADER_DISTANCE_REQUEST_CMD:
		handler.DistanceRequested(sender, cmd.Source)
	case LEADER_DISTANCE_RESPONSE_CMD:
		handler.DistanceReceived(sender, cmd.Iarg[0])
	case JOIN_RING_CMD:
		handler.RingJoinRequested(sender, cmd.Source)
	case LEAVE_RING_CMD:
		handler.RingLeaveAnnounced(sender, cmd.Source)
	case GET_REQUEST_CMD:
		handler.ValueGetRequested(sender, cmd.Source)
	case GET_RESPONSE_CMD:
		handler.GotValue(sender, cmd.Iarg[0])
	case SET_REQUEST_CMD:
		handler.ValueSetRequested(sender, cmd.Source, cmd.Iarg[0])
	case SET_RESPONSE_CMD:
		handler.ValueSetConfirmed(sender)
	}
}
