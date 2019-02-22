package blink

// RouteTable : A routing table strucutre of map[string][]*Link (mapping addresses to their corresponding links)
type RouteTable map[string][]*Link

// Initialize (RouteTable) : Initializes a routing table given a list of all nodes in the system
func (RoutingTable *RouteTable) Initialize(nodesTable []string) {
	var linksArray []*Link
	RouteTable := *RoutingTable

	for _, value := range nodesTable {
		var newLink Link
		newLink.Addr = value
		linksArray = append(linksArray, &newLink)
	}

	for nodeIndex, nodeValue := range nodesTable {
		RouteTable[nodeValue] = make([]*Link, 0)
		for linkIndex, linkValue := range linksArray {
			if nodeIndex != linkIndex {
				RouteTable[nodeValue] = append(RouteTable[nodeValue], linkValue)
			}
		}
	}
}
