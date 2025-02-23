package resolver

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/netip"
)

const (
	A     = 0x01
	NS    = 0x02
	CNAME = 0x05
	SOA   = 0x06
	PTR   = 0x0c
	MX    = 0x0f
	TXT   = 0x10
	AAAA  = 0x1c
)

const (
	IN = 0x01
	CS = 0x02
	CH = 0x03
	HS = 0x04
)

// GetResourceRecordType returns the string representation of the resource record type
func GetResourceRecordType(recordType uint16) string {
	switch recordType {
	case A:
		return "A"
	case NS:
		return "NS"
	case CNAME:
		return "CNAME"
	case SOA:
		return "SOA"
	case PTR:
		return "PTR"
	case MX:
		return "MX"
	case TXT:
		return "TXT"
	case AAAA:
		return "AAAA"
	default:
		return "Unknown"
	}
}

// GetResourceRecordClass returns the
func GetResourceRecordClass(class uint16) string {
	switch class {
	case IN:
		return "IN"
	case CS:
		return "CS"
	case CH:
		return "CH"
	case HS:
		return "HS"
	default:
		return "Unknown"
	}
}

// ResourceRecord represents a DNS resource record
type ResourceRecord struct {
	Name       string
	RecordType uint16
	Class      uint16
	TTL        uint32
	DataLength uint16
	Data       string
}

// ParseResourceRecord creates a new resource record from the byte slice
func ParseResourceRecord(b []byte, count, offset int) (rr []ResourceRecord, newOffset int) {
	for i := 0; i < count; i++ {
		name, nOffset := DecodeDomainName(b, offset)
		recordType := binary.BigEndian.Uint16(b[nOffset : nOffset+2])
		class := binary.BigEndian.Uint16(b[nOffset+2 : nOffset+4])
		ttl := binary.BigEndian.Uint32(b[nOffset+4 : nOffset+8])
		dataLength := binary.BigEndian.Uint16(b[nOffset+8 : nOffset+10])
		var data string
		switch recordType {
		case 0x01:
			ip, ok := netip.AddrFromSlice(b[nOffset+10 : nOffset+14])
			if ok {
				data = ip.String()
			}
		case 0x02, 0x05:
			data, _ = DecodeDomainName(b, nOffset+10)
		case 0x1c:
			ip, ok := netip.AddrFromSlice(b[nOffset+10 : nOffset+26])
			if ok {
				data = ip.String()
			}
		default:
			data = hex.EncodeToString(b[nOffset+10 : nOffset+10+int(dataLength)])
		}
		rr = append(rr, ResourceRecord{name, recordType, class, ttl, dataLength, data})
		nOffset += 10 + int(dataLength)
		offset = nOffset
	}
	return rr, offset
}

// Print returns the string representation of the message
func (rr ResourceRecord) Print() string {
	if rr.RecordType == 0x01 {
		return fmt.Sprintf("\tName: %s \n\tType: %s \n\tClass: %s \n\tTTL: %d \n\tLength: %d \n\tAddress: %s\n",
			rr.Name, GetResourceRecordType(rr.RecordType), GetResourceRecordClass(rr.Class), int(rr.TTL), int(rr.DataLength), rr.Data)
	}
	return fmt.Sprintf("\tName: %s \n\tType: %s \n\tClass: %s \n\tTTL: %d \n\tLength: %d \n\tName server: %s\n",
		rr.Name, GetResourceRecordType(rr.RecordType), GetResourceRecordClass(rr.Class), int(rr.TTL), int(rr.DataLength), rr.Data)
}
