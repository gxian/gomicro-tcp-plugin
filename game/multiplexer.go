package game

import (
	"bytes"
	"io"
	"sync"
)

// SendFunc ...
type SendFunc func(Message) error

// Send ...
func (s SendFunc) Send(m Message) error {
	return s(m)
}

// Sender ...
type Sender interface {
	Send(Message) error
}

// Handler ...
type Handler interface {
	Handle(Message, Sender)
}

// Multiplexer ...
type Multiplexer struct {
	codec    Codec
	handlers map[int]Handler
	bufs     map[io.Writer]*bytes.Buffer
	mu       *sync.Mutex
}

// NewMultiplexer ...
func NewMultiplexer(c Codec) *Multiplexer {
	return &Multiplexer{
		codec:    c,
		handlers: make(map[int]Handler),
		bufs:     make(map[io.Writer]*bytes.Buffer),
		mu:       &sync.Mutex{},
	}
}

func (m *Multiplexer) wrapSender(w io.Writer) SendFunc {
	return func(message Message) error {
		b, err := m.codec.Encode(message)
		if err != nil {
			return err
		}
		_, err = w.Write(b)
		return err
	}
}

// HandleFunc ...
func (m *Multiplexer) HandleFunc(id int, h Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[id] = h
}

// Read ...
func (m *Multiplexer) Read(w io.Writer, b []byte) (int, error) {
	// write to buf
	m.mu.Lock()
	buf, ok := m.bufs[w]
	if !ok {
		buf = bytes.NewBuffer([]byte{})
		m.bufs[w] = buf
	}
	m.mu.Unlock()

	n, err := buf.Write(b)
	if err != nil {
		return n, err
	}

	n, msg, err := m.codec.Decode(buf.Bytes())
	if err != nil {
		// kick
		return n, err
	}

	if n != 0 {
		// handle message
		m.mu.Lock()
		defer m.mu.Unlock()
		h, ok := m.handlers[int(msg.ID())]
		if !ok {
			return n, nil
		}
		h.Handle(msg, m.wrapSender(w))

		// offset buf
		buf.Next(n)
	}

	return 0, nil
}

// Close ...
func (m *Multiplexer) Close(w io.Writer) error {
	// clear writer buf
	return nil
}
