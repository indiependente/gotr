package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/indiependente/gotr/traceroute"
)

func isSupportedProto(p string) bool {
	switch p {
	case "icmp", "udp":
		return true
	default:
		return false
	}
}

func main() {
	// default values
	maxTTL := 30
	argslen := len(os.Args)
	proto := "icmp"

	if argslen != 2 && argslen != 3 && argslen != 4 {
		log.Fatal("Usage: ./gotr hostname [TTL] [PROTOCOL]")
	}
	address := os.Args[1]
	if argslen == 3 {
		var err error
		maxTTL, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
	} else if argslen == 4 {
		proto = os.Args[3]
		if !isSupportedProto(proto) {
			log.Fatalf("Unsupported protocol %s", proto)
		}
	}
	for logMessage := range traceroute.Traceroute(address, maxTTL, proto) {
		fmt.Println(logMessage)
	}
}
