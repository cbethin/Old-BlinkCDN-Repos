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
  loss int
  delay float64
}



func main() {

  if len(os.Args) != 2 {
    fmt.Println("Not enough arguments")
    os.Exit(1)
  }

  // Resolve addresses
  serverAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
  checkError(err)
  clientAddr, err := net.ResolveUDPAddr("udp", ":4000")

  // Set up connection
  conn, err := net.DialUDP("udp", clientAddr, serverAddr)
  checkError(err)

  var RoutingTable [10]link
  var link1 link
  link1.name = "c1"
  link1.delay = Ping(conn)
  link1.loss = ProbeLoss(conn)
  RoutingTable[0] = link1

  fmt.Println(RoutingTable)
}

func ProbeLoss(conn *net.UDPConn) int {
  buffer := make([]byte, 1024)
  copy(buffer[:], "10")
  _, err := conn.Write(buffer)
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
  return lossCount
}

func Ping(conn *net.UDPConn) float64 {
  // Setup ping buffer
  buffer := make([]byte, 1024)
  copy(buffer[:], "11")
  // Write buffer to connectiont
  _, err := conn.Write(buffer)
  timeSent := time.Now()
  checkError(err)
  // Get response
  _, _, err = conn.ReadFromUDP(buffer)
  timeTaken := time.Since(timeSent)
  timeTakenSec := timeTaken.Seconds()
  delay := timeTakenSec / 2
  checkError(err)

  //fmt.Println("Delay:", oneWayTime, "seconds")
  return delay
}

func TestLatency(conn *net.UDPConn) {
  latency := Ping(conn)
  fmt.Println("--,", latency)

  for {
    time.Sleep(time.Millisecond * 400)
    newDelay := Ping(conn)
    latency = 0.9 * latency + 0.1 * newDelay
    fmt.Println(newDelay, ",", latency)
  }
}

func checkError(err error) {
  if err != nil {
    fmt.Println("ERROR:", err)
    os.Exit(1)
  }
}
