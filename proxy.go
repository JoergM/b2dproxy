package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	proxyport("192.168.59.103", 8080)
	//create API for use in main

	//udp connections?
}

func proxyport(host string, port int64) {
	portstr := fmt.Sprintf(":%d", port)

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

		go proxyconnection(host, port, conn)
	}
}

func proxyconnection(host string, port int64, upstream net.Conn) {

	log.Printf("Proxy %v to %s\n", upstream.RemoteAddr(), host)

	downstream, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}

	log.Println(downstream)

	io.Copy(downstream, upstream)

	upstream.Close()
	downstream.Close()

	//copy data from incoming to outgoing

	//handling responses

	//handling multiple connections (test with ab)
}
