package main

import (
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	// Listen on TCP port 35353 on all available unicast & anycast IP addresses.
	l, err := net.Listen("tcp", ":35353")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go proxy(conn)
	}
}

func proxy(conn net.Conn) {
	defer conn.Close()

	upstream, err := net.Dial("tcp", "cloudflare-dns.com:853")
	if err != nil {
		log.Print(err)
		return
	}
	defer upstream.Close()

	go io.Copy(upstream, conn)
	io.Copy(conn, upstream)
}

func copyToStderr(conn net.Conn) {
	// copyToStderr - echo all incoming data with number of bytes.
	defer conn.Close()
	for {
		// Defining a buffer of bytes.
		var buff [128]byte
		// Close connection if no data has been received within 5 seconds since last packet came into.
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, err := conn.Read(buff[:])
		if err != nil {
			log.Print(err)
			return
		}
		os.Stderr.Write(buff[:n])
	}
}