package main

import (
  "fmt"
  "os"
  "./blink"
  "strconv"
  "net"
  "time"
)

func main()  {

  /* 4 Arguments (3 plus Arg[0]) needed to run the program (1) local address (2) blink node address
     (3) destination address */
  if len(os.Args) != 4 {
    fmt.Println("Incorrect number of arguments!!!")
    os.Exit(1)
  }

  /* Set up all needed addresses. The information will be sent from this node (localAddr)
     to the blinkNode (blinkNodeAddr), which will then forward the packet to the final
     destination (destAddr).
  */
  localAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
  blink.CheckError(err)

  oracleAddr, err := net.ResolveUDPAddr("udp", os.Args[2])
  blink.CheckError(err)

  destAddr, err := net.ResolveUDPAddr("udp", os.Args[3])
  blink.CheckError(err)

  conn, err := net.ListenUDP("udp", localAddr)
  blink.CheckError(err)

  // Setup With Oracle
  buffer := blink.MakeBlinkPacket("0000000100000000", localAddr, destAddr, blink.HelloOracle, []byte(""))
  _, err = conn.WriteToUDP(buffer, oracleAddr)
  blink.CheckError(err)
  fmt.Println("Sent to Oracle")

  // Read Oracle's Response
  _, _, err = conn.ReadFromUDP(buffer)
  blink.CheckError(err)

  packetData := blink.ExtractPacketData(buffer)
  fmt.Println("Received From Oracle", string(packetData))

  SID := string(packetData[:16])
  var blinkAddrString string

  for i:=16; i <=1024; i++ {
    if string(packetData[i]) != "!" {
      blinkAddrString += string(packetData[i])
    } else {
      break
    }
  }

  fmt.Println("Connected to Oracle. SID:", SID, "Addr:", blinkAddrString)

  blinkAddr, err := net.ResolveUDPAddr("udp", blinkAddrString)
  blink.CheckError(err)

  /* Loops through from 0 - 10,000 sending each number to the blinkNode. Each packet
     is wrapped in a blink header, which will be used by the blink node to properly
     forward the packet
  */

  time.Sleep(3*time.Second)
  for i := 1; i < 10000; i++ {
    buffer := blink.MakeBlinkPacket(SID, localAddr, destAddr, blink.Iframe, []byte(strconv.Itoa(i)))
    _, err := conn.WriteToUDP(buffer, blinkAddr)
    blink.CheckError(err)
    fmt.Println("Sent:", i)
    time.Sleep(33333 * time.Microsecond)
  }
}
