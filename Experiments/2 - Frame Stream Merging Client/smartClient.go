/*

*/

package main

import (
  "fmt"
  "net"
  "os"
  "strconv"
  "time"
)

func main() {

    if len(os.Args) != 3 {
      fmt.Println("Improper number of arguments.")
      os.Exit(0)
    }

    // Initialize Server/Client Addresses
    server1Addr, err := net.ResolveUDPAddr("udp", os.Args[1])
    checkError(err)
    server2Addr, err := net.ResolveUDPAddr("udp", os.Args[2])
    checkError(err)
    client1Addr, err := net.ResolveUDPAddr("udp", ":0")
    checkError(err)
    // client2Addr, err := net.ResolveUDPAddr("udp", ":1")
    // checkError(err)

    // Set up Connection Objects
    conn1, err := net.DialUDP("udp", client1Addr, server1Addr)
    checkError(err)
    conn2, err := net.DialUDP("udp", client1Addr, server2Addr)
    checkError(err)
    defer conn1.Close()
    defer conn2.Close()

    // Set up playout buffer, bitmask, and channel
    bitMask := make([]bool, 700)
    file, err := os.Create("playout.txt")
    checkError(err)
    defer file.Close()
    c := fanIn(handleServer(conn1), handleServer(conn2))

    for {
        receivedBuffer := <-c
        sequenceNumber, err := strconv.Atoi(string(receivedBuffer[:3]))
        receivedMessage := string(receivedBuffer[4:10])
        checkError(err)

        if bitMask[sequenceNumber] == false {
          bitMask[sequenceNumber] = true
          // Write to file
          currentTime := time.Now().String()
          fileLine := strconv.Itoa(sequenceNumber) + " " + receivedMessage + " " + currentTime + "\n"
          file.WriteString(fileLine)
          fmt.Println("Wrote:", fileLine)
        } else {
          fmt.Println("Discarded:", sequenceNumber)
        }
    }

    fmt.Println("You're both boring; I'm leaving.")

}

func fanIn(input1, input2 <-chan []byte) <-chan []byte {
    c := make(chan []byte)
    go func() { for { c <- <-input1 } }()
    go func() { for { c <- <-input2 } }()
    return c
}

func handleServer(conn *net.UDPConn) <-chan []byte { // Returns receive-only channel of strings.
    // Create channel (to pass values btwn functions) and buffer (for transmissions)
    c := make(chan []byte)
    buffer := make([]byte, 1024)

    // Initialize server messages by sending play message
    copy(buffer[:], "play100")
    _, err := conn.Write(buffer)
    checkError(err)

    buffer = make([]byte, 1024) //Reset Buffer

    go func() { // We launch the goroutine from inside the function.
        for i := 0; ; i++ {
            _, _, err = conn.ReadFromUDP(buffer)
            checkError(err)
            c <- buffer
        }
    }()
    return c // Return the channel to the caller.
}

func checkError(err error) {
  if err != nil {
    fmt.Println("ERROR:", err)
  }
}
