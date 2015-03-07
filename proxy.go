package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func proxyPort(host string, port int) {
	portstr := fmt.Sprintf(":%d", port)
	log.Printf("New Proxy on Port %s for host %s", portstr, host)

	ln, err := net.Listen("tcp", portstr)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		pc := proxyConnection{
			upstream:       conn,
			downstreamHost: host,
			downstreamPort: port,
		}

		go pc.establish()
	}

	//todo return value that can be used to end proxyPort
}

type proxyConnection struct {
	upstream       net.Conn
	downstream     net.Conn
	downstreamHost string
	downstreamPort int
}

func (pc *proxyConnection) establish() {

	defer pc.upstream.Close()

	log.Printf("Proxy %v to %s\n", pc.upstream.RemoteAddr(), pc.downstreamHost)

	var err error
	pc.downstream, err = net.Dial("tcp", fmt.Sprintf("%s:%d", pc.downstreamHost, pc.downstreamPort))
	if err != nil {
		panic(err)
	}
	defer pc.downstream.Close()

	//in parallel copy responses back
	done := make(chan bool)
	go copyContent(pc.downstream, pc.upstream, done)
	go copyContent(pc.upstream, pc.downstream, done)

	//wait for one channel to finish
	<-done
}

func copyContent(in net.Conn, out net.Conn, done chan bool) {
	io.Copy(out, in)
	done <- true
}
