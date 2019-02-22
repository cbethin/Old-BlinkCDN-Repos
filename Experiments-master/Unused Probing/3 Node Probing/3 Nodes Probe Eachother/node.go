package main

import (
  "fmt"
  "net"
  "os"
  "time"
  "strconv"
)

type link struct {
  name string
  loss float64
  latency float64
  address *net.UDPAddr
  sourceAddr *net.UDPAddr
}

var clientDialAddr *net.UDPAddr
var clientListenAddr *net.UDPAddr

func main() {

  if len(os.Args) != 6 {
    fmt.Println("Not enough arguments")
    os.Exit(1)
  }

  // Resolve addresses
  clientDialAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
  clientListenAddr, err := net.ResolveUDPAddr("udp", os.Args[2])
  nodeAddr1, err := net.ResolveUDPAddr("udp", os.Args[3])
  checkError(err)
  nodeAddr2, err := net.ResolveUDPAddr("udp", os.Args[4])
  checkError(err)
  oracleAddr, err := net.ResolveUDPAddr("udp", os.Args[5])

  // Set up listener
  listener, err := net.ListenUDP("udp", clientListenAddr)
  checkError(err)
  go listenForProbe(listener)

  // Set up routing table
  var RoutingTable [3]link
  time.Sleep(10*time.Second)
  var link1 link
  link1.name = "S1"
  link1.latency = ProbeLatency(nodeAddr1)
  link1.loss = ProbeLoss(nodeAddr1)
  link1.address = nodeAddr1
  link1.sourceAddr = clientListenAddr
  RoutingTable[0] = link1

  var link2 link
  link2.name = "S2"
  link2.latency = ProbeLatency(nodeAddr2)
  link2.loss = ProbeLoss(nodeAddr2)
  link2.address = nodeAddr2
  link2.sourceAddr = clientListenAddr
  RoutingTable[1] = link2

  buffer := make([]byte, 1024)

  fmt.Println(RoutingTable)
  conn, err := net.DialUDP("udp", clientDialAddr, oracleAddr)
  checkError(err)
  buf1 := LinkToByteArray(link1)
  buf2 := LinkToByteArray(link2)
  copy(buffer[0:98], buf1)
  copy(buffer[98:196], buf2)
  _, err = conn.Write(buffer)
  checkError(err)
  fmt.Println("Routing Table Sent")

  // // Wait for decision
  // _, _, err = conn.ReadFromUDP(buffer)
  // checkError(err)
  // conn.Close()
  //
  // indexToUse, err := strconv.Atoi(string(buffer[0]))
  // checkError(err)
  // fmt.Println("Received:", indexToUse)

  time.Sleep(25*time.Second)
  listener.Close()
}

func LinkToByteArray(link link) []byte {
  byteArray := make([]byte, 192)
  copy(byteArray[:4], link.name)
  copy(byteArray[4:19], strconv.FormatFloat(link.loss, 'f', 15, 64))
  copy(byteArray[19:34], strconv.FormatFloat(link.latency, 'f', 15, 64))
  copy(byteArray[34:66], link.address.String() + "!")
  copy(byteArray[66:98], link.sourceAddr.String() + "!")
  fmt.Println(string(byteArray))
  return byteArray
}

func ProbeLoss(serverAddr *net.UDPAddr) float64 {
  // Set up connection
  conn, err := net.DialUDP("udp", clientDialAddr, serverAddr)
  checkError(err)

  buffer := make([]byte, 1024)
  copy(buffer[:], "10")
  _, err = conn.Write(buffer)
  checkError(err)


  lossCount := 0
  receivedNumber := 0
  expectedNumber := 100

  for expectedNumber != 199 {
    // Set connection so it times out after 1 second
    currentTime := time.Now()
    deadline := currentTime.Add(1000*time.Millisecond)
    conn.SetReadDeadline(deadline)

    //read from connection, check for error
    _, _, err := conn.ReadFromUDP(buffer)
    if err != nil {
      restLost := 199 - receivedNumber
      lossCount += restLost
      break
    }

    //Pull number from the message
    receivedNumber, err = strconv.Atoi(string(buffer[0:3]))
    checkError(err)

    if receivedNumber != expectedNumber {
      //Add 1 to the error count
      lossCount += (receivedNumber - expectedNumber)
    }

    expectedNumber = expectedNumber + (receivedNumber - expectedNumber) + 1
  }

  fmt.Println("Loss:", lossCount)
  conn.Close()
  return float64(lossCount / 100.)
}

func ProbeLatency(serverAddr *net.UDPAddr) float64 {
  // Set up connection
  conn, err := net.DialUDP("udp", clientDialAddr, serverAddr)
  checkError(err)

  // Setup ping buffer
  buffer := make([]byte, 1024)
  copy(buffer[:], "11")
  // Write buffer to connectiont
  _, err = conn.Write(buffer)
  timeSent := time.Now()
  checkError(err)
  // Get response
  _, _, err = conn.ReadFromUDP(buffer)
  timeTaken := time.Since(timeSent)
  timeTakenSec := timeTaken.Seconds()
  delay := timeTakenSec / 2
  checkError(err)

  conn.Close()
  return delay
}

func listenForProbe(conn *net.UDPConn) {
  buffer := make([]byte, 1024)

  for {
    _, addr, err := conn.ReadFromUDP(buffer)
    checkError(err)

    if string(buffer[0:2]) == "11" {
      copy(buffer[:], "11")
      _, err := conn.WriteToUDP(buffer, addr)
      checkError(err)
    } else if string(buffer[0:2]) == "10" {
      for i := 100; i < 200; i++ {
        copy(buffer[:], strconv.Itoa(i))
        _, err := conn.WriteToUDP(buffer, addr)
        checkError(err)
      }
    }
  }
}

func checkError(err error) {
  if err != nil {
    fmt.Println("ERROR:", err)
    os.Exit(1)
  }
}
