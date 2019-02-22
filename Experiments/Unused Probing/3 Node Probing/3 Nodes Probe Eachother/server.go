package main

import (
  "fmt"
  "net"
  "os"
  "strconv"
)

func main() {
  if len(os.Args) != 2 {
    fmt.Println("Not enough arguments")
    os.Exit(1)
  }

  // Resovle addresses
  serverAddr, err := net.ResolveUDPAddr("udp", ":"+os.Args[1])
  checkError(err)

  // Set up listener
  conn, err := net.ListenUDP("udp", serverAddr)
  checkError(err)

  // Listen and respond
  buffer := make([]byte, 1024)
  fmt.Println("Listening...")

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
    fmt.Println("Error:", err)
    os.Exit(1)
  }
}
