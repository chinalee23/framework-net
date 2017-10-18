package tcp

import (
	"fmt"
	"net"
	"network/message"
	"sync"
	"sync/atomic"
)

type TcpServer struct {
	svr    *net.TCPListener
	conns  map[uint32]*stTcp
	nextId uint32
	chmsg  chan *message.Message
	wg     sync.WaitGroup
	eh     func(error)
	dh     func(uint32)
}

func NewServer(addr string, eh func(error), dh func(uint32)) *TcpServer {
	laddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println("tcp addr[", addr, "] error:", err)
		return nil
	}
	svr, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		fmt.Println("tcp listen[", addr, "] error:", err)
		return nil
	}
	return &TcpServer{
		svr:    svr,
		conns:  make(map[uint32]*stTcp),
		nextId: 0,
		chmsg:  make(chan *message.Message),
		eh:     eh,
		dh:     dh,
	}
}

func (p *TcpServer) Start() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			conn, err := p.svr.AcceptTCP()
			if err != nil {
				fmt.Println("tcp[", p.svr.Addr(), "] accept error:", err)
				p.eh(err)
				return
			}
			p.accept(conn)
		}
	}()
}

func (p *TcpServer) Close() {
	p.svr.Close()
	p.wg.Wait()
	for _, conn := range p.conns {
		conn.close()
	}
}

func (p *TcpServer) Read() []*message.Message {
	msgs := make([]*message.Message, 0)
	for {
		select {
		case msg := <-p.chmsg:
			msgs = append(msgs, msg)
			break
		default:
			return msgs
		}
	}
}

func (p *TcpServer) Write(connId uint32, msgType int, msg []byte) bool {
	conn, ok := p.conns[connId]
	if !ok {
		return false
	}
	conn.write(msgType, msg)
	return true
}

func (p *TcpServer) accept(conn *net.TCPConn) {
	atomic.AddUint32(&p.nextId, 1)
	tcp := newTcp(p.nextId, conn, p.chmsg, p.dh)
	p.conns[p.nextId] = tcp
	tcp.start()
}
