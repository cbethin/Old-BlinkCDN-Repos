package main

import (
	"fmt"
	"net"
	"os"
	"time"
	// Blink library import needs "./"
	"./blink"
)

var (
	RoutingTable map[string]*blink.Link
	localNode     *net.UDPAddr
	oracleAddr 		*net.UDPAddr
	nextNode			*blink.Link
)

func main() {

	if len(os.Args) != 6 {
		fmt.Println("Incorrect Amount of Inputs")
		os.Exit(1)
	}
	defer fmt.Println(time.Now())

	// Set Up Addresses to use
	localNodeAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
	blink.CheckError(err)

	nodeAddr1, err := net.ResolveUDPAddr("udp", os.Args[2])
	blink.CheckError(err)

	nodeAddr2, err := net.ResolveUDPAddr("udp", os.Args[3])
	blink.CheckError(err)

	nodeAddr3, err := net.ResolveUDPAddr("udp", os.Args[4])
	blink.CheckError(err)

	oracleAddr, err = net.ResolveUDPAddr("udp", os.Args[5])
	blink.CheckError(err)

	// Setup listening connection
	// conn, err := net.ListenUDP("udp", localNode)
	// blink.CheckError(err)

	// Setup Routing Table and It's Links
	RoutingTable = make(map[string]*blink.Link)

	var link1 blink.Link
	link1.Addr = nodeAddr1.String()
	link1.Latency = 0
	link1.TTLExpiration = time.Now().Add(blink.ROUTE_TTL * time.Millisecond)
	RoutingTable[nodeAddr1.String()] = &link1

	var link2 blink.Link
	link2.Addr = nodeAddr2.String()
	link2.Latency = 0
	link2.TTLExpiration = time.Now().Add(blink.ROUTE_TTL * time.Millisecond)
	RoutingTable[nodeAddr2.String()] = &link2

	var link3 blink.Link
	link3.Addr = nodeAddr3.String()
	link3.Latency = 0
	link3.TTLExpiration = time.Now().Add(blink.ROUTE_TTL * time.Millisecond)
	RoutingTable[nodeAddr3.String()] = &link3

	// Initialize Decision Table
	DecisionTable := make(map[string][3]string)
	DecisionTable[""] = [3]string{"k", "l", "m"}
	blink.SetDecisionTable(DecisionTable)
	fmt.Println(blink.DecisionTable)
	// Let Blink Know What Variables are What
	blink.SetRoutingTable(RoutingTable)
	blink.SetLocalAddr(localNodeAddr)
	blink.SetOracleAddr(oracleAddr)

	blink.SetupListener()

	/* Go Check Expiration of Routing Table and Listen for Incoming Data
	 	 This is the actual functionality of the nodes, with 2 go routines Checking
	   to see if they should initalize a probe to another node (1 node per routine)
	   the third function being the active listener, whihc will determine what to do
	   with or in response to the incoming packets*/

	//StartProber(conn)
	blink.StartProber()
	blink.GetData()

	time.Sleep(30 * time.Second)

}
