package udp

import (
	"fmt"
	"net"
)

type stUdp struct {
	id   uint32
	addr net.UDPAddr
}

func newUdp(id uint32, addr net.UDPAddr) *stUdp {
	return &stUdp{
		id:   id,
		addr: addr,
	}
}

func (p *stUdp) start() {

}

func (p *stUdp) close() {

}

func (p *stUdp) write(msgType int, msg []byte) {

}
