package common

import "fmt"

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

var cmdNames = []string {
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
	Op int
	Sarg []string
	Iarg []int
}

func NewPeerInfoRequestCommand() TcpCommand {
	return TcpCommand{Op: PEER_INFO_REQUEST_CMD}
}

func NewPeerInfoResponseCommand(name string, leader string) TcpCommand {
	return TcpCommand{Op: PEER_INFO_RESPONSE_CMD, Sarg: []string{name, leader}}
}

func NewSetLeaderCommand(leader string) TcpCommand {
	return TcpCommand{Op: SET_LEADER_CMD, Sarg: []string{leader}}
}

func NewSetLinkedPeersCommand(loPeer string, hiPeer string) TcpCommand {
	return TcpCommand{Op: SET_LINKED_PEERS_CMD, Sarg: []string{loPeer, hiPeer}}
}

func NewReportPeerCommand(peer string) TcpCommand {
	return TcpCommand{Op: REPORT_PEER_CMD, Sarg: []string{peer}}
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

func NewGetResponseCommand() TcpCommand {
	return TcpCommand{Op: GET_RESPONSE_CMD}
}

func NewSetRequestCommand(value int) TcpCommand {
	return TcpCommand{Op: SET_REQUEST_CMD, Iarg: []int{value}}
}

func NewSetResponseCommand() TcpCommand {
	return TcpCommand{Op: SET_RESPONSE_CMD}
}

func (cmd *TcpCommand) String() string {
	if cmd.Op >= 0 && cmd.Op < MAX_CMD {
		return cmdNames[cmd.Op]
	} else {
		return fmt.Sprintf("UNKNOWN_CMD(%d)", cmd.Op)
	}
}