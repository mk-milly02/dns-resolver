package main

import (
	"dns-resolver/resolver"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
)

func main() {
	domain := flag.String("d", "example.com", "domain to resolve")
	flag.Parse()

	// Create a new message
	mRequest := resolver.NewMessage(*domain)
	query := mRequest.BuildQuery()

	// Send the query to the name server
	hex_query, err := hex.DecodeString(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Querying %s for %s\n", resolver.DEFAULT_NAME_SERVER, *domain)
	response, err := resolver.SendRequest(hex_query)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the response is valid
	if !mRequest.ValidateResponse(response) {
		log.Fatal("Invalid response")
	}

	// Parse the response
	var mResponse resolver.Message
	mResponse.ParseResponse(response)

	for len(mResponse.Answer) == 0 {
		recursiveQuery(&mResponse, hex_query, domain)
	}

	fmt.Println(mResponse.Print())
}

func recursiveQuery(m *resolver.Message, hex_query []byte, domain *string) {
	if len(m.Answer) == 0 && len(m.Authority) != 0 {
		var nameServer string
		for _, rr := range m.Additional {
			if rr.RecordType == 0x01 {
				nameServer = rr.Data
				break
			}
		}
		if nameServer != "" {
			fmt.Printf("Querying %s for %s\n", nameServer, *domain)
			response, err := resolver.SendRequestTo(hex_query, nameServer)
			if err != nil {
				log.Fatal(err)
			}
			if !m.ValidateResponse(response) {
				log.Fatal("Invalid response")
			}
			m.ParseResponse(response)
		}
	}
}
