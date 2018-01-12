package common

import (
	"fmt"
	"net"
	"strings"
)

type PeerAddr string

type PeerInfo struct {
	Addr PeerAddr
	Name string
}


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

func GetInterfaceIPv4Addr(iface *net.Interface) (net.IP, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	for i, _ := range addrs {
		addr, _, err := net.ParseCIDR(addrs[i].String())
		if err != nil {
			return nil, fmt.Errorf("%v: %v", addrs[i].String(), err)
		}
		v4 := addr.To4()
		if v4 != nil {
			return v4, nil
		}
	}
	return nil, fmt.Errorf("%v has no IPv4 address", iface.Name)
}

func GetTransmissionInfoString(source, from, to, destination PeerAddr) string {
	strip := func(s PeerAddr) string {
		ip := strings.Split(string(s), ":")
		ip = strings.Split(ip[0], ".")
		return ip[len(ip) - 1]
	}
	res := strip(source) + "-"
	if from != source {
		res += "-" + strip(from)
	}
	if to != destination {
		res += ".." + strip(to)
	}
	res += "->" + strip(destination)
	return res
}