package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

// Setup Link Structure
// Each link will be a row in the routing table
type link struct {
  addr string
  latency float64
  loss float64
  TTLexpiration time.Time
}

const (
	//PROBE_TIME = 3000 // 1 probe every 3s
	ROUTE_TTL = 2000 //60seconds
	PROBE_TIMEOUT = 1000 // 2 SECONDS
)

func main() {

	// Must input the local node address and address of other node
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

  // Probe
  for {
    // If current time is before experation time do nothing, else probe
    if time.Now().Before(RoutingTable[0].TTLexpiration) == false {
      RoutingTable[0].SendProbe(conn)
      fmt.Print(" EWMA: ", RoutingTable[0].latency , "\n")
    }
  }
}


func (link1 *link) SendProbe(conn *net.UDPConn) {

  // Update the link's expiration time to be ROUT_TTL milliseconds from now
  link1.TTLexpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)

	// Resolve link's node address
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
		fmt.Println("Timeout:", err)
	}

	timeElapsed := time.Since(probeTime) // Calculate time elapsed since send

  // Write 3rd Leg
	_, err = conn.WriteToUDP(buffer, nodeAddr)
	CheckError(err)

  // Update Link's latency value
  latencyInMilliseconds := timeElapsed.Seconds() * 1000 / 2
  fmt.Print("Delay: ", timeElapsed/2)

	// Update link's latency
  if link1.latency <= 0 {
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
