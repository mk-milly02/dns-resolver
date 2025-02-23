package resolver

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

// ResourceRecord represents a DNS resource record
type ResourceRecord struct {
	name       string
	recordType uint16
	class      uint16
	ttl        uint32
	dataLength uint16
	data       string
}

// ParseResourceRecord creates a new resource record from the byte slice
func ParseResourceRecord(b []byte, count, offset int) (rr []ResourceRecord, newOffset int) {
	for i := 0; i < count; i++ {
		name, offset := DecodeDomainName(b, offset)
		recordType := binary.BigEndian.Uint16(b[offset : offset+2])
		class := binary.BigEndian.Uint16(b[offset+2 : offset+4])
		ttl := binary.BigEndian.Uint32(b[offset+4 : offset+8])
		dataLength := binary.BigEndian.Uint16(b[offset+8 : offset+10])
		var data string
		switch recordType {
		case 0x01:
			data = fmt.Sprintf("%d.%d.%d.%d", b[offset+10], b[offset+11], b[offset+12], b[offset+13])
		case 0x05:
			data, _ = DecodeDomainName(b, offset+10)
		default:
			data = hex.EncodeToString(b[offset+10 : offset+10+int(dataLength)])
		}
		rr = append(rr, ResourceRecord{name, recordType, class, ttl, dataLength, data})
		offset += 10 + int(dataLength)
	}
	return rr, offset
}

// String returns the string representation of the message
func (rr ResourceRecord) String() string {
	if rr.recordType == 0x01 {
		return fmt.Sprintf("Name: %s \nType: %04x \nClass: %04x \nTTL:%08x \nLength: %04x \nAddress: %s",
			rr.name, rr.recordType, rr.class, rr.ttl, rr.dataLength, rr.data)
	}
	return fmt.Sprintf("Name: %s \nType: %04x \nClass: %04x \nTTL:%08x \nLength: %04x \nName server: %s",
		rr.name, rr.recordType, rr.class, rr.ttl, rr.dataLength, rr.data)
}
