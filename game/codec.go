package game

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

type defaultCodec struct {
}

func (d *defaultCodec) Encode(m Message) ([]byte, error) {
	return []byte{}, nil
}

func (d *defaultCodec) Decode(b []byte) (int, Message, error) {
	return 0, nil, nil
}

// NewCodec ...
func NewCodec() Codec {
	return &defaultCodec{}
}
