package main

import (
	"dns-resolver/resolver"
	"encoding/hex"
	"fmt"
	"log"
)

func main() {
	request_messgae := resolver.NewMessage()
	query := request_messgae.BuildQuery()
	fmt.Println("Query:", query)
	hex_query, err := hex.DecodeString(query)
	if err != nil {
		log.Fatal(err)
	}
	response_bytes, err := resolver.SendRequest(hex_query)
	if err != nil {
		log.Fatal(err)
	}
	resolver.DumpResponse(response_bytes)
	fmt.Println("Response Validity:", request_messgae.ValidateResponse(response_bytes))
	response_message := resolver.ParseResponse(response_bytes, request_messgae)
	for _, record := range response_message.Answer {
		fmt.Println(record.String())
	}
}
