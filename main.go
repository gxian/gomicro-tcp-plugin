package main

import (
	"flag"
	"fmt"
	"gomicro-tcp-plugin/game"
	"gomicro-tcp-plugin/tcp"
)

var (
	addr string
)

func init() {
	flag.StringVar(&addr, "address", "", "listen address")
}

func main() {
	flag.Parse()
	codec := game.NewCodec()
	mux := game.NewMultiplexer(codec)
	srv := tcp.NewServer(addr, mux)
	err := srv.Run()
	if err != nil {
		fmt.Println(err)
	}
}
