package main

import (
  "fmt"
  "net"
  "os"
  "time"
  "./blink"
)

type link struct {
	addr            string
	latency         float64
	loss            int
}

func main() {
  if len(os.Args) != 2 {
    fmt.Println("Not enough arguments")
    os.Exit(1)
  }

  // Set Up listening connection
  oracleAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
  blink.CheckError(err)

  conn, err := net.ListenUDP("udp", oracleAddr)
  blink.CheckError(err)
  fmt.Println("Listening...")

  // Set up Read Buffer
  buffer := make([]byte, 1024)
  //var RoutingTable [2]link

  /* Read incoming data, extract the packet data from the received Blink Packet
     and then print out the information in the received routing table
  */
  for {
    _, addr, err := conn.ReadFromUDP(buffer)
    blink.CheckError(err)
    packetData := blink.ExtractPacketData(buffer)
    link := blink.ByteArrayToLink(packetData)
    fmt.Println("----------\nfrom", addr.String())
    fmt.Println("Latency:", link.Latency)
    fmt.Println("Received at:", time.Now().String())
  }
}
