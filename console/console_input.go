package console

import (
	"bufio"
	"os"
	"strings"

	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/logger"
)

func ListenConsole(ctx *common.Context, stop *chan struct{}) {
	ch := make(chan string)
	go func(ch chan string) {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')
			if err != nil {
				close(ch)
				return
			}
			s = strings.TrimRight(s, "\n \t")
			ch <- s
		}
	}(ch)

	for {
		select {
		case str, ok := <-ch:
			if !ok {
				return
			} else {
				handleAction(str, ctx.State)
			}
		case <-*stop:
			return
		}
	}
}
func handleAction(action string, handler common.ActionHandler) {
	if action == "start ch ro" {
		handler.ActionStartChRo()
	} else {
		logger.Warn("Unknown command: %v", action)
	}
}
