package resolver

import (
	"encoding/hex"
	"fmt"
)

// Message represents a DNS message
type Message struct {
	Header     Header
	Question   []Question
	Answer     []ResourceRecord
	Authority  []ResourceRecord
	Additional []ResourceRecord
}

// NewMessage creates a new DNS message with default values.
// The message includes a header, a question section with a query for the given name,
// and empty answer, authority, and additional sections.
func NewMessage(name string) Message {
	return Message{
		Header:   NewHeader(),
		Question: []Question{{EncodeDomainName(name), 0x01, 0x01}},
	}
}

// BuildQuery builds the DNS query from the message
func (m Message) BuildQuery() string {
	return m.Header.String() + m.Question[0].String()
}

// ValidateResponse validates the response from the name server
func (m Message) ValidateResponse(response []byte) bool {
	return hex.EncodeToString(response[:2]) == fmt.Sprintf("%04x", m.Header.id)
}

// IsAResponse checks if the message is a response
func (m Message) IsAResponse() bool {
	return m.Header.flags&0x8000 == 0x8000
}

// ParseResponse parses the response from the name server
func (m *Message) ParseResponse(response []byte) {
	header, offset := ParseHeader(response)
	m.Header = header
	questions, newOffset := ParseQuestion(response, int(header.queryCount), offset)
	m.Question = questions
	answers, newOffset := ParseResourceRecord(response, int(header.answerCount), newOffset)
	m.Answer = answers
	authority, newOffset := ParseResourceRecord(response, int(header.authorityCount), newOffset)
	m.Authority = authority
	additional, _ := ParseResourceRecord(response, int(header.additionalCount), newOffset)
	m.Additional = additional
}

// Print returns the string representation of the message
func (m Message) Print() string {
	var str string
	str += "Header:\n" + m.Header.Print() + "\n"
	str += "Query:\n"
	for _, q := range m.Question {
		str += q.Print() + "\n"
	}
	str += "Answer:\n"
	for _, a := range m.Answer {
		str += a.Print() + "\n"
	}
	str += "Authority:\n"
	for _, a := range m.Authority {
		str += a.Print() + "\n"
	}
	str += "Additional:\n"
	for _, a := range m.Additional {
		str += a.Print() + "\n"
	}
	return str
}
