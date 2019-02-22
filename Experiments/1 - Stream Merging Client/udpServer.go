/*
	FUNCITONALITY: This program acts as a server designed to send a given number of
	packets to a client. It works by listening for a "play" message from a client,
	which is to be sent in the format "playX" where X is a 3 digit number that
	identifies the message # with which the server should start sending messages.
	The server will send one message for every number between X and NUM_PACKETS
	(defined below).

	USAGE: The program is to be called with 1 argument, the argument being the port
	number the user wishes to connect the server to.

	NOTES: The address of the client can be obtained from the ReadFromUDP
	command every time the server reads an incoming message. Thus, the server
	can respond to any client that wishes to connect to it. For the system to
	work, two instances of the server must be set up for the client to connect
	to, each server having a different address.

	VARIABLES:
		serverAddrString (string) : UDP address of the server using the port given by user

		serverAddr (*net.UDPAddr) : address of server resolved into type *net.UDPAddr

		conn (net.Conn): object that allows any incoming and outgoing connections
				to/from the server.

		err (error):  variable that is used to hold any errors thrown by various
				functions used throughout the program.

		buffer ([]byte) : all messages must be passed to the buffer before writing to
				a connection or reading from a connection, as the connection object will
				only handle byte streams.
*/

package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"strconv"
)

const (
	NUM_PACKETS = 999
)

func main() {

	// Make sure 2 arguments are inputted to the program.
	// The single argument used will be the server's port.
	if len(os.Args) != 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	// Obtain the port from the program argument and resolve into
	// type UDP address. Will be the address the server listens to.
	serverAddrString := ":" + os.Args[1]
	serverAddr, err := net.ResolveUDPAddr("udp", serverAddrString)
	checkError(err)

	// Setup our connection object as a listener. Make sure to
	// close the connection before the program exits
	conn, err := net.ListenUDP("udp", serverAddr)
	checkError(err)
	defer conn.Close()
	fmt.Println("Listening for message...")

	// Create the buffer of 1024 bytes. This will be initialized to a slice
	// of 1024 bytes, each containing 0.
	buffer := make([]byte, 1024)

	for {

		// Listen for incoming messages. When a play message is received, check to see
		// If the message is the play message. The client's address will be pulled from
		// each message and stored in the addr variable
		_, addr, err := conn.ReadFromUDP(buffer)
		checkError(err)

		if string(buffer[0:4]) == "play" {

			// If the play message is received, pull the frame number and convert to an integer.
			fmt.Println("------------")
		  fmt.Println("Received:", string(buffer[:10]), "\nSize:", len(buffer), "bytes | From:", addr.String())

			frameNumber, err := strconv.Atoi(string(buffer[4:7]))
			checkError(err)

			buffer = make([]byte, 1024)
			fmt.Println("Sending:")

			// Starting on the desired frame number, begin writing messages to the client
			// containing the frame number until frame 700 is reached. Wait for
			// 5 milliseconds between each message transmission.
			for i := frameNumber+1; i <= NUM_PACKETS; i++ {
				message := strconv.Itoa(i)
				fmt.Print(message, " ")
			  copy(buffer[:], message)
				_, err = conn.WriteToUDP(buffer, addr)
				checkError(err)
				time.Sleep(time.Microsecond * 33333)
			}

			// When all the messages have been transmitted, send a "done" messages
			// to the client and exit the program.
			copy(buffer[:], "done")
			_, err = conn.WriteToUDP(buffer, addr)
			os.Exit(1)

		} else if string(buffer[0:4]) == "done" {
			// if the server receives a "done" message from the client, close
			// the connection and exit the program. Your work is done.
			conn.Close()
			os.Exit(2)
		}
	}
}

// If an error is found print the error and exit program
func checkError(err error) {
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}
