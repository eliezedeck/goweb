package http

import (
	"log"
	"net"
)

// GetTCPListener creates and binds a listener to the specified address (in the
// for of host:port) and returns that listener. It is handy in combination with
// consumers such as the goji.ServeListener().
func GetTCPListener(address string) *net.TCPListener {
	addr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		log.Fatalf("Could not parse listen-address: %s", err)
	}
	l, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Fatalf("Could not open TCP for listening: %s", err)
	}
	return l
}
