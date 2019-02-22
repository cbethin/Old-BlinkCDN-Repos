package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const (
	//  PROBE_TIME = 3000 // Probe 20 times every minute
	ROUTE_TTL     = 15000 // 10 seconds
	PROBE_TIMEOUT = 5000  // 5 seconds
)

var (
	RoutingTable map[string]*link
	lossCounter int
	probeCounter int
	hasReceived bool
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

	lossCounter = 0
	probeCounter = 0

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

	probeCounter++

	fmt.Println("Sent IP")
	fmt.Println("Loss:", lossCounter, "Probes:", probeCounter)
	fmt.Println("-------------")
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

	probeCounter++

	link1.timerStartTime = time.Now()
	link1.timeoutDeadline = link1.timerStartTime.Add(time.Millisecond * PROBE_TIMEOUT)
	go link1.CheckForTimeout()

	fmt.Println("Sent R1")
	fmt.Println("Loss:", lossCounter, "Probes:", probeCounter)
	fmt.Println("-------------")
}

func (link1 *link) sendResponse2(conn *net.UDPConn) {
	hasReceived = true

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

	//fmt.Println("Loss", lossCounter, "Probes:", probeCounter, "EWMA: ", link1.latency)
	fmt.Println("Sent R2")
	fmt.Println("Loss:", lossCounter, "Probes:", probeCounter)
	fmt.Println("-------------")
}

func (link1 *link) recieveResponse2(conn *net.UDPConn) {
	hasReceived = true

	// Takes time elapsed
	timeElapsed := time.Since(link1.timerStartTime)

	// Calculates Latency
	link1.latency = 0.9*link1.latency + 0.1*((timeElapsed/2).Seconds()*1000)

	probeCounter++

	link1.loss = lossCounter
	//fmt.Println("Loss", lossCounter, "Probes:", probeCounter, "EWMA: ", link1.latency)
	fmt.Println("Received R2")
	fmt.Println("Loss:", lossCounter, "Probes:", probeCounter)
	fmt.Println("-------------")
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

func (link1 link) CheckForTimeout() {
	// Store the index of the bitMask's index which this function should check for timeout

	for {
		if hasReceived == true {
			hasReceived = false
			link1.timeoutDeadline = time.Now().Add(5000 * time.Second)
			return
		} else if time.Now().Before(link1.timeoutDeadline) == false {
			// Update Loss
			fmt.Println("------------\nTIMEOUT\n-------------")
			lossCounter++
			hasReceived = false
			link1.timeoutDeadline = time.Now().Add(5000 * time.Second)
			continue
		}

		//fmt.Println("Checking Timeout")
		time.Sleep(200 * time.Millisecond)
	}

	return
}

func getData(conn *net.UDPConn) {

	buffer := make([]byte, 64)

	for {
		fmt.Println("Reading...")
		_, addr, err := conn.ReadFromUDP(buffer)
		CheckError(err)
		fmt.Println(string(buffer[:2]))

		if string(buffer[:2]) == "IP" {
			RoutingTable[addr.String()].sendResponse1(conn)
		} else if string(buffer[:2]) == "R1" {
			RoutingTable[addr.String()].sendResponse2(conn)
		} else if string(buffer[:2]) == "R2" {
			RoutingTable[addr.String()].recieveResponse2(conn)
		} else {
			fmt.Println("Packet recieved is an undefined type")
		}
		fmt.Println("Loss:", lossCounter, "Probes:", probeCounter)
		fmt.Println("-------------")
	}
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Err:", err)
		os.Exit(1)
	}
}
