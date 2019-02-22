package blink

import "net"

// MakeBlinkPacket : Creates a Blink Packet, which is simply an array of bytes, where the first 66 bytes contain information about the source address, destination address, type of packet being sent. The remaining 958 bytes are used for data transmission
func MakeBlinkPacket(SID string, srcAddr *net.UDPAddr, finalDestAddr *net.UDPAddr, packetType string, buf []byte) []byte {

	// This copies the src address and destination adress into bytes 0-32 and 32-64
	// respectively. The addresses are inputted as a string followed by a ! to let the program
	// know where the string ends. Bytes 64-66 are filled with the packet type as a string, and
	// the rest of the program are filled with the actual packet to send

	outBuf := make([]byte, 1024)
	copy(outBuf[:16], []byte(SID+"!"))
	copy(outBuf[16:48], []byte(srcAddr.String()+"!"))
	copy(outBuf[48:80], []byte(finalDestAddr.String()+"!"))
	copy(outBuf[80:82], []byte(packetType))
	copy(outBuf[82:], buf)
	return outBuf
}

// UnwrapHeader : Extract all information from the Blink Packet. Return src Addr, Destination Addr, Packet Type, and Packet Data (in that order)
func UnwrapHeader(inBuf []byte) (string, *net.UDPAddr, *net.UDPAddr, string, []byte) {
	SID := ExtractSID(inBuf)
	srcAddr := ExtractSrcAddr(inBuf)
	finalDestAddr := ExtractFinalDestAddr(inBuf)
	packetType := ExtractPacketType(inBuf)
	packetData := ExtractPacketData(inBuf)

	return SID, srcAddr, finalDestAddr, packetType, packetData
}

// ExtractSID : Extract SID from the header of the Blink Packet. Returns SID as a string
func ExtractSID(inBuf []byte) string {
	SID := string(inBuf[:16])
	return SID
}

// ExtractSrcAddr : Extract Source Address from the header of the Blink Packet. Returns address as a pointer to a resolved net.UDPAddr
func ExtractSrcAddr(inBuf []byte) *net.UDPAddr {

	// Pull in the header bytes corresponding to the src address (0-32)
	addrBuf := inBuf[16:48]
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

// ExtractFinalDestAddr : Extract Destination Address from the header of the Blink Packet. Returns address as a pointer to a resolved net.UDPAddr
func ExtractFinalDestAddr(inBuf []byte) *net.UDPAddr {

	// Pull in the header bytes corresponding to the destination address (32-64)

	addrBuf := inBuf[48:80]
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

// ExtractPacketType : From an inputted blink packet (an array of bytes), extract the packet type from the header. Returns value as a string type
func ExtractPacketType(inBuf []byte) string {
	return string(inBuf[80:82])
}

// ExtractPacketData : From an inputted blink packet (an array of bytes), extract the packet's data from the packet
func ExtractPacketData(inBuf []byte) []byte {
	return inBuf[82:]
}
