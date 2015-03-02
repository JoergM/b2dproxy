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

		go proxyconnection(host, conn)
	}
}

func proxyconnection(host string, conn net.Conn) {

	log.Printf("Proxy %v to %s\n", conn.RemoteAddr(), host)
	io.Copy(conn, conn)

	conn.Close()
	//dial a remote port and send data

	//copy data from incoming to outgoing

	//handling responses

	//handling multiple connections (test with ab)
}
