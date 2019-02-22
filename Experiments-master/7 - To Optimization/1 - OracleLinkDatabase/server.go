package main

import (
  "fmt"
  "os"
  "./blink"
  "net"
  "time"
  "io/ioutil"
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

  // SEND VIDEO

  filename := "test.m4v"
  videoBuffer, err := ioutil.ReadFile(filename)
  blink.CheckError(err)

  outputBuf := make([]byte, 942)

  for i:=0; i < len(videoBuffer); i += 942 {
    if i + 943 < len(videoBuffer) {
      copy(outputBuf[:], videoBuffer[i:i+942])
    } else {
      finalValue := len(videoBuffer)-1
      outputBuf = fillByteArray(videoBuffer[i:finalValue], 942)
    }

    blinkPacket := blink.MakeBlinkPacket(SID, localAddr, destAddr, blink.Iframe, outputBuf)

    _, err = conn.WriteToUDP(blinkPacket, blinkAddr)
    blink.CheckError(err)
    packetData := blink.ExtractPacketData(outputBuf)
    fmt.Println("Buffer:", packetData)
    time.Sleep(1 * time.Millisecond)
  }

    fmt.Println("Length:", len(buffer))
    fmt.Println("Amount Extra:", len(buffer)%942)
}

func fillByteArray(buffer []byte, toLength int) []byte {
  outBuf := make([]byte, toLength)
  copy(outBuf[:], buffer[:])
	for i:=len(buffer); i < toLength; i++ {
		outBuf[i] = ":"[0]
	}

	return outBuf
}
