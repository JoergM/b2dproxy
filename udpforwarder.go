package main

import (
	"fmt"
	"log"
	"net"
)

type UDPForwarder struct {
	stopChan   chan bool
	closed     bool
	readChan   chan []byte
	downstream *net.UDPAddr
	host       string
	portstr    string
	conn       *net.UDPConn
}

func NewUDPForwarder(host string, port int) *UDPForwarder {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}

	return &UDPForwarder{
		stopChan:   make(chan bool),
		closed:     false,
		readChan:   make(chan []byte),
		downstream: addr,
		host:       host,
		portstr:    fmt.Sprintf(":%d", port),
	}
}

func (uf *UDPForwarder) start() error {

	log.Printf("New UDP Forwarder on Port %s, for host %s", uf.portstr, uf.host)

	udpAddr, _ := net.ResolveUDPAddr("udp", uf.portstr)
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Println(err)
		return err
	}
	uf.conn = conn

	//just a guess, that 10k should be a sane upper limit of a package
	conn.SetReadBuffer(10240)

	go uf.readPackage()
	go uf.listenLoop()

	return nil
}

func (uf *UDPForwarder) listenLoop() {
	for {
		select {
		case <-uf.stopChan:
			uf.closed = true
			uf.conn.Close()
			break
		case pkg := <-uf.readChan:
			go uf.forwardPackage(pkg)
		}
	}
}

func (uf *UDPForwarder) Stop() {
	log.Printf("Stop listening to port %v\n", uf.portstr[1:])
	uf.stopChan <- true
}

func (uf *UDPForwarder) readPackage() {
	buf := make([]byte, 10240)
	for !uf.closed {
		n, _, err := uf.conn.ReadFromUDP(buf)
		if err != nil {
			if !uf.closed {
				log.Println(err)
			}
			continue
		}
		uf.readChan <- buf[:n]
	}
}

func (uf *UDPForwarder) forwardPackage(pkg []byte) {
	conn, err := net.DialUDP("udp", nil, uf.downstream)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(pkg)
	if err != nil {
		log.Println(err)
	}
}
