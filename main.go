package main

import (
	"fmt"

	"github.com/indiependente/gotr/traceroute"
	"github.com/integrii/flaggy"
)

type tracerouteInput struct {
	Address string
	MaxTTL  int
}

func main() {
	input := flagParse()

	tr := traceroute.NewTracer(input.Address)
	tr.Traceroute(input.MaxTTL)
	for hopMessage := range tr.Hops() {
		fmt.Println(hopMessage)
	}
}

func flagParse() tracerouteInput {
	var version = "0.0.1"
	flaggy.SetName("GoTR")
	flaggy.SetDescription("Golang implementation of the traceroute command (ICMP based) (root access required)")
	var address string
	flaggy.AddPositionalValue(&address, "host", 1, true, "hostname to traceroute")

	maxTTL := 30
	flaggy.Int(&maxTTL, "t", "ttl", "Number of allowed hops")

	flaggy.SetVersion(version)
	flaggy.Parse()

	return tracerouteInput{
		Address: address,
		MaxTTL:  maxTTL,
	}
}
