package state

import (
	"runtime"

	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/logger"
)

type NullState struct{}

func (h *NullState) NewPeer(sender common.PeerAddr, addr common.PeerAddr,
	name string, shouldReply bool) {
	notImplementedWarning()
}

func (h *NullState) LeaderChanged(sender common.PeerAddr, leader common.PeerAddr) {
	notImplementedWarning()
}
func (h *NullState) LinkedPeersChanged(sender common.PeerAddr,
	loPeer common.PeerAddr, hiPeer common.PeerAddr) {
	notImplementedWarning()
}
func (h *NullState) PeerReported(sender common.PeerAddr, reportedPeer common.PeerAddr) {
	notImplementedWarning()
}
func (h *NullState) DistanceRequested(sender common.PeerAddr, source common.PeerAddr) {
	notImplementedWarning()
}
func (h *NullState) DistanceReceived(sender common.PeerAddr, distance int) {
	notImplementedWarning()
}
func (h *NullState) SyncPeers(sender common.PeerAddr, values []string) {
	notImplementedWarning()
}

func (h *NullState) ChRoIdReceived(sender common.PeerAddr, id int) {
	notImplementedWarning()
}

func (h *NullState) RingJoinRequested(sender common.PeerAddr, source common.PeerAddr) {
	notImplementedWarning()
}
func (h *NullState) RingLeaveAnnounced(sender common.PeerAddr, source common.PeerAddr) {
	notImplementedWarning()
}
func (h *NullState) ValueGetRequested(sender common.PeerAddr, source common.PeerAddr) {
	notImplementedWarning()
}
func (h *NullState) GotValue(sender common.PeerAddr, value int) {
	notImplementedWarning()
}
func (h *NullState) ValueSetRequested(sender common.PeerAddr,
	source common.PeerAddr, value int) {
	notImplementedWarning()
}
func (h *NullState) ValueSetConfirmed(sender common.PeerAddr) {
	notImplementedWarning()
}

func notImplementedWarning() {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	eventF := runtime.FuncForPC(pc[0])

	runtime.Callers(3, pc)
	callerF := runtime.FuncForPC(pc[0])
	file, line := callerF.FileLine(pc[0])

	logger.Warn("%v event can't be handled in current state. Called from %s:%d", eventF.Name(), file, line)
}
