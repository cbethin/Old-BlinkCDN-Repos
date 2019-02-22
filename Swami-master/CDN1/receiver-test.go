package main

import (
	"fmt"
	"net"

	"./blink"
)

func main() {
	localAddr, err := net.ResolveUDPAddr("udp", ":8001")
	blink.CheckError(err)

	// hostAddr, err := net.ResolveUDPAddr("udp", ":8001")
	// blink.CheckError(err)

	conn, err := net.ListenUDP("udp", localAddr)
	blink.CheckError(err)

	fmt.Println("Connected to", localAddr.String())

	b := make([]byte, 1024)
	_, _, err = conn.ReadFromUDP(b)
	blink.CheckError(err)
	fmt.Println("Packet:", blink.ByteArrayToJSONToPacket(b))
	fmt.Println(blink.ByteArrayToJSONToPacket(b).Number)
	fmt.Println("T1:", !blink.ByteArrayToJSONToPacket(b).Time1.IsZero() && !blink.ByteArrayToJSONToPacket(b).Time2.IsZero())
	// fmt.Println("T2:", !blink.ByteArrayToJsonToPacket(b).Time2.IsZero())
}
