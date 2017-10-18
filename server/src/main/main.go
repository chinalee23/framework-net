package main

import (
	"fmt"
	"network/tcp"
)

func onSvrError(err error) {

}

func onDisconnect(id uint32) {

}

func main() {
	fmt.Println("main.start...")

	tcp.NewServer("192.168.142.140:12345", onSvrError, onDisconnect)
}
