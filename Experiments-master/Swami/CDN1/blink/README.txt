PACKAGE DOCUMENTATION

package blink
    import "./blink"


CONSTANTS

const (
    ROUTE_TTL     = 15000 // 15 seconds
    PROBE_TIMEOUT = 10000 // 10 seconds
)

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
)
    Possible Blink packet types

VARIABLES

var (
    RoutingTable map[string]*Link

    LocalAddrString string

    DecisionTable map[string][]string

    PacketTrackTable map[string]map[int]Packet
)
    Global variables used by Blink nodes and oracle

FUNCTIONS

func Bounce(buf []byte)
    Bounce : This function will forward an inputted buffer to the desired
    destination address. If the current node is the desired destination
    address, the system will not Bounce the blink packet

func ByteArrayToDecisionTable(inBuf []byte) map[string][]string
    ByteArrayToDecisionTable : Converts an array of bytes into a decision
    table. Currently only supports one session

func CheckError(err error)
    CheckError : Checks to see if inputted error is empty. If not (meaning
    there is an error) then the error is printed and the program quits. */

func CheckExpiration(conn *net.UDPConn, link1 *Link)
    CheckExpiration : Every 100 milliseconds the program will check the
    TTLExpiration of the designated link. If the link has expired, the
    program will initialize the probing process by sending an initial probe
    to the corresponding node

func DecisionTableToByteArray(decisionTable map[string][]string) []byte
    DecisionTableToByteArray : Converts a decision table into an array of
    bytes, currently only supports one session

func ExtractFinalDestAddr(inBuf []byte) *net.UDPAddr
    ExtractFinalDestAddr : Extract Destination Address from the header of
    the Blink Packet. Returns address as a pointer to a resolved net.UDPAddr

func ExtractPacketData(inBuf []byte) []byte
    ExtractPacketData : From an inputted blink packet (an array of bytes),
    extract the packet's data from the packet

func ExtractPacketType(inBuf []byte) string
    ExtractPacketType : From an inputted blink packet (an array of bytes),
    extract the packet type from the header. Returns value as a string type

func ExtractSID(inBuf []byte) string
    ExtractSID : Extract SID from the header of the Blink Packet. Returns
    SID as a string

func ExtractSrcAddr(inBuf []byte) *net.UDPAddr
    ExtractSrcAddr : Extract Source Address from the header of the Blink
    Packet. Returns address as a pointer to a resolved net.UDPAddr

func GetData()
    GetData : Reads data flowing in to the listening connection that is
    passed into this function. The function then checks the type of packet
    being sent and responds accordingly.

func GetLatencyForPaths()
    GetLatencyForPaths : calculates the measured latency of each path in
    pathsArray using oracleRoutingTable values. Stores in global
    latencyArray.

func LinkToByteArray(link Link) []byte
    LinkToByteArray : Converts a link into a byte. The byte array is
    formatted to allow conversion using the ByteArrayToLink format

func MakeBlinkPacket(SID string, srcAddr *net.UDPAddr, finalDestAddr *net.UDPAddr, packetType string, buf []byte) []byte
    MakeBlinkPacket : Creates a Blink Packet, which is simply an array of
    bytes, where the first 66 bytes contain information about the source
    address, destination address, type of packet being sent. The remaining
    958 bytes are used for data transmission

func PacketToJSONToByteArray(p Packet) []byte
    PacketToJSONToByteArray : converts a Packet type to a JSON object
    encoded and returned as a byteArray

func ReceiveResponse2(link1 *Link)
    ReceiveResponse2 : Calculates the latency after a ResponseTwo is
    received and stores the latency in the corresponding link. Sends routing
    table to oracle when it's done calculating

func SelectPath(conn *net.UDPConn)
    SelectPath : function selects the best of two paths (based on lowest
    latency), updates the decision table accordingly. The oracle will then
    distribute this decision table to all of the nodes.

func SendInitialProbe(link1 *Link)
    SendInitialProbe : Sends the initial probe packet to another node via a
    link, which is passed into the function.

func SendResponse1(link1 *Link)
    SendResponse1 : Sends the Response 1 packet to another node via a link,
    which is passed into the function. */

func SendResponse2(link1 *Link)
    SendResponse2 : Sends the Response 2 packet to another node via a link,
    which is passed into the function. Will also calculate and store the
    calculated Latency of the link. Sends routing table to oracle when it's
    done calculating

func SetDecisionTable(decisionTable map[string][]string)
    SetDecisionTable : Sets the value of the global decision table

func SetLocalAddr(localNodeAddr *net.UDPAddr)
    SetLocalAddr : identifies the local UDP address in blink

func SetOracleAddr(oracleNodeAddr *net.UDPAddr)
    SetOracleAddr : identifies the oracle's UDP address in blink

func SetRoutingTable(routingTable map[string]*Link)
    SetRoutingTable : idenitfies the RoutingTable variable

func SetupListener()
    SetupListener : Sets up the listening connection to be used on the node.
    Sets the connection object as the global variable conn.

func SetupServer(serverAddr *net.UDPAddr, conn *net.UDPConn, buf []byte)
    SetupServer : This function initializes a data-sending server that lies
    outside of the Blink Node System. It responds to a HelloOracle message
    (sent by the server) by setting up a SessionID as well as an initial
    series of hops for the session. The oracle then sends the SessionID and
    the first blink node's address to the Server. The oracle then sends
    Decision Table updates to all of the nodes

func StartNode()
    StartNode : Starts a node instance and it's corresponding probing and
    data retreive processes.

func StartOracle(oracleAddrString string)
    StartOracle : This function initializes main functionality of the
    oracle. It reads any incoming data and decides how to appropriately
    respond to it.

func StartProber()
    StartProber : starts the probing functionality of the node. It must be
    started after the RoutingTable has been initialized and must be
    restarted every time the Routing Table is updated.

func UnwrapHeader(inBuf []byte) (string, *net.UDPAddr, *net.UDPAddr, string, []byte)
    UnwrapHeader : Extract all information from the Blink Packet. Return src
    Addr, Destination Addr, Packet Type, and Packet Data (in that order)

func UpdateDecisionTable(inBuf []byte)
    UpdateDecisionTable : Updates Global Decision Table at a given node
    based on an incoming blink packet.

func UpdateRoutingTable(addr *net.UDPAddr, inBuf []byte)
    UpdateRoutingTable : This function updates the GlobalRoutingTable based
    on information from receieved Routing Table

TYPES

type Link struct {
    Addr            string
    Latency         float64
    Loss            int
    TTLExpiration   time.Time // When was the last time I probed you (Creepy!)
    TimeoutDeadline time.Time
    TimerStartTime  time.Time
}
    Link : structure containing information about a give connection between
    two Blink nodes

func ByteArrayToLink(byteArray []byte) Link
    ByteArrayToLink : Converts a byte array into a link. The byte array must
    be formatted properly using the blink LinkToByteArray() function.

type Packet struct {
    Number int
    Time1  time.Time
    Time2  time.Time
    Data   []byte
}
    Packet : structure containing information about a given packet track

func ByteArrayToJSONToPacket(b []byte) Packet
    ByteArrayToJSONToPacket : converts a []byte and decodes it as a JSON
    following the Packet type structure


