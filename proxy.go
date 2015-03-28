package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func NewSinglePortProxy(host string, port int, ptype string) (*SinglePortProxy, error) {
	portstr := fmt.Sprintf(":%d", port)

	if ptype == "udp" {
		startUDPForwarder(host, portstr)
	}

	ln, err := net.Listen("tcp", portstr)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Printf("New TCP Proxy on Port %s, for host %s", portstr, host)

	spp := SinglePortProxy{
		host:     host,
		port:     port,
		listener: ln,
	}

	go spp.listenForConnections()

	return &spp, nil
}

type SinglePortProxy struct {
	host     string
	port     int
	listener net.Listener
	stopped  bool
}

func (proxy *SinglePortProxy) listenForConnections() {
	for {
		defer func() {
			if r := recover(); r != nil {
				//it's ok, Accept probably stopped because port is closed
			}
		}()

		conn, err := proxy.listener.Accept()
		if err != nil {
			if proxy.stopped {
				break
			} else {
				panic(err)
			}
		}

		pc := proxyConnection{
			upstream:       conn,
			downstreamHost: proxy.host,
			downstreamPort: proxy.port,
		}

		go pc.establish()
	}
}

func (proxy SinglePortProxy) stopListen() {
	log.Printf("Stop listening to port %v\n", proxy.port)
	proxy.stopped = true
	proxy.listener.Close()
}

type proxyConnection struct {
	upstream       net.Conn
	downstream     net.Conn
	downstreamHost string
	downstreamPort int
}

func (pc *proxyConnection) establish() {

	defer pc.upstream.Close()

	var err error
	pc.downstream, err = net.Dial("tcp", fmt.Sprintf("%s:%d", pc.downstreamHost, pc.downstreamPort))
	if err != nil {
		//TODO if it is connection refused pass this on to the client
		log.Println(err)
		return
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

func startUDPForwarder(host, portstr string) {

	log.Printf("New UDP Forwarder on Port %s, for host %s", portstr, host)
	//TODO implement
	//	udpAddr, _ := net.ResolveUDPAddr("udp", portstr)
	//	ln, err = net.ListenUDP("udp", udpAddr)
	//	if err != nil {
	//		log.Println(err)
	//		return nil, err
	//	}
}
