package network

import (
	"errors"
	"sync"
)

var (
	ERR_EXIT = errors.New("exit")
)

type Socket interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}
type Connection struct {
	conn   Socket
	wg     sync.WaitGroup
	mutex  sync.Mutex
	chexit chan bool
	chsend chan []byte
	chrecv chan []byte
}

func NewConnection(conn Socket, maxcount int) *Connection {
	count := maxcount
	if count < 1024 {
		count = 1024
	}
	return &Connection{
		conn:   conn,
		chsend: make(chan []byte, count),
		chexit: make(chan bool),
		chrecv: make(chan []byte, count),
	}
}

func (p *Connection) Start() {
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

func (p *Connection) Close() {
	close(p.chexit)
	p.conn.Close()
	p.wg.Wait()
}

func (p *Connection) Write(data []byte) error {
	select {
	case <-p.chexit:
		return ERR_EXIT
	case p.chsend <- data:
		break
	}
	return nil
}

func (p *Connection) send() {
	for {
		select {
		case <-p.chexit:
			return
		case data := <-p.chsend:
			_, err := p.conn.Write(data)
			if err != nil {
				return
			}
		}
	}
}

func (p *Connection) recv() (err error) {
	defer func() {
		if err != nil {
			select {
			case <-p.chexit:
				err = nil
			default:
				return
			}
		}
	}()
	for {
		select {
		case <-p.chexit:
			return nil
		default:
			break
		}
		data := make([]byte, 1024)
		sz, err := p.conn.Read(data)
		if err != nil {
			return err
		}
		p.chrecv <- data[0:sz]
	}
	return nil
}
