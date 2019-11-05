package main

import "sync"

// Encoder ...
type Encoder interface {
	Encode(Message) ([]byte, error)
}

// Decoder ...
type Decoder interface {
	Decode([]byte) (Message, error)
}

// Codec ...
type Codec interface {
	Encoder
	Decoder
}

// Sender ...
type Sender interface {
	Send(Message)
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
	mu       *sync.Mutex
}

// NewMultiplexer ...
func NewMultiplexer(c Codec) *Multiplexer {
	return &Multiplexer{
		codec:    c,
		handlers: make(map[int]Handler),
		mu:       &sync.Mutex{},
	}
}

// HandleFunc ...
func (m *Multiplexer) HandleFunc(id int, h Handler) {
	m.handlers[id] = h
}

// Read ...
func (m *Multiplexer) Read(b []byte) (int, error) {
	// 断包处理
	return 0, nil
}

// Close ...
func (m *Multiplexer) Close() error {
	return nil
}
