package game

// Message ...
type Message interface {
	HeaderLen() uint16
	ID() int32
	BodyLen() int32
	Header() []byte
	Body() []byte
}

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
