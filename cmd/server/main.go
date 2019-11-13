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

type msg struct {
}

// 测试是否connection accept之后超时会被踢掉;包体异常大小处理
// 测试百级并发连接下消息处理速度和时延
// 完善日志信息用于正式联调

func init() {
	flag.StringVar(&addr, "address", "", "server address")
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
