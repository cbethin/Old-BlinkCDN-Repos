package blink

import (
	"fmt"
	"net"
	"time"

	"./blinkEstimators"
)

// LatencyEstimator : a funciton taking in an old latency estimate (float64) and new latency (float64), and producing a new latency estimate (float64)
type LatencyEstimator func(float64, float64) float64

// SendInitialProbe : Sends the initial probe packet to another node via a link, which is passed into the function.
func SendInitialProbe(link1 *Link) {
	// Update the link expiration as the first thing
	link1.TTLExpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)

	nodeAddr, err := net.ResolveUDPAddr("udp", link1.Addr)
	CheckError(err)

	// Create a Buffer now
	buffer := MakeBlinkPacket(ProbeID, localAddr, nodeAddr, InitialProbe, []byte(""))

	_, err = conn.WriteToUDP(buffer, nodeAddr)
	CheckError(err)

	link1.TimerStartTime = time.Now()
	link1.TimeoutDeadline = link1.TimerStartTime.Add(time.Millisecond * PROBE_TIMEOUT)
}

// SendResponse1 : Sends the Response 1 packet to another node via a link, which is passed into the function. */
func SendResponse1(link1 *Link) {
	link1.TTLExpiration = time.Now().Add(ROUTE_TTL * time.Millisecond)

	nodeAddr, err := net.ResolveUDPAddr("udp", link1.Addr)
	CheckError(err)

	// Create a Buffer now
	buffer := MakeBlinkPacket(ProbeID, localAddr, nodeAddr, ResponseOne, []byte(""))

	_, err = conn.WriteToUDP(buffer, nodeAddr)
	CheckError(err)

	link1.TimerStartTime = time.Now()
	link1.TimeoutDeadline = link1.TimerStartTime.Add(time.Millisecond * PROBE_TIMEOUT)
}

// SendResponse2 : Sends the Response 2 packet to another node via a link, which is passed into the function. Will also calculate and store the calculated Latency of the link. Sends routing table to oracle when it's done calculating
func SendResponse2(link1 *Link) {

	// Calculate time elapsed since the probe was sent
	timeElapsed := time.Since(link1.TimerStartTime)

	// Resolve other node's address, create a Blink Response 2 packet
	// and send that buffer over the inputted connection object
	nodeAddr, err := net.ResolveUDPAddr("udp", link1.Addr)
	CheckError(err)

	buffer := MakeBlinkPacket(ProbeID, localAddr, nodeAddr, ResponseTwo, []byte(""))

	_, err = conn.WriteToUDP(buffer, nodeAddr)
	CheckError(err)

	// Calculates Latency
	// link1.Latency = 0.9*link1.Latency + 0.1*((timeElapsed/2).Seconds()*1000)
	link1.Latency = createLatencyEstimate(link1.Latency, (timeElapsed/2).Seconds()*1000)

	// Send Routing Table To Oracle
	for _, value := range RoutingTable {
		// Create a blink packet for
		buffer := MakeBlinkPacket(ProbeID, localAddr, oracleAddr, RoutingTableType, LinkToByteArray(*value))

		_, err := conn.WriteToUDP(buffer, oracleAddr)
		CheckError(err)
	}

	// Save time sent
	saveLatency((timeElapsed/2).Seconds()*1000, link1.Addr)
}

// ReceiveResponse2 : Calculates the latency after a ResponseTwo is received and stores the latency in the corresponding link. Sends routing table to oracle when it's done calculating
func ReceiveResponse2(link1 *Link) {

	// Takes time elapsed since probe was sent
	timeElapsed := time.Since(link1.TimerStartTime)

	// Calculates Latency
	link1.Latency = 0.9*link1.Latency + 0.1*((timeElapsed/2).Seconds()*1000)

	// Send Routing Table To Oracle
	for _, value := range RoutingTable {
		// Create a blink packet to send the link to the oracle.
		buffer := MakeBlinkPacket(ProbeID, localAddr, oracleAddr, RoutingTableType, LinkToByteArray(*value))
		// Send the new Blink Packet to the oracle.
		_, err := conn.WriteToUDP(buffer, oracleAddr)
		CheckError(err)

	}

	// Save time recieved
	saveLatency((timeElapsed/2).Seconds()*1000, link1.Addr)
}

// CheckExpiration : Every 100 milliseconds the program will check the TTLExpiration of the designated link. If the link has expired, the program will initialize the probing process by sending an initial probe to the corresponding node
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

func createLatencyEstimate(latency float64, timeElapsed float64) float64 {
	if EstimatorFunction == nil {
		return blinkEstimators.WeightedAverageEstimator(latency, timeElapsed)
	}

	return EstimatorFunction(latency, timeElapsed)
}
