package resolver

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// Message represents a DNS message
type Message struct {
	Header     header
	Question   question
	Answer     [6]uint16
	Authority  [6]uint16
	Additional [6]uint16
}

// Header represents the DNS message header
type header struct {
	id              uint16
	flags           uint16
	queryCount      uint16
	answerCount     uint16
	authorityCount  uint16
	additionalCount uint16
}

// String returns the string representation of the header
func (h header) String() string {
	return fmt.Sprintf("%04x%04x%04x%04x%04x%04x", h.id, h.flags, h.queryCount, h.answerCount, h.authorityCount, h.additionalCount)
}

// Question represents a DNS question section
type question struct {
	name       string
	recordType uint16
	class      uint16
}

// String returns the string representation of the question
func (q question) String() string {
	return fmt.Sprintf("%s%04x%04x", q.name, q.recordType, q.class)
}

// NewMessage creates a new DNS message with default values
// The message includes a header, a question section with a default query for "dns.google.com",
// and empty answer, authority, and additional sections.
func NewMessage() Message {
	return Message{
		Header: header{
			id:         0x16,
			flags:      0x100,
			queryCount: 0x01,
		},
		Question: question{EncodeURL("dns.google.com"), 0x01, 0x01},
	}
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

// BuildQuery builds the DNS query from the message
func (m Message) BuildQuery() string {
	return m.Header.String() + m.Question.String()
}

// ValidateResponse validates the response from the name server
func (m Message) ValidateResponse(response []byte) bool {
	return hex.EncodeToString(response[:2]) == fmt.Sprintf("%04x", m.Header.id)
}
