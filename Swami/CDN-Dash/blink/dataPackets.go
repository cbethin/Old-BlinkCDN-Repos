package blink
//
import (
  "time"
  "bytes"
	"encoding/json"
)

type ByteArray []byte

type Packet struct {
  Number  int
  Data    []byte
	Time1  time.Time
	Time2  time.Time
	Time3  time.Time
  Filename string
}

func (p Packet) ToByteStream() ([]byte) {
  b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(p)
	return b.Bytes()
}

func (b ByteArray) ToPacket() (Packet) {
  var p Packet
   err := json.NewDecoder(bytes.NewReader(b)).Decode(&p)
   CheckError(err)
   return p
}
