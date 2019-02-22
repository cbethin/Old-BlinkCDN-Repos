package blink

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

// SetupListener : Sets up the listening connection to be used on the node. Sets the connection object as the global variable conn.
func SetupListener() {
	var err error
	conn, err = net.ListenUDP("udp", localAddr)
	CheckError(err)
	fmt.Println("Setup Listner on", localAddr.String())
}

// StartNode : Starts a node instance and it's corresponding probing and data retreive processes.
func StartNode() {
	StartProber()
	GetData()
}

// StartProber : starts the probing functionality of the node. It must be started after the RoutingTable has been initialized and must be restarted every time the Routing Table is updated.
func StartProber() {
	for _, value := range RoutingTable {
		go CheckExpiration(conn, value)
	}
}

// GetData : Reads data flowing in to the listening connection that is passed into this function.  The function then checks the type of packet being sent and responds accordingly.
func GetData() {

	// Create File
	fileName = strconv.Itoa(localAddr.Port) + ".txt"
	file, err := os.Create(fileName)
	CheckError(err)
	defer file.Close()

	// GET DATA
	buffer := make([]byte, 1024)

	for {
		// Read from connection, then extract packet type using Blink's built in function
		_, addr, err := conn.ReadFromUDP(buffer)
		CheckError(err)

		packetType := ExtractPacketType(buffer)

		fmt.Println("---------------------")
		// Chech the packet type, and call the proper function for each packet type
		switch packetType {
		case InitialProbe:
			fmt.Println("Received IP, Sending R1")
			SendResponse1(RoutingTable[addr.String()])
		case ResponseOne:
			fmt.Println("Received R1, Sending R2")
			SendResponse2(RoutingTable[addr.String()])
		case ResponseTwo:
			fmt.Println("Received R2")
			ReceiveResponse2(RoutingTable[addr.String()])
		case Iframe:
			fmt.Println("Handling I Frame")
			Bounce(buffer)
		case Pframe:
			fmt.Println("Handling P Frame")
			Bounce(buffer)
		case Bframe:
			fmt.Println("Handling B Frame")
			Bounce(buffer)
		case DecisionTableType:
			fmt.Println("Handling Decision Table")
			UpdateDecisionTable(buffer)
		default:
			fmt.Println("Packet recieved is an undefined type")
		}
		fmt.Println("---------------------")
	}
}

// Bounce : This function will forward an inputted buffer to the desired destination address. If the current node is the desired destination address, the system will not Bounce the blink packet
func Bounce(buf []byte) {

	packetType := ExtractPacketType(buf)
	SID := ExtractSID(buf)
	packetNumber := ByteArrayToJSONToPacket(ExtractPacketData(buf)).Number
	destAddr := ExtractFinalDestAddr(buf)
	receivedTime := time.Now()

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	CheckError(err)
	printString := SID + " " + packetType + " " + strconv.Itoa(packetNumber) + " " + time.Now().String() + "\n"
	_, err = file.WriteString(printString)
	CheckError(err)
	file.Close()

	fmt.Println("Bouncing", packetNumber)

	// Checks to see if SID is present in the decision table. If it is
	// the funciton will proceed to find the next node to bounce to and
	// bounce the data to that node. If the SID is not found, nothing will be
	// done with the packet.

	nextNodeAddrString := ""
	hopNumber := -1
	if _, ok := DecisionTable[SID]; ok {
		// Is the current node the destination address?
		if destAddr.String() != LocalAddrString {
			// If thisNode != destNode then determine next node and foward data
			// to next node
			for i, value := range DecisionTable[SID] {
				if value == LocalAddrString {
					hopNumber = i
					if i == len(DecisionTable[SID])-1 {
						fmt.Println("Receiving", packetNumber)
						trackPacketDelay(packetNumber, hopNumber, SID)
						return
						// destAddr := ExtractFinalDestAddr(buf)
						// nextNodeAddrString = destAddr.String()
					}

					nextNodeAddrString = DecisionTable[SID][i+1]
					// fmt.Println("Next Node:", nextNodeAddrString)
				}
			}

			// If current node was not found, break off. Bounce was incorrect
			if nextNodeAddrString == "" {
				fmt.Println("Next node not found")
				return
			}

			// Resolve the NextNode's address and bounce it off.
			nextNodeAddr, err := net.ResolveUDPAddr("udp", nextNodeAddrString)
			CheckError(err)
			fmt.Println("Actual Next Node:", nextNodeAddr.String())
			trackPacketDelay(packetNumber, hopNumber, SID)

			_, err = conn.WriteToUDP(buf, nextNodeAddr)
			CheckError(err)

			// Send time update to oracle
			message := packetType + strconv.Itoa(packetNumber) + " received-at: " + LocalAddrString + " at: " + receivedTime.String()
			blinkPacket := MakeBlinkPacket(ProbeID, localAddr, oracleAddr, BounceUpdate, []byte(message))

			_, err = conn.WriteToUDP(blinkPacket, oracleAddr)
			CheckError(err)
		} else {
			fmt.Println("Receiving", packetNumber)
		}
	}

	CheckError(err)

	// fmt.Println("Tracking at hop", hopNumber)
	// trackPacketDelay(ByteArrayToJSONToPacket(ExtractPacketData(buf)).Number, hopNumber, SID)

	/* Make sure the destination address is not the current node. If the dest addr
	   is the current node, simply print out the message. If not, forward the message
	   to the desired address. */
}

// trackPacketDelay : Tracks packet delay in a bounced node
func trackPacketDelay(packetNumber int, hopNumber int, SID string) {
	t := time.Now()

	// Sends time 1 if on the first node, sends time 2 if on the last node
	var p Packet
	if hopNumber == 0 {
		p = Packet{Number: packetNumber, Time1: t}
	} else if hopNumber == 1 {
		p = Packet{Number: packetNumber, Time2: t}
	} else if hopNumber == 2 {
		p = Packet{Number: packetNumber, Time3: t}
	}

	data := PacketToJSONToByteArray(p)
	buff := MakeBlinkPacket(SID, localAddr, oracleAddr, PacketTrack, data)
	conn.WriteToUDP(buff, oracleAddr)
	fmt.Println("Packet tracked:", p.Number)
}
