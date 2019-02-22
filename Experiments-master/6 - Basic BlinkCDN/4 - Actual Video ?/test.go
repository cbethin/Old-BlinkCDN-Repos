package main

import (
  "fmt"
  "os"
  "io/ioutil"
  "net"
  "time"
)

var (
  localAddr *net.UDPAddr
  otherAddr *net.UDPAddr
)

func main() {

  var err error
  localAddr, err = net.ResolveUDPAddr("udp", "127.0.0.1:4000")
  CheckError(err)
  otherAddr, err = net.ResolveUDPAddr("udp", "127.0.0.1:5000")
  CheckError(err)

  conn, err := net.ListenUDP("udp", localAddr)
  CheckError(err)
  defer conn.Close()

  filename := "test.m4v"
  sendVideo(filename, conn)

  time.Sleep(10*time.Second)



  // filename := "test.m4v"
  // buf, err := ioutil.ReadFile(filename)
  // if err != nil {
  //   fmt.Println("Error:", err)
  // }
  // fmt.Println("GOT IT:", buf[:20])
  // fmt.Println("Length:", len(buf))
  //
  // fmt.Println("Writing...")
  // err = ioutil.WriteFile("test2.png", buf, os.ModeDevice)
}

func sendVideo(filename string, conn *net.UDPConn) {
  // Open File In Buffer
  buffer, err := ioutil.ReadFile(filename)
  CheckError(err)

  outputBuf := make([]byte, 1024)

  for i:=0; i < len(buffer); i += 1024 {
    if i + 1025 < len(buffer) {
      copy(outputBuf[:], buffer[i:i+1024])
    } else {
      finalValue := len(buffer)-1
      outputBuf = fillByteArray(buffer[i:finalValue], 1024)
    }

    _, err = conn.WriteToUDP(outputBuf, otherAddr)
    CheckError(err)
    fmt.Println("Buffer:", outputBuf)
    time.Sleep(100 * time.Microsecond)
  }

    fmt.Println("Length:", len(buffer))
    fmt.Println("Amount Extra:", len(buffer)%1024)
}

func CheckError(err error) {
  if err != nil {
    fmt.Println("Err:", err)
    os.Exit(1)
  }
}

func fillByteArray(buffer []byte, toLength int) []byte {
  outBuf := make([]byte, toLength)
  copy(outBuf[:], buffer[:])
	for i:=len(buffer); i < toLength; i++ {
		outBuf[i] = ":"[0]
	}

  fmt.Println("Output:", outBuf)
	return outBuf
}
