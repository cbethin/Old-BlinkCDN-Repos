package blink

import (
  "fmt"
  "net"
  "os"
)

func PrintHello() {
  fmt.Println("Hey There how are you")
  return
}

func WrapHeader(srcAddr *net.UDPAddr, finalDstAddr *net.UDPAddr, packetType string, buf []byte) []byte {

  // This copies the source address and destination adress into bytes 0-32 and 32-64
  // respectively. The addresses are inputted as a string followed by a ! to let the program
  // know where the string ends. Bytes 64-66 are filled with the packet type as a string, and
  // the rest of the program are filled with the actual packet to send

  outBuf := make([]byte, 1024)
  copy(outBuf[:32], []byte(srcAddr.String() + "!"))
  copy(outBuf[32:64], []byte(finalDstAddr.String() + "!"))
  copy(outBuf[64:66], []byte(packetType))
  copy(outBuf[66:], buf)
  return outBuf
}

// Extract all information from the Blink Packet. Return Source Addr, Destination Addr,
// Packet Type, and Packet Data (in that order)
func UnwrapHeader(inBuf []byte) (*net.UDPAddr, *net.UDPAddr, string, []byte) {
  srcAddr := ExtractSrcAddr(inBuf)
  finalDstAddr := ExtractFinalDstAddr(inBuf)
  packetType := ExtractPacketType(inBuf)
  packetData := ExtractPacketData(inBuf)

  return srcAddr, finalDstAddr, packetType, packetData
}

// Extract Source Address from the header of the Blink Packet.

func ExtractSrcAddr(inBuf []byte) *net.UDPAddr {

  // Pull in the header bytes corresponding to the source address (0-32)
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

//Extract Destination Address from the header of the Blink Packet.
func ExtractFinalDstAddr(inBuf []byte) *net.UDPAddr {

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

// Extract the packet type from the Blink Packet header, return as type string
func ExtractPacketType(inBuf []byte) string {
  return string(inBuf[64:66])
}

// Extract Packet Data from the Blink Packet header
func ExtractPacketData(inBuf []byte) []byte {
  return inBuf[66:]
}

func CheckError(err error) {
  if err != nil {
    fmt.Println("Error:", err)
    os.Exit(1)
  }
}
