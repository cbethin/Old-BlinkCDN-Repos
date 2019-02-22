package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"./blink"
)

var (
	clientAddr *net.UDPAddr
	// conn       *net.UDPConn
	node1Addr   *net.UDPAddr
	node2Addr   *net.UDPAddr
	node3Addr   *net.UDPAddr
	node4Addr   *net.UDPAddr
	currentConn *net.UDPConn
)

func main() {
	if len(os.Args) != 6 {
		fmt.Println("ERROR: Improper number of arguments")
		os.Exit(0)
	}

	// Set as UDP we can always change it to TCP if needed
	clientAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
	blink.CheckError(err)

	node1Addr, err = net.ResolveUDPAddr("udp", os.Args[2])
	blink.CheckError(err)

	node2Addr, err = net.ResolveUDPAddr("udp", os.Args[3])
	blink.CheckError(err)

	node3Addr, err = net.ResolveUDPAddr("udp", os.Args[4])
	blink.CheckError(err)

	node4Addr, err = net.ResolveUDPAddr("udp", os.Args[5])
	blink.CheckError(err)

	// // Open Listening Connection
	// conn, err := net.ListenUDP("udp", clientAddr)
	// blink.CheckError(err)
	// defer conn.Close()
	// currentConn = conn

	go startHttpListening()
	for {
		conn, err := net.ListenUDP("udp", clientAddr)
		blink.CheckError(err)
		defer conn.Close()
		currentConn = conn

		startListening(conn)
	}

}

func startListening(conn *net.UDPConn) {
	// Make Read Buffer
	shouldStop := false
	readBuffer := make([]byte, 1024)

	for {
		if shouldStop == true {
			return
		}

		// Read from Listening Connection
		_, _, err := conn.ReadFromUDP(readBuffer)
		if err != nil {
			shouldStop = true
			continue
		}
		// blink.CheckError(err)

		// If statement checks which node to send data to
		if string(readBuffer[0:5]) == "node0" {
			shouldStop = clientSendIt(conn, node1Addr)

		} else if string(readBuffer[0:5]) == "node1" {
			shouldStop = clientSendIt(conn, node2Addr)

		} else if string(readBuffer[0:5]) == "node2" {
			shouldStop = clientSendIt(conn, node3Addr)

		} else if string(readBuffer[0:5]) == "node3" {
			shouldStop = clientSendIt(conn, node4Addr)
		}
	}
}

func startHttpListening() {
	http.HandleFunc("/", handleHttpResponse)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}

func clientSendIt(conn *net.UDPConn, nodeAddr *net.UDPAddr) bool {

	// number of packets sent
	numberOfPackets := 1000

	packetNumber := 0
	nullAddr, err := net.ResolveUDPAddr("udp", ":8080")
	blink.CheckError(err)

	for i := packetNumber + 1; i < numberOfPackets; i++ {
		p := blink.Packet{Number: i}
		sendBuffer := blink.PacketToJSONToByteArray(p)

		blinkPacket := blink.MakeBlinkPacket("0001000100010001", clientAddr, nullAddr, blink.Iframe, sendBuffer)

		_, err := conn.WriteToUDP(blinkPacket, nodeAddr)
		if err != nil {
			return true
		}

		blink.CheckError(err)
		fmt.Println(p, " ")

		time.Sleep(time.Millisecond * 1000)
	}

	return false
}

func restartConnection() {
	// Closes connection and then re-opens it

	currentConn.Close()
	fmt.Println("Connection Closed from Close Connection Function")
}

func handleHttpResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	if r.Method == "POST" {
		if r.URL.Path == "/stoptest" {
			fmt.Println("Stopping test")
			restartConnection()
			w.Write([]byte("success"))
		}
	}
}
