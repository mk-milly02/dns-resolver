package resolver

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/netip"
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

// parseHeader creates a new header from the byte slice
// The header is 12 bytes long
func parseHeader(b []byte) header {
	return header{
		id:              binary.BigEndian.Uint16(b[:2]),
		flags:           binary.BigEndian.Uint16(b[2:4]),
		queryCount:      binary.BigEndian.Uint16(b[4:6]),
		answerCount:     binary.BigEndian.Uint16(b[6:8]),
		authorityCount:  binary.BigEndian.Uint16(b[8:10]),
		additionalCount: binary.BigEndian.Uint16(b[10:12]),
	}
}

// parseResourceRecord creates a new resource record from the byte slice
func parseResourceRecord(b []byte, domainName string, offset int) (rr resourceRecord, newOffset int) {
	if hex.EncodeToString(b[offset:offset+1]) == "c0" {
		offset += 2
		rr.name = domainName
	}
	rr.recordType = binary.BigEndian.Uint16(b[offset : offset+2])
	offset += 2
	rr.class = binary.BigEndian.Uint16(b[offset : offset+2])
	offset += 2
	rr.ttl = binary.BigEndian.Uint32(b[offset : offset+4])
	offset += 4
	rr.dataLength = binary.BigEndian.Uint16(b[offset : offset+2])
	offset += 2
	if rr.recordType == 0x01 {
		ip, _ := netip.AddrFromSlice(b[offset : offset+int(rr.dataLength)])
		rr.data = ip.String()
	}
	offset += int(rr.dataLength)
	return rr, offset
}

// ParseResponse parses the response from the name server
func ParseResponse(response []byte, req Message) (m Message) {
	offset := 0
	m.Header = parseHeader(response[offset:12]) // 12 bytes for the header
	offset += 12
	m.Question = req.Question
	offset += len(req.Question.name)/2 + 4 // 4 bytes for the record type and class
	for i := 0; i < int(m.Header.answerCount); i++ {
		rr, nextOffset := parseResourceRecord(response, req.Question.name, offset)
		offset = nextOffset
		m.Answer = append(m.Answer, rr)
	}
	for i := 0; i < int(m.Header.authorityCount); i++ {
		rr, nextOffset := parseResourceRecord(response, req.Question.name, offset)
		offset = nextOffset
		m.Authority = append(m.Authority, rr)
	}
	for i := 0; i < int(m.Header.additionalCount); i++ {
		rr, newOffset := parseResourceRecord(response, req.Question.name, offset)
		offset = newOffset
		m.Additional = append(m.Additional, rr)
	}
	return m
}

// IsAResponse checks if the message is a response
func (m Message) IsAResponse() bool {
	return m.Header.flags&0x8000 == 0x8000
}

// String returns the string representation of the message
func (rr resourceRecord) String() string {
	return fmt.Sprintf("Name: %s\n Type: %d\n Class: %d\n TTL: %d\n Address: %s", decodeURL(rr.name), rr.recordType, rr.class, int(rr.ttl), rr.data)
}
