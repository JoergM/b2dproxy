package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	//listen to a port and print incoming data
	listenonport(8080)

	//dial a remote port and send data

	//copy data from incoming to outgoing

	//handling responses

	//handling multiple connections (test with ab)

	//create API for use in main

	//udp connections?
}

func listenonport(port int64) {
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

		go func(c net.Conn) {
			log.Println(c.RemoteAddr())
			io.Copy(c, c)
			c.Close()
		}(conn)
	}
}
