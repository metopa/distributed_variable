package state

import (
	"fmt"

	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/net"
)

type LinkedState struct {
	DiscoveryState
}
func (s *LinkedState) Start() {

}
func (s *LinkedState) GotValue(sender common.PeerAddr, value int) {
	fmt.Printf("Value = %v\n", value)
}
func (s *LinkedState) ValueSetConfirmed(sender common.PeerAddr) {
	fmt.Printf("Value is updated\n")
}

func (s *LinkedState) Name() string {
	return "Linked state state"
}

func (s *LinkedState) ActionStartChRo() {
	fmt.Println("Peer has already joined the ring")
}
func (s *LinkedState)  ActionSetValue(value int) {
	fmt.Println("Set requested")
	net.SendToRingLeader(s.Ctx, common.NewSetRequestCommand(value))
}
func (s *LinkedState)  ActionGetValue() {
	fmt.Println("Get requested")
	net.SendToRingLeader(s.Ctx, common.NewGetRequestCommand())
}
