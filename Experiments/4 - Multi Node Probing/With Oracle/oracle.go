package main

import (
  "fmt"
  "net"
  "os"
  "strconv"
  "time"
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

  // Set Up listener
  oracleAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
  checkError(err)

  conn, err := net.ListenUDP("udp", oracleAddr)
  checkError(err)
  fmt.Println("Listening...")

  buffer := make([]byte, 1024)
  //var RoutingTable [2]link

  for {
    _, _, err = conn.ReadFromUDP(buffer)
    checkError(err)
    fmt.Println("----------\n", ByteArrayToLink(buffer), "\n", time.Now().String())
  }
}

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
  checkError(err)

  newLink.addr = address.String()
  newLink.latency, _ = strconv.ParseFloat(latencyString, 64)
  return newLink
}


func checkError(err error) {
  if err != nil {
    fmt.Println("ERROR:", err)
    os.Exit(1)
  }
}
