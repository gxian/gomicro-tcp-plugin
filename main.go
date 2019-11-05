package main

import (
	"flag"
	"fmt"
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
	srv := tcp.NewServer(addr)
	err := srv.Run()
	if err != nil {
		fmt.Println(err)
	}
}
