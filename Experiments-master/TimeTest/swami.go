/*
- Listens on UDP port and print data received
from nodes in text file

- Data Received will be an array like below
swamiData := [nodeAddr, Time1, Time2]

*/
package main

import (
	"fmt"
	"net"
	"os"
	// Import swamiLib Library
	"./swamiLib"
)

func main() {

	// Check for Correct Number of arguments
	if len(os.Args) != 2 {
		fmt.Println("Wrong Number of Arguments")
		os.Exit(1)
	}

	// Get Server IP Address
	serverAddrString := ":" + os.Args[1]
	serverAddr, err := net.ResolveUDPAddr("udp", serverAddrString)
	swamiLib.CheckError(err)

	// Establish Listening Connection on Server IP
	conn, err := net.ListenUDP("udp", serverAddr)
	swamiLib.CheckError(err)
	defer conn.Close()
	fmt.Println("Listening for packets...")

	// Create the Buffer
	buffer := make([]byte, 1024)

	// Create the Data File
	file, err := os.Create("Swami.txt")
	swamiLib.CheckError(err)
	defer file.Close()

	for {
		// Listen for incoming messages & write to local buffer
		_, _, err := conn.ReadFromUDP(buffer)
		swamiLib.CheckError(err)

		//Insert conversion (convert buffer to string)
		rcvdString := string(buffer[:200]) + "\n"

		// Write buffer to text file and print in console
		file.WriteString(rcvdString)
		fmt.Println("Wrote:", string(buffer[:100]))
	}

	fmt.Println("Main Function done running")
}
