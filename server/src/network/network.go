package network

import (
	"fmt"
	"net"
	"sync/atomic"
)

var conns map[uint32]*Connection
var connId uint32

func acceptClient(conn net.Conn) {
	connection := NewConnection(conn, 1024)
	atomic.AddUint32(&connId, 1)
	conns[connId] = connection
}

func Start() {
	fmt.Println("network.start...")

	connId = 0
}
