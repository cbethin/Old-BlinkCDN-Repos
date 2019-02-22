package main

import (
        "fmt"
        "log"
        "net"
        "bufio"
        "os"
)


const(
  play = "play"
)

func handleUDPConnection(conn *net.UDPConn) {

        // here is where you want to do stuff like read or write to client
        defer conn.Close()
        buffer := make([]byte, 1024)

        n, addr, err := conn.ReadFromUDP(buffer)

        fmt.Println("UDP client : ", addr)
        fmt.Println("Received from UDP client :  ", string(buffer[:n]))

        if err != nil {
                log.Fatal(err)
        }

        // NOTE : Need to specify client address in WriteToUDP() function
        //        otherwise, you will get this error message
        //        write udp : write: destination address required if you use Write() function instead of WriteToUDP()

        // write message back to client
        message := []byte("Hello UDP client!")
        _, err = conn.WriteToUDP(message, addr)

        if err != nil {
                log.Println(err)
        }
        s := string(buffer[0:n])
        fmt.Println(s)
        if s[0:4] == play{
          //open file
          file, err := os.Open("ourfile.txt");
          if err != nil {
            fmt.Println("Could not open file")
          }
          defer file.Close()
          scanner := bufio.NewScanner(file)
          for scanner.Scan() {
              //log.Println(scanner.Text())
              //fmt.Println(scanner.Text())
              message := []byte(fillString(scanner.Text(), 1024))
              _, err = conn.WriteToUDP(message, addr)
          }
          if err = scanner.Err(); err != nil {
            log.Fatal(err)
          }
          message := []byte("done")
            _, err = conn.WriteToUDP(message, addr)

          // fmt.Println("A client has connected!")
        	// file, err := os.Open("ourfile.txt")
        	// if err != nil {
        	// 	fmt.Println(err)
        	// 	return
        	// }
        	// fileInfo, err := file.Stat()
        	// if err != nil {
        	// 	fmt.Println(err)
        	// 	return
        	// }
        	// fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
        	// fileName := fillString(fileInfo.Name(), 64)
        	// fmt.Println("Sending filename and filesize!")
        	// conn.WriteToUDP([]byte(fileSize), addr)
        	// conn.WriteToUDP([]byte(fileName), addr)
        	// sendBuffer := make([]byte, 1024)
        	// fmt.Println("Start sending file!")
        	// for {
        	// 	_,_, err = conn.ReadFromUDP(sendBuffer)
        	// 	if err == io.EOF {
        	// 		break
        	// 	}
        	// 	conn.WriteToUDP(sendBuffer, addr)
        	// }
        	// fmt.Println("File has been sent, closing connection!")
        	// return
        }


}

func main() {
        hostName := "155.246.66.150"
        portNum := "6000"
        service := hostName + ":" + portNum

        udpAddr, err := net.ResolveUDPAddr("udp4", service)

        if err != nil {
                log.Fatal(err)
        }

        // setup listener for incoming UDP connection
        ln, err := net.ListenUDP("udp", udpAddr)

        if err != nil {
                log.Fatal(err)
        }

        fmt.Println("UDP server up and listening on port 6000")

        defer ln.Close()

        for {
                // wait for UDP client to connect
                handleUDPConnection(ln)
        }

}



func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}
