/*
Sends 100 packets to next node
Each packet is sent at a 50 Millisecond Increment
*/

package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"./swamiLib"
)

func main() {

	// Check for Proper Number of Arguments
	if len(os.Args) != 3 {
		fmt.Println("Improper number of Arguments")
		os.Exit(0)
	}

	// Setup Server/Client Address Connections
	clientAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
	swamiLib.CheckError(err)

	nextNodeAddr, err := net.ResolveUDPAddr("udp", os.Args[2])
	swamiLib.CheckError(err)

	// Set up Connection Objects
	conn, err := net.ListenUDP("udp", clientAddr)
	swamiLib.CheckError(err)
	defer conn.Close()

	// Make Buffer
	buffer := make([]byte, 1024)

	// Initialize Packet Counter
	packetNumber := 1

	// // For loop that sends 100 Packets then stops
	// for i := packetNumber + 1; i <= 100; i++ {
	// 	// Create Message
	// 	message := strconv.Itoa(i)
	//
	// 	// Copy Message to buffer
	// 	copy(buffer[:], message)
	//
	// 	// Send buffer to next node
	// 	_, err := conn.WriteToUDP(buffer, nextNodeAddr)
	// 	swamiLib.CheckError(err)
	// 	fmt.Print(message, " ")
	//
	// 	// Sleep for 25 Milliseconds
	// 	time.Sleep(time.Millisecond * 2000)
	// }

	// For loop that sends a packet every 2 seconds
	for {
		// Create message
		message := strconv.Itoa(packetNumber)

		// Copy Message to buffer
		copy(buffer[:], message)

		// Send buffer to next node
		_, err := conn.WriteToUDP(buffer, nextNodeAddr)
		swamiLib.CheckError(err)
		fmt.Print(message, " ")

		// Sleep for 25 Milliseconds
		time.Sleep(time.Millisecond * 2000)

		packetNumber += 1
	}

	fmt.Println("")
}
