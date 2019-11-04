package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

var (
	host string
	port string
)

// Codec ...
type Codec interface {
	MsgID() int32
}

type rawMsg struct {
	header rawMsgHeader
}

type rawMsgHeader struct {
	bodySize  uint16
	checksum  uint16
	timestamp int32
	msgID     int32
}

func init() {
	flag.StringVar(&host, "host", "", "listen host")
	flag.StringVar(&port, "port", "35000", "listen port")
}

func main() {
	flag.Parse()

	var l net.Listener
	var err error
	l, err = net.Listen("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + host + ":" + port)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
			os.Exit(1)
		}
		//logs an incoming message
		fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	for {
		io.Copy(conn, conn)
	}
}
