package blink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Packet : structure containing information about a given packet track
type Packet struct {
	Number int
	Time1  time.Time
	Time2  time.Time
	Time3  time.Time
	Data   []byte
}

// PacketToJSONToByteArray : converts a Packet type to a JSON object encoded and returned as a byteArray
func PacketToJSONToByteArray(p Packet) []byte {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(p)
	return b.Bytes()
}

// ByteArrayToJSONToPacket : converts a []byte and decodes it as a JSON following the Packet type structure
func ByteArrayToJSONToPacket(b []byte) Packet {
	var p Packet
	err := json.NewDecoder(bytes.NewReader(b)).Decode(&p)
	CheckError(err)
	return p
}

// extractPacketTrackLatency : from a given Packet, subtracts Time2 from Time1 and saves latency to database
func extractPacketTrackLatency(p Packet) { // SET BACK TO MAP
	// Takes t1 and t2 and finds t total then sends that, packet number and service time to function firebaseData
	t1 := p.Time1
	t2 := p.Time2
	t3 := p.Time3
	pktNum := p.Number

	// Total Trip time
	totalTripTime := t3.Sub(t1)
	if totalTripTime < 0 {
		fmt.Println("Error in trip time calculation:", p)
		return
	}

	latencyData := make(map[string]float64)
	latencyData["packetNumber"] = float64(pktNum)
	latencyData["time"] = totalTripTime.Seconds()

	fmt.Println("Updated packet", pktNum, " | Latency:", latencyData["time"])
	// firebase_T1_T3(latencyData)

	// Hop 1
	t1_t2 := t2.Sub(t1)
	if t1_t2 < 0 {
		fmt.Println("Error in trip time calculation:", p)
	}

	latencyData = make(map[string]float64)
	latencyData["packetNumber"] = float64(pktNum)
	latencyData["time"] = t1_t2.Seconds()

	fmt.Println("Updated packet", pktNum, " | Latency:", latencyData["time"])
	// firebase_T1_T2(latencyData)

	// Hop 2
	t2_t3 := t3.Sub(t2)
	if t2_t3 < 0 {
		fmt.Println("Error in trip time calculation:", p)
	}

	latencyData = make(map[string]float64)
	latencyData["packetNumber"] = float64(pktNum)
	latencyData["time"] = t2_t3.Seconds()

	fmt.Println("Updated packet", pktNum, " | Latency:", latencyData["time"])
	// firebase_T2_T3(latencyData)
}
