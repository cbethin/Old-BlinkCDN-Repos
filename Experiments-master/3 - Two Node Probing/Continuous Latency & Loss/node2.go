package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const (
	//PROBE_FREQUENCY = 20 // 20 every minute
	// ROUTE_TTL = 60000 //60seconds
	PROBE_TIMEOUT = 100 // 2 SECONDS
)

var EWMA float64

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Not Enough Arguements")
		os.Exit(1)
	}

	localNode, err := net.ResolveUDPAddr("udp", os.Args[1])
	CheckError(err)

	conn, err := net.ListenUDP("udp", localNode)
	CheckError(err)
  defer conn.Close()
  EWMA = 0

  for {
	   RespondToProbe(conn)
  }
}

func RespondToProbe(conn *net.UDPConn) {

	buffer := make([]byte, 64)

	_, addr, err := conn.ReadFromUDP(buffer)
	CheckError(err)

	copy(buffer[:], "Hello")

	_, err = conn.WriteToUDP(buffer, addr)
	CheckError(err)

	probeTime := time.Now()

	_, _, err = conn.ReadFromUDP(buffer)
	CheckError(err)

	timeElapsed := time.Since(probeTime)
	fmt.Print("\nLatency: ", timeElapsed/2)
  if EWMA == 0 {
    EWMA = (timeElapsed.Seconds() * 1000 / 2)
  } else {
    EWMA = 0.9*EWMA + 0.1*(timeElapsed.Seconds() * 1000 / 2)
  }
  fmt.Print(" EWMA: ", EWMA)

}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Err:", err)
		os.Exit(1)
	}
}
