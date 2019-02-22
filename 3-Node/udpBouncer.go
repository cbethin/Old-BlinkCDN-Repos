package main

import (
	//"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	play = "play"
)

func main() {

	// Setup Sending/Receiving Addresses
  if len(os.Args) != 3 {
    fmt.Println("Arguemtn Error")
    os.Exit(0)
  }
  sendingHostName := os.Args[1] + ":6000"
  receivingHostName := os.Args[2] + ":6000"
  selfHostName := ":6000"

	selfUdpAddr, err := net.ResolveUDPAddr("udp4", selfHostName)

	if err != nil {
		log.Fatal(err)
	}

	// setup listener for incoming UDP connection
	listener, err := net.ListenUDP("udp", selfUdpAddr)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("UDP server up and listening on port 6000")

	defer listener.Close()

	for {
		// wait for UDP client to connect
		waitForPlay(listener, sendingHostName, receivingHostName)
	}

}

// This function takes a UDP listener and waits for the "play" Command
// to be sent from the Client.
// Once the message is received, it forwards the message to the Server
// it then begins to wait for messages to come through
func waitForPlay(conn *net.UDPConn, sendingHostName string, receivingHostName string) {

	// here is where you want to do stuff like read or write to client
	defer conn.Close()
	buffer := make([]byte, 1024)

	n, addr, err := conn.ReadFromUDP(buffer)

	fmt.Println("UDP client : ", addr)
	fmt.Println("Received from UDP client :  ", string(buffer[:n]))

	if err != nil {
		log.Fatal(err)
	}

	// NOTE : Need to specify client address in WriteToUDP() function
	//        otherwise, you will get this error message
	//        write udp : write: destination address required if you use Write() function instead of WriteToUDP()
  sendingHostAddr, err := net.ResolveUDPAddr("udp", sendingHostName)
  checkError(err)


	if err != nil {
		log.Println(err)
	}
	s := string(buffer[0:n])
	fmt.Println(s)
	if s[0:4] == play {
		//open file
		  _, err = conn.WriteToUDP([]byte("play"), sendingHostAddr)
  }

  waitForMessage(conn, sendingHostName, receivingHostName)
}


// This function takes a listener and waits for information to begin getting sent
// When information of any kind is received by the listener it bounces it
// to the Client..
func waitForMessage(conn *net.UDPConn, sendingHostName string, receivingHostName string) {
  for {
    buffer := make([]byte, 1024)
    n, addr, err := conn.ReadFromUDP(buffer)
    checkError(err)

		sendingHostAddr, err := net.ResolveUDPAddr("udp", sendingHostName)
		checkError(err)

		receivingHostAddr, err := net.ResolveUDPAddr("udp", receivingHostName)
    checkError(err)
    //fmt.Println("Sending to:", receivingHostName)

		fmt.Println(addr.String(), sendingHostAddr.String())
		if (addr.String() == sendingHostAddr.String()) {
	    if string(buffer[:n]) != "done" {
	      // convert data
	      ///fmt.Println(string(buffer[0:15]))
	      _, err = conn.WriteToUDP(buffer[:n], receivingHostAddr)

	    } else {
	      _, err = conn.WriteToUDP([]byte("done"), receivingHostAddr)
	      fmt.Println("Finished File Transfer.")
	      os.Exit(1)
	    }
		}
  }
}

func checkError(err error) {
  if err != nil {
    fmt.Println("ERROR:", err)
    os.Exit(2)
  }
}
