package main

import (
	"fmt"
	"net"
	"time"

	"./blink"
)

func main() {
	localAddr, err := net.ResolveUDPAddr("udp", ":8000")
	blink.CheckError(err)

	hostAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8001")
	blink.CheckError(err)

	conn, err := net.ListenUDP("udp", localAddr)
	blink.CheckError(err)

	fmt.Println("Connected to", localAddr.String())

	p := blink.Packet{Number: 1, Time1: time.Now(), Time2: time.Now()}
	b := blink.PacketToJsonToByteArray(p)

	// for {
	_, err = conn.WriteToUDP(b, hostAddr)
	blink.CheckError(err)
	fmt.Println("Sent")
	// }

}
