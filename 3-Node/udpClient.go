package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	//"bufio"
	//"io"
	"os"
	//"strings"
	//"strconv"
	"math/rand"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Not enough args")
	}
	hostName := os.Args[1]
	portNum := "6000"

	service := hostName + ":" + portNum

	RemoteAddr, err := net.ResolveUDPAddr("udp", service)

	//LocalAddr := nil
	// see https://golang.org/pkg/net/#DialUDP

	conn, err := net.DialUDP("udp", nil, RemoteAddr)
	defer conn.Close()

	myUDP, err := net.ResolveUDPAddr("udp", ":6000")
	checkError(err)

	ln, err := net.ListenUDP("udp", myUDP)
	checkError(err)
	// note : you can use net.ResolveUDPAddr for LocalAddr as well
	//        for this tutorial simplicity sake, we will just use nil

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Established connection to %s \n", service)
	log.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
	log.Printf("Local UDP client address : %s \n", conn.LocalAddr().String())

	_, err = conn.Write([]byte("play"))
	checkError(err)
	fmt.Println("----------\nPlaying...")

	// Create the input buffer
	buffer := make([]byte, 1024)
	//Create the File
	file, err := os.Create("ourfile.txt")
	checkError(err)

	var SNC float64 = 0 // Sequence Number Counter

	//Receive until "done" is received
	for {
		//fmt.Println("running")
		n, _, err := ln.ReadFromUDP(buffer)
		checkError(err)

		if string(buffer[:n]) != "done" {
			// convert data
			toWrite := string(buffer[:n]) + "\n"
			file.WriteString(toWrite)
			fmt.Println(string(buffer[0:15]), "added.")
			scanList := strings.Fields(toWrite) // Parses String toWrite into a list of 3 strings
			sequenceNumber, _ := strconv.ParseFloat(scanList[2], 64)

			if SNC != sequenceNumber {
				SNC += 1
				switchConnection(&RemoteAddr, conn, sequenceNumber)
			} else {
				SNC = sequenceNumber
			}
		} else {
			fmt.Println("Finished File Transfer.")
			os.Exit(1)
		}
	}

	file.Close()

}

func switchConnection(addr **net.UDPAddr, currentConn net.Conn, seqNum float64) net.Conn {

	_, err := currentConn.Write([]byte("stop")) // Sends Stop Signa'
	checkError(err)

	currentConn.Close() // Closes Connection

	addrS2 := selectConnection()

	RemoteAddr, err := net.ResolveUDPAddr("udp", addrS2)

	//LocalAddr := nil
	// see https://golang.org/pkg/net/#DialUDP

	conn, err := net.DialUDP("udp", nil, RemoteAddr)
	defer conn.Close()

	// myUDP, err := net.ResolveUDPAddr("udp", ":6000")
	// checkError(err)
	//
	// ln, err := net.ListenUDP("udp", myUDP)
	// checkError(err)
	// // note : you can use net.ResolveUDPAddr for LocalAddr as well
	// //        for this tutorial simplicity sake, we will just use nil

	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write([]byte("play"))
	checkError(err)
	fmt.Println("----------\nPlaying...")

	return conn
}

func selectConnection() string {
	// Contains a slice of server connections avaliable
	// Generates a random number within the range of the slice
	// choose that number within the slice

	/* EX
	   slice1 = [1,2,3,4,5]
	   randNum = rand.generate(range of slice1)
	   connection = slice1[randNum]
	*/

	connectionSlice := []string{"155.246.66.150"}
	rangeOfSlice := len(connectionSlice)
	randNum := rand.Intn(rangeOfSlice)
	connection := connectionSlice[randNum]
	return connection
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
