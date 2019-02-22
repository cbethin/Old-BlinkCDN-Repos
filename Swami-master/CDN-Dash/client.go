package main

import (
	"fmt"
	"io/ioutil"
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
	serverAddr  *net.UDPAddr
)

const (
	SESSION_ID_TEST = "1001100110011001"
	// SERVER_ADDR     = "155.246.45.43:9000"
	SERVER_ADDR = "18.221.128.102"
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

	serverAddr, err = net.ResolveUDPAddr("udp", SERVER_ADDR)
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
	readBuffer := make([]byte, blink.MaxPacketSize)

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

		p := blink.ByteArrayToJSONToPacket(readBuffer)
		packetNumber := p.Number
		filename := string(p.Data)
		fmt.Println(packetNumber, ' ', filename)

		// If statement checks which node to send data to
		switch packetNumber {
		case 0:
			fmt.Println("Sending to:", node1Addr.String())
			shouldStop = clientSendIt(conn, node1Addr, filename)
		case 1:
			fmt.Println("Sending to:", node2Addr.String())
			shouldStop = clientSendIt(conn, node2Addr, filename)
		case 2:
			fmt.Println("Sending to:", node3Addr.String())
			shouldStop = clientSendIt(conn, node3Addr, filename)
		case 3:
			fmt.Println("Sending to:", node4Addr.String())
			shouldStop = clientSendIt(conn, node4Addr, filename)
		default:
			fmt.Println("Client not found")
		}
	}
}

func startHttpListening() {
	http.HandleFunc("/", handleHttpResponse)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}

func clientSendIt(conn *net.UDPConn, nodeAddr *net.UDPAddr, filename string) bool {

	// nullAddr, err := net.ResolveUDPAddr("udp", ":0000")
	// blink.CheckError(err)

	// SEND VIDEO
	// filename := "video/test.mp4"

	if _, err := os.Stat(filename); err == nil {

		videoBuffer, err := ioutil.ReadFile(filename)
		blink.CheckError(err)

		var textBuffer []byte

		stepSize := blink.MaxPacketSize - blink.BlinkPacketSize - 450
		outputBuf := make([]byte, 0)

		var blinkPacket []byte
		packetCount := 0
		for i := 0; i < len(videoBuffer); i += stepSize {
			if i+stepSize+1 < len(videoBuffer) {
				outputBuf = append(outputBuf, videoBuffer[i:i+stepSize]...)
			} else {
				// finalValue := len(videoBuffer) - 1
				outputBuf = append(outputBuf, videoBuffer[i:]...)
				// outputBuf = fillByteArray(videoBuffer[i:finalValue], stepSize)
			}

			for j := 0; j < len(outputBuf); j++ {
				textBuffer = append(textBuffer, outputBuf[j])
			}

			// Make packet out of video outputBuf
			outputPacket := blink.Packet{Number: packetCount, Data: outputBuf, Filename: filename}
			outputPacketBuffer := blink.PacketToJSONToByteArray(outputPacket)
			blinkPacket = blink.MakeBlinkPacket(SESSION_ID_TEST, clientAddr, serverAddr, blink.Iframe, outputPacketBuffer)

			_, err = conn.WriteToUDP(blinkPacket, nodeAddr)
			if err != nil {
				return true
			}

			packetCount++

			fmt.Println(outputPacket.Number)
			// fmt.Println("Buffer:", outputPacket.Number)
			// fmt.Println("Size:", len(outputPacketBuffer))
			time.Sleep(40 * time.Millisecond)

			outputBuf = make([]byte, 0)
		}

		// Send transfer complete message
		outputBuf = append(outputBuf, []byte("TRANSFER COMPLETE")...)
		outputPacket := blink.Packet{Number: packetCount, Data: outputBuf, Filename: filename}
		outputPacketBuffer := blink.PacketToJSONToByteArray(outputPacket)
		blinkPacket = blink.MakeBlinkPacket(SESSION_ID_TEST, clientAddr, serverAddr, blink.Iframe, outputPacketBuffer)
		_, err = conn.WriteToUDP(blinkPacket, nodeAddr)
		if err != nil {
			return true
		}

		fmt.Println("Length:", len(videoBuffer))
		fmt.Println("Amount Extra:", len(videoBuffer)%(blink.MaxPacketSize-blink.BlinkPacketSize))
		fmt.Println("Packets Sent", packetCount)

		_ = ioutil.WriteFile("test-og.txt", textBuffer, os.ModeAppend)
		fmt.Println("DONE Writing File")
		return false
	} else {
		fmt.Println("NOT HERE:", filename)
		return true
	}
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

func fillByteArray(buffer []byte, toLength int) []byte {
	outBuf := make([]byte, toLength)
	copy(outBuf[:], buffer[:])
	for i := len(buffer); i < toLength; i++ {
		outBuf[i] = ":"[0]
	}

	return outBuf
}
