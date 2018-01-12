package console

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/logger"
)

var setValueRegex = regexp.MustCompile("set\\s+([0-9]+)")

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
	if action == "?" {
		fmt.Print("Available commands:\n\tstart\n\tget\n\tset %d\n\tleave\n\force-leave\n\treconnect\n")
	} else if action == "start" {
		go handler.ActionStartChRo()
	} else if action == "get" {
		go handler.ActionGetValue()
	} else {
		m := setValueRegex.Find([]byte(action))

		if m != nil {
			n, err := strconv.Atoi(string(m[1]))
			if err != nil {
				fmt.Printf("Failed to parse number: %v->%v\n", string(m[1]), err)
			} else {
				go handler.ActionSetValue(n)
			}
		} else if action == "leave" {
			go handler.ActionLeave()
		} else if action == "force-leave" {
			go handler.ActionDisconnect()
		} else if action == "reconnect" {
			go handler.ActionReconnect()
		} else {
			logger.Warn("Unknown command: %v", action)
		}
	}
}
