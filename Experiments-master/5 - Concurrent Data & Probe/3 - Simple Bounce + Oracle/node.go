package main

import (
	"fmt"
	"net"
	"os"
	"time"
	// Blink library import needs "./"
	"./blink"
	"strconv"
)

const (
	//  PROBE_TIME = 3000 // Probe 20 times every minute
	ROUTE_TTL     = 15000 // 15 seconds
	PROBE_TIMEOUT = 10000 // 10 seconds
)

var (
	RoutingTable map[string]*link
	nodeAddr     *net.UDPAddr
	localNode    *net.UDPAddr
	// nodeAddr3 *net.UDPAddr
	oracleAddr 		*net.UDPAddr
)

type link struct {
	addr            string
	latency         float64
	loss            int
	TTLExpiration   time.Time // When was the last time I probed you (Creepy!)
	timeoutDeadline time.Time
	timerStartTime  time.Time
	lossBitMask     [100]bool
	lossIndex       int
}

func main() {

	if len(os.Args) != 4 {
		fmt.Println("Incorrect Amount of Inputs")
		os.Exit(1)
	}
	defer fmt.Println(time.Now())

	// Set Up Addresses to use
	localNodeAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
	blink.CheckError(err)

	localNode = localNodeAddr

	nodeAddr, err = net.ResolveUDPAddr("udp", os.Args[2])
	blink.CheckError(err)

	oracleAddr, err = net.ResolveUDPAddr("udp", os.Args[3])
	blink.CheckError(err)

	// nodeAddr2, err = net.ResolveUDPAddr("udp", os.Args[3])
	// blink.CheckError(err)

	// nodeAddr3, err = net.ResolveUDPAddr("udp", os.Args[4])
	// blink.CheckError(err)

	conn, err := net.ListenUDP("udp", localNode)
	blink.CheckError(err)

	RoutingTable = make(map[string]*link)

	// Setup Routing Table and It's Links
	var link1 link
	link1.addr = nodeAddr.String()
	link1.latency = 0
	link1.loss = 0
	link1.lossIndex = 0
	link1.TTLExpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)
	RoutingTable[nodeAddr.String()] = &link1

	// var link2 link
	// link2.addr = nodeAddr2.String()
	// link2.latency = 0
	// link2.loss = 0
	// link2.lossIndex = 0
	// link2.TTLExpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)
	// RoutingTable[nodeAddr2.String()] = &link2

	// var link3 link
	// link3.addr = nodeAddr3.String()
	// link3.latency = 0
	// link3.loss = 0
	// link3.lossIndex = 0
	// link3.TTLExpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)
	// RoutingTable[nodeAddr3.String()] = &link3

	/* Go Check Expiration of Routing Table and Listen for Incoming Data
	 	 This is the actual functionality of the nodes, with 2 go routines Checking
	   to see if they should initalize a probe to another node (1 node per routine)
	   the third function being the active listener, whihc will determine what to do
	   with or in response to the incoming packets*/

	go RoutingTable[nodeAddr.String()].checkExpiration(conn)
	// go RoutingTable[nodeAddr2.String()].checkExpiration(conn)
	// go RoutingTable[nodeAddr3.String()].checkExpiration(conn)
	getData(conn)
	time.Sleep(30 * time.Second)
}




/* Sends the initial probe packet to another node via a link, which is
   passed into the function.*/

func (link1 *link) sendInitialProbe(conn *net.UDPConn) {
	// Update the link expiration as the first thing
	link1.TTLExpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)

	nodeAddr, err := net.ResolveUDPAddr("udp", link1.addr)
	blink.CheckError(err)

	// Create a Buffer now
	buffer := blink.MakeBlinkPacket(localNode, nodeAddr, blink.InitialProbe, []byte(""))

	_, err = conn.WriteToUDP(buffer, nodeAddr)
	blink.CheckError(err)

	link1.timerStartTime = time.Now()
	link1.timeoutDeadline = link1.timerStartTime.Add(time.Millisecond * PROBE_TIMEOUT)
}




/* Sends the Response 1 packet to another node via a link, which is
   passed into the function. */

func (link1 *link) sendResponse1(conn *net.UDPConn) {
	link1.TTLExpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)

	// Create a Buffer now
	buffer := blink.MakeBlinkPacket(localNode, nodeAddr, blink.ResponseOne, []byte(""))
	nodeAddr, err := net.ResolveUDPAddr("udp", link1.addr)
	blink.CheckError(err)

	_, err = conn.WriteToUDP(buffer, nodeAddr)
	blink.CheckError(err)

	link1.timerStartTime = time.Now()
	link1.timeoutDeadline = link1.timerStartTime.Add(time.Millisecond * PROBE_TIMEOUT)
}




/* Sends the Response 2 packet to another node via a link, which is
   passed into the function. Will also calculate and store the calculated
	 latency of the link. Sends routing table to oracle when it's done calculating
*/

func (link1 *link) sendResponse2(conn *net.UDPConn) {

	// Calculate time elapsed since the probe was sent
	timeElapsed := time.Since(link1.timerStartTime)

	// Resolve other node's address, create a Blink Response 2 packet
	// and send that buffer over the inputted connection object
	nodeAddr, err := net.ResolveUDPAddr("udp", link1.addr)
	blink.CheckError(err)

	buffer := blink.MakeBlinkPacket(localNode, nodeAddr, blink.ResponseTwo, []byte(""))

	_, err = conn.WriteToUDP(buffer, nodeAddr)
	blink.CheckError(err)

	// Calculates Latency
	link1.latency = 0.9*link1.latency + 0.1*((timeElapsed/2).Seconds()*1000)

	// Send Routing Table To Oracle
	for _, value := range RoutingTable {
		// Create a blink packet for
		buffer := blink.MakeBlinkPacket(localNode, oracleAddr, blink.RoutingTable, LinkToByteArray(*value))

		_, err := conn.WriteToUDP(buffer, oracleAddr)
		blink.CheckError(err)
	}

}




/* Calculates the latency after a ResponseTwo is received and stores the latency in
	 the corresponding link. Sends routing table to oracle when it's done calculating
*/
func (link1 *link) receiveResponse2(conn *net.UDPConn) {

	// Takes time elapsed since probe was sent
	timeElapsed := time.Since(link1.timerStartTime)

	// Calculates Latency
	link1.latency = 0.9*link1.latency + 0.1*((timeElapsed/2).Seconds()*1000)

	// Send Routing Table To Oracle
	for _, value := range RoutingTable {
		// Create a blink packet to send the link to the oracle.
		buffer := blink.MakeBlinkPacket(localNode, oracleAddr, blink.RoutingTable, LinkToByteArray(*value))
		// Send the new Blink Packet to the oracle.
		_, err := conn.WriteToUDP(buffer, oracleAddr)
		blink.CheckError(err)
	}
}




/* Every 100 milliseconds the program will check the TTLExpiration of the designated
   link. If the link has expired, the program will initialize the probing process by
	 sending an initial probe to the corresponding node */

func (link1 *link) checkExpiration(conn *net.UDPConn) {

	for {
		if time.Now().Before(link1.TTLExpiration) == false {
			// If your Current time is before the expiation time
			// You do Nothing!
			// FOR FUTURE : The link that requires the probing funtion
			//              i.e. the one that times out, can update the
			//              linkNumber and call the SendProbe
			//              RoutingTable[linkNumber].SendProbe(nodeAddr,conn)
			link1.sendInitialProbe(conn)
		}
		time.Sleep(100 * time.Millisecond)
	}
}




/* Read data flowing in to the listening connection that is passed into this function
   The function then checks the type of packet being sent and responds accordingly.
*/

func getData(conn *net.UDPConn) {

	buffer := make([]byte, 1024)

	for {
		// Read from connection, then extract packet type using Blink's built in function
		_, addr, err := conn.ReadFromUDP(buffer)
		blink.CheckError(err)

		packetType := blink.ExtractPacketType(buffer)
		fmt.Println(packetType, "from", addr.String())

		// Chech the packet type, and call the proper function for each packet type
		switch packetType {
		case blink.InitialProbe:
			go RoutingTable[addr.String()].sendResponse1(conn)
		case blink.ResponseOne:
			go RoutingTable[addr.String()].sendResponse2(conn)
		case blink.ResponseTwo:
			go RoutingTable[addr.String()].receiveResponse2(conn)
		case blink.Iframe:
			fmt.Println("Received I Frame")
			go Bounce(conn, buffer)
		case blink.Pframe:
			fmt.Println("Received an P Frame")
			go Bounce(conn, buffer)
		case blink.Bframe:
			fmt.Println("Received an B Frame")
			go Bounce(conn, buffer)
		default:
			fmt.Println("Packet recieved is an undefined type")
		}
		fmt.Println("------------------")
		fmt.Println("Link 1 -- Lat:", RoutingTable[nodeAddr.String()].latency)
		// fmt.Println("Link 2 -- Lat:", RoutingTable[nodeAddr2.String()].latency)
	}
}




/* This function will forward an inputted buffer to the desired destination address.
	 if the current node is the desired destination address, the system will not Bounce
	 the blink packet
*/

func Bounce(conn *net.UDPConn, buf []byte) {
	// Extract the destination address from the Blink packet
	destAddr := blink.ExtractFinalDestAddr(buf)

	/* Make sure the destination address is not the current node. If the dest addr
	   is the current node, simply print out the message. If not, forward the message
	   to the desired address. */
	if destAddr.String() != localNode.String() {
		_, err := conn.WriteToUDP(buf, destAddr)
		blink.CheckError(err)

		fmt.Println("-------------------------")
		fmt.Println("Bouncing", string(buf[66:75]))
	} else {
		fmt.Println("-------------------------")
		fmt.Println("Receiving", string(buf[66:75]))
	}
}

func LinkToByteArray(link link) []byte {
  byteArray := make([]byte, 192)
  copy(byteArray[0:16], strconv.FormatFloat(link.latency, 'f', 15, 64))
  copy(byteArray[16:48], link.addr + "!")
  fmt.Println(string(byteArray))
  return byteArray
}
