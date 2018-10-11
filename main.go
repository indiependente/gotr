package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/indiependente/gotr/traceroute"
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
	tr := traceroute.NewTracer(address)
	tr.Traceroute(maxTTL)
	for logMessage := range tr.Out {
		fmt.Println(logMessage)
	}
}
