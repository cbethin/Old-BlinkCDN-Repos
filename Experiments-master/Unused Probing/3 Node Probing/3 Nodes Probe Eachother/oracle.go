package main

import (
  "fmt"
  "net"
  "os"
  "strconv"
)

type link struct {
  name string
  loss float64
  latency float64
  address *net.UDPAddr
  sourceAddr *net.UDPAddr
}

type connection struct {
  cost float64
  addr1 *net.UDPAddr
  addr2 *net.UDPAddr
}

type routingTable struct {
  sourceAddr *net.UDPAddr
  routes [2]link
}

func main() {
  if len(os.Args) != 2 {
    fmt.Println("Not enough arguments")
    os.Exit(1)
  }

  // Set Up listener
  oracleAddr, err := net.ResolveUDPAddr("udp", ":"+os.Args[1])
  checkError(err)
  conn, err := net.ListenUDP("udp", oracleAddr)
  checkError(err)
  fmt.Println("Listening...")
  //var clientAddr *net.UDPAddr

  GlobalRoutingTable := getGlobalRoutingTable(conn)
  fmt.Println(GlobalRoutingTable)
  // serverToUse := MakeDecision(RoutingTable)
  // copy(buffer[:], strconv.Itoa(serverToUse))
  // _, err = conn.WriteToUDP(buffer, clientAddr)
  // checkError(err)

  MakePathDecision(GlobalRoutingTable)
}

func MakePathDecision(GlobalRoutingTable [3]routingTable) {

  var ConnectionsTable [3]connection

  for i:=0; i < 3; i++ { // loops through global routing table
    currentNodeAddr := GlobalRoutingTable[i].sourceAddr
    for j:=0; j < 2; j++ { // loops through each link of every routing table in global
      otherNodeAddr := GlobalRoutingTable[i].routes[j].address
      for k:=0; k < 3; k++ { // loops through all other routing tables in global
        if GlobalRoutingTable[k].sourceAddr.String() == otherNodeAddr.String() {
          for l:=0; l < 2; l++ { // loops through each link of those other r.tables
            if GlobalRoutingTable[k].routes[l].address.String() == currentNodeAddr.String() {
              for m:=0; m < 3; m++ { // loops through connections table to see if link exists
                if ConnectionsTable[m].addr1.String() != otherNodeAddr.String() || ConnectionsTable[m].addr1.String() == currentNodeAddr.String() {
                  averageLoss := GlobalRoutingTable[k].routes[l].loss * 0.5 + GlobalRoutingTable[i].routes[j].loss
                  averageLatency := GlobalRoutingTable[k].routes[l].latency * 0.5 + GlobalRoutingTable[i].routes[j].latency
                  averageCost := 0.7 * averageLoss + 0.3 * averageLatency
                  ConnectionsTable[i].addr1 = currentNodeAddr
                  ConnectionsTable[i].addr2 = otherNodeAddr
                  ConnectionsTable[i].cost = averageCost
                }
              }
            }
          }
        }
      }
    }
  }

  fmt.Println(ConnectionsTable[0].addr1.String(), ConnectionsTable[0].addr2.String())
  fmt.Println(ConnectionsTable[1].addr1.String(), ConnectionsTable[1].addr2.String())
  fmt.Println(ConnectionsTable[2].addr1.String(), ConnectionsTable[2].addr2.String())
}

func getGlobalRoutingTable(conn *net.UDPConn) [3]routingTable {
  buffer := make([]byte, 1024)
  var GlobalRoutingTable [3]routingTable
  var RoutingTable routingTable

  for i := 0; i < 3; i++ {
    _, _, err := conn.ReadFromUDP(buffer)
    checkError(err)

    link := ByteArrayToLink(buffer[0:98])
    link2 := ByteArrayToLink(buffer[98:196])
    RoutingTable.sourceAddr = link.sourceAddr
    RoutingTable.routes[0] = link
    RoutingTable.routes[1] = link2
    GlobalRoutingTable[i] = RoutingTable
  }

  return GlobalRoutingTable
}

func ByteArrayToLink(byteArray []byte) link {
  var newLink link
  newLink.name = string(byteArray[:4])
  lossString := string(byteArray[4:19])
  latencyString := string(byteArray[19:34])
  addressBytes := byteArray[34:66]
  sourceAddrBytes := byteArray[66:98]
  addressString := ""
  sourceAddressString := ""

  for i:= 0; i < 32; i++ {
    if string(addressBytes[i]) != "!" {
      addressString += string(addressBytes[i])
    } else {
      break
    }
  }

  for i:= 0; i < 32; i++ {
    if string(sourceAddrBytes[i]) != "!" {
      sourceAddressString += string(sourceAddrBytes[i])
    } else {
      break
    }
  }

  address, err := net.ResolveUDPAddr("udp", addressString)
  checkError(err)
  sourceAddr, err := net.ResolveUDPAddr("udp", sourceAddressString)
  checkError(err)

  newLink.sourceAddr = sourceAddr
  newLink.address = address
  newLink.loss, _ = strconv.ParseFloat(lossString, 64)
  newLink.latency, _ = strconv.ParseFloat(latencyString, 64)
  return newLink
}


func MakeDecision(RoutingTable [2]link) int {
  linkCost1 := 0.7 * RoutingTable[0].loss + 0.3 * RoutingTable[0].latency
  linkCost2 := 0.7 * RoutingTable[1].loss + 0.3 * RoutingTable[1].latency
  if linkCost1 <= linkCost2 {
    return 0
  } else {
    return 1
  }
}

func checkError(err error) {
  if err != nil {
    fmt.Println("ERROR:", err)
    os.Exit(1)
  }
}
