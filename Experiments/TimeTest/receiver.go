/*
  Receives incoming packets from central nodes.
  Does nothing with them
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

  localAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
  swamiLib.CheckError(err)

  conn, err := net.ListenUDP("udp", localAddr)
  swamiLib.CheckError(err)

  buffer := make([]byte, 1024)
  fmt.Println("Listening for packets..")

  for {
    _, _, err := conn.ReadFromUDP(buffer)
    swamiLib.CheckError(err)
    fmt.Print(string(buffer[0]), " ")
  }

  fmt.Println("")

}
