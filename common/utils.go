package common

import (
	"net"
	"strings"
)

type PeerAddr string

func IsTimeoutError(err error) bool {
	nerr, ok := err.(net.Error)
	return ok && nerr.Timeout()
}

func GetInterfaceNames() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	var ifaceNames []string
	for _, i := range ifaces {
		ifaceNames = append(ifaceNames, i.Name)
	}
	return strings.Join(ifaceNames, ", ")
}
