package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const (
	//  PROBE_TIME = 3000 // Probe 20 times every minute
	ROUTE_TTL     = 200 // 3 seconds
	PROBE_TIMEOUT = 100  // 2 SECONDS
)

var (
	RoutingTable map[string]*link
)

type link struct {
	addr            string
	latency         float64
	loss            int
	TTLExpiration   time.Time // When was the last time I probed you (Creepy!)
	timeoutDeadline time.Time
	timerStartTime  time.Time
	lossBitMask			[100]bool
	lossIndex				int
}

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Not Enough Arguements")
		os.Exit(1)
	}

	// Set Up Addresses to use
	localNode, err := net.ResolveUDPAddr("udp", os.Args[1])
	CheckError(err)

	nodeAddr, err := net.ResolveUDPAddr("udp", os.Args[2])
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

	// Go Check Expiration of Routing Table and Listen for Incoming Data
	go RoutingTable[nodeAddr.String()].checkExpiration(conn)
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
	go link1.CheckForTimeout()
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
	go link1.CheckForTimeout()
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


	// Update Bitmask
	link1.lossBitMask[link1.lossIndex] = true
	link1.lossIndex = (link1.lossIndex + 1) % 100
	lossCount := 0

		// Count up losses in the bitmask
	for _, value := range link1.lossBitMask {
		if value == false {
			lossCount++
		}
	}
		// Store lossCount into the link's loss count
	link1.loss = lossCount
	fmt.Println("Loss", link1.loss, "EWMA: ", link1.latency)

}

func (link1 *link) recieveResponse2(conn *net.UDPConn) {
	// Takes time elapsed
	timeElapsed := time.Since(link1.timerStartTime)

	// Calculates Latency
	link1.latency = 0.9*link1.latency + 0.1*((timeElapsed/2).Seconds()*1000)

	// Update Bitmask
	link1.lossBitMask[link1.lossIndex] = true
	link1.lossIndex = (link1.lossIndex + 1) % 100
	lossCount := 0

		// Count up losses in the bitmask
	for _, value := range link1.lossBitMask {
		if value == false {
			lossCount++
		}
	}
		// Store lossCount into the link's loss count
	link1.loss = lossCount
	fmt.Println("Loss", link1.loss, "EWMA Lat: ", link1.latency)
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

func (link1 link) CheckForTimeout() {
	// Store the index of the bitMask's index which this function should check for timeout
	lossIndex := link1.lossIndex

	for {
		if link1.lossBitMask[lossIndex] == true {
			return
		}

		if time.Now().Before(link1.timeoutDeadline) == false {
			// Update Loss
			link1.lossBitMask[lossIndex] = false
			link1.lossIndex = (lossIndex + 1) % 100
			return
		}
	}

}

func getData(conn *net.UDPConn) {

	buffer := make([]byte, 64)

	for {
		_, addr, err := conn.ReadFromUDP(buffer)
		CheckError(err)

		if string(buffer[:2]) == "IP" {
			RoutingTable[addr.String()].sendResponse1(conn)
		} else if string(buffer[:2]) == "R1" {
			RoutingTable[addr.String()].sendResponse2(conn)
		} else if string(buffer[:2]) == "R2" {
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
