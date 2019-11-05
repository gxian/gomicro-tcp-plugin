package main

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
	Handle(Message)
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
}

// NewMultiplexer ...
func NewMultiplexer(c Codec) *Multiplexer {
	return &Multiplexer{}
}

// HandleFunc ...
func (m *Multiplexer) HandleFunc(msg int, s Sender, h Handler) {

}

// Read ...
func (m *Multiplexer) Read(b []byte) (int, error) {
	return 0, nil
}

// Close ...
func (m *Multiplexer) Close() error {
	return nil
}
