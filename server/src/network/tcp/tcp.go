package tcp

import (
	"fmt"
	"net"
	"network/message"
	"sync"
)

type stTcp struct {
	id     uint32
	conn   *net.TCPConn
	buff   []byte
	chsend chan []byte
	chexit chan bool
	wg     sync.WaitGroup
	chmsg  chan *message.Message
	dh     func(uint32)
}

func newTcp(id uint32, conn *net.TCPConn, chmsg chan *message.Message, dh func(uint32)) *stTcp {
	return &stTcp{
		id:     id,
		conn:   conn,
		buff:   make([]byte, 0),
		chsend: make(chan []byte),
		chexit: make(chan bool),
		chmsg:  chmsg,
		dh:     dh,
	}
}

func (p *stTcp) start() {
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

func (p *stTcp) close() {
	close(p.chexit)
	p.conn.Close()
	p.wg.Wait()
}

func (p *stTcp) write(msgType int, msg []byte) {
	data := p.pack(msgType, msg)
	p.chsend <- data
}

func (p *stTcp) recv() (err error) {
	defer func() {
		if err != nil {
			select {
			case <-p.chexit:
				err = nil
			default:
				p.dh(p.id)
			}
		}
	}()

	for {
		select {
		case <-p.chexit:
			return
		default:
			data := make([]byte, 1024)
			sz, err := p.conn.Read(data)
			if err != nil {
				fmt.Println("tcp read err:", err)
				return err
			}
			p.buff = append(p.buff, data[:sz]...)
			p.unpack()
		}
	}
}

func (p *stTcp) send() {
	for {
		select {
		case <-p.chexit:
			return
		case data := <-p.chsend:
			p.conn.Write(data)
		}
	}
}

func (p *stTcp) pack(msgType int, msg []byte) []byte {
	sz := len(msg)
	data := make([]byte, 0)
	if sz > 127 {
		data = append(data, byte(((sz>>8)&0x7f)|0x80))
		data = append(data, byte(sz&0xff))
	} else {
		data = append(data, byte(sz))
	}
	data = append(data, message.ToBytes(msgType)...)
	data = append(data, msg...)

	return data
}

func (p *stTcp) unpack() {
	var lenBuff int
	var sz int
	var lenSz int
	for {
		lenBuff = len(p.buff)
		if lenBuff == 0 {
			break
		}
		if p.buff[0] > 127 {
			if lenBuff < 2 {
				break
			}
			sz = int((p.buff[0]&0x7f))*256 + int(p.buff[1])
			lenSz = 2
		} else {
			sz = int(p.buff[0])
			lenSz = 1
		}
		if sz+lenSz+2 > lenBuff {
			break
		}
		p.buff = p.buff[lenSz:]

		msgType := message.ToInt(p.buff[:2])
		p.buff = p.buff[2:]
		msg := &message.Message{
			ConnId:  p.id,
			MsgType: msgType,
			Data:    p.buff[:sz],
		}
		p.chmsg <- msg
	}
}
