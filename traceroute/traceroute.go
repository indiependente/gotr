package traceroute

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const readTimeoutSec = 10

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// Tracer implements the Traceroute operation.
type Tracer struct {
	address string
	Out     chan string
}

// NewTracer returns a new Tracer.
func NewTracer(addr string) Tracer {
	return Tracer{
		address: addr,
		Out:     make(chan string),
	}
}

// Traceroute returns a string channel to range over in order to get the trace results
func (tr Tracer) Traceroute(maxTTL int) {
	go iCMPTraceroute(tr.address, maxTTL, tr.Out)
}

func iCMPTraceroute(address string, maxTTL int, outCh chan string) {
	// Resolve to IP
	ipaddr, err := net.ResolveIPAddr("ip4", address)
	checkError(err)

	c, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	checkError(err)
	defer c.Close()

	pc := ipv4.NewPacketConn(c)
	defer pc.Close()

	outCh <- fmt.Sprintf("%s Launching traceroute against %s (%s)", time.Now().Format("2006-01-02 15:04:05"), address, ipaddr.IP.String())
	outCh <- "\t#HOP\tREMOTE IP\t\tMSGLENGTH\tNAMES"
	t1 := time.Now()
	var i int
	for i = 1; i <= maxTTL; i++ {
		pc.SetTTL(i)
		sendPacket(pc, ipaddr)
		if readPacket(pc, outCh) {
			break
		}
	}
	if i >= maxTTL {
		outCh <- fmt.Sprintf(">>> Host unreachable in %d hops!", maxTTL)
	}

	t2 := time.Since(t1)
	outCh <- fmt.Sprintf("Time elapsed : %02dms", t2/time.Millisecond)
	close(outCh)
}

func sendPacket(pc *ipv4.PacketConn, ipaddress *net.IPAddr) {
	p, err := craftPacket()
	checkError(err)
	_, err = pc.WriteTo(p, nil, ipaddress)
	checkError(err)
}

func readPacket(pc *ipv4.PacketConn, out chan string) bool {
	buff := make([]byte, 256)
	pc.SetReadDeadline(time.Now().Add(readTimeoutSec * time.Second))
	n, _, peer, err := pc.ReadFrom(buff)
	if err != nil {
		fmt.Printf("Request %ds Timeout\n", readTimeoutSec)
		return false
	}
	m, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), buff[:n])
	checkError(err)
	ttl, err := pc.TTL()
	checkError(err)
	switch m.Type {
	case ipv4.ICMPTypeTimeExceeded: // hop
		logHop(ttl, peer, m, out)
	case ipv4.ICMPTypeDestinationUnreachable:
		out <- fmt.Sprintf("\t-:\t%v\t\t[%d bytes]\tDESTINATION UNREACHABLE\n", peer, m.Body.Len(1))
	case ipv4.ICMPTypeEchoReply: // destination
		logHop(ttl, peer, m, out)
		return true
	default:
		out <- fmt.Sprintf("%v: %+v", peer, m)
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

func logHop(ttl int, peer net.Addr, m *icmp.Message, outCh chan string) {
	names, _ := net.LookupAddr(peer.String())
	if len(names) > 0 {
		outCh <- fmt.Sprintf("\t%d:\t%v\t\t[%d bytes]\t%+v\n", ttl, peer, m.Body.Len(1), names)
	} else {
		outCh <- fmt.Sprintf("\t%d:\t%v\t\t[%d bytes]\n", ttl, peer, m.Body.Len(1))
	}
}
