package tcp

import (
	"net"
	"sync"
)

// SendFunc ...
type SendFunc func(int, []byte) error

// Multiplexer ...
type Multiplexer interface {
	Init(SendFunc)
	Read(int, []byte) (int, error)
	Close(int) error
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
	conns map[int]*connection
	host  string
	port  string
	mux   Multiplexer
	gen   *idgen
	done  chan bool
	mu    *sync.Mutex
}

// NewServer ...
func NewServer(host, port string) *Server {
	return &Server{
		conns: make(map[int]*connection),
		host:  host,
		port:  port,
		gen:   newIDGen(),
		done:  make(chan bool, 0),
		mu:    &sync.Mutex{},
	}
}

// Run ...
func (s *Server) Run() error {
	var l net.Listener
	var err error
	l, err = net.Listen("tcp", s.host+":"+s.port)
	if err != nil {
		return err
	}

	for {
		select {
		case <-s.done:
			// clean up
			return nil
		default:
			conn, err := l.Accept()
			if err != nil {
				return err
			}
			// TODO: accept timeout处理
			s.handleRequest(conn)
		}
	}
}

func (s *Server) muxReader(id int) readerFunc {
	return func(b []byte) (int, error) {
		return s.mux.Read(id, b)
	}
}

func (s *Server) muxCloser(id int) closerFunc {
	// del s.conns
	return func() error {
		return s.mux.Close(id)
	}
}

func (s *Server) handleRequest(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.gen.Get()
	c := newConnection(conn, s.muxReader(id), s.muxCloser(id))
	s.conns[id] = c
	go c.Start()
}
