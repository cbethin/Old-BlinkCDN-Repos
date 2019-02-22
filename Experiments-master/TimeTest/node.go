/*
Bounces data to client

While doing this, it records the time it receives
the data (Time1) and records the time it sends the
data (Time2)

Sends this data to Swami server in the following
formated array

swamiData := [nodeAddr, Time1, Time2]

*/
package main

import (
	"fmt"
	"net"
	"os"
	"time"
	// Import swamiLib Library
	"./swamiLib"
)

func main() {

	if len(os.Args) != 4 {
		fmt.Println()
	}

	/*
	  Resolve UDP addresses for the the local node,
	  next node and swami node
	*/
	localNodeAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
	swamiLib.CheckError(err)

	nextNodeAddr, err := net.ResolveUDPAddr("udp", os.Args[2])
	swamiLib.CheckError(err)

	swamiAddr, err := net.ResolveUDPAddr("udp", os.Args[3])
	swamiLib.CheckError(err)

	// Set up a listening connection on the local node
	conn, err := net.ListenUDP("udp", localNodeAddr)
	swamiLib.CheckError(err)

	// Create the Buffer
	buffer := make([]byte, 1024)
	fmt.Println("Listening for packets")
	rcvdCount := 0

	for {
		// Read buffer from the listening port
		_, _, err := conn.ReadFromUDP(buffer)
		swamiLib.CheckError(err)
		rcvdCount++
		fmt.Println(rcvdCount)

		// Record time data is recieved
		time1 := time.Now().String()

		// Send buffer over to the next node
		_, err = conn.WriteToUDP(buffer, nextNodeAddr)
		swamiLib.CheckError(err)

		// Record time data is sent
		time2 := time.Now().String()

		// Create Data array to send to swami

		//*** Need to convert all things to string and then to byte array
		localNodeAddrString := localNodeAddr.String()
		swamiData := [3]string{localNodeAddrString, time1, time2}
		swamiDataString := swamiData[0] + "," + swamiData[1] + "," + swamiData[2] + ","

		swamiPacket := make([]byte, 1024)
		copy(swamiPacket[:1024], []byte(swamiDataString))

		// Send swami data to swami server
		_, err = conn.WriteToUDP(swamiPacket, swamiAddr)
		swamiLib.CheckError(err)

	}

}
