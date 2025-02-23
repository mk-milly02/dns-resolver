package resolver

import (
	"fmt"
	"net"
)

const DEFAULT_NAME_SERVER = "198.41.0.4"
const GOOGLE_DNS_SERVER = "8.8.8.8"
const DEFAULT_PORT = "53"

// SendRequest sends the DNS query to the name server
func SendRequest(query []byte) ([]byte, error) {
	conn, err := net.Dial("udp", DEFAULT_NAME_SERVER+":"+DEFAULT_PORT)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	_, err = conn.Write(query)
	if err != nil {
		return nil, err
	}
	response := make([]byte, 512)
	_, err = conn.Read(response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// ReceiveResponse receives and parses the response from the name server
func ReceiveResponse(response []byte) {
	// Parsing logic here
	fmt.Println("Response received:", response)
}

// DumpResponse outputs the response in a readable format
func DumpResponse(response []byte) {
	// Dumping logic here
	fmt.Printf("Response: %x\n", response)
}
