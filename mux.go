package main

import (
	"bytes"
	"io"
	"sync"
)

// Encoder ...
type Encoder interface {
	Encode(Message) ([]byte, error)
}

// Decoder ...
type Decoder interface {
	Decode([]byte) (int, Message, error)
}

// Codec ...
type Codec interface {
	Encoder
	Decoder
}

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

// Header ...
type Header interface {
	BodySize() uint16
	Checksum() uint16
	Timestamp() int32
	MsgID() int32
}

// Message ...
type Message interface {
	Header() Header
	Body() []byte
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
	// 断包处理
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
		// 包体完整
		m.mu.Lock()
		defer m.mu.Unlock()
		h, ok := m.handlers[int(msg.Header().MsgID())]
		if !ok {
			return n, nil
		}
		h.Handle(msg, m.wrapSender(w))

		// 偏移buf
		buf.Next(n)
	}

	return 0, nil
}

// Close ...
func (m *Multiplexer) Close(w io.Writer) error {
	// clear writer buf
	return nil
}
