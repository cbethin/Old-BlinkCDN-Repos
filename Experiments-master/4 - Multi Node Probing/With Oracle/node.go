package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"strconv"
)

const (
	//  PROBE_TIME = 3000 // Probe 20 times every minute
	ROUTE_TTL     = 15000 // 10 seconds
	PROBE_TIMEOUT = 10000 // 5 seconds
)

var (
	RoutingTable map[string]*link
	nodeAddr  *net.UDPAddr
	nodeAddr2 *net.UDPAddr
	nodeAddr3 *net.UDPAddr
	oracleAddr *net.UDPAddr
)

type link struct {
	addr            string
	latency         float64
	loss            int
	TTLExpiration   time.Time // When was the last time I probed you (Creepy!)
	timeoutDeadline time.Time
	timerStartTime  time.Time
	lossBitMask     [100]bool
	lossIndex       int
}

func main() {

	if len(os.Args) != 5 {
		fmt.Println("Not Enough Arguements")
		os.Exit(1)
	}
	defer fmt.Println(time.Now())

	// Set Up Addresses to use
	localNode, err := net.ResolveUDPAddr("udp", os.Args[1])
	CheckError(err)

	nodeAddr, err = net.ResolveUDPAddr("udp", os.Args[2])
	CheckError(err)

	nodeAddr2, err = net.ResolveUDPAddr("udp", os.Args[3])
	CheckError(err)

	// nodeAddr3, err = net.ResolveUDPAddr("udp", os.Args[4])
	// CheckError(err)

	oracleAddr, err = net.ResolveUDPAddr("udp", os.Args[4])
	CheckError(err)

	conn, err := net.ListenUDP("udp", localNode)
	CheckError(err)

	RoutingTable = make(map[string]*link)

	// Setup Routing Table and It's Links
	var link1 link
	link1.addr = nodeAddr.String()
	link1.latency = 0
	link1.loss = 0
	link1.lossIndex = 0
	link1.TTLExpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)
	RoutingTable[nodeAddr.String()] = &link1

	var link2 link
	link2.addr = nodeAddr2.String()
	link2.latency = 0
	link2.loss = 0
	link2.lossIndex = 0
	link2.TTLExpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)
	RoutingTable[nodeAddr2.String()] = &link2

	// Go Check Expiration of Routing Table and Listen for Incoming Data
	go RoutingTable[nodeAddr.String()].checkExpiration(conn)
	go RoutingTable[nodeAddr2.String()].checkExpiration(conn)
	getData(conn)
	time.Sleep(30 * time.Second)
}

func (link1 *link) sendInitialProbe(conn *net.UDPConn) {
	// Update the link expiration as the first thing
	link1.TTLExpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)

	// Create a Buffer now
	buffer := make([]byte, 64)
	copy(buffer[:], "IP")

	nodeAddr, err := net.ResolveUDPAddr("udp", link1.addr)
	CheckError(err)

	_, err = conn.WriteToUDP(buffer, nodeAddr)
	CheckError(err)

	link1.timerStartTime = time.Now()
	link1.timeoutDeadline = link1.timerStartTime.Add(time.Millisecond * PROBE_TIMEOUT)
	// go link1.CheckForTimeout()
	//
	// probeCounter++

	// fmt.Println("Sent IP")
	// fmt.Println("Loss:", lossCounter, "Probes:", probeCounter)
	// fmt.Println("-------------")
}

func (link1 *link) sendResponse1(conn *net.UDPConn) {
	link1.TTLExpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)

	// Create a Buffer now
	buffer := make([]byte, 64)
	copy(buffer[:], "R1")

	nodeAddr, err := net.ResolveUDPAddr("udp", link1.addr)
	CheckError(err)

	_, err = conn.WriteToUDP(buffer, nodeAddr)
	CheckError(err)

	// probeCounter++

	link1.timerStartTime = time.Now()
	link1.timeoutDeadline = link1.timerStartTime.Add(time.Millisecond * PROBE_TIMEOUT)
	// go link1.CheckForTimeout()

	// fmt.Println("Sent R1")
	// fmt.Println("Loss:", lossCounter, "Probes:", probeCounter)
	// fmt.Println("-------------")
}

func (link1 *link) sendResponse2(conn *net.UDPConn) {
	// hasReceived = true

	// Takes time elapsed
	timeElapsed := time.Since(link1.timerStartTime)

	// Create a Buffer & sends response 2
	buffer := make([]byte, 64)
	copy(buffer[:], "R2")

	nodeAddr, err := net.ResolveUDPAddr("udp", link1.addr)
	CheckError(err)

	_, err = conn.WriteToUDP(buffer, nodeAddr)
	CheckError(err)

	// Calculates Latency
	link1.latency = 0.9*link1.latency + 0.1*((timeElapsed/2).Seconds()*1000)

	// Send Routing Table To Oracle
	for _, value := range RoutingTable {
		buffer := make([]byte, 50)

		copy(buffer[:], LinkToByteArray(*value))
		_, err := conn.WriteToUDP(buffer, oracleAddr)
		CheckError(err)
	}

}

func (link1 *link) recieveResponse2(conn *net.UDPConn) {
	// Takes time elapsed
	timeElapsed := time.Since(link1.timerStartTime)

	// Calculates Latency
	link1.latency = 0.9*link1.latency + 0.1*((timeElapsed/2).Seconds()*1000)

	// Send Routing Table To Oracle
	for _, value := range RoutingTable {
		buffer := make([]byte, 50)

		copy(buffer[:], LinkToByteArray(*value))
		_, err := conn.WriteToUDP(buffer, oracleAddr)
		CheckError(err)
	}

}

func (link1 *link) checkExpiration(conn *net.UDPConn) {

	for {
		//
		if time.Now().Before(link1.TTLExpiration) == false {
			// If your Current time is before the expiation time
			// You do Nothing!
			// FOR FUTURE : The link that requires the probing funtion
			//              i.e. the one that times out, can update the
			//              linkNumber and call the SendProbe
			//              RoutingTable[linkNumber].SendProbe(nodeAddr,conn)
			link1.sendInitialProbe(conn)
		}
		//fmt.Println("Checking Expiration")
		time.Sleep(100 * time.Millisecond)
	}
}

func getData(conn *net.UDPConn) {

	buffer := make([]byte, 64)

	for {
		_, addr, err := conn.ReadFromUDP(buffer)
		CheckError(err)
		fmt.Println("-----------------")
		fmt.Println(string(buffer[:2]), "from", addr.String())

		if string(buffer[:2]) == "IP" {
			RoutingTable[addr.String()].sendResponse1(conn)
		} else if string(buffer[:2]) == "R1" {
			RoutingTable[addr.String()].sendResponse2(conn)
		} else if string(buffer[:2]) == "R2" {
			RoutingTable[addr.String()].recieveResponse2(conn)
		} else {
			fmt.Println("Packet recieved is an undefined type")
		}
		fmt.Println("------------------")
		fmt.Println("Link 1 -- Lat:", RoutingTable[nodeAddr.String()].latency)
		fmt.Println("Link 2 -- Lat:", RoutingTable[nodeAddr2.String()].latency)
	}
}

func LinkToByteArray(link link) []byte {
  byteArray := make([]byte, 192)
  copy(byteArray[0:16], strconv.FormatFloat(link.latency, 'f', 15, 64))
  copy(byteArray[16:48], link.addr + "!")
  fmt.Println(string(byteArray))
  return byteArray
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Err:", err)
		os.Exit(1)
	}
}
