package udp

import (
	"fmt"
	"net"
	"network/netdef"
	"sync"
	"sync/atomic"
)

type stSendPackage struct {
	addr net.Addr
	data []byte
}

type stUdpServer struct {
	svr      *net.UDPConn
	conns    map[string]*stUdp
	id2conn  map[uint32]*stUdp
	nextId   uint32
	wg       sync.WaitGroup
	chsend   chan *stSendPackage
	chexit   chan bool
	eh       func(error)
	chClient chan netdef.Connection
}

func NewServer(addr string, eh func(error)) (stUdpServer, error) {
	var udpSvr stUdpServer

	laddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println("udp addr[", addr, "] error:", err)
		return udpSvr, err
	}
	svr, err := net.ListenUDP("udp", laddr)
	if err != nil {
		fmt.Println("udp listen[", addr, "] error:", err)
		return udpSvr, err
	}
	return stUdpServer{
		svr:     svr,
		conns:   make(map[string]*stUdp),
		id2conn: make(map[uint32]*stUdp),
		nextId:  0,
		chsend:  make(chan *stSendPackage, 1024),
		chexit:  make(chan bool),
		eh:      eh,
	}, nil
}

func (p stUdpServer) Start(chClient chan netdef.Connection) {
	p.chClient = chClient

	p.wg.Add(2)
	go func() {
		defer p.wg.Done()
		p.recv()
	}()

	go func() {
		defer p.wg.Done()
		p.send()
	}()
}

func (p stUdpServer) Close() {
	p.svr.Close()
	close(p.chexit)
	p.wg.Wait()
	for _, conn := range p.conns {
		conn.Close()
	}
}

func (p stUdpServer) recv() (err error) {
	defer func() {
		if err != nil {
			select {
			case <-p.chexit:
				err = nil
			default:
				p.eh(err)
			}
		}
	}()
	for {
		select {
		case <-p.chexit:
			return
		default:
			data := make([]byte, 512)
			sz, addr, err := p.svr.ReadFrom(data)
			if err != nil {
				fmt.Println("udp[", p.svr.LocalAddr(), "] read error:", err)
				return err
			}
			fmt.Println("receive data from", addr, "data:", data[:sz])
			conn, ok := p.conns[addr.String()]
			if ok {
				conn.recv(data[:sz])
			} else {
				p.accept(addr)
			}
		}
	}
}

func (p stUdpServer) send() {
	for {
		select {
		case <-p.chexit:
			break
		case pkg := <-p.chsend:
			fmt.Println("send data:", pkg.data)
			p.svr.WriteTo(pkg.data, pkg.addr)
			break
		}
	}
}

func (p stUdpServer) accept(addr net.Addr) {
	atomic.AddUint32(&p.nextId, 1)
	udp := newUdp(p.nextId, addr, p.chsend)
	p.conns[addr.String()] = udp
	p.id2conn[p.nextId] = udp

	p.chClient <- udp
}
