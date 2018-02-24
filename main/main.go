package main

import (
	"goexcrs/traceroute"
	"log"
	"os"
	"strconv"
)

func main() {
	maxTTL := 30
	argslen := len(os.Args)
	if argslen != 2 && argslen != 3 {
		log.Fatal("Usage: ./gotr hostname [TTL]")
	}
	address := os.Args[1]
	if argslen == 3 {
		var err error
		maxTTL, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
	}
	traceroute.Traceroute(address, maxTTL)
}
