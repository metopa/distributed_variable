package common

import (
	"fmt"
)

const (
	PEER_INFO_REQUEST_CMD = iota
	PEER_INFO_RESPONSE_CMD
	SET_LEADER_CMD
	SET_LINKED_PEERS_CMD
	REPORT_PEER_CMD
	LEADER_DISTANCE_REQUEST_CMD
	LEADER_DISTANCE_RESPONSE_CMD
	SYNC_PEERS_CMD
	CHANG_ROBERTS_ID_CMD
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
	"SYNC_PEERS_CMD",
	"CHANG_ROBERTS_ID_CMD",
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
	Source      PeerAddr
	Destination PeerAddr
	Ttl         int
}

func NewPeerInfoRequestCommand(name string) TcpCommand {
	return TcpCommand{Op: PEER_INFO_REQUEST_CMD, Sarg: []string{name}}
}

func NewPeerInfoResponseCommand(name string, leader PeerAddr) TcpCommand {
	return TcpCommand{Op: PEER_INFO_RESPONSE_CMD, Sarg: []string{name, string(leader)}}
}

func NewSetLeaderCommand(leader PeerAddr) TcpCommand {
	return TcpCommand{Op: SET_LEADER_CMD, Sarg: []string{string(leader)}}
}

func NewSetLinkedPeersCommand(loPeer PeerAddr, hiPeer PeerAddr) TcpCommand {
	return TcpCommand{Op: SET_LINKED_PEERS_CMD, Sarg: []string{string(loPeer), string(hiPeer)}}
}

func NewReportPeerCommand(peer PeerAddr) TcpCommand {
	return TcpCommand{Op: REPORT_PEER_CMD, Sarg: []string{string(peer)}}
}

func NewLeaderDistanceRequestCommand() TcpCommand {
	return TcpCommand{Op: LEADER_DISTANCE_REQUEST_CMD}
}

func NewLeaderDistanceResponseCommand(distance int) TcpCommand {
	return TcpCommand{Op: LEADER_DISTANCE_RESPONSE_CMD, Iarg: []int{distance}}
}

func NewSyncPeersCmd(ctx *Context) TcpCommand {
	cmd := TcpCommand{Op: SYNC_PEERS_CMD}
	ctx.Sync.Lock()
	for _, v := range ctx.KnownPeers {
		cmd.Sarg = append(cmd.Sarg, string(v.Addr), v.Name)
	}
	ctx.Sync.Unlock()
	return cmd
}

func NewChangRobertIdCmd(id int) TcpCommand {
	return TcpCommand{Op: CHANG_ROBERTS_ID_CMD, Iarg: []int{id}}
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


func (cmd TcpCommand) String() string {
	if cmd.Op >= 0 && cmd.Op < MAX_CMD {
		return fmt.Sprintf("{%s}", cmdNames[cmd.Op])
	} else {
		return fmt.Sprintf("{???(%d)}", cmd.Op)
	}
}

func DispatchCommand(handler CommandHandler, sender PeerAddr, cmd TcpCommand) {
	switch cmd.Op {
	case PEER_INFO_REQUEST_CMD:
		handler.NewPeer(sender, cmd.Source, cmd.Sarg[0], true)
	case PEER_INFO_RESPONSE_CMD:
		handler.NewPeer(sender, cmd.Source, cmd.Sarg[0], false)
		if cmd.Sarg[1] != "" {
			handler.LeaderChanged(sender, PeerAddr(cmd.Sarg[1]))
		}
	case SET_LEADER_CMD:
		handler.LeaderChanged(sender, PeerAddr(cmd.Sarg[0]))
	case SET_LINKED_PEERS_CMD:
		handler.LinkedPeersChanged(sender,
			PeerAddr(cmd.Sarg[0]), PeerAddr(cmd.Sarg[1]))
	case REPORT_PEER_CMD:
		handler.PeerReported(sender, PeerAddr(cmd.Sarg[0]))
	case LEADER_DISTANCE_REQUEST_CMD:
		handler.DistanceRequested(sender, cmd.Source)
	case LEADER_DISTANCE_RESPONSE_CMD:
		handler.DistanceReceived(sender, cmd.Iarg[0])
	case SYNC_PEERS_CMD:
		handler.SyncPeers(sender, cmd.Sarg)
	case CHANG_ROBERTS_ID_CMD:
		handler.ChRoIdReceived(sender, cmd.Iarg[0])
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
