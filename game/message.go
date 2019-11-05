package game

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
