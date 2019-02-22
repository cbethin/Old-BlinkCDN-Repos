package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"./blink"
)

var (
	clientAddr   *net.UDPAddr
	conn         *net.UDPConn
	currentConn  *net.UDPConn
	serverAddr   *net.UDPAddr
	finalBuffers = make(map[string][]byte)
)

const (
// SESSION_ID_TEST = "1001100110011001"
// SERVER_ADDR = "155.246.45.33:8080"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: Improper number of arguments")
		os.Exit(0)
	}

	// Set as UDP we can always change it to TCP if needed
	clientAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
	blink.CheckError(err)

	// // Open Listening Connection
	// conn, err := net.ListenUDP("udp", clientAddr)
	// blink.CheckError(err)
	// defer conn.Close()
	// currentConn = conn

	conn, err := net.ListenUDP("udp", clientAddr)
	blink.CheckError(err)
	defer conn.Close()
	currentConn = conn

	// Make Read Buffer
	readBuffer := make([]byte, blink.MaxPacketSize)

	for {
		// Read from connection, then extract packet type using Blink's built in function
		_, _, err := conn.ReadFromUDP(readBuffer)
		blink.CheckError(err)
		addFile(readBuffer)
	}
}

func addFile(buf []byte) {
	sid := blink.ExtractSID(buf)
	packetDataAsBuffer := blink.ByteArray(blink.ExtractPacketData(buf))
	packetAsPacket := packetDataAsBuffer.ToPacket()
	filepath := "server/" + packetAsPacket.Filename
	// packetAsPacket := blink.ByteArrayToJSONToPacket(packetDataAsBuffer)
	// packetNumber := packetAsPacket.Number
	videoBuffer := packetAsPacket.Data
	var fileHasEnded bool

	fileHasEnded = (string(videoBuffer) == "TRANSFER COMPLETE")
	if !fileHasEnded {
		if len(finalBuffers[sid]) < 1 {
			finalBuffers[sid] = make([]byte, 0)
		}

		finalBuffers[sid] = append(finalBuffers[sid], videoBuffer...)
	} else {
		dirPath := ""
		foldersPaths := strings.Split(filepath, "/")
		for i := 0; i < len(foldersPaths)-1; i++ {
			dirPath += foldersPaths[i] + "/"
		}

		if _, err := os.Stat(filepath); err != nil {
			os.MkdirAll(dirPath, os.ModePerm)
			file, err := os.Create(filepath)
			blink.CheckError(err)
			file.Close()
		}

		file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0666)
		blink.CheckError(err)
		_, err = file.WriteString(string(finalBuffers[sid]))
		file.Close()
		fmt.Println("DONE Writing File:", filepath)
		finalBuffers[sid] = make([]byte, 0)
	}

}
