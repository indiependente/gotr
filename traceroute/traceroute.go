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

// Traceroute returns a string channel to range over in order to get the trace results
func Traceroute(address string, maxTTL int) <-chan string {
	outCh := make(chan string)
	go iCMPTraceroute(address, maxTTL, outCh)
	return outCh
}

func iCMPTraceroute(address string, maxTTL int, outCh chan string) {
	// Resolve to IP
	ipaddr, err := net.ResolveIPAddr("ip4", address)
	checkError(err)

	c, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	checkError(err)
	defer c.Close()

	p := ipv4.NewPacketConn(c)

	outCh <- fmt.Sprintf("%s Launching traceroute against %s (%s)", time.Now().Format("2006-01-02 15:04:05"), address, ipaddr.IP.String())
	outCh <- "\t\t\t#HOP\tREMOTE IP\t\tMSGLENGTH"
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
		outCh <- fmt.Sprintf(">>> Host unreachable in %d hops!", maxTTL)
	}

	t2 := time.Now().Sub(t1)
	outCh <- fmt.Sprintf("Time elapsed : %02dms", t2/time.Millisecond)
	close(outCh)
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
