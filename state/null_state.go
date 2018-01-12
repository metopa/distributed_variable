package state

import (
	"runtime"
	"strings"

	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/logger"
)

type NullState struct{}

func (h *NullState) NewPeer(sender common.PeerAddr, addr common.PeerAddr,
	name string, shouldReply bool) {
	NotHandled()
}

func (h *NullState) LeaderChanged(sender common.PeerAddr, leader common.PeerAddr) {
	NotHandled()
}
func (h *NullState) LinkedPeersChanged(sender common.PeerAddr,
	loPeer common.PeerAddr, hiPeer common.PeerAddr) {
	NotHandled()
}
func (h *NullState) PeerReported(sender common.PeerAddr, reportedPeer common.PeerAddr) {
	NotHandled()
}
func (h *NullState) DistanceRequested(sender common.PeerAddr, source common.PeerAddr) {
	NotHandled()
}
func (h *NullState) DistanceReceived(sender common.PeerAddr, distance int) {
	NotHandled()
}
func (h *NullState) SyncPeers(sender common.PeerAddr, values []string) {
	NotHandled()
}

func (h *NullState) ChRoIdReceived(sender common.PeerAddr, id int) {
	NotHandled()
}

func (h *NullState) RingJoinRequested(sender common.PeerAddr, source common.PeerAddr) {
	NotHandled()
}
func (h *NullState) RingLeaveAnnounced(sender common.PeerAddr, source common.PeerAddr) {
	NotHandled()
}
func (h *NullState) ValueGetRequested(sender common.PeerAddr, source common.PeerAddr) {
	NotHandled()
}
func (h *NullState) GotValue(sender common.PeerAddr, value int) {
	NotHandled()
}
func (h *NullState) ValueSetRequested(sender common.PeerAddr,
	source common.PeerAddr, value int) {
	NotHandled()
}
func (h *NullState) ValueSetConfirmed(sender common.PeerAddr) {
	NotHandled()
}
func (h *NullState) ActionSetValue(value int) {
	NotHandled()
}
func (h *NullState) ActionGetValue() {
	NotHandled()
}
func (h *NullState) ActionStartChRo() {
	NotHandled()
}
func (h *NullState) ActionLeave() {
	NotHandled()
}
func (h *NullState) ActionDisconnect() {
	NotHandled()
}
func (h *NullState) ActionReconnect() {
	NotHandled()
}
func (h *NullState) Name() string {
	return "Null state"
}

func NotHandled() {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	eventF := runtime.FuncForPC(pc[0])

	runtime.Callers(3, pc)
	callerF := runtime.FuncForPC(pc[0])
	file, line := callerF.FileLine(pc[0])
	eventFName := strings.Split(eventF.Name(), ".")

	logger.Warn("%v event can't be handled in current state. Called from %s:%d",
		eventFName[len(eventFName)-1], file, line)
}
