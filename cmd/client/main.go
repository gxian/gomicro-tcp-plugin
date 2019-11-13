package main

import (
	"flag"
	"fmt"
	"gomicro-tcp-plugin/tcp"
	"gomicro-tcp-plugin/tcp/message"
)

var (
	addr string
)

func init() {
	flag.StringVar(&addr, "address", "", "listen address")
}

func main() {
	flag.Parse()
	srv := tcp.NewServer(
		addr,
		tcp.NewMultiplexer(
			message.NewGateCodec(),
		),
	)
	err := srv.Run()
	if err != nil {
		fmt.Println(err)
	}
}
