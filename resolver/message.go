package resolver

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// Message represents a DNS message
type Message struct {
	Header     Header
	Question   Question
	Answer     []ResourceRecord
	Authority  []ResourceRecord
	Additional []ResourceRecord
}

// NewMessage creates a new DNS message with default values.
// The message includes a header, a question section with a query for the given name,
// and empty answer, authority, and additional sections.
func NewMessage(name string) Message {
	return Message{
		Header: Header{
			id:         0x16,
			flags:      0x00,
			queryCount: 0x01,
		},
		Question: Question{EncodeDomainName(name), 0x01, 0x01},
	}
}

// encodeDomainName encodes a domain name into the DNS format
func EncodeDomainName(name string) string {
	var encoded string
	labels := strings.Split(name, ".")
	for _, l := range labels {
		encoded += fmt.Sprintf("%02x", len(l))
		encoded += hex.EncodeToString([]byte(l))
	}
	encoded += fmt.Sprintf("%02x", 0)
	return encoded
}

// decodeDomainName decodes a domain name from the DNS format
func DecodeDomainName(encoded string) string {
	var decoded string
	encoded_bytes, _ := hex.DecodeString(encoded)
	encoded = string(encoded_bytes)
	for i := 0; i < len(encoded); i++ {
		length := encoded[i]
		if length == '0' {
			break
		}
		decoded += encoded[i+1 : i+1+int(length)]
		i += int(length)
		if i+2 < len(encoded) {
			decoded += "."
		}
	}
	return decoded
}

// BuildQuery builds the DNS query from the message
func (m Message) BuildQuery() string {
	return m.Header.String() + m.Question.String()
}

// ValidateResponse validates the response from the name server
func (m Message) ValidateResponse(response []byte) bool {
	return hex.EncodeToString(response[:2]) == fmt.Sprintf("%04x", m.Header.id)
}

// IsAResponse checks if the message is a response
func (m Message) IsAResponse() bool {
	return m.Header.flags&0x8000 == 0x8000
}
