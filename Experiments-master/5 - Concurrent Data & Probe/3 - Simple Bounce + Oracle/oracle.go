package main

import (
  "fmt"
  "net"
  "os"
  "strconv"
  "time"
  "./blink"
)

type link struct {
	addr            string
	latency         float64
	loss            int
}

func main() {
  if len(os.Args) != 2 {
    fmt.Println("Not enough arguments")
    os.Exit(1)
  }

  // Set Up listening connection
  oracleAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
  blink.CheckError(err)

  conn, err := net.ListenUDP("udp", oracleAddr)
  blink.CheckError(err)
  fmt.Println("Listening...")

  // Set up Read Buffer
  buffer := make([]byte, 1024)
  //var RoutingTable [2]link

  /* Read incoming data, extract the packet data from the received Blink Packet
     and then print out the information in the received routing table
  */
  for {
    _, addr, err := conn.ReadFromUDP(buffer)
    blink.CheckError(err)
    packetData := blink.ExtractPacketData(buffer)
    link := ByteArrayToLink(packetData)
    fmt.Println("----------\n", link, "from", addr.String())
    fmt.Println("Latency:", link.latency)
    fmt.Println("Received at:", time.Now().String())
  }
}




/* Converts a byte array into a link. The byte array must be formatted properly
   using the blink LinkToByteArray() function.
*/
func ByteArrayToLink(byteArray []byte) link {
  var newLink link
  latencyString := string(byteArray[0:16])
  addressBytes := byteArray[16:48]
  addressString := ""

  for i:= 0; i < 32; i++ {
    if string(addressBytes[i]) != "!" {
      addressString += string(addressBytes[i])
    } else {
      break
    }
  }

  address, err := net.ResolveUDPAddr("udp", addressString)
  blink.CheckError(err)

  newLink.addr = address.String()
  newLink.latency, _ = strconv.ParseFloat(latencyString, 64)
  return newLink
}
