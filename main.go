package main

import (
	"flag"
)

var (
	addr string
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
	flag.StringVar(&addr, "address", "", "listen address")
}

func main() {
	flag.Parse()
}
