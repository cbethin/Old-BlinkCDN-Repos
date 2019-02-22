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
  if len(os.Args) != 4{
    fmt.Println("Incorrect number of arguments!!!")
    os.Exit(1)
  }

  /* Set up all needed addresses. The information will be sent from this node (localAddr)
     to the blinkNode (blinkNodeAddr), which will then forward the packet to the final
     destination (destAddr).
  */
  localAddr, err := net.ResolveUDPAddr("udp", os.Args[1])
  blink.CheckError(err)

  blinkNodeAddr, err := net.ResolveUDPAddr("udp", os.Args[2])
  blink.CheckError(err)

  destAddr, err := net.ResolveUDPAddr("udp", os.Args[3])
  blink.CheckError(err)

  conn, err := net.ListenUDP("udp", localAddr)
  blink.CheckError(err)

  /* Loops through from 0 - 10,000 sending each number to the blinkNode. Each packet
     is wrapped in a blink header, which will be used by the blink node to properly
     forward the packet
  */
  for i := 0; i < 10000; i++ {
    buffer := blink.MakeBlinkPacket(localAddr, destAddr, blink.Iframe, []byte(strconv.Itoa(i)))
    _, err := conn.WriteToUDP(buffer, blinkNodeAddr)
    blink.CheckError(err)
    time.Sleep(333 * time.Microsecond)
  }
}
