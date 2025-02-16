package main

import (
	"dns-resolver/resolver"
	"encoding/hex"
	"fmt"
	"log"
)

func main() {
	default_message := resolver.NewMessage()
	query := default_message.BuildQuery()
	fmt.Println("Query:", query)
	hex_query, err := hex.DecodeString(query)
	if err != nil {
		log.Fatal(err)
	}
	response, err := resolver.SendRequest(hex_query)
	if err != nil {
		log.Fatal(err)
	}
	resolver.DumpResponse(response)
	fmt.Println("Response Validity:", default_message.ValidateResponse(response))
}
