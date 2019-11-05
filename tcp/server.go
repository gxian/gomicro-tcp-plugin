package tcp

import (
	"io"
	"net"
	"sync"
	"time"
)

// Multiplexer ...
type Multiplexer interface {
	Read(io.Writer, []byte) (int, error)
	Close(io.Writer) error
}

type readerFunc func([]byte) (int, error)

func (r readerFunc) Read(b []byte) (int, error) {
	return r(b)
}

type closerFunc func() error

func (r closerFunc) Close() error {
	return r()
}

// Server ...
type Server struct {
	conns  map[io.Writer]*connection
	accept map[io.Writer]int64
	addr   string
	mux    Multiplexer
	done   chan bool
	mu     *sync.Mutex
}

// NewServer ...
func NewServer(addr string, mux Multiplexer) *Server {
	return &Server{
		conns: make(map[io.Writer]*connection),
		addr:  addr,
		done:  make(chan bool, 0),
		mux:   mux,
		mu:    &sync.Mutex{},
	}
}

// Run ...
func (s *Server) Run() error {
	var l net.Listener
	var err error
	l, err = net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	for {
		select {
		case <-s.done:
			l.Close()
			// clean up
			return nil
		default:
			conn, err := l.Accept()
			if err != nil {
				return err
			}
			// TODO: accept timeout
			s.handleAccept(conn)
		}
	}
}

func (s *Server) muxReader(w io.Writer) readerFunc {
	delete(s.accept, w)
	return func(b []byte) (int, error) {
		return s.mux.Read(w, b)
	}
}

func (s *Server) muxCloser(w io.Writer) closerFunc {
	// del s.conns
	delete(s.conns, w)
	delete(s.accept, w)
	return func() error {
		return s.mux.Close(w)
	}
}

func (s *Server) handleAccept(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c := newConnection(conn, s.muxReader(conn), s.muxCloser(conn))
	s.conns[conn] = c
	s.accept[conn] = time.Now().Unix()
	go c.Start()
}
