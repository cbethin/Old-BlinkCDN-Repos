package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const (
	//  PROBE_TIME = 3000 // Probe 20 times every minute
	ROUTE_TTL     = 2000 // 3 seconds
	PROBE_TIMEOUT = 500  // 2 SECONDS
)

var (
	lossBitMask  [100]bool
	bitMaskIndex int
	lossCounter  int
	RoutingTable map[string]*link
)

type link struct {
	addr            string
	latency         float64
	loss            int
	TTLExpiration   time.Time // When was the last time I probed you (Creepy!)
	timeoutDeadline time.Time
	timerStartTime  time.Time
}

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Not Enough Arguements")
		os.Exit(1)
	}

	localNode, err := net.ResolveUDPAddr("udp", os.Args[1])
	CheckError(err)

	nodeAddr, err := net.ResolveUDPAddr("udp", os.Args[2])
	CheckError(err)

	conn, err := net.ListenUDP("udp", localNode)
	CheckError(err)

	// Initialize Bitmask
	bitMaskIndex = 0
	lossCounter = 0

	RoutingTable = make(map[string]*link)

	// Routing Table Block
	var link1 link
	link1.addr = nodeAddr.String()
	link1.latency = 0
	link1.loss = 0
	link1.TTLExpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)

	RoutingTable[nodeAddr.String()] = &link1

	go RoutingTable[nodeAddr.String()].checkExpiration(conn)
	getData(conn)
	time.Sleep(30 * time.Second)
}

// func (link1 *link) SendProbe(conn *net.UDPConn) {
//
// 	// Read From Bufferint
// 	_, _, err := conn.ReadFromUDP(buffer)
// 	if err != nil {
// 		fmt.Println("You Lost a Packet")
// 		lossBitMask[bitMaskIndex] = false
// 	} else {
// 		lossBitMask[bitMaskIndex] = true
// 	}
//
// 	// INCREMENT bitMaskIndex
// 	bitMaskIndex = (bitMaskIndex + 1) % 100
//
// 	lossCounter = 0
// 	for _, value := range lossBitMask {
// 		if value == false {
// 			lossCounter++
// 		}
// 	}
//
// 	fmt.Println("Loss Count: ", lossCounter)
// 	link1.loss = lossCounter
//
// 	timeElapsed := time.Since(probeTime)
// 	fmt.Println("Current Latency:", timeElapsed/2)
//
// 	_, err = conn.WriteToUDP(buffer, nodeAddr)
// 	CheckError(err)
//
// }

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

	link1.timerStartTime = time.Now()
	link1.timeoutDeadline = link1.timerStartTime.Add(time.Millisecond * PROBE_TIMEOUT)
}

func (link1 *link) sendResponse2(conn *net.UDPConn) {

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
	fmt.Println("1 Current Latency:", timeElapsed/2)
	fmt.Println(" 2 EWMA Latency: ", link1.latency)

}

func (link1 *link) recieveResponse2(conn *net.UDPConn) {
	// Takes time elapsed
	timeElapsed := time.Since(link1.timerStartTime)
	// Calculates Latency
	link1.latency = 0.9*link1.latency + 0.1*((timeElapsed/2).Seconds()*1000)
	fmt.Println("  3 Current Latency:", timeElapsed/2)
	fmt.Println("   4 EWMA Latency: ", link1.latency)
}

func (link1 *link) checkExpiration(conn *net.UDPConn) {
	for {
		//
		if time.Now().Before(link1.TTLExpiration) {
			// If your Current time is before the expiation time
			// You do Nothing!
			// FOR FUTURE : The link that requires the probing funtion
			//              i.e. the one that times out, can update the
			//              linkNumber and call the SendProbe
			//              RoutingTable[linkNumber].SendProbe(nodeAddr,conn)
		} else {
			link1.sendInitialProbe(conn)
		}
	}
}

func getData(conn *net.UDPConn) {

	buffer := make([]byte, 64)

	for {
		_, addr, err := conn.ReadFromUDP(buffer)
		CheckError(err)

		if string(buffer[:2]) == "IP" {
			fmt.Println("Recieved IP")
			RoutingTable[addr.String()].sendResponse1(conn)
		} else if string(buffer[:2]) == "R1" {
			fmt.Println("Recieved R1")
			RoutingTable[addr.String()].sendResponse2(conn)
		} else if string(buffer[:2]) == "R2" {
			fmt.Println("Recieved R2")
			RoutingTable[addr.String()].recieveResponse2(conn)
		} else {
			fmt.Println("Packet recieved is an undefined type")
		}
	}
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Err:", err)
		os.Exit(1)
	}
}
