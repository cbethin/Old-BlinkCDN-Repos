package main

import (
  "fmt"
  "net"
  "io/ioutil"
  "os"
  "./blink"
)

func main() {

  myAddress, _ := net.ResolveUDPAddr("udp", "127.0.0.1:9009")

  conn, _ := net.ListenUDP("udp", myAddress)

  buffer := make([]byte, 1024)

  var finalBuffer []byte
  var fileHasEnded = false
  var toWriteToFile []byte

  for fileHasEnded == false {
    conn.ReadFromUDP(buffer)
    packetData := blink.ExtractPacketData(buffer)
    fmt.Println("Data:", packetData)
    toWriteToFile, fileHasEnded = unfillByteArray(packetData, 942)
    for _, value := range toWriteToFile {
      finalBuffer = append(finalBuffer, value)
    }
  }

  //fmt.Println("Output:", finalBuffer)
  _ = ioutil.WriteFile("test2.m4v", []byte(finalBuffer), os.ModeAppend)
  fmt.Println("DONE")
}

func unfillByteArray(inBuf []byte, fromLength int) ([]byte, bool) {

  finalOutBuf := make([]byte, 942)
  outBuf := make([]byte, 942)
  fileHasEnded := false

  for i:=0; i < len(inBuf); i++ {
    if ((i+6) < len(inBuf)) && (string(inBuf[i:i+6]) == "::::::") {
      fmt.Println("FOUND IT")
      outBuf2 := make([]byte, i)
      copy(outBuf2[:], inBuf[:i])
      finalOutBuf = outBuf2
      fileHasEnded = true
      break
    } else {
      outBuf[i] = inBuf[i]
      finalOutBuf = outBuf
      fileHasEnded = false
    }
  }

  return finalOutBuf, fileHasEnded
}

func CheckError(err error) {
  if err != nil {
    fmt.Println("Err:", err)
    os.Exit(1)
  }
}
