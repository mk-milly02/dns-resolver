package resolver

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

// Question represents a DNS question section
type Question struct {
	name       string
	recordType uint16
	class      uint16
}

// Print returns the string representation of the question
func (q Question) Print() string {
	return fmt.Sprintf("\tName: %s \n\tType: %s \n\tClass: %s", q.name, GetResourceRecordType(q.recordType), GetResourceRecordClass(q.class))
}

// String returns the string representation of the question
func (q Question) String() string {
	return fmt.Sprintf("%s%04x%04x", q.name, q.recordType, q.class)
}

// ParseQuestion creates a new question from the byte slice
func ParseQuestion(b []byte, qcount, offset int) ([]Question, int) {
	var questions []Question
	for i := 0; i < qcount; i++ {
		name, newOffset := DecodeDomainName(b, offset)
		offset = newOffset
		recordType := binary.BigEndian.Uint16(b[offset : offset+2])
		class := binary.BigEndian.Uint16(b[offset+2 : offset+4])
		questions = append(questions, Question{name, recordType, class})
		offset += 4
	}
	return questions, offset
}

// EncodeDomainName encodes a domain name into the DNS format
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

// DecodeDomainName parses a domain name from the byte slice
func DecodeDomainName(b []byte, offset int) (string, int) {
	var name string
	for {
		length := int(b[offset])
		if length == 0 {
			break
		}
		name += string(b[offset+1 : offset+1+length])
		offset += length + 1
		if b[offset] != 0 {
			name += "."
		}
	}
	return name, offset + 1
}
