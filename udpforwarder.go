package main

import (
	"fmt"
	"log"
	"net"
)

type UDPForwarder struct {
	closed  chan bool
	host    string
	portstr string
}

func NewUDPForwarder(host string, port int) *UDPForwarder {
	return &UDPForwarder{host: host,
		portstr: fmt.Sprintf(":%d", port),
		closed:  make(chan bool),
	}
}

func (uf *UDPForwarder) start() error {

	log.Printf("New UDP Forwarder on Port %s, for host %s", uf.portstr, uf.host)

	udpAddr, _ := net.ResolveUDPAddr("udp", uf.portstr)
	_, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Println(err)
		return err
	}

	readChan := make(chan []byte)

	for {
		select {
		case <-uf.closed:
			break
		case pkg := <-readChan:
			go forwardPackage(pkg)
		}
	}

	return nil
}

func (uf *UDPForwarder) Stop() {
	uf.closed <- true
}

func forwardPackage(pkg []byte) {
}
