package main

import (
  "fmt"
  "net"
  "os"
  "strconv"
)

type link struct {
  name string
  loss float64
  latency float64
  address *net.UDPAddr
}

func main() {
  if len(os.Args) != 2 {
    fmt.Println("Not enough arguments")
    os.Exit(1)
  }

  // Set Up listener
  oracleAddr, err := net.ResolveUDPAddr("udp", ":"+os.Args[1])
  checkError(err)
  conn, err := net.ListenUDP("udp", oracleAddr)
  checkError(err)
  fmt.Println("Listening...")
  var clientAddr *net.UDPAddr

  buffer := make([]byte, 1024)
  var RoutingTable [2]link

  _, clientAddr, err = conn.ReadFromUDP(buffer)
  checkError(err)

  link := ByteArrayToLink(buffer[0:66])
  link2 := ByteArrayToLink(buffer[66:132])
  RoutingTable[0] = link
  RoutingTable[1] = link2


  serverToUse := MakeDecision(RoutingTable)
  copy(buffer[:], strconv.Itoa(serverToUse))
  _, err = conn.WriteToUDP(buffer, clientAddr)
  checkError(err)
}

func ByteArrayToLink(byteArray []byte) link {
  var newLink link
  newLink.name = string(byteArray[:4])
  lossString := string(byteArray[4:19])
  latencyString := string(byteArray[19:34])
  addressBytes := byteArray[34:66]
  addressString := ""

  for i:= 0; i < 32; i++ {
    if string(addressBytes[i]) != "!" {
      addressString += string(addressBytes[i])
    } else {
      break
    }
  }

  address, err := net.ResolveUDPAddr("udp", addressString)
  checkError(err)

  newLink.address = address
  newLink.loss, _ = strconv.ParseFloat(lossString, 64)
  newLink.latency, _ = strconv.ParseFloat(latencyString, 64)
  return newLink
}


func MakeDecision(RoutingTable [2]link) int {
  linkCost1 := 0.7 * RoutingTable[0].loss + 0.3 * RoutingTable[0].latency
  linkCost2 := 0.7 * RoutingTable[1].loss + 0.3 * RoutingTable[1].latency
  if linkCost1 <= linkCost2 {
    return 0
  } else {
    return 1
  }
}

func checkError(err error) {
  if err != nil {
    fmt.Println("ERROR:", err)
    os.Exit(1)
  }
}
