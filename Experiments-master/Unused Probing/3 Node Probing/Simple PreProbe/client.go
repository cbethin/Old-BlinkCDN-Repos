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
}

var clientAddr *net.UDPAddr

func main() {

  if len(os.Args) != 3 {
    fmt.Println("Not enough arguments")
    os.Exit(1)
  }

  // Resolve addresses
  serverAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
  checkError(err)
  serverAddr2, err := net.ResolveUDPAddr("udp", os.Args[2])
  checkError(err)
  clientAddr, err = net.ResolveUDPAddr("udp", ":4000")

  // Set up routing table
  var RoutingTable [2]link
  var link1 link
  link1.name = "S1"
  link1.latency = ProbeLatency(serverAddr)
  link1.loss = ProbeLoss(serverAddr)
  link1.address = serverAddr
  RoutingTable[0] = link1

  var link2 link
  link2.name = "S2"
  link2.latency = ProbeLatency(serverAddr2)
  link2.loss = ProbeLoss(serverAddr2)
  link2.address = serverAddr2
  RoutingTable[1] = link2

  fmt.Println(RoutingTable)
  serverToUse := decisionMaker(RoutingTable)
  fmt.Println("Server Chosen:", serverToUse)

  conn, err := net.DialUDP("udp", clientAddr, serverToUse)
  checkError(err)

  buffer := make([]byte, 1024)
  copy(buffer[:], "11")
  _, err = conn.Write(buffer)
  checkError(err)

  _, addr, err := conn.ReadFromUDP(buffer)
  checkError(err)
  fmt.Println("Received from:", addr)
}

func decisionMaker(RoutingTable [2]link) *net.UDPAddr {
  linkCost1 := 0.7 * RoutingTable[0].loss + 0.3 * RoutingTable[0].latency
  linkCost2 := 0.7 * RoutingTable[1].loss + 0.3 * RoutingTable[1].latency
  if linkCost1 <= linkCost2 {
    return RoutingTable[0].address
  } else {
    return RoutingTable[1].address
  }
}

func ProbeLoss(serverAddr *net.UDPAddr) float64 {
  // Set up connection
  conn, err := net.DialUDP("udp", clientAddr, serverAddr)
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
  conn, err := net.DialUDP("udp", clientAddr, serverAddr)
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

func checkError(err error) {
  if err != nil {
    fmt.Println("ERROR:", err)
    os.Exit(1)
  }
}
