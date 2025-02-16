package main

import (
	"dns-resolver/resolver"
	"fmt"
)

func main() {
	default_message := resolver.NewMessage()
	fmt.Println(default_message.BuildQuery())
}

//0016 0100 0001 0000 0000 0000 03646e7306676f6f676c6503636f6d00 0001 0001
//0016 0100 0001 0000 0000 0000 03646e7306676f6f676c6503636f6d00 0001 0001 000000000000000000000000000000000000000000000000000000000000000000000000