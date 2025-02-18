package resolver

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

// Message represents a DNS message
type Message struct {
	Header     header
	Question   question
	Answer     []resourceRecord
	Authority  []resourceRecord
	Additional []resourceRecord
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

// Question represents a DNS question section
type question struct {
	name       string
	recordType uint16
	class      uint16
}

// ResourceRecord represents a DNS resource record
type resourceRecord struct {
	name       string
	recordType uint16
	class      uint16
	ttl        uint32
	dataLength uint16
	data       string
}

// String returns the string representation of the header
func (h header) String() string {
	return fmt.Sprintf("%04x%04x%04x%04x%04x%04x", h.id, h.flags, h.queryCount, h.answerCount, h.authorityCount, h.additionalCount)
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
		Question: question{encodeURL("dns.google.com"), 0x01, 0x01},
	}
}

// encodeURL encodes a domain name into the DNS format
func encodeURL(name string) string {
	var encoded string
	labels := strings.Split(name, ".")
	for _, l := range labels {
		encoded += fmt.Sprintf("%02x", len(l))
		encoded += hex.EncodeToString([]byte(l))
	}
	encoded += fmt.Sprintf("%02x", 0)
	return encoded
}

// decodeURL decodes a domain name from the DNS format
func decodeURL(encoded string) string {
	var decoded string
	for i := 0; i < len(encoded); i++ {
		length := encoded[i]
		if length == '0' {
			break
		}
		decoded += encoded[i+1 : i+1+int(length)]
		i += int(length)
		if i+1 < len(encoded) {
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

// newHeader creates a new header from the byte slice
// The header is 12 bytes long
func newHeader(b []byte) header {
	return header{
		id:              binary.BigEndian.Uint16(b[:2]),
		flags:           binary.BigEndian.Uint16(b[2:4]),
		queryCount:      binary.BigEndian.Uint16(b[4:6]),
		answerCount:     binary.BigEndian.Uint16(b[6:8]),
		authorityCount:  binary.BigEndian.Uint16(b[8:10]),
		additionalCount: binary.BigEndian.Uint16(b[10:12]),
	}
}

// newQuestion creates a new question from the byte slice
//
// The question is variable length
func newQuestion(b []byte, offset int) (question, int) {
	return question{
		name:       decodeURL(string(b[:offset])),
		recordType: binary.BigEndian.Uint16(b[offset : offset+2]),
		class:      binary.BigEndian.Uint16(b[offset+2 : offset+4]),
	}, offset + 4
}

// newResourceRecord creates a new resource record from the byte slice
func newResourceRecord(b []byte, domainName string, offset int) (rr resourceRecord, nextOffset int) {
	rr.name = domainName
	rr.recordType = binary.BigEndian.Uint16(b[offset : offset+2])
	rr.class = binary.BigEndian.Uint16(b[offset+2 : offset+4])
	rr.ttl = binary.BigEndian.Uint32(b[offset+4 : offset+8])
	rr.dataLength = binary.BigEndian.Uint16(b[offset+8 : offset+10])
	rr.data = string(b[offset+10 : offset+10+int(rr.dataLength)])
	return rr, offset + 10 + int(rr.dataLength)
}

// ParseResponse parses the response from the name server
func ParseResponse(response []byte, req Message) (m Message) {
	m.Header = newHeader(response)
	q, offset := newQuestion(response[12:], len(req.Question.name))
	m.Question = q
	for i := 0; i < int(m.Header.answerCount); i++ {
		rr, nextOffset := newResourceRecord(response, m.Question.name, offset)
		m.Answer = append(m.Answer, rr)
		offset = nextOffset
	}
	return m
}
