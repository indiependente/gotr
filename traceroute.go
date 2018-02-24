package traceroute

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func Traceroute(address string, maxTTL int) {
	// Resolve to IP
	ipaddr, err := net.ResolveIPAddr("ip4", address)
	checkError(err)

	c, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	checkError(err)
	defer c.Close()

	p := ipv4.NewPacketConn(c)

	log.Printf("Launching traceroute against %s (%s)\n", address, ipaddr.IP.String())
	fmt.Printf("\t\t\t#HOP\tREMOTE IP\t\tMSGLENGTH\n")
	t1 := time.Now()
	var done bool

	for i := 1; i < maxTTL; i++ {
		p.SetTTL(i)
		sendPacket(p, ipaddr)
		checkError(err)
		done = readPacket(p)
		if done {
			break
		}
	}
	if !done {
		fmt.Printf(">>> Host unreachable in %d hops!\n", maxTTL)
	}

	t2 := time.Now().Sub(t1)
	log.Printf("Time elapsed : %02dms", t2/time.Millisecond)
}

func sendPacket(pc *ipv4.PacketConn, ipaddress *net.IPAddr) {
	p, err := craftPacket()
	checkError(err)
	_, err = pc.WriteTo(p, nil, ipaddress)
	checkError(err)
}

func readPacket(pc *ipv4.PacketConn) bool {
	buff := make([]byte, 1500)
	n, _, peer, err := pc.ReadFrom(buff)
	checkError(err)
	m, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), buff[:n])
	checkError(err)
	ttl, err := pc.TTL()
	checkError(err)
	switch m.Type {
	case ipv4.ICMPTypeTimeExceeded: // hop
		logHop(ttl, peer, m)
	case ipv4.ICMPTypeEchoReply: // destination
		logHop(ttl, peer, m)
		return true
	default:
		log.Printf("%v: %+v", peer, m)
	}
	return false
}

func craftPacket() ([]byte, error) {
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("GoFF"),
		},
	}
	wb, err := wm.Marshal(nil)
	return wb, err
}

func logHop(ttl int, peer net.Addr, m *icmp.Message) {
	names, _ := net.LookupAddr(peer.String())
	if len(names) > 0 {
		log.Printf("\t%d:\t%v\t\t[%d bytes]\t%+v\n", ttl, peer, m.Body.Len(1), names)
	} else {
		log.Printf("\t%d:\t%v\t\t[%d bytes]\n", ttl, peer, m.Body.Len(1))
	}
}
