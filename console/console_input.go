package console

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/metopa/distributed_variable/common"
)

var setValueRegex = regexp.MustCompile("set\\s+([0-9]+)")

func ListenConsole(ctx *common.Context, stdInChan chan string) {
	for {
		str, ok := <-stdInChan
		if !ok {
			return
		} else {
			if handleAction(str, ctx.GetState()) {
				return
			}
		}
	}
}

func handleAction(action string, handler common.ActionHandler) bool {
	if action == "" {
		return false
	}
	if action == "?" {
	} else if action == "run" {
		go handler.ActionStartChRo()
	} else if action == "get" {
		go handler.ActionGetValue()
	} else if action == "sync" {
		go handler.ActionSync()
	} else {
		m := setValueRegex.FindStringSubmatch(action)

		if m != nil {
			n, err := strconv.Atoi(string(m[1]))
			if err != nil {
				fmt.Printf("Failed to parse number: %v->%v\n", string(m[1]), err)
			} else {
				go handler.ActionSetValue(n)
			}
		} else if action == "leave" {
			if handler.ActionLeave() {
				return true
			}
		} else if action == "disconnect" {
			return true
		} else {
			fmt.Printf("Unknown command: %v\n", action)
			fmt.Print("Valid commands:\n\trun\n\tleave\n\tdisconnect\n\tsync\n\tget\n\tset %d\n")
		}
	}
	return false
}
