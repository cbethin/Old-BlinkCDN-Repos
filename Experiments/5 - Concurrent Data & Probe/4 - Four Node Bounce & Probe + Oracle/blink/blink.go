package blink

import (
  "fmt"
  "net"
  "os"
  "strconv"
  "time"
)

const (
	//  PROBE_TIME = 3000 // Probe 20 times every minute
	ROUTE_TTL     = 15000 // 15 seconds
	PROBE_TIMEOUT = 10000 // 10 seconds
)

const (
  InitialProbe = "IP"
  ResponseOne  = "R1"
  ResponseTwo  = "R2"
  Iframe       = "IF"
  Pframe       = "PF"
  Bframe       = "BF"
  RoutingTableType = "RT"
)

var (
	RoutingTable  map[string]*Link
	localAddr     *net.UDPAddr
	oracleAddr 		*net.UDPAddr
	nextNode			*Link
  conn          *net.UDPConn
)


type Link struct {
	Addr            string
	Latency         float64
	Loss            int
	TTLExpiration   time.Time // When was the last time I probed you (Creepy!)
	TimeoutDeadline time.Time
	TimerStartTime  time.Time
}

/////////////////////////////
//// HEADER FUNCTIONALITY ///
/////////////////////////////

/* Creates a Blink Packet, which is simply an array of bytes, where the first 66
   bytes contain information about the source address, destination address, type of
   packet being sent. The remaining 958 bytes are used for data transmission
*/
func MakeBlinkPacket(srcAddr *net.UDPAddr, finalDestAddr *net.UDPAddr, packetType string, buf []byte) []byte {

  // This copies the src address and destination adress into bytes 0-32 and 32-64
  // respectively. The addresses are inputted as a string followed by a ! to let the program
  // know where the string ends. Bytes 64-66 are filled with the packet type as a string, and
  // the rest of the program are filled with the actual packet to send

  outBuf := make([]byte, 1024)
  copy(outBuf[:32], []byte(srcAddr.String() + "!"))
  copy(outBuf[32:64], []byte(finalDestAddr.String() + "!"))
  copy(outBuf[64:66], []byte(packetType))
  copy(outBuf[66:], buf)
  return outBuf
}




/* Extract all information from the Blink Packet. Return src Addr, Destination Addr,
   Packet Type, and Packet Data (in that order) */
func UnwrapHeader(inBuf []byte) (*net.UDPAddr, *net.UDPAddr, string, []byte) {
  srcAddr := ExtractSrcAddr(inBuf)
  finalDestAddr := ExtractFinalDestAddr(inBuf)
  packetType := ExtractPacketType(inBuf)
  packetData := ExtractPacketData(inBuf)

  return srcAddr, finalDestAddr, packetType, packetData
}




/* Extract Source Address from the header of the Blink Packet.
   Returns address as a pointer to a resolved net.UDPAddr
*/

func ExtractSrcAddr(inBuf []byte) *net.UDPAddr {

  // Pull in the header bytes corresponding to the src address (0-32)
  addrBuf := inBuf[:32]
  addrString := ""

  // Loop through the characters in that buffer and append each character to
  // an address string until we encounter the exclamation mark, which tells us
  // we have reached the end of the address.

  for _, value := range addrBuf {
    if string(value) != "!" {
      addrString += string(value)
    } else {
      break
    }
  }

  // Resolve the address string into a UDP address and return
  addr, err := net.ResolveUDPAddr("udp", addrString)
  CheckError(err)

  return addr
}




/* Extract Destination Address from the header of the Blink Packet.
   Returns address as a pointer to a resolved net.UDPAddr
*/
func ExtractFinalDestAddr(inBuf []byte) *net.UDPAddr {

  // Pull in the header bytes corresponding to the destination address (32-64)

  addrBuf := inBuf[32:64]
  addrString := ""

  // Loop through the characters in that buffer and append each character to
  // an address string until we encounter the exclamation mark, which tells us
  // we have reached the end of the address.

  for _, value := range addrBuf {
    if string(value) != "!" {
      addrString += string(value)
    } else {
      break
    }
  }

  // Resolve the address string into a UDP address and return
  addr, err := net.ResolveUDPAddr("udp", addrString)
  CheckError(err)

  return addr
}




/* From an inputted blink packet (an array of bytes), extract the packet type
   from the header. Returns value as a string type
*/
func ExtractPacketType(inBuf []byte) string {
  return string(inBuf[64:66])
}




/* From an inputted blink packet (an array of bytes), extract the packet's data
   from the packet
*/
func ExtractPacketData(inBuf []byte) []byte {
  return inBuf[66:]
}




/* Checks to see if inputted error is empty. If not (meaning there is an error)
   then the error is printed and the program quits. */
func CheckError(err error) {
  if err != nil {
    fmt.Println("Error:", err)
    os.Exit(1)
  }
}


/////////////////////////////
///// NODE FUNCTIONALITY ////
/////////////////////////////


/* Sends the initial probe packet to another node via a link, which is
   passed into the function.*/

func SendInitialProbe(link1 *Link) {
	// Update the link expiration as the first thing
	link1.TTLExpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)

	nodeAddr, err := net.ResolveUDPAddr("udp", link1.Addr)
	CheckError(err)

	// Create a Buffer now
	buffer := MakeBlinkPacket(localAddr, nodeAddr, InitialProbe, []byte(""))

	_, err = conn.WriteToUDP(buffer, nodeAddr)
	CheckError(err)

	link1.TimerStartTime = time.Now()
	link1.TimeoutDeadline = link1.TimerStartTime.Add(time.Millisecond * PROBE_TIMEOUT)
}




/* Sends the Response 1 packet to another node via a link, which is
   passed into the function. */

func SendResponse1(link1 *Link) {
	link1.TTLExpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)

	nodeAddr, err := net.ResolveUDPAddr("udp", link1.Addr)
	CheckError(err)

	// Create a Buffer now
	buffer := MakeBlinkPacket(localAddr, nodeAddr, ResponseOne, []byte(""))

	_, err = conn.WriteToUDP(buffer, nodeAddr)
	CheckError(err)

	link1.TimerStartTime = time.Now()
	link1.TimeoutDeadline = link1.TimerStartTime.Add(time.Millisecond * PROBE_TIMEOUT)
}




/* Sends the Response 2 packet to another node via a link, which is
   passed into the function. Will also calculate and store the calculated
	 Latency of the link. Sends routing table to oracle when it's done calculating
*/

func SendResponse2(link1 *Link) {

	// Calculate time elapsed since the probe was sent
	timeElapsed := time.Since(link1.TimerStartTime)

	// Resolve other node's address, create a Blink Response 2 packet
	// and send that buffer over the inputted connection object
	nodeAddr, err := net.ResolveUDPAddr("udp", link1.Addr)
	CheckError(err)

	buffer := MakeBlinkPacket(localAddr, nodeAddr, ResponseTwo, []byte(""))

	_, err = conn.WriteToUDP(buffer, nodeAddr)
	CheckError(err)

	// Calculates Latency
	link1.Latency = 0.9*link1.Latency + 0.1*((timeElapsed/2).Seconds()*1000)

	// Send Routing Table To Oracle
	for _, value := range RoutingTable {
		// Create a blink packet for
		buffer := MakeBlinkPacket(localAddr, oracleAddr, RoutingTableType, LinkToByteArray(*value))

		_, err := conn.WriteToUDP(buffer, oracleAddr)
		CheckError(err)
	}

}




/* Calculates the latency after a ResponseTwo is received and stores the latency in
	 the corresponding link. Sends routing table to oracle when it's done calculating
*/
func ReceiveResponse2(link1 *Link) {

	// Takes time elapsed since probe was sent
	timeElapsed := time.Since(link1.TimerStartTime)

	// Calculates Latency
	link1.Latency = 0.9*link1.Latency + 0.1*((timeElapsed/2).Seconds()*1000)

	// Send Routing Table To Oracle
	for _, value := range RoutingTable {
		// Create a blink packet to send the link to the oracle.
		buffer := MakeBlinkPacket(localAddr, oracleAddr, RoutingTableType, LinkToByteArray(*value))
		// Send the new Blink Packet to the oracle.
		_, err := conn.WriteToUDP(buffer, oracleAddr)
		CheckError(err)
	}
}




/* Every 100 milliseconds the program will check the TTLExpiration of the designated
   link. If the link has expired, the program will initialize the probing process by
	 sending an initial probe to the corresponding node */

func CheckExpiration(conn *net.UDPConn, link1 *Link) {

	for {
		if time.Now().Before(link1.TTLExpiration) == false {
			// If your Current time is before the expiation time
			// You do Nothing!
			// FOR FUTURE : The link that requires the probing funtion
			//              i.e. the one that times out, can update the
			//              linkNumber and call the SendProbe
			//              RoutingTable[linkNumber].SendProbe(nodeAddr,conn)
			SendInitialProbe(link1)
      fmt.Println("Sent IP to:", link1.Addr)
		}
		time.Sleep(100 * time.Millisecond)
	}
}




/* This function will forward an inputted buffer to the desired destination address.
	 if the current node is the desired destination address, the system will not Bounce
	 the blink packet
*/

func Bounce(buf []byte) {
	// Extract the destination address from the Blink packet
	destAddr := ExtractFinalDestAddr(buf)

	// Resolve the address of the next node to bounce to
  fmt.Println("PRinting Next Node:", nextNode)
	nextNodeAddr, err := net.ResolveUDPAddr("udp", nextNode.Addr)
	CheckError(err)

	/* Make sure the destination address is not the current node. If the dest addr
	   is the current node, simply print out the message. If not, forward the message
	   to the desired address. */
	if destAddr.String() != localAddr.String() {
		_, err := conn.WriteToUDP(buf, nextNodeAddr)
		CheckError(err)

		fmt.Println("-------------------------")
		fmt.Println("Bouncing", string(buf[66:75]))
	} else {
		fmt.Println("-------------------------")
		fmt.Println("Receiving", string(buf[66:75]))
	}
}


///////////////////////////////
//// USABLE NODE FUNCTIONS ////
///////////////////////////////

/* StartProber starts the probing functionality of the node. It must be started after the
   RoutingTable has been initialized and must be restarted every time the Routing Table is
   updated.
*/

func StartProber() {
  // go blink.RoutingTable[nodeAddr1.String()].checkExpiration(conn)
	// go blink.RoutingTable[nodeAddr2.String()].checkExpiration(conn)
	// go blink.RoutingTable[nodeAddr3.String()].checkExpiration(conn)
  // fmt.Println("Length of Routing Table:", len(RoutingTable))
  //
  for _, value := range RoutingTable {
    go CheckExpiration(conn, value)
  }

}


/* Read data flowing in to the listening connection that is passed into this function
   The function then checks the type of packet being sent and responds accordingly.
*/

func GetData() {

	buffer := make([]byte, 1024)

	for {
		// Read from connection, then extract packet type using Blink's built in function
		_, addr, err := conn.ReadFromUDP(buffer)
		CheckError(err)

		packetType := ExtractPacketType(buffer)
		fmt.Println(packetType, "from", addr.String())

		// Chech the packet type, and call the proper function for each packet type
		switch packetType {
		case InitialProbe:
			go SendResponse1(RoutingTable[addr.String()])
		case ResponseOne:
			go SendResponse2(RoutingTable[addr.String()])
		case ResponseTwo:
			go ReceiveResponse2(RoutingTable[addr.String()])
		case Iframe:
			fmt.Println("Received I Frame")
			go Bounce(buffer)
		case Pframe:
			fmt.Println("Received an P Frame")
			go Bounce(buffer)
		case Bframe:
			fmt.Println("Received an B Frame")
			go Bounce(buffer)
		default:
			fmt.Println("Packet recieved is an undefined type")
		}
	}
}




////////////////
///// Links ////
////////////////

func LinkToByteArray(link Link) []byte {
  byteArray := make([]byte, 192)
  copy(byteArray[0:16], strconv.FormatFloat(link.Latency, 'f', 15, 64))
  copy(byteArray[16:48], link.Addr + "!")
  return byteArray
}




/* Converts a byte array into a link. The byte array must be formatted properly
   using the blink LinkToByteArray() function.
*/
func ByteArrayToLink(byteArray []byte) Link {
  var newLink Link
  latencyString := string(byteArray[0:16])
  addressBytes := byteArray[16:48]
  addressString := ""

  for i:= 0; i < 32; i++ {
    if string(addressBytes[i]) != "!" {
      addressString += string(addressBytes[i])
    } else {
      break
    }
  }

  address, err := net.ResolveUDPAddr("udp", addressString)
  CheckError(err)

  newLink.Addr = address.String()
  newLink.Latency, _ = strconv.ParseFloat(latencyString, 64)
  return newLink
}





/////////////////////////////////////
///// Routing Tables & Addresses ////
/////////////////////////////////////

func SetRoutingTable(routingTable map[string]*Link) {
  RoutingTable = routingTable
}

func SetLocalAddr(localNodeAddr *net.UDPAddr) {
  localAddr = localNodeAddr
}

func SetOracleAddr(oracleNodeAddr *net.UDPAddr) {
  oracleAddr = oracleNodeAddr
}

func SetNextNode(NextNode *Link) {
  nextNode = NextNode
}

func SetupListener() {
  var err error
  conn, err = net.ListenUDP("udp", localAddr)
  CheckError(err)
  fmt.Println("Setup Listner on", localAddr.String())
}
