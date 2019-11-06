package game

// Message ...
type Message interface {
	ID() int32
	Body() []byte
	Bytes() []byte
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
