package main

import (
	"dns-resolver/resolver"
	"fmt"
)

func main() {
	default_message := resolver.NewMessage()
	fmt.Println(default_message.ConvertToHexString())
}
