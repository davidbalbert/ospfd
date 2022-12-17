package main

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"time"
)

// https://www.rfc-editor.org/rfc/rfc2328.html

var allSPFRouters = netip.MustParseAddr("224.0.0.5")
var allDRouters = netip.MustParseAddr("224.0.0.6")

// OSPF capabilities
const (
	capE  = 1 << 1
	capMC = 1 << 2
	capNP = 1 << 3
	capEA = 1 << 4
	capDC = 1 << 5
)

func toNetAddr(addr netip.Addr) net.Addr {
	return &net.IPAddr{IP: addr.AsSlice()}
}

func to4(addr netip.Addr) []byte {
	b := addr.As4()
	return b[:]
}

func mustAddrFromSlice(b []byte) netip.Addr {
	addr, ok := netip.AddrFromSlice(b)
	if !ok {
		panic("mustAddrFromSlice: slice should be either 4 or 16 bytes, but got " + fmt.Sprint(len(b)))
	}
	return addr
}

func tickImmediately(d time.Duration) <-chan time.Time {
	c := make(chan time.Time)

	go func() {
		c <- time.Now()
		for t := range time.Tick(d) {
			c <- t
		}
	}()

	return c
}

func main() {
	fmt.Printf("Starting ospfd with uid %d\n", os.Getuid())

	config, err := NewConfig("192.168.200.1")
	if err != nil {
		panic(err)
	}

	if err := config.AddNetwork("192.168.105.0/24", "0.0.0.0"); err != nil {
		panic(err)
	}

	if err := config.AddInterface("bridge100", "0.0.0.0", networkPointToMultipoint, 10, 40, 5); err != nil {
		panic(err)
	}

	instance, err := NewInstance(config)
	if err != nil {
		panic(err)
	}

	instance.Run()
}
