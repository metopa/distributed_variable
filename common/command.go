package common

import (
	"fmt"
)

const (
	PEER_INFO_REQUEST_CMD        = iota
	PEER_INFO_RESPONSE_CMD
	SET_LEADER_CMD
	REPORT_PEER_CMD
	REMOVE_PEER_CMD
	LEADER_DISTANCE_REQUEST_CMD
	LEADER_DISTANCE_RESPONSE_CMD
	SYNC_PEERS_CMD
	PING_CMD
	PONG_CMD
	CHANG_ROBERTS_ID_CMD
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
	"REPORT_PEER_CMD",
	"REMOVE_PEER_CMD",
	"LEADER_DISTANCE_REQUEST_CMD",
	"LEADER_DISTANCE_RESPONSE_CMD",
	"SYNC_PEERS_CMD",
	"PING_CMD",
	"PONG_CMD",
	"CHANG_ROBERTS_ID_CMD",
	"GET_REQUEST_CMD",
	"GET_RESPONSE_CMD",
	"SET_REQUEST_CMD",
	"SET_RESPONSE_CMD",
}

type Command struct {
	Op          int
	Sarg        []string
	Iarg        []int
	Source      PeerAddr
	Destination PeerAddr
	From        PeerAddr
	Ttl         int
	Clock       LamportClock
}

func NewPeerInfoRequestCommand(name string, leader PeerAddr) Command {
	return Command{Op: PEER_INFO_REQUEST_CMD, Sarg: []string{name, string(leader)}}
}

func NewPeerInfoResponseCommand(name string) Command {
	return Command{Op: PEER_INFO_RESPONSE_CMD, Sarg: []string{name}}
}

func NewSetLeaderCommand(leader PeerAddr, value int) Command {
	return Command{Op: SET_LEADER_CMD, Sarg: []string{string(leader)}, Iarg: []int{value}}
}

func NewReportPeerCommand(peer PeerAddr) Command {
	return Command{Op: REPORT_PEER_CMD, Sarg: []string{string(peer)}}
}

func NewRemovePeerCommand(peer PeerAddr, direction int) Command {
	return Command{Op: REMOVE_PEER_CMD, Sarg: []string{string(peer)}, Iarg: []int{direction}}
}

func NewLeaderDistanceRequestCommand() Command {
	return Command{Op: LEADER_DISTANCE_REQUEST_CMD}
}

func NewLeaderDistanceResponseCommand(distance int, direction int) Command {
	return Command{Op: LEADER_DISTANCE_RESPONSE_CMD, Iarg: []int{distance, direction}}
}

func NewSyncPeersCmd(peers map[PeerAddr]PeerInfo) Command {
	cmd := Command{Op: SYNC_PEERS_CMD}

	for _, v := range peers {
		cmd.Sarg = append(cmd.Sarg, string(v.Addr), v.Name)
	}

	return cmd
}

func NewPingCmd() Command {
	return Command{Op: PING_CMD}
}

func NewPongCmd() Command {
	return Command{Op: PONG_CMD}
}

func NewChangRobertIdCmd(id int) Command {
	return Command{Op: CHANG_ROBERTS_ID_CMD, Iarg: []int{id}}
}

func NewGetRequestCommand() Command {
	return Command{Op: GET_REQUEST_CMD}
}

func NewGetResponseCommand(value int) Command {
	return Command{Op: GET_RESPONSE_CMD, Iarg: []int{value}}
}

func NewSetRequestCommand(value int) Command {
	return Command{Op: SET_REQUEST_CMD, Iarg: []int{value}}
}

func NewSetResponseCommand() Command {
	return Command{Op: SET_RESPONSE_CMD}
}

func (cmd Command) String() string {
	if cmd.Op >= 0 && cmd.Op < MAX_CMD {
		return fmt.Sprintf("{%s}", cmdNames[cmd.Op])
	} else {
		return fmt.Sprintf("{???(%d)}", cmd.Op)
	}
}

func DispatchCommand(handler CommandHandler, sender PeerAddr, cmd Command) {
	switch cmd.Op {
	case PEER_INFO_REQUEST_CMD:
		handler.NewPeer(sender, cmd.Source, cmd.Sarg[0], true)
		if cmd.Sarg[1] != "" {
			handler.LeaderChanged(sender, PeerAddr(cmd.Sarg[1]), 0)
		}
	case PEER_INFO_RESPONSE_CMD:
		handler.NewPeer(sender, cmd.Source, cmd.Sarg[0], false)
	case SET_LEADER_CMD:
		handler.LeaderChanged(sender, PeerAddr(cmd.Sarg[0]), cmd.Iarg[0])
	case REPORT_PEER_CMD:
		handler.PeerReported(PeerAddr(cmd.Sarg[0]))
	case REMOVE_PEER_CMD:
		handler.PeerRemoved(sender, PeerAddr(cmd.Sarg[0]), cmd.Iarg[0])
	case LEADER_DISTANCE_REQUEST_CMD:
		handler.DistanceRequested(sender, cmd.Source)
	case LEADER_DISTANCE_RESPONSE_CMD:
		handler.DistanceReceived(sender, cmd.Iarg[0], cmd.Iarg[1])
	case SYNC_PEERS_CMD:
		handler.SyncPeers(sender, cmd.Sarg)
	case PING_CMD:
		handler.Ping(sender, cmd.Source)
	case PONG_CMD:
		handler.Pong(sender, cmd.Source)
	case CHANG_ROBERTS_ID_CMD:
		handler.ChRoIdReceived(sender, cmd.Iarg[0])
	case GET_REQUEST_CMD:
		handler.ValueGetRequested(sender, cmd.Source)
	case GET_RESPONSE_CMD:
		handler.GotValue(sender, cmd.Iarg[0])
	case SET_REQUEST_CMD:
		handler.ValueSetRequested(sender, cmd.Source, cmd.Iarg[0])
	case SET_RESPONSE_CMD:
		handler.ValueSetConfirmed(sender)
	default:
		fmt.Errorf("Unhandled command: %v\n", cmd)
	}
}
