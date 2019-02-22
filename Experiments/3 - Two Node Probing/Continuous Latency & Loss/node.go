package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

type link struct {
  addr string
  latency float64
  loss float64
  TTLexpiration time.Time
}

const (
	//PROBE_TIME = 3000 // 1 probe every 3s
	ROUTE_TTL = 100 //60seconds
	PROBE_TIMEOUT = 50 // 2 SECONDS
)

var (
	lossMask [100]bool
	maskIndex int
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Not Enough Arguements")
		os.Exit(1)
	}

  // Initialize Addresses
	nodeAddr, err := net.ResolveUDPAddr("udp", os.Args[2])
	CheckError(err)

	localNode, err := net.ResolveUDPAddr("udp", os.Args[1])
	CheckError(err)

	conn, err := net.ListenUDP("udp", localNode)
	CheckError(err)

  // Initialize Our Routing Table
  var RoutingTable [1]*link
  var link1 link
  link1.addr = nodeAddr.String()
  link1.TTLexpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)
  RoutingTable[0] = &link1

	// Initialize Loss Count
	maskIndex = 0
	for i:=0; i < 100; i++ {
		lossMask[i] = true
	}

  // Probe
  for {
    // If current time is before experation time do nothing, else probe
    if time.Now().Before(RoutingTable[0].TTLexpiration) == false {
      RoutingTable[0].SendProbe(conn)
      fmt.Print("Mask Index: ", maskIndex, " Loss: ", RoutingTable[0].loss, " EWMA: ", RoutingTable[0].latency , "\n")
    }
  }
}


func (link1 *link) SendProbe(conn *net.UDPConn) {

  // Update the link's expiration time to be ROUT_TTL milliseconds from now
  link1.TTLexpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)

  nodeAddr, err := net.ResolveUDPAddr("udp", link1.addr)
  CheckError(err)

  // Write initial probe
	buffer := make([]byte, 64)
	copy(buffer[:], "Hello")
	_, err = conn.WriteToUDP(buffer, nodeAddr)
	CheckError(err)

	probeTime := time.Now()
	deadline := probeTime.Add(time.Millisecond * PROBE_TIMEOUT)
	conn.SetReadDeadline(deadline)

  // Read 2nd Leg of Handshake & Record Latency
	_, _, err = conn.ReadFromUDP(buffer)
	if err != nil {
		//fmt.Println("Timeout:", err)
		lossMask[maskIndex] = false
	} else  {
		lossMask[maskIndex] = true
	}

	timeElapsed := time.Since(probeTime)

  // Write 3rd Leg
	_, err = conn.WriteToUDP(buffer, nodeAddr)
	CheckError(err)

	// Update Loss Values
	maskIndex = (maskIndex + 1) %  100 // Increment maskIndex every time you send a probe, loop back to 0 when you get to 100
	var lossCount float64 = 0

	for _, value := range lossMask { // Loop through lossMask and count the number of packets lost in the past 100
		if value == false {
			lossCount++
		}
	}

	link1.loss = lossCount

  // Update Link's latency value
  latencyInMilliseconds := timeElapsed.Seconds() * 1000 / 2

  if link1.latency == 0 {
    link1.latency = latencyInMilliseconds
  } else {
    link1.latency = 0.9 * link1.latency + 0.1 * (latencyInMilliseconds)
  }
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Err:", err)
		os.Exit(1)
	}
}
