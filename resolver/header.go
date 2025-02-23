package resolver

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

// Header represents the DNS message header
type Header struct {
	id              uint16
	flags           uint16
	queryCount      uint16
	answerCount     uint16
	authorityCount  uint16
	additionalCount uint16
}

// String returns the string representation of the header
func (h Header) String() string {
	return fmt.Sprintf("%04x%04x%04x%04x%04x%04x", h.id, h.flags,
		h.queryCount, h.answerCount, h.authorityCount, h.additionalCount)
}

// ParseHeader creates a new header from the byte slice
// The header is 12 bytes long
func ParseHeader(b []byte) (Header, int) {
	if len(b) < 12 {
		log.Fatal("Invalid DNS header")
	}
	return Header{
		id:              binary.BigEndian.Uint16(b[:2]),
		flags:           binary.BigEndian.Uint16(b[2:4]),
		queryCount:      binary.BigEndian.Uint16(b[4:6]),
		answerCount:     binary.BigEndian.Uint16(b[6:8]),
		authorityCount:  binary.BigEndian.Uint16(b[8:10]),
		additionalCount: binary.BigEndian.Uint16(b[10:12]),
	}, 12
}

// generate_transaction_id generates a random transaction ID
func generate_transaction_id() uint16 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := uint16(r.Intn(math.MaxUint16))
	return num
}

// NewHeader creates a new DNS header with default values
func NewHeader() Header {
	return Header{
		id:         generate_transaction_id(),
		flags:      0x00,
		queryCount: 0x01,
	}
}
