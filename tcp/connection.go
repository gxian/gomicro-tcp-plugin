package tcp

import (
	"io"
	"net"
)

const (
	bufSize = 1024
)

type connection struct {
	conn   net.Conn
	reader io.Reader
	closer io.Closer
	buf    []byte
	stop   chan bool
}

func newConnection(conn net.Conn, r io.Reader, c io.Closer) *connection {
	return &connection{
		conn:   conn,
		buf:    make([]byte, 0, bufSize),
		stop:   make(chan bool, 0),
		reader: r,
		closer: c,
	}
}

func (c *connection) Start() {
	// go read
	go c.read()
	// go write
}

func (c *connection) Stop(err error) {
	c.closer.Close()
}

func (c *connection) read() {
	var n int
	var err error
	for {
		select {
		case <-c.stop:
			c.Stop(nil)
			return
		default:
			n, err = c.reader.Read(c.buf[0:n])
			if err != nil {
				c.Stop(err)
			}
		}
	}
}

func (c *connection) Send(b []byte) error {
	_, err := c.conn.Write(b)
	if err != nil {
		return err
	}
	return nil
}
