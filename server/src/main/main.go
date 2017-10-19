package main

import (
	"fmt"
	"network/udp"
)

func onSvrError(err error) {

}

func onDisconnect(id uint32) {

}

func main() {
	fmt.Println("main.start...")

	udp.NewServer("192.168.142.140:12345", onSvrError)
}
