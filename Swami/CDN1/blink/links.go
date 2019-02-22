package blink

import (
	"strconv"
	"time"
)

// Link : structure containing information about a give connection between two Blink nodes
type Link struct {
	Addr            string
	Latency         float64
	Loss            int
	TTLExpiration   time.Time // When was the last time I probed you (Creepy!)
	TimeoutDeadline time.Time
	TimerStartTime  time.Time
}

// LinkToByteArray : Converts a link into a byte. The byte array is formatted to allow conversion using the ByteArrayToLink format
func LinkToByteArray(link Link) []byte {
	byteArray := make([]byte, 192)
	copy(byteArray[0:16], strconv.FormatFloat(link.Latency, 'f', 15, 64))
	copy(byteArray[16:48], link.Addr+"!")
	return byteArray
}

// ByteArrayToLink : Converts a byte array into a link. The byte array must be formatted properly using the blink LinkToByteArray() function.
func ByteArrayToLink(byteArray []byte) Link {
	var newLink Link
	latencyString := string(byteArray[0:16])
	addressBytes := byteArray[16:48]
	addressString := ""

	for i := 0; i < 32; i++ {
		if string(addressBytes[i]) != "!" {
			addressString += string(addressBytes[i])
		} else {
			break
		}
	}

	newLink.Addr = addressString
	newLink.Latency, _ = strconv.ParseFloat(latencyString, 64)
	return newLink
}
