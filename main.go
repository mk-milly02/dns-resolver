package main

import (
	"dns-resolver/resolver"
	"flag"
	"fmt"
	"log"
)

func main() {
	domain := flag.String("domain", "example.com", "domain to resolve")
	flag.Parse()

	// Create a new message
	m := resolver.NewMessage(*domain)
	query := m.BuildQuery()
	fmt.Println("Query:", query)

	// Send the query to the name server
	response, err := resolver.SendRequest([]byte(query))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Response:", response)

	// Check if the response is valid
	if !m.ValidateResponse(response) {
		log.Fatal("Invalid response")
	}
	fmt.Println("Response is valid")

	// Parse the response
	var rMessage resolver.Message
	rMessage.ParseResponse(response)
	fmt.Println(rMessage.String())
}
