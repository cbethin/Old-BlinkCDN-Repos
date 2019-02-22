/*
Sends 100 packets to next node
Each packet is sent at a 50 Millisecond Increment
*/

package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"./swamiLib"
)

func main() {

	// Check for Proper Number of Arguments
	if len(os.Args) != 4 {
		fmt.Println("Improper number of Arguments")
		os.Exit(0)
	}

	// Setup Server/Client Address Connections
	clientAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
	swamiLib.CheckError(err)

	nextNodeAddr, err := net.ResolveUDPAddr("udp", os.Args[2])
	swamiLib.CheckError(err)

  swamiAddr, err := net.ResolveUDPAddr("udp", os.Args[3])
  swamiLib.CheckError(err)

	// Set up Connection Objects
	conn, err := net.ListenUDP("udp", clientAddr)
	swamiLib.CheckError(err)
	defer conn.Close()

	// Make Buffer
	buffer := make([]byte, 1024)
  // Copy Message to buffer
  copy(buffer[:], "ok")

	// Initialize Packet Counter
	packetNumber := 1

	// For loop that sends a packet every 2 seconds
	for {

		// Send buffer to next node
		_, err := conn.WriteToUDP(buffer, nextNodeAddr)

    // Record time data is sent
		time1 := time.Now().String()

		swamiLib.CheckError(err)
		fmt.Print(packetNumber, " ")

    // Record time 3
		_, _, err = conn.ReadFromUDP(buffer)
		time3 := time.Now().String()

		//*** Need to convert all things to string and then to byte array
		localNodeAddrString := clientAddr.String()
		swamiData := [3]string{localNodeAddrString, time1, time3}
		swamiDataString := swamiData[0] + "," + swamiData[1] + "," + swamiData[2] + ","

		swamiPacket := make([]byte, 1024)
		copy(swamiPacket[:1024], []byte(swamiDataString))

		// Send swami data to swami server
		_, err = conn.WriteToUDP(swamiPacket, swamiAddr)
		swamiLib.CheckError(err)

		packetNumber++
	}


	fmt.Println("Done.")
}
