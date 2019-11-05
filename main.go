package main

import (
	"flag"
)

var (
	addr string
)

func init() {
	flag.StringVar(&addr, "address", "", "listen address")
}

func main() {
	flag.Parse()
}
