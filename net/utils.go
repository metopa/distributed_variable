package net

import (
	"sync/atomic"
	"time"

	"github.com/metopa/distributed_variable/common"
)

func StartChRoTimer(ctx *common.Context) {
	if atomic.CompareAndSwapInt32(&ctx.StartedChRoTimer, 0, 1) {
		go func() {
			time.Sleep(ctx.ChRoTimerDur / 2)
			SendToHi(ctx, common.NewChangRobertIdCmd(ctx.PeerId), false)
		}()
	}
}
