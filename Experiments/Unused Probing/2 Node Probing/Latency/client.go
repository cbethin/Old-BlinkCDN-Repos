package main

import (
  "fmt"
  "net"
  "os"
  "time"
)

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

  TestLatency(conn)

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
    fmt.Println("Error:", err)
    os.Exit(1)
  }
}
