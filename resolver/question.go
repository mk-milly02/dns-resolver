package resolver

import "fmt"

// Question represents a DNS question section
type Question struct {
	name       string
	recordType uint16
	class      uint16
}

// String returns the string representation of the question
func (q Question) String() string {
	return fmt.Sprintf("%s%04x%04x", q.name, q.recordType, q.class)
}
