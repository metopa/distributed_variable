package state

import (
	"runtime"
	"strings"

	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/logger"
)

type NullState struct{}

func (s *NullState) Init() {
	logger.Info("Current state: NULL")
}

func (s *NullState) NewPeer(sender common.PeerAddr, addr common.PeerAddr,
	name string, shouldReply bool) {
	NotHandled()
}

func (s *NullState) LeaderChanged(sender common.PeerAddr, leader common.PeerAddr) {
	NotHandled()
}
func (s *NullState) LinkedPeersChanged(sender common.PeerAddr,
	loPeer common.PeerAddr, hiPeer common.PeerAddr) {
	NotHandled()
}
func (s *NullState) PeerReported(sender common.PeerAddr, reportedPeer common.PeerAddr) {
	NotHandled()
}
func (s *NullState) DistanceRequested(sender common.PeerAddr, source common.PeerAddr) {
	NotHandled()
}
func (s *NullState) DistanceReceived(sender common.PeerAddr, distance int) {
	NotHandled()
}
func (s *NullState) SyncPeers(sender common.PeerAddr, values []string) {
	NotHandled()
}

func (s *NullState) ChRoIdReceived(sender common.PeerAddr, id int) {
	NotHandled()
}

func (s *NullState) RingJoinRequested(sender common.PeerAddr, source common.PeerAddr) {
	NotHandled()
}
func (s *NullState) RingLeaveAnnounced(sender common.PeerAddr, source common.PeerAddr) {
	NotHandled()
}
func (s *NullState) ValueGetRequested(sender common.PeerAddr, source common.PeerAddr) {
	NotHandled()
}
func (s *NullState) GotValue(sender common.PeerAddr, value int) {
	NotHandled()
}
func (s *NullState) ValueSetRequested(sender common.PeerAddr,
	source common.PeerAddr, value int) {
	NotHandled()
}
func (s *NullState) ValueSetConfirmed(sender common.PeerAddr) {
	NotHandled()
}
func (s *NullState) ActionSetValue(value int) {
	NotHandled()
}
func (s *NullState) ActionGetValue() {
	NotHandled()
}
func (s *NullState) ActionStartChRo() {
	NotHandled()
}
func (s *NullState) ActionLeave() {
	NotHandled()
}
func (s *NullState) ActionDisconnect() {
	NotHandled()
}
func (s *NullState) ActionReconnect() {
	NotHandled()
}
func (s *NullState)ActionReportPeer(addr common.PeerAddr) {
	NotHandled()
}
func (s *NullState) Name() string {
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
