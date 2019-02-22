package blink

import (
	"fmt"
	"net"
)

const (
	ROUTE_TTL     = 15000 // 15 seconds
	PROBE_TIMEOUT = 10000 // 10 seconds
)

// Possible Blink packet types
const (
	InitialProbe      = "IP"
	ResponseOne       = "R1"
	ResponseTwo       = "R2"
	Iframe            = "IF"
	Pframe            = "PF"
	Bframe            = "BF"
	RoutingTableType  = "RT"
	DecisionTableType = "DT"
	HelloOracle       = "HO"
	BounceUpdate      = "BU"
	ProbeID           = "0000000000000000"
	PacketTrack       = "PT"
	BlinkSignal       = "BS"
)

// Global variables used by Blink nodes and oracle
var (
	RoutingTable      map[string]*Link
	localAddr         *net.UDPAddr
	LocalAddrString   string
	oracleAddr        *net.UDPAddr
	DecisionTable     map[string][]string
	conn              *net.UDPConn
	fileName          string
	latencyArray      = make([]float64, 2)
	pathsArray        = make([][]string, 2)
	PacketTrackTable  map[string]map[int]Packet
	EstimatorFunction LatencyEstimator
)

// SetRoutingTable : idenitfies the RoutingTable variable
func SetRoutingTable(routingTable map[string]*Link) {
	RoutingTable = routingTable
}

// SetLocalAddr : identifies the local UDP address in blink
func SetLocalAddr(localNodeAddr *net.UDPAddr) {
	localAddr = localNodeAddr
}

// SetOracleAddr : identifies the oracle's UDP address in blink
func SetOracleAddr(oracleNodeAddr *net.UDPAddr) {
	oracleAddr = oracleNodeAddr
}

// SetLatencyEstimatorFunction : Sets the estimator funciton to be used by the prober
func SetLatencyEstimatorFunction(fn LatencyEstimator) {
	EstimatorFunction = fn
}

// CheckError : Checks to see if inputted error is empty. If not (meaning there is an error) then the error is printed and the program quits. */
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		// os.Exit(1)
	}
}
