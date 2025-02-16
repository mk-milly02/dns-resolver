package resolver

import (
	"encoding/hex"
	"fmt"
	"strings"
)

/*
+---------------------+
|        Header       |
+---------------------+
|       Question      | the question for the name server
+---------------------+
|        Answer       | RRs answering the question
+---------------------+
|      Authority      | RRs pointing toward an authority
+---------------------+
|      Additional     | RRs holding additional information
+---------------------+
*/
type Message struct {
	Header     [6]uint16
	Question   QMessage
	Answer     [6]uint16
	Authority  [6]uint16
	Additional [6]uint16
}

// QMessage represents a DNS question section
type QMessage struct {
	Name       string
	RecordType uint16
	Class      uint
}

// NewMessage creates a new DNS message with default values
// The message includes a header, a question section with a default query for "dns.google.com",
// and empty answer, authority, and additional sections.
func NewMessage() Message {
	return Message{
		Header:     [6]uint16{0x16, 0x100, 0x01, 0x00, 0x00, 0x00},
		Question:   QMessage{EncodeURL("dns.google.com"), 0x01, 0x01},
		Answer:     [6]uint16{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		Authority:  [6]uint16{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		Additional: [6]uint16{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
}

// HexStringFromUInt16 converts an array of uint16 to a hex string
func HexStringFromUInt16(arr []uint16) (result string) {
	for i := 0; i < len(arr); i++ {
		result += fmt.Sprintf("%04x", arr[i])
	}
	return
}

// EncodeURL encodes a domain name into the DNS format
func EncodeURL(name string) string {
	var encoded string
	labels := strings.Split(name, ".")
	for _, l := range labels {
		encoded += fmt.Sprintf("%02x", len(l))
		encoded += hex.EncodeToString([]byte(l))
	}
	encoded += fmt.Sprintf("%02x", 0)
	return encoded
}

func (m Message) BuildQuery() string {
	return HexStringFromUInt16(m.Header[:]) + m.Question.Name + HexStringFromUInt16([]uint16{m.Question.RecordType, uint16(m.Question.Class)}) +
		HexStringFromUInt16(m.Answer[:]) + HexStringFromUInt16(m.Authority[:]) + HexStringFromUInt16(m.Additional[:])
}
