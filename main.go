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
	fmt.Println(mResponse.Print())
}
