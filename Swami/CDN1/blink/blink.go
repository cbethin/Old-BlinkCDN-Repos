package blink

// const (
// 	ROUTE_TTL     = 15000 // 15 seconds
// 	PROBE_TIMEOUT = 10000 // 10 seconds
// )
//
// // Possible Blink packet types
// const (
// 	InitialProbe      = "IP"
// 	ResponseOne       = "R1"
// 	ResponseTwo       = "R2"
// 	Iframe            = "IF"
// 	Pframe            = "PF"
// 	Bframe            = "BF"
// 	RoutingTableType  = "RT"
// 	DecisionTableType = "DT"
// 	HelloOracle       = "HO"
// 	BounceUpdate      = "BU"
// 	ProbeID           = "0000000000000000"
// 	PacketTrack       = "PT"
// )
//
// // Global variables used by Blink nodes and oracle
// var (
// 	RoutingTable     map[string]*Link
// 	localAddr        *net.UDPAddr
// 	LocalAddrString  string
// 	oracleAddr       *net.UDPAddr
// 	DecisionTable    map[string][]string
// 	conn             *net.UDPConn
// 	fileName         string
// 	latencyArray     = make([]float64, 3)
// 	pathsArray       = make([][]string, 3)
// 	PacketTrackTable map[string]map[int]Packet
// )

// // Link : structure containing information about a give connection between two Blink nodes
// type Link struct {
// 	Addr            string
// 	Latency         float64
// 	Loss            int
// 	TTLExpiration   time.Time // When was the last time I probed you (Creepy!)
// 	TimeoutDeadline time.Time
// 	TimerStartTime  time.Time
// }
//
// // Packet : structure containing information about a given packet track
// type Packet struct {
// 	Number int
// 	Time1  time.Time
// 	Time2  time.Time
// 	Data   []byte
// }

// ////////////////////////////////////
// ///  MARK: Setting Global Values ///
// ////////////////////////////////////
//
// // SetRoutingTable : idenitfies the RoutingTable variable
// func SetRoutingTable(routingTable map[string]*Link) {
// 	RoutingTable = routingTable
// }
//
// // SetLocalAddr : identifies the local UDP address in blink
// func SetLocalAddr(localNodeAddr *net.UDPAddr) {
// 	localAddr = localNodeAddr
// }
//
// // SetOracleAddr : identifies the oracle's UDP address in blink
// func SetOracleAddr(oracleNodeAddr *net.UDPAddr) {
// 	oracleAddr = oracleNodeAddr
// }

// // SetupListener : Sets up the listening connection to be used on the node. Sets the connection object as the global variable conn.
// func SetupListener() {
// 	var err error
// 	conn, err = net.ListenUDP("udp", localAddr)
// 	CheckError(err)
// 	fmt.Println("Setup Listner on", localAddr.String())
// }
//
// ////////////////////
// /// MARK: Links ////
// ////////////////////
//
// // LinkToByteArray : Converts a link into a byte. The byte array is formatted to allow conversion using the ByteArrayToLink format
// func LinkToByteArray(link Link) []byte {
// 	byteArray := make([]byte, 192)
// 	copy(byteArray[0:16], strconv.FormatFloat(link.Latency, 'f', 15, 64))
// 	copy(byteArray[16:48], link.Addr+"!")
// 	return byteArray
// }
//
// // ByteArrayToLink : Converts a byte array into a link. The byte array must be formatted properly using the blink LinkToByteArray() function.
// func ByteArrayToLink(byteArray []byte) Link {
// 	var newLink Link
// 	latencyString := string(byteArray[0:16])
// 	addressBytes := byteArray[16:48]
// 	addressString := ""
//
// 	for i := 0; i < 32; i++ {
// 		if string(addressBytes[i]) != "!" {
// 			addressString += string(addressBytes[i])
// 		} else {
// 			break
// 		}
// 	}
//
// 	newLink.Addr = addressString
// 	newLink.Latency, _ = strconv.ParseFloat(latencyString, 64)
// 	return newLink
// }

///////////////////////////////////
//// MARK: HEADER FUNCTIONALITY ///
///////////////////////////////////

// // MakeBlinkPacket : Creates a Blink Packet, which is simply an array of bytes, where the first 66 bytes contain information about the source address, destination address, type of packet being sent. The remaining 958 bytes are used for data transmission
// func MakeBlinkPacket(SID string, srcAddr *net.UDPAddr, finalDestAddr *net.UDPAddr, packetType string, buf []byte) []byte {
//
// 	// This copies the src address and destination adress into bytes 0-32 and 32-64
// 	// respectively. The addresses are inputted as a string followed by a ! to let the program
// 	// know where the string ends. Bytes 64-66 are filled with the packet type as a string, and
// 	// the rest of the program are filled with the actual packet to send
//
// 	outBuf := make([]byte, 1024)
// 	copy(outBuf[:16], []byte(SID+"!"))
// 	copy(outBuf[16:48], []byte(srcAddr.String()+"!"))
// 	copy(outBuf[48:80], []byte(finalDestAddr.String()+"!"))
// 	copy(outBuf[80:82], []byte(packetType))
// 	copy(outBuf[82:], buf)
// 	return outBuf
// }
//
// // UnwrapHeader : Extract all information from the Blink Packet. Return src Addr, Destination Addr, Packet Type, and Packet Data (in that order)
// func UnwrapHeader(inBuf []byte) (string, *net.UDPAddr, *net.UDPAddr, string, []byte) {
// 	SID := ExtractSID(inBuf)
// 	srcAddr := ExtractSrcAddr(inBuf)
// 	finalDestAddr := ExtractFinalDestAddr(inBuf)
// 	packetType := ExtractPacketType(inBuf)
// 	packetData := ExtractPacketData(inBuf)
//
// 	return SID, srcAddr, finalDestAddr, packetType, packetData
// }
//
// // ExtractSID : Extract SID from the header of the Blink Packet. Returns SID as a string
// func ExtractSID(inBuf []byte) string {
// 	SID := string(inBuf[:16])
// 	return SID
// }
//
// // ExtractSrcAddr : Extract Source Address from the header of the Blink Packet. Returns address as a pointer to a resolved net.UDPAddr
// func ExtractSrcAddr(inBuf []byte) *net.UDPAddr {
//
// 	// Pull in the header bytes corresponding to the src address (0-32)
// 	addrBuf := inBuf[16:48]
// 	addrString := ""
//
// 	// Loop through the characters in that buffer and append each character to
// 	// an address string until we encounter the exclamation mark, which tells us
// 	// we have reached the end of the address.
//
// 	for _, value := range addrBuf {
// 		if string(value) != "!" {
// 			addrString += string(value)
// 		} else {
// 			break
// 		}
// 	}
//
// 	// Resolve the address string into a UDP address and return
// 	addr, err := net.ResolveUDPAddr("udp", addrString)
// 	CheckError(err)
//
// 	return addr
// }
//
// // ExtractFinalDestAddr : Extract Destination Address from the header of the Blink Packet. Returns address as a pointer to a resolved net.UDPAddr
// func ExtractFinalDestAddr(inBuf []byte) *net.UDPAddr {
//
// 	// Pull in the header bytes corresponding to the destination address (32-64)
//
// 	addrBuf := inBuf[48:80]
// 	addrString := ""
//
// 	// Loop through the characters in that buffer and append each character to
// 	// an address string until we encounter the exclamation mark, which tells us
// 	// we have reached the end of the address.
//
// 	for _, value := range addrBuf {
// 		if string(value) != "!" {
// 			addrString += string(value)
// 		} else {
// 			break
// 		}
// 	}
//
// 	// Resolve the address string into a UDP address and return
// 	addr, err := net.ResolveUDPAddr("udp", addrString)
// 	CheckError(err)
//
// 	return addr
// }
//
// // ExtractPacketType : From an inputted blink packet (an array of bytes), extract the packet type from the header. Returns value as a string type
// func ExtractPacketType(inBuf []byte) string {
// 	return string(inBuf[80:82])
// }
//
// // ExtractPacketData : From an inputted blink packet (an array of bytes), extract the packet's data from the packet
// func ExtractPacketData(inBuf []byte) []byte {
// 	return inBuf[82:]
// }

////////////////////////////////////
///// MARK: PROBE FUNCTIONALITY ////
////////////////////////////////////

// // SendInitialProbe : Sends the initial probe packet to another node via a link, which is passed into the function.
// func SendInitialProbe(link1 *Link) {
// 	// Update the link expiration as the first thing
// 	link1.TTLExpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)
//
// 	nodeAddr, err := net.ResolveUDPAddr("udp", link1.Addr)
// 	CheckError(err)
//
// 	// Create a Buffer now
// 	buffer := MakeBlinkPacket(ProbeID, localAddr, nodeAddr, InitialProbe, []byte(""))
//
// 	_, err = conn.WriteToUDP(buffer, nodeAddr)
// 	CheckError(err)
//
// 	link1.TimerStartTime = time.Now()
// 	link1.TimeoutDeadline = link1.TimerStartTime.Add(time.Millisecond * PROBE_TIMEOUT)
// }
//
// // SendResponse1 : Sends the Response 1 packet to another node via a link, which is passed into the function. */
// func SendResponse1(link1 *Link) {
// 	link1.TTLExpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)
//
// 	nodeAddr, err := net.ResolveUDPAddr("udp", link1.Addr)
// 	CheckError(err)
//
// 	// Create a Buffer now
// 	buffer := MakeBlinkPacket(ProbeID, localAddr, nodeAddr, ResponseOne, []byte(""))
//
// 	_, err = conn.WriteToUDP(buffer, nodeAddr)
// 	CheckError(err)
//
// 	link1.TimerStartTime = time.Now()
// 	link1.TimeoutDeadline = link1.TimerStartTime.Add(time.Millisecond * PROBE_TIMEOUT)
// }
//
// // SendResponse2 : Sends the Response 2 packet to another node via a link, which is passed into the function. Will also calculate and store the calculated Latency of the link. Sends routing table to oracle when it's done calculating
// func SendResponse2(link1 *Link) {
//
// 	// Calculate time elapsed since the probe was sent
// 	timeElapsed := time.Since(link1.TimerStartTime)
//
// 	// Resolve other node's address, create a Blink Response 2 packet
// 	// and send that buffer over the inputted connection object
// 	nodeAddr, err := net.ResolveUDPAddr("udp", link1.Addr)
// 	CheckError(err)
//
// 	buffer := MakeBlinkPacket(ProbeID, localAddr, nodeAddr, ResponseTwo, []byte(""))
//
// 	_, err = conn.WriteToUDP(buffer, nodeAddr)
// 	CheckError(err)
//
// 	// Calculates Latency
// 	link1.Latency = 0.9*link1.Latency + 0.1*((timeElapsed/2).Seconds()*1000)
//
// 	// Send Routing Table To Oracle
// 	for _, value := range RoutingTable {
// 		// Create a blink packet for
// 		buffer := MakeBlinkPacket(ProbeID, localAddr, oracleAddr, RoutingTableType, LinkToByteArray(*value))
//
// 		_, err := conn.WriteToUDP(buffer, oracleAddr)
// 		CheckError(err)
// 	}
//
// 	// Save time sent
// 	saveLatency((timeElapsed/2).Seconds()*1000, link1.Addr)
// }
//
// // ReceiveResponse2 : Calculates the latency after a ResponseTwo is received and stores the latency in the corresponding link. Sends routing table to oracle when it's done calculating
// func ReceiveResponse2(link1 *Link) {
//
// 	// Takes time elapsed since probe was sent
// 	timeElapsed := time.Since(link1.TimerStartTime)
//
// 	// Calculates Latency
// 	link1.Latency = 0.9*link1.Latency + 0.1*((timeElapsed/2).Seconds()*1000)
//
// 	// Send Routing Table To Oracle
// 	for _, value := range RoutingTable {
// 		// Create a blink packet to send the link to the oracle.
// 		buffer := MakeBlinkPacket(ProbeID, localAddr, oracleAddr, RoutingTableType, LinkToByteArray(*value))
// 		// Send the new Blink Packet to the oracle.
// 		_, err := conn.WriteToUDP(buffer, oracleAddr)
// 		CheckError(err)
//
// 	}
//
// 	// Save time recieved
// 	saveLatency((timeElapsed/2).Seconds()*1000, link1.Addr)
// }
//
// // CheckExpiration : Every 100 milliseconds the program will check the TTLExpiration of the designated link. If the link has expired, the program will initialize the probing process by sending an initial probe to the corresponding node
// func CheckExpiration(conn *net.UDPConn, link1 *Link) {
//
// 	for {
// 		if time.Now().Before(link1.TTLExpiration) == false {
// 			// If your Current time is before the expiation time
// 			// You do Nothing!
// 			// FOR FUTURE : The link that requires the probing funtion
// 			//              i.e. the one that times out, can update the
// 			//              linkNumber and call the SendProbe
// 			//              RoutingTable[linkNumber].SendProbe(nodeAddr,conn)
// 			SendInitialProbe(link1)
// 			fmt.Println("Sent IP to:", link1.Addr)
// 		}
// 		time.Sleep(100 * time.Millisecond)
// 	}
// }

//////////////////////////////////////
////// MARK: NODE FUNCTIONALITY //////
//////////////////////////////////////

// // StartNode : Starts a node instance and it's corresponding probing and data retreive processes.
// func StartNode() {
// 	StartProber()
// 	GetData()
// }
//
// // StartProber : starts the probing functionality of the node. It must be started after the RoutingTable has been initialized and must be restarted every time the Routing Table is updated.
// func StartProber() {
// 	for _, value := range RoutingTable {
// 		go CheckExpiration(conn, value)
// 	}
// }
//
// // GetData : Reads data flowing in to the listening connection that is passed into this function.  The function then checks the type of packet being sent and responds accordingly.
// func GetData() {
//
// 	// Create File
// 	fileName = strconv.Itoa(localAddr.Port) + ".txt"
// 	file, err := os.Create(fileName)
// 	CheckError(err)
// 	defer file.Close()
//
// 	// GET DATA
// 	buffer := make([]byte, 1024)
//
// 	for {
// 		// Read from connection, then extract packet type using Blink's built in function
// 		_, addr, err := conn.ReadFromUDP(buffer)
// 		CheckError(err)
//
// 		packetType := ExtractPacketType(buffer)
//
// 		fmt.Println("---------------------")
// 		// Chech the packet type, and call the proper function for each packet type
// 		switch packetType {
// 		case InitialProbe:
// 			fmt.Println("Received IP, Sending R1")
// 			SendResponse1(RoutingTable[addr.String()])
// 		case ResponseOne:
// 			fmt.Println("Received R1, Sending R2")
// 			SendResponse2(RoutingTable[addr.String()])
// 		case ResponseTwo:
// 			fmt.Println("Received R2")
// 			ReceiveResponse2(RoutingTable[addr.String()])
// 		case Iframe:
// 			fmt.Println("Handling I Frame")
// 			Bounce(buffer)
// 		case Pframe:
// 			fmt.Println("Handling P Frame")
// 			Bounce(buffer)
// 		case Bframe:
// 			fmt.Println("Handling B Frame")
// 			Bounce(buffer)
// 		case DecisionTableType:
// 			fmt.Println("Handling Decision Table")
// 			UpdateDecisionTable(buffer)
// 		default:
// 			fmt.Println("Packet recieved is an undefined type")
// 		}
// 		fmt.Println("---------------------")
// 	}
// }
//
// // Bounce : This function will forward an inputted buffer to the desired destination address. If the current node is the desired destination address, the system will not Bounce the blink packet
// func Bounce(buf []byte) {
//
// 	packetType := ExtractPacketType(buf)
// 	SID := ExtractSID(buf)
// 	packetNumber := ByteArrayToJSONToPacket(ExtractPacketData(buf)).Number
// 	destAddr := ExtractFinalDestAddr(buf)
// 	receivedTime := time.Now()
//
// 	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
// 	CheckError(err)
// 	printString := SID + " " + packetType + " " + strconv.Itoa(packetNumber) + " " + time.Now().String() + "\n"
// 	_, err = file.WriteString(printString)
// 	CheckError(err)
// 	file.Close()
//
// 	fmt.Println("Bouncing", packetNumber)
//
// 	// Checks to see if SID is present in the decision table. If it is
// 	// the funciton will proceed to find the next node to bounce to and
// 	// bounce the data to that node. If the SID is not found, nothing will be
// 	// done with the packet.
//
// 	nextNodeAddrString := ""
// 	hopNumber := -1
// 	if _, ok := DecisionTable[SID]; ok {
// 		// Is the current node the destination address?
// 		if destAddr.String() != LocalAddrString {
// 			// If thisNode != destNode then determine next node and foward data
// 			// to next node
// 			for i, value := range DecisionTable[SID] {
// 				if value == LocalAddrString {
// 					hopNumber = i
// 					if i == len(DecisionTable[SID])-1 {
// 						fmt.Println("Receiving", packetNumber)
// 						trackPacketDelay(packetNumber, hopNumber, SID)
// 						return
// 						// destAddr := ExtractFinalDestAddr(buf)
// 						// nextNodeAddrString = destAddr.String()
// 					}
//
// 					nextNodeAddrString = DecisionTable[SID][i+1]
// 					// fmt.Println("Next Node:", nextNodeAddrString)
// 				}
// 			}
//
// 			// If current node was not found, break off. Bounce was incorrect
// 			if nextNodeAddrString == "" {
// 				fmt.Println("Next node not found")
// 				return
// 			}
//
// 			// Resolve the NextNode's address and bounce it off.
// 			nextNodeAddr, err := net.ResolveUDPAddr("udp", nextNodeAddrString)
// 			CheckError(err)
// 			fmt.Println("Actual Next Node:", nextNodeAddr.String())
//
// 			_, err = conn.WriteToUDP(buf, nextNodeAddr)
// 			CheckError(err)
//
// 			// Send time update to oracle
// 			message := packetType + strconv.Itoa(packetNumber) + " received-at: " + LocalAddrString + " at: " + receivedTime.String()
// 			blinkPacket := MakeBlinkPacket(ProbeID, localAddr, oracleAddr, BounceUpdate, []byte(message))
//
// 			_, err = conn.WriteToUDP(blinkPacket, oracleAddr)
// 			CheckError(err)
// 		} else {
// 			fmt.Println("Receiving", packetNumber)
// 		}
// 	}
//
// 	CheckError(err)
//
// 	if hopNumber == 0 || hopNumber == 2 {
// 		fmt.Println("Tracking at hop", hopNumber)
// 		trackPacketDelay(packetNumber, hopNumber, SID)
// 	}
//
// 	/* Make sure the destination address is not the current node. If the dest addr
// 	   is the current node, simply print out the message. If not, forward the message
// 	   to the desired address. */
// }
//
// // trackPacketDelay : Tracks packet delay in a bounced node
// func trackPacketDelay(packetNumber int, hopNumber int, SID string) {
// 	t := time.Now()
//
// 	// Sends time 1 if on the first node, sends time 2 if on the last node
// 	var p Packet
// 	if hopNumber == 0 {
// 		p = Packet{Number: packetNumber, Time1: t}
// 	} else if hopNumber == 2 {
// 		p = Packet{Number: packetNumber, Time2: t}
// 	}
//
// 	data := PacketToJSONToByteArray(p)
// 	buff := MakeBlinkPacket(SID, localAddr, oracleAddr, PacketTrack, data)
// 	conn.WriteToUDP(buff, oracleAddr)
// 	fmt.Println("Packet tracked:", p.Number)
// }

///////////////////////////////////////////
////////// MARK: DECISION TABLES //////////
///////////////////////////////////////////

// // UpdateDecisionTable : Updates Global Decision Table at a given node based on an incoming blink packet.
// func UpdateDecisionTable(inBuf []byte) {
// 	tableData := ExtractPacketData(inBuf)
// 	decisionTable := ByteArrayToDecisionTable(tableData)
// 	SetDecisionTable(decisionTable)
// 	fmt.Println("Decision Update:", DecisionTable)
// }
//
// // SetDecisionTable : Sets the value of the global decision table
// func SetDecisionTable(decisionTable map[string][]string) {
// 	DecisionTable = decisionTable
// }
//
// // DecisionTableToByteArray : Converts a decision table into an array of bytes, currently only supports one session
// func DecisionTableToByteArray(decisionTable map[string][]string) []byte {
// 	outBuf := make([]byte, 100)
// 	for key, value := range decisionTable {
// 		copy(outBuf[:16], []byte(key))
// 		// Add Hops Array to byte arra
// 		hops := value[0] + "!" + value[1] + "!" + value[2] + "!"
// 		copy(outBuf[16:], []byte(hops))
// 	}
//
// 	return outBuf
// }
//
// // ByteArrayToDecisionTable : Converts an array of bytes into a decision table. Currently only supports one session
// func ByteArrayToDecisionTable(inBuf []byte) map[string][]string {
//
// 	decisionTable := make(map[string][]string)
// 	SID := string(inBuf[:16])
// 	hopsArray := inBuf[16:]
// 	hopStringsArray := []string{}
// 	addrCount := 0
// 	addrString := ""
//
// 	for _, value := range hopsArray {
// 		if addrCount < 3 {
// 			if string(value) != "!" {
// 				addrString += string(value)
// 			} else {
// 				hopStringsArray = append(hopStringsArray, addrString)
// 				// hopStringsArray[addrCount] = addrString
// 				addrCount++
// 				addrString = ""
// 			}
// 		} else {
// 			break
// 		}
// 	}
//
// 	decisionTable[SID] = hopStringsArray
// 	return decisionTable
// }

///////////////////////////////////
///  MARK: ORACLE FUNCTIONALITY ///
///////////////////////////////////

// var (
// 	nodesTable         []string // array containing addresses of all blink nodes
// 	locations          map[string]string
// 	oracleRoutingTable map[string][]*Link
// )

// // StartOracle : This function initializes main functionality of the oracle. It reads any incoming data and decides how to appropriately respond to it.
// func StartOracle(oracleAddrString string) {
//
// 	// Update IP addresses accordingly for nodes in use
// 	nodesTable = []string{"52.53.177.194:8001", "18.184.225.196:8001", "35.176.239.10:8001", "13.115.224.27:8001"}
// 	go startHTTPListening()
//
// 	// Set Nodes A,B,C,D = corresponding index in nodesTable array
//
// 	// Setup Oracle
// 	var err error
// 	oracleAddr, err = net.ResolveUDPAddr("udp", oracleAddrString)
// 	CheckError(err)
//
// 	conn, err = net.ListenUDP("udp", oracleAddr)
// 	CheckError(err)
// 	fmt.Println("Listening...")
//
// 	// Setup Routing Table (THIS IS HARD CODED IN RIGHT NOW)
// 	var newLink1 Link
// 	newLink1.Addr = nodesTable[0]
// 	var newLink2 Link
// 	newLink2.Addr = nodesTable[1]
// 	var newLink3 Link
// 	newLink3.Addr = nodesTable[2]
// 	var newLink4 Link
// 	newLink4.Addr = nodesTable[3]
//
// 	newRoutingTable := make(map[string][]*Link)
// 	newRoutingTable[nodesTable[0]] = []*Link{&newLink2, &newLink3, &newLink4}
// 	newRoutingTable[nodesTable[1]] = []*Link{&newLink1, &newLink3, &newLink4}
// 	newRoutingTable[nodesTable[2]] = []*Link{&newLink1, &newLink2, &newLink4}
// 	newRoutingTable[nodesTable[3]] = []*Link{&newLink1, &newLink2, &newLink3}
// 	oracleRoutingTable = newRoutingTable
//
// 	// Setup PacketTrackTable
// 	PacketTrackTable = make(map[string]map[int]Packet)
//
// 	//Setup the data file
// 	fileName = strconv.Itoa(oracleAddr.Port) + ".txt"
// 	file, err := os.Create(fileName)
// 	CheckError(err)
// 	defer file.Close()
//
// 	// Set up Read Buffer
// 	buffer := make([]byte, 1024)
//
// 	// Read the Data and decide what to do
// 	for {
// 		_, addr, err := conn.ReadFromUDP(buffer)
// 		CheckError(err)
// 		packetType := ExtractPacketType(buffer)
//
// 		// fmt.Println("---------------------")
// 		switch packetType {
// 		case RoutingTableType:
// 			// fmt.Println("Routing Table Received")
// 			UpdateRoutingTable(addr, buffer)
// 			// SelectPath(conn)
// 			//fmt.Println(oracleRoutingTable)
// 		case HelloOracle:
// 			SetupServer(addr, conn, buffer)
// 			// fmt.Println("Received Hello Oracle")
// 		case BounceUpdate:
// 			handleBounceUpdate(buffer)
// 		case PacketTrack:
// 			handlePacketTrackUpdate(buffer)
// 		default:
// 			// fmt.Println("Packet type unrecognized")
// 		}
// 		// fmt.Println("---------------------")
//
// 	}
// }
//
// func startHTTPListening() {
// 	locations = make(map[string]string)
// 	for i := 0; i < len(nodesTable); i++ {
// 		locations[nodesTable[i]] = nodesTable[i]
// 	}
//
// 	// Handle HTTP Functionality to convey info
// 	http.HandleFunc("/", handleHTTPResponse)
// 	if err := http.ListenAndServe(":8081", nil); err != nil {
// 		panic(err)
// 	}
// }
//
// // handleHTTPResponse : handles http responses to the oracle
// func handleHTTPResponse(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println(r.URL.Path)
// 	if r.Method == "GET" {
// 		if r.URL.Path == "/getpaths" {
// 			query := r.URL.Query()
// 			destination := query["dest"][0]
// 			source := query["source"][0]
// 			// sessionID := query["sid"][0]
//
// 			formPath(source, destination)
// 			GetLatencyForPaths()
//
// 			message := ""
// 			for i := 0; i < len(latencyArray); i++ {
// 				message += strconv.FormatFloat(latencyArray[i], 'f', -1, 32)
// 				message += " "
// 			}
//
// 			w.Write([]byte(message))
// 		}
// 	} else if r.Method == "POST" {
// 		if r.URL.Path == "/setpath" {
// 			path := r.URL.Query()["path"][0]
// 			pathNumber, err := strconv.Atoi(path)
// 			CheckError(err)
// 			sessionID := "0001000100010001"
// 			fmt.Println(pathsArray)
// 			selectedPath := []string{pathsArray[pathNumber][0], pathsArray[pathNumber][1], pathsArray[pathNumber][2]}
// 			setPath(selectedPath, sessionID)
// 			tellClientToStart(pathsArray[pathNumber][0])
// 			w.Write([]byte("success"))
// 		} else if r.URL.Path == "/resettest" {
// 			wipeDatabase()
// 			w.Write([]byte("success"))
// 		}
// 	}
// }
//
// // SetupServer : This function initializes a data-sending server that lies outside of the Blink Node System. It responds to a HelloOracle message (sent by the server) by setting up a SessionID as well as an initial series of hops for the session. The oracle then sends the SessionID and the first blink node's address to the Server. The oracle then sends Decision Table updates to all of the nodes
// func SetupServer(serverAddr *net.UDPAddr, conn *net.UDPConn, buf []byte) {
// 	destAddr := ExtractFinalDestAddr(buf)
// 	fmt.Println("Dest Addr:", destAddr.String())
// 	SID := "1001100110011001"
//
// 	decisionTable := make(map[string][]string)
// 	decisionTable[SID] = []string{nodesTable[0], nodesTable[2], destAddr.String()}
// 	DecisionTable = decisionTable
//
// 	// Send Reply to Server
// 	message := SID + decisionTable[SID][0] + "!"
//
// 	blinkPacket := MakeBlinkPacket(SID, oracleAddr, serverAddr, HelloOracle, []byte(message))
//
// 	_, err := conn.WriteToUDP(blinkPacket, serverAddr)
// 	fmt.Println("Contacted by Server")
// 	CheckError(err)
//
// 	// Send Decision Tables to Nodes
// 	for _, value := range nodesTable {
// 		addr, err := net.ResolveUDPAddr("udp", value)
// 		CheckError(err)
//
// 		SelectPath(conn) // Update the decision table before sending
// 		buffer := MakeBlinkPacket("0000000000000011", oracleAddr, serverAddr, DecisionTableType, DecisionTableToByteArray(decisionTable))
// 		_, err = conn.WriteToUDP(buffer, addr)
// 		CheckError(err)
// 		fmt.Println("Sent DT to:", addr.String())
// 		fmt.Println(DecisionTable)
// 	}
// }
//
// // UpdateRoutingTable : This function updates the GlobalRoutingTable based on information from receieved Routing Table
// func UpdateRoutingTable(addr *net.UDPAddr, inBuf []byte) {
// 	// See if the routing table from this source exists
// 	for sourceAddress, routingTable := range oracleRoutingTable {
// 		// If it does exist.
// 		if sourceAddress == addr.String() {
// 			//Extract packet data
// 			packetData := ExtractPacketData(inBuf)
// 			//Convert Byte array to link
// 			newLink := ByteArrayToLink(packetData)
// 			//Store that link.
// 			for i, value := range routingTable {
// 				if value.Addr == newLink.Addr {
// 					routingTable[i] = &newLink
// 					oracleRoutingTable[sourceAddress] = routingTable
// 				}
// 			}
// 		}
// 	}
// }
//
// // SelectPath : function selects the best of two paths (based on lowest latency), updates the decision table accordingly. The oracle will then distribute this decision table to all of the nodes.
// func SelectPath(conn *net.UDPConn) {
// 	latency1to2 := oracleRoutingTable[nodesTable[1]][0].Latency
// 	latency2to4 := oracleRoutingTable[nodesTable[3]][1].Latency
// 	latency1to3 := oracleRoutingTable[nodesTable[2]][0].Latency
// 	latency3to4 := oracleRoutingTable[nodesTable[3]][2].Latency
// 	// fmt.Println("1to2:", latency1to2, "2to4:", latency2to4)
// 	// fmt.Println("1to3:", latency1to3, "3to4:", latency3to4)
//
// 	latencyBlue := latency1to3 + latency3to4
// 	latencyRed := latency1to2 + latency2to4
//
// 	// fmt.Println("Red:", latencyRed, "Blue:", latencyBlue)
//
// 	bluePath := []string{nodesTable[0], nodesTable[2], nodesTable[3]}
// 	redPath := []string{nodesTable[0], nodesTable[1], nodesTable[3]}
//
// 	if latencyBlue <= latencyRed {
// 		// fmt.Println("Blue Path Selected")
// 		for key := range DecisionTable {
// 			fmt.Println("Key", key)
// 			DecisionTable[key] = bluePath
// 		}
// 	} else {
// 		// fmt.Println("Red Path Selected")
// 		for key := range DecisionTable {
// 			fmt.Println("SID:", key)
// 			DecisionTable[key] = redPath
// 		}
// 	}
//
// 	packetData := DecisionTableToByteArray(DecisionTable)
// 	blinkPacket := MakeBlinkPacket("0001000100010001", oracleAddr, oracleAddr, DecisionTableType, packetData)
//
// 	// Send the decision table to all the nodes
// 	for _, nodeAddrString := range nodesTable {
// 		addr, err := net.ResolveUDPAddr("udp", nodeAddrString)
// 		CheckError(err)
//
// 		_, err = conn.WriteToUDP(blinkPacket, addr)
// 		CheckError(err)
// 	}
//
// }
//
// // setPath : Sets the path for a given session in the decision table. Takes an inputted array of address strings and sets them as the decision table. It then distributes decision table to all Blink nodes.
// func setPath(path []string, sessionID string) {
// 	decisionTable := make(map[string][]string)
// 	decisionTable["0001000100010001"] = path
// 	DecisionTable = decisionTable
//
// 	packetData := DecisionTableToByteArray(DecisionTable)
// 	blinkPacket := MakeBlinkPacket("0001000100010001", oracleAddr, oracleAddr, DecisionTableType, packetData)
//
// 	// Send the decision table to all the nodes
// 	for _, nodeAddrString := range nodesTable {
// 		addr, err := net.ResolveUDPAddr("udp", nodeAddrString)
// 		CheckError(err)
//
// 		_, err = conn.WriteToUDP(blinkPacket, addr)
// 		CheckError(err)
// 	}
//
// 	fmt.Println("Paths set")
// }
//
// // formPath : given two Blink node addresses, forms any 2 or 3 node path between the two nodes. Assigns this path to the global pathsArray variable.
// func formPath(nodeA string, nodeB string) {
//
// 	// Forms 3 paths from source(A) to destination(B)
//
// 	var bounceNode string
// 	pathIndex := 0
// 	for i := range pathsArray {
// 		pathsArray[i] = make([]string, 3)
// 	}
//
// 	for i := 0; i < len(nodesTable); i++ {
// 		/*
// 			If middle node is not equal to A then set middle node equal to bounce node
// 			This will result in 3 paths (assuming 4 total nodes)
// 			Example ( A -> C -> B) (A -> D -> B) (A -> B)
// 		*/
// 		if nodesTable[i] != nodeA {
// 			bounceNode = nodesTable[i]
// 			// Assign Items in the Array
// 			pathsArray[pathIndex][0] = nodeA
// 			pathsArray[pathIndex][1] = bounceNode
// 			pathsArray[pathIndex][2] = nodeB
// 			pathIndex += 1
// 		}
// 	}
//
// }
//
// // GetLatencyForPaths : calculates the measured latency of each path in pathsArray using oracleRoutingTable values. Stores in global latencyArray.
// func GetLatencyForPaths() {
// 	/*
// 		Gets latency for each path and puts it in latencyArray at corresponding index
// 		The 3 paths will then be showed to the user
// 	*/
// 	latencyIndex := 0
//
// 	for i := 0; i < len(pathsArray); i++ {
// 		// Latency from the start of the path to the mid point node
// 		route1oracle := oracleRoute(pathsArray[i][0], pathsArray[i][1])
//
// 		routeStartToHalf := oracleRoutingTable[pathsArray[i][0]][route1oracle].Latency
//
// 		route2oracle := oracleRoute(pathsArray[i][1], pathsArray[i][2])
//
// 		// Latency from t he mid point node to the destination node
// 		routeHalfToEnd := 0.0
// 		if pathsArray[i][1] != pathsArray[i][2] {
// 			routeHalfToEnd = oracleRoutingTable[pathsArray[i][1]][route2oracle].Latency
// 		}
//
// 		// Total route latency by adding both delays together
// 		routeLatency := routeStartToHalf + routeHalfToEnd
// 		latencyArray[latencyIndex] = routeLatency
// 		latencyIndex++
// 	}
//
// 	sortLatencyAndPathArrays()
// }
//
// // sortLatencyAndPathArrays : using global latencyArray, this function will sort the latencyArray (lowest->highest), sorting the pathsArray along with it to assure the latencyArray index always matches its corresponding pathsArray index.
// func sortLatencyAndPathArrays() {
// 	fmt.Println("Before:", latencyArray)
// 	fmt.Println(pathsArray)
//
// 	for i := 0; i < len(pathsArray)-1; i++ {
// 		for j := 0; j < len(pathsArray)-1; j++ {
// 			if latencyArray[j] > latencyArray[j+1] {
// 				tmpLat := latencyArray[j]
// 				tmpPath := pathsArray[j]
//
// 				latencyArray[j] = latencyArray[j+1]
// 				pathsArray[j] = pathsArray[j+1]
//
// 				latencyArray[j+1] = tmpLat
// 				pathsArray[j+1] = tmpPath
// 			}
// 		}
// 	}
//
// 	fmt.Println("After:", latencyArray)
// 	fmt.Println("-------")
// 	fmt.Println(pathsArray)
// }
//
// // oracleRoute : returns the index of the link in the oracleRoutingTable given source and destination
// func oracleRoute(A string, B string) int {
// 	for j := 0; j < len(oracleRoutingTable[A]); j++ {
// 		if oracleRoutingTable[A][j].Addr == B {
// 			return j
// 		}
// 	}
//
// 	return -1
// }
//
// func tellClientToStart(nodeAddr string) {
// 	nodeIndex := 0
// 	for i := 0; i < len(nodesTable); i++ {
// 		if nodesTable[i] == nodeAddr {
// 			nodeIndex = i
// 		}
// 	}
//
// 	clientAddr, err := net.ResolveUDPAddr("udp", "18.216.192.154:8000")
// 	CheckError(err)
//
// 	message := []byte("node" + strconv.Itoa(nodeIndex))
// 	_, err = conn.WriteToUDP(message, clientAddr)
// }
//
// func handleBounceUpdate(inBuf []byte) {
//
// 	// Print bounce update
// 	packetData := ExtractPacketData(inBuf)
// 	sourceAddr := ExtractSrcAddr(inBuf)
//
// 	// Write Bounce update to file
// 	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
// 	CheckError(err)
// 	printString := string(packetData) + " from " + sourceAddr.String() + "\n"
// 	_, err = file.WriteString(printString)
// 	file.Close()
//
// 	fmt.Println(string(packetData) + " from " + sourceAddr.String())
// }
//
// func saveLatency(latency float64, addr string) {
// 	timeStamp := time.Now().String()
//
// 	file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_WRONLY, 0600)
// 	CheckError(err)
//
// 	latencyString := strconv.FormatFloat(latency, 'f', 15, 64)
// 	printString := latencyString + ", " + addr + ", " + timeStamp + ", \n"
// 	_, err = file.WriteString(printString)
// 	file.Close()
// }
//
//
// func handlePacketTrackUpdate(buffer []byte) {
// 	SID := ExtractSID(buffer)
// 	data := ExtractPacketData(buffer)
// 	packet := ByteArrayToJSONToPacket(data)
//
// 	// fmt.Println("Received packet update:", packet)
//
// 	if _, ok := PacketTrackTable[SID]; ok {
// 		if _, ok2 := PacketTrackTable[SID][packet.Number]; ok2 {
// 		} else {
// 			PacketTrackTable[SID][packet.Number] = packet
// 		}
// 	} else {
// 		packetMap := make(map[int]Packet)
// 		PacketTrackTable[SID] = packetMap
// 	}
//
// 	if !packet.Time1.IsZero() {
// 		p := Packet{Number: packet.Number, Time1: packet.Time1, Time2: PacketTrackTable[SID][packet.Number].Time2}
// 		PacketTrackTable[SID][packet.Number] = p
// 	}
//
// 	if !packet.Time2.IsZero() {
// 		p := Packet{Number: packet.Number, Time1: PacketTrackTable[SID][packet.Number].Time1, Time2: packet.Time2}
// 		PacketTrackTable[SID][packet.Number] = p
// 	}
//
// 	if !PacketTrackTable[SID][packet.Number].Time1.IsZero() && !PacketTrackTable[SID][packet.Number].Time2.IsZero() {
// 		// fmt.Println("PACKET TRACK UPDATING:", PacketTrackTable[SID][packet.Number])
// 		extractPacketTrackLatency(PacketTrackTable[SID][packet.Number])
// 	}
// }

//////////////////////////////
/// MARK: HELPER FUNCTIONS ///
//////////////////////////////

// // PacketToJSONToByteArray : converts a Packet type to a JSON object encoded and returned as a byteArray
// func PacketToJSONToByteArray(p Packet) []byte {
// 	b := new(bytes.Buffer)
// 	json.NewEncoder(b).Encode(p)
// 	return b.Bytes()
// }
//
// // ByteArrayToJSONToPacket : converts a []byte and decodes it as a JSON following the Packet type structure
// func ByteArrayToJSONToPacket(b []byte) Packet {
// 	var p Packet
// 	err := json.NewDecoder(bytes.NewReader(b)).Decode(&p)
// 	CheckError(err)
// 	return p
// }
//
// // extractPacketTrackLatency : from a given Packet, subtracts Time2 from Time1 and saves latency to database
// func extractPacketTrackLatency(p Packet) { // SET BACK TO MAP
// 	// Takes t1 and t2 and finds t total then sends that, packet number and service time to function firebaseData
// 	t1 := p.Time1
// 	t2 := p.Time2
// 	pktNum := p.Number
// 	// serviceTime := time.service
//
// 	totalTripTime := t2.Sub(t1)
// 	if totalTripTime < 0 {
// 		fmt.Println("Error in trip time calculation:", p)
// 		return
// 	}
//
// 	latencyData := make(map[string]float64)
// 	latencyData["packetNumber"] = float64(pktNum)
// 	latencyData["time"] = totalTripTime.Seconds()
//
// 	fmt.Println("Updated packet", pktNum, " | Latency:", latencyData["time"])
// 	firebaseData(latencyData)
// }

// // firebaseData : stores an inputted map[string]float64 to firebase database
// func firebaseData(latencyData map[string]float64) {
// 	sa := option.WithCredentialsFile("./blink/swami-database-firebase-adminsdk-6e01e-6d14ba6b69.json")
//
// 	app, err := firebase.NewApp(context.Background(), nil, sa)
// 	CheckError(err)
//
// 	client, err := app.Firestore(context.Background())
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
//
// 	_, err = client.Collection("Latency_Data").Doc("Latency").Set(context.Background(), latencyData)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
//
// 	// fmt.Println(result)
// 	defer client.Close()
// }
//
// // wipeDatabase : wipes the current dataset in the firebase databasse
// func wipeDatabase() {
// 	fmt.Println("Wiping database")
// 	sa := option.WithCredentialsFile("./blink/swami-database-firebase-adminsdk-6e01e-6d14ba6b69.json")
//
// 	app, err := firebase.NewApp(context.Background(), nil, sa)
// 	CheckError(err)
//
// 	client, err := app.Firestore(context.Background())
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
//
// 	latencyData := make(map[string]float64)
// 	latencyData["packetNumber"] = float64(0)
// 	latencyData["time"] = -1
// 	fmt.Println(latencyData)
//
// 	result, err := client.Collection("Latency_Data").Doc("Latency").Set(context.Background(), latencyData)
// 	fmt.Println(result)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
//
// }
