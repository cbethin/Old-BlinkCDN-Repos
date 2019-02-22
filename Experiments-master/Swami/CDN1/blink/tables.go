package blink

import "fmt"

// MARK: DECISION TABLE

// UpdateDecisionTable : Updates Global Decision Table at a given node based on an incoming blink packet.
func UpdateDecisionTable(inBuf []byte) {
	tableData := ExtractPacketData(inBuf)
	decisionTable := ByteArrayToDecisionTable(tableData)
	SetDecisionTable(decisionTable)
	fmt.Println("Decision Update:", DecisionTable)
}

// SetDecisionTable : Sets the value of the global decision table
func SetDecisionTable(decisionTable map[string][]string) {
	DecisionTable = decisionTable
}

// DecisionTableToByteArray : Converts a decision table into an array of bytes, currently only supports one session
func DecisionTableToByteArray(decisionTable map[string][]string) []byte {
	outBuf := make([]byte, 100)
	for key, value := range decisionTable {
		copy(outBuf[:16], []byte(key))
		// Add Hops Array to byte arra
		hops := value[0] + "!" + value[1] + "!" + value[2] + "!"
		copy(outBuf[16:], []byte(hops))
	}

	return outBuf
}

// ByteArrayToDecisionTable : Converts an array of bytes into a decision table. Currently only supports one session
func ByteArrayToDecisionTable(inBuf []byte) map[string][]string {

	decisionTable := make(map[string][]string)
	SID := string(inBuf[:16])
	hopsArray := inBuf[16:]
	hopStringsArray := []string{}
	addrCount := 0
	addrString := ""

	for _, value := range hopsArray {
		if addrCount < 3 {
			if string(value) != "!" {
				addrString += string(value)
			} else {
				hopStringsArray = append(hopStringsArray, addrString)
				// hopStringsArray[addrCount] = addrString
				addrCount++
				addrString = ""
			}
		} else {
			break
		}
	}

	decisionTable[SID] = hopStringsArray
	return decisionTable
}
