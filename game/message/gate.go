package message

import (
	"encoding/binary"

	"gomicro-tcp-plugin/game"
)

const (
	// HeaderLen header length
	HeaderLen      = 12
	bodySizeBegin  = 0
	bodySizeEnd    = 2
	checksumBegin  = bodySizeEnd
	checksumEnd    = 4
	timestampBegin = checksumEnd
	timestampEnd   = 8
	msgIDBegin     = timestampEnd
	msgIDEnd       = 12
)

// Gate ...
type Gate struct {
	BodySize  uint16
	Checksum  uint16
	Timestamp int32
	MsgID     int32
	Payload   []byte
}

// ID ...
func (g *Gate) ID() int32 {
	return g.MsgID
}

// Body ...
func (g *Gate) Body() []byte {
	return g.Payload
}

// Bytes ...
func (g *Gate) Bytes() []byte {
	b := make([]byte, HeaderLen+g.BodySize)
	binary.LittleEndian.PutUint16(b[:bodySizeBegin], g.BodySize)
	binary.LittleEndian.PutUint16(b[checksumBegin:checksumEnd], g.Checksum)
	binary.LittleEndian.PutUint32(b[timestampBegin:timestampEnd], uint32(g.Timestamp))
	binary.LittleEndian.PutUint32(b[msgIDBegin:msgIDEnd], uint32(g.MsgID))
	copy(b[HeaderLen:], g.Payload)
	return b
}

type gateCodec struct {
}

// NewGateCodec ...
func NewGateCodec() game.Codec {
	return &gateCodec{}
}

// Encode ...
func (g *gateCodec) Encode(m game.Message) ([]byte, error) {
	return m.Bytes(), nil
}

// Decode ...
func (g *gateCodec) Decode(b []byte) (int, game.Message, error) {
	total := int32(len(b))
	if total < HeaderLen {
		return 0, nil, nil
	}
	bodySize := binary.LittleEndian.Uint16(b[:bodySizeEnd])
	msgLen := int(HeaderLen + bodySize)
	if total < int32(msgLen) {
		// incomplete
		return 0, nil, nil
	}
	checksum := binary.LittleEndian.Uint16(b[checksumBegin:checksumEnd])
	timestamp := int32(binary.LittleEndian.Uint32(b[timestampBegin:timestampEnd]))
	msgID := int32(binary.LittleEndian.Uint32(b[msgIDBegin:msgIDEnd]))
	if timestamp == 0 {
		msgID = 0
	} else {
		msgID = msgID / (timestamp%10000 + 1)
	}
	msg := &Gate{
		BodySize:  bodySize,
		Checksum:  checksum,
		Timestamp: timestamp,
		MsgID:     msgID,
	}

	msg.Payload = make([]byte, msg.BodySize)
	copy(msg.Payload, b[HeaderLen:msgLen])
	return msgLen, msg, nil
}
