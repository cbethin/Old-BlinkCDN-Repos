package blink

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	// firebase "firebase.google.com/go"
)

var (
	nodesTable         []string // array containing addresses of all blink nodes
	oracleRoutingTable map[string][]*Link
	processes          = make(map[string]*Process)
	packetTrackMutex   = sync.RWMutex{}
	firebaseCtx        = context.Background()
	app                *firebase.App
	client             *firestore.Client
)

const (
	SESSION_ID_TEST  = "1001100110011001"
	clientAddrString = "18.216.192.154:9000"
	// clientAddrString = "155.246.113.15:8000"
	nodeAddrString = "155.246.45.43"
)

// StartOracle : This function initializes main functionality of the oracle. It reads any incoming data and decides how to appropriately respond to it.
func StartOracle(oracleAddrString string) {

	// Update IP addresses accordingly for nodes in use
	nodesTable = []string{"52.53.177.194:8001", "18.184.225.196:8001", "35.176.239.10:8001", "13.115.224.27:8001"}
	// nodesTable = []string{nodeAddrString + ":8001", nodeAddrString + ":8002", nodeAddrString + ":8003", nodeAddrString + ":8004"}
	go startHTTPListening()

	// Set Nodes A,B,C,D = corresponding index in nodesTable array
	sa := option.WithCredentialsFile("./blink/swami-database-firebase-adminsdk-6e01e-6d14ba6b69.json")
	var err error
	app, err = firebase.NewApp(firebaseCtx, nil, sa)
	CheckError(err)
	client, err = app.Firestore(firebaseCtx)
	CheckError(err)

	// Setup Oracle
	oracleAddr, err = net.ResolveUDPAddr("udp", oracleAddrString)
	CheckError(err)

	conn, err = net.ListenUDP("udp", oracleAddr)
	CheckError(err)
	fmt.Println("Listening...")

	// Setup Routing Table (THIS IS HARD CODED IN RIGHT NOW)
	newRoutingTable := make(RouteTable)
	newRoutingTable.Initialize(nodesTable)
	oracleRoutingTable = newRoutingTable

	// Setup PacketTrackTable
	PacketTrackTable = make(map[string]map[int]Packet)

	//Setup the data file
	// fileName = strconv.Itoa(oracleAddr.Port) + ".txt"
	// file, err := os.Create(fileName)
	// CheckError(err)
	// defer file.Close()

	// Set up Read Buffer
	buffer := make([]byte, MaxPacketSize)

	// Read the Data and decide what to do
	for {
		_, addr, err := conn.ReadFromUDP(buffer)
		CheckError(err)
		packetType := ExtractPacketType(buffer)

		// fmt.Println("---------------------")
		switch packetType {
		case RoutingTableType:
			// fmt.Println("Routing Table Received")
			UpdateRoutingTable(addr, buffer)
			// SelectPath(conn)
			//fmt.Println(oracleRoutingTable)
		case HelloOracle:
			SetupServer(addr, conn, buffer)
			// fmt.Println("Received Hello Oracle")
		case BounceUpdate:
			handleBounceUpdate(buffer)
		case PacketTrack:
			handlePacketTrackUpdate(buffer)
		case BlinkSignal:
			handleBlinkSignal(buffer)
		default:
			// fmt.Println("Packet type unrecognized")
		}

	}
}

func handleBlinkSignal(buf []byte) {
	if string(ExtractPacketData(buf)[:7]) == "success" {
		SID := ExtractSID(buf)
		fn := processes[SID].Callback
		fn()
		processes[SID].HasCompleted = true
		clearSession(SID)
	}
}

func startHTTPListening() {

	// Handle HTTP Functionality to convey info
	http.HandleFunc("/", handleFuncWrapper)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}

func handleFuncWrapper(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got HTTP Request: ", r.URL.Path)
	c := make(chan bool)
	go handleHTTPResponse(w, r, c)
	x := <-c
	if x == true {
		fmt.Println("Process " + r.URL.Path + " completed")
	}
}

// handleHTTPResponse : handles http responses to the oracle
func handleHTTPResponse(w http.ResponseWriter, r *http.Request, c chan bool) {
	if r.Method == "GET" {
		if r.URL.Path == "/getpaths" {
			query := r.URL.Query()
			destination := query["dest"][0]
			source := query["source"][0]
			// sessionID := query["sid"][0]

			formPath(source, destination)
			GetLatencyForPaths()

			message := ""
			for i := 0; i < len(latencyArray); i++ {
				message += strconv.FormatFloat(latencyArray[i], 'f', -1, 32)
				message += " "
			}

			w.Write([]byte(message))
			c <- true
		}
	} else if r.Method == "POST" {
		if r.URL.Path == "/setpath" {
			path := r.URL.Query()["path"][0]
			pathNumber, err := strconv.Atoi(path)
			CheckError(err)
			sessionID := SESSION_ID_TEST
			selectedPath := []string{pathsArray[pathNumber][0], pathsArray[pathNumber][1], pathsArray[pathNumber][2]}
			setPath(selectedPath, sessionID)
			// tellClientToStart(pathsArray[pathNumber][0], "")
			w.Write([]byte(selectedPath[0] + " " + selectedPath[1] + " " + selectedPath[2]))
			fmt.Println("Paths Set")
			c <- true
		} else if r.URL.Path == "/resettest" {
			// wipeDatabase()
			w.Write([]byte("success"))
			c <- true
		} else if r.URL.Path == "/getfile" {
			filename := r.URL.Query()["file"][0]
			sid := r.URL.Query()["sid"][0]
			fn := func() {
				fmt.Println("Sending file complete.")
				w.Write([]byte("success22"))
			}

			processes[sid] = &Process{Callback: fn, HasCompleted: false}

			// Tell Client to start sending the file, and then wait for the process to finish
			tellClientToStart(DecisionTable[SESSION_ID_TEST][0], filename)
			for {
				if processes[sid].HasCompleted {
					c <- true
					return
				}

				time.Sleep(5 * time.Millisecond)
			}
			// w.Write([]byte("success"))
		}
	}
}

// SetupServer : This function initializes a data-sending server that lies outside of the Blink Node System. It responds to a acle message (sent by the server) by setting up a SessionID as well as an initial series of hops for the session. The oracle then sends the SessionID and the first blink node's address to the Server. The oracle then sends Decision Table updates to all of the nodes
func SetupServer(serverAddr *net.UDPAddr, conn *net.UDPConn, buf []byte) {
	destAddr := ExtractFinalDestAddr(buf)
	fmt.Println("Dest Addr:", destAddr.String())
	// SID := "1001100110011001"
	SID := SESSION_ID_TEST

	decisionTable := make(map[string][]string)
	decisionTable[SID] = []string{nodesTable[0], nodesTable[2], destAddr.String()}
	DecisionTable = decisionTable

	// Send Reply to Server
	message := SID + decisionTable[SID][0] + "!"

	blinkPacket := MakeBlinkPacket(SID, oracleAddr, serverAddr, HelloOracle, []byte(message))

	_, err := conn.WriteToUDP(blinkPacket, serverAddr)
	fmt.Println("Contacted by Server")
	CheckError(err)

	// Send Decision Tables to Nodes
	for _, value := range nodesTable {
		addr, err := net.ResolveUDPAddr("udp", value)
		CheckError(err)

		SelectPath(conn) // Update the decision table before sending
		buffer := MakeBlinkPacket(SESSION_ID_TEST, oracleAddr, serverAddr, DecisionTableType, DecisionTableToByteArray(decisionTable))
		_, err = conn.WriteToUDP(buffer, addr)
		CheckError(err)
		fmt.Println("Sent DT to:", addr.String())
		// fmt.Println(DecisionTable)
	}
}

// UpdateRoutingTable : This function updates the GlobalRoutingTable based on information from receieved Routing Table
func UpdateRoutingTable(addr *net.UDPAddr, inBuf []byte) {
	// See if the routing table from this source exists
	for sourceAddress, routingTable := range oracleRoutingTable {
		// If it does exist.
		if sourceAddress == addr.String() {
			//Extract packet data
			packetData := ExtractPacketData(inBuf)
			//Convert Byte array to link
			newLink := ByteArrayToLink(packetData)
			//Store that link.
			for i, value := range routingTable {
				if value.Addr == newLink.Addr {
					routingTable[i] = &newLink
					oracleRoutingTable[sourceAddress] = routingTable
				}
			}
		}
	}
}

// SelectPath : function selects the best of two paths (based on lowest latency), updates the decision table accordingly. The oracle will then distribute this decision table to all of the nodes.
func SelectPath(conn *net.UDPConn) {
	latency1to2 := oracleRoutingTable[nodesTable[1]][0].Latency
	latency2to4 := oracleRoutingTable[nodesTable[3]][1].Latency
	latency1to3 := oracleRoutingTable[nodesTable[2]][0].Latency
	latency3to4 := oracleRoutingTable[nodesTable[3]][2].Latency
	// fmt.Println("1to2:", latency1to2, "2to4:", latency2to4)
	// fmt.Println("1to3:", latency1to3, "3to4:", latency3to4)

	latencyBlue := latency1to3 + latency3to4
	latencyRed := latency1to2 + latency2to4

	// fmt.Println("Red:", latencyRed, "Blue:", latencyBlue)

	bluePath := []string{nodesTable[0], nodesTable[2], nodesTable[3]}
	redPath := []string{nodesTable[0], nodesTable[1], nodesTable[3]}

	if latencyBlue <= latencyRed {
		// fmt.Println("Blue Path Selected")
		for key := range DecisionTable {
			fmt.Println("Key", key)
			DecisionTable[key] = bluePath
		}
	} else {
		// fmt.Println("Red Path Selected")
		for key := range DecisionTable {
			fmt.Println("SID:", key)
			DecisionTable[key] = redPath
		}
	}

	packetData := DecisionTableToByteArray(DecisionTable)
	blinkPacket := MakeBlinkPacket("0001000100010001", oracleAddr, oracleAddr, DecisionTableType, packetData)

	// Send the decision table to all the nodes
	for _, nodeAddrString := range nodesTable {
		addr, err := net.ResolveUDPAddr("udp", nodeAddrString)
		CheckError(err)

		_, err = conn.WriteToUDP(blinkPacket, addr)
		CheckError(err)
	}
}

// setPath : Sets the path for a given session in the decision table. Takes an inputted array of address strings and sets them as the decision table. It then distributes decision table to all Blink nodes.
func setPath(path []string, sessionID string) {
	decisionTable := make(map[string][]string)
	decisionTable[SESSION_ID_TEST] = path
	DecisionTable = decisionTable

	packetData := DecisionTableToByteArray(DecisionTable)
	blinkPacket := MakeBlinkPacket(SESSION_ID_TEST, oracleAddr, oracleAddr, DecisionTableType, packetData)

	// Send the decision table to all the nodes
	for _, nodeAddrString := range nodesTable {
		addr, err := net.ResolveUDPAddr("udp", nodeAddrString)
		CheckError(err)

		_, err = conn.WriteToUDP(blinkPacket, addr)
		CheckError(err)
	}
}

// formPath : given two Blink node addresses, forms any 2 or 3 node path between the two nodes. Assigns this path to the global pathsArray variable.
func formPath(nodeA string, nodeB string) {

	// Forms 3 paths from source(A) to destination(B)
	doesNodeAExist := false
	doesNodeBExist := false
	for i := 0; i < len(nodesTable); i++ {
		if nodesTable[i] == nodeA {
			doesNodeAExist = true
		}

		if nodesTable[i] == nodeB {
			doesNodeBExist = true
		}
	}

	if !doesNodeAExist || !doesNodeBExist {
		fmt.Println("Node A or B does not exist")
	}

	var bounceNode string
	pathIndex := 0
	for i := range pathsArray {
		pathsArray[i] = make([]string, 3)
	}

	for i := 0; i < len(nodesTable); i++ {
		/*
			If middle node is not equal to A then set middle node equal to bounce node
			This will result in 3 paths (assuming 4 total nodes)
			Example ( A -> C -> B) (A -> D -> B) (A -> B)
		*/
		if nodesTable[i] != nodeA && nodesTable[i] != nodeB {
			bounceNode = nodesTable[i]
			// Assign Items in the Array
			pathsArray[pathIndex][0] = nodeA
			pathsArray[pathIndex][1] = bounceNode
			pathsArray[pathIndex][2] = nodeB
			pathIndex += 1
		}
	}
}

// GetLatencyForPaths : calculates the measured latency of each path in pathsArray using oracleRoutingTable values. Stores in global latencyArray.
func GetLatencyForPaths() {
	/*
		Gets latency for each path and puts it in latencyArray at corresponding index
		The 3 paths will then be showed to the user
	*/
	latencyIndex := 0

	for i := 0; i < len(pathsArray); i++ {
		if pathsArray[i][1] == pathsArray[i][2] {
			latencyArray[latencyIndex] = -1
		} else {
			// Latency from the start of the path to the mid point node
			route1oracle := oracleRoute(pathsArray[i][0], pathsArray[i][1])

			routeStartToHalf := oracleRoutingTable[pathsArray[i][0]][route1oracle].Latency

			route2oracle := oracleRoute(pathsArray[i][1], pathsArray[i][2])

			// Latency from t he mid point node to the destination node
			routeHalfToEnd := 0.0
			if pathsArray[i][1] != pathsArray[i][2] {
				routeHalfToEnd = oracleRoutingTable[pathsArray[i][1]][route2oracle].Latency
			}

			// Total route latency by adding both delays together
			routeLatency := routeStartToHalf + routeHalfToEnd
			latencyArray[latencyIndex] = routeLatency
			latencyIndex++
		}
	}

	sortLatencyAndPathArrays()
}

// sortLatencyAndPathArrays : using global latencyArray, this function will sort the latencyArray (lowest->highest), sorting the pathsArray along with it to assure the latencyArray index always matches its corresponding pathsArray index.
func sortLatencyAndPathArrays() {
	for i := 0; i < len(pathsArray)-1; i++ {
		for j := 0; j < len(pathsArray)-1; j++ {
			if latencyArray[j] > latencyArray[j+1] {
				tmpLat := latencyArray[j]
				tmpPath := pathsArray[j]

				latencyArray[j] = latencyArray[j+1]
				pathsArray[j] = pathsArray[j+1]

				latencyArray[j+1] = tmpLat
				pathsArray[j+1] = tmpPath
			}
		}
	}

	// fmt.Println("After:", latencyArray)
	// fmt.Println("-------")
	// fmt.Println(pathsArray)
}

// oracleRoute : returns the index of the link in the oracleRoutingTable given source and destination
func oracleRoute(A string, B string) int {
	for j := 0; j < len(oracleRoutingTable[A]); j++ {
		if oracleRoutingTable[A][j].Addr == B {
			return j
		}
	}

	return -1
}

func tellClientToStart(nodeAddr string, filename string) {
	nodeIndex := 0
	for i := 0; i < len(nodesTable); i++ {
		if nodesTable[i] == nodeAddr {
			nodeIndex = i
		}
	}

	clientAddr, err := net.ResolveUDPAddr("udp", clientAddrString)
	CheckError(err)

	p := Packet{Number: nodeIndex, Data: []byte(filename)}
	buffer := PacketToJSONToByteArray(p)
	_, err = conn.WriteToUDP(buffer, clientAddr)
	fmt.Println("Sent:", filename, " to ", nodeAddr)
}

func handleBounceUpdate(inBuf []byte) {

	// Print bounce update
	packetData := ExtractPacketData(inBuf)
	sourceAddr := ExtractSrcAddr(inBuf)

	// Write Bounce update to file
	// file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	// CheckError(err)
	// printString := string(packetData) + " from " + sourceAddr.String() + "\n"
	// _, err = file.WriteString(printString)
	// file.Close()

	fmt.Println("Bounce Update:", string(packetData)+" from "+sourceAddr.String())
}

func saveLatency(latency float64, addr string) {
	// timeStamp := time.Now().String()

	// file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_WRONLY, 0600)
	// CheckError(err)

	// latencyString := strconv.FormatFloat(latency, 'f', 15, 64)
	// printString := latencyString + ", " + addr + ", " + timeStamp + ", \n"
	// _, err = file.WriteString(printString)
	// file.Close()
}

func handlePacketTrackUpdate(buffer []byte) {
	SID := ExtractSID(buffer)
	data := ExtractPacketData(buffer)
	packet := ByteArrayToJSONToPacket(data)

	if _, ok := PacketTrackTable[SID]; ok {
		// packetTrackMutex.Lock()
		_, ok2 := PacketTrackTable[SID][packet.Number]
		// packetTrackMutex.Unlock()

		if !ok2 {
			for i := packet.Number; i < packet.Number+1000; i++ {
				// packetTrackMutex.Lock()
				PacketTrackTable[SID][i] = Packet{Number: i}
				// packetTrackMutex.Unlock()
			}
		}
	} else {
		packetMap := make(map[int]Packet)
		// packetTrackMutex.Lock()
		PacketTrackTable[SID] = packetMap
		// packetTrackMutex.Unlock()

		for i := packet.Number; i < packet.Number+1000; i++ {
			// packetTrackMutex.Lock()
			PacketTrackTable[SID][i] = Packet{Number: i}
			// packetTrackMutex.Unlock()
		}
		// for i := packet.Number; i < packet.Number+500; i++ {
		// 	p := Packet{Number: i}
		// 	PacketTrackTable[SID][i] = &p
		// }
	}

	if !packet.Time1.IsZero() {
		// packetTrackMutex.Lock()
		p := Packet{
			Number: packet.Number,
			Time1:  packet.Time1,
			Time2:  PacketTrackTable[SID][packet.Number].Time2,
			Time3:  PacketTrackTable[SID][packet.Number].Time3,
		}
		PacketTrackTable[SID][packet.Number] = p
		// packetTrackMutex.Unlock()
	}

	if !packet.Time2.IsZero() {
		// packetTrackMutex.Lock()
		p := Packet{
			Number: packet.Number,
			Time1:  PacketTrackTable[SID][packet.Number].Time1,
			Time2:  packet.Time2,
			Time3:  PacketTrackTable[SID][packet.Number].Time3,
		}
		PacketTrackTable[SID][packet.Number] = p
		// packetTrackMutex.Unlock()
	}

	if !packet.Time3.IsZero() {
		// packetTrackMutex.Lock()
		p := Packet{
			Number: packet.Number,
			Time1:  PacketTrackTable[SID][packet.Number].Time1,
			Time2:  PacketTrackTable[SID][packet.Number].Time2,
			Time3:  packet.Time3,
		}

		PacketTrackTable[SID][packet.Number] = p
		// packetTrackMutex.Unlock()
	}

	// packetTrackMutex.Lock()
	p := PacketTrackTable[SID][packet.Number]
	if !p.Time1.IsZero() && !p.Time2.IsZero() && !p.Time3.IsZero() {
		go extractPacketTrackLatency(p)
	}
	// packetTrackMutex.Unlock()

}

// sends data from t1 to t3 to firebase
func firebase_T1_T3(latencyData map[string]float64) {

	client, err := app.Firestore(firebaseCtx)
	CheckError(err)

	_, err = client.Collection("Latency_Data").Doc("t1_t3").Set(context.Background(), latencyData)
	if err != nil {
		log.Fatalln(err)
	}

	// _, err = client.Collection("Latency_Data").Doc("t1_t3_Stored").Update(firebaseCtx, []firestore.Update{
	// 	{Path: strconv.FormatFloat(latencyData["packetNumber"], 'f', -1, 32), Value: latencyData["time"]},
	// })
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	defer client.Close()
}

// sends data from t1 to t2 to firebase
func firebase_T1_T2(latencyData map[string]float64) {

	client, err := app.Firestore(firebaseCtx)
	CheckError(err)

	_, err = client.Collection("Latency_Data").Doc("t1_t2").Set(context.Background(), latencyData)
	if err != nil {
		log.Fatalln(err)
	}

	// _, err = client.Collection("Latency_Data").Doc("t1_t2_Stored").Update(firebaseCtx, []firestore.Update{
	// 	{Path: strconv.FormatFloat(latencyData["packetNumber"], 'f', -1, 32), Value: latencyData["time"]},
	// })
	// CheckError(err)

	defer client.Close()
}

// sends data from t2 to t3 to firebase
func firebase_T2_T3(latencyData map[string]float64) {

	client, err := app.Firestore(firebaseCtx)
	CheckError(err)

	_, err = client.Collection("Latency_Data").Doc("t2_t3").Set(context.Background(), latencyData)
	if err != nil {
		log.Fatalln(err)
	}

	// _, err = client.Collection("Latency_Data").Doc("t2_t3_Stored").Update(firebaseCtx, []firestore.Update{
	// 	{Path: strconv.FormatFloat(latencyData["packetNumber"], 'f', -1, 32), Value: latencyData["time"]},
	// })
	// CheckError(err)

	defer client.Close()
}

// wipeDatabase : wipes the current dataset in the firebase databasse
func wipeDatabase() {
	latencyData := make(map[string]float64)
	latencyData["packetNumber"] = float64(0)
	latencyData["time"] = -1

	_, err := client.Collection("Latency_Data").Doc("Latency").Set(firebaseCtx, latencyData)
	CheckError(err)
}

func clearSession(SID string) {
	// Wipe the packet map data
	packetMap := make(map[int]Packet)
	packetTrackMutex.Lock()
	PacketTrackTable[SID] = packetMap
	packetTrackMutex.Unlock()

	for i := 0; i < 1000; i++ {
		packetTrackMutex.Lock()
		PacketTrackTable[SID][i] = Packet{Number: i}
		packetTrackMutex.Unlock()
	}
}
