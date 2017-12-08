package network

import (
	"network/netdef"
	"network/tcp"
	"network/udp"
)

func NewTcpServer(addr string, eh func(error)) (netdef.Server, error) {
	return tcp.NewServer(addr, eh)
}

func NewUdpServer(addr string, eh func(error)) (netdef.Server, error) {
	return udp.NewServer(addr, eh)
}
