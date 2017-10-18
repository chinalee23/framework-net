package udp

import (
	"fmt"
	"net"
	"network/message"
	"sync"
)

type UdpServer struct {
	svr     *net.UDPConn
	wg      sync.WaitGroup
	addr2id map[net.UDPAddr]uint32
	id2addr map[uint32]net.UDPAddr
}

func NewServer(addr string, eh func(error), dh func(uint32)) *UdpServer {
	laddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println("udp addr[", addr, "] error:", err)
		return nil
	}
	svr, err := net.ListenUDP("udp", laddr)
	if err != nil {
		fmt.Println("udp listen[", addr, "] error:", err)
		return nil
	}
	return &UdpServer{
		svr: svr,
	}
}

func (p *UdpServer) Start() {
	p.wg.Add(2)
}

func (p *UdpServer) Close() {
	p.svr.Close()
	p.wg.Wait()
}

func (p *UdpServer) Read() []*message.Message {
	return nil
}

func (p *UdpServer) Write(connId uint32, msgType int, msg []byte) bool {
	return true
}

func (p *UdpServer) recv() {
	for {
		data := make([]byte, 512)
		sz, addr, err := p.svr.ReadFrom(data)
		if err != nil {
			fmt.Println("udp[", p.svr.LocalAddr(), "] read error:", err)
			return
		}

	}
}

func (p *UdpServer) send() {

}
