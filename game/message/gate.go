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

// HeaderLen ...
func (g *Gate) HeaderLen() uint16 {
	return HeaderLen
}

// ID ...
func (g *Gate) ID() int32 {
	return 0
}

// BodyLen ...
func (g *Gate) BodyLen() int32 {
	return 0
}

// Header ...
func (g *Gate) Header() []byte {
	return []byte{}
}

// Body ...
func (g *Gate) Body() []byte {
	return []byte{}
}

type gateCodec struct {
}

// NewGateCodec ...
func NewGateCodec() game.Codec {
	return &gateCodec{}
}

// Encode ...
func (g *gateCodec) Encode(m game.Message) ([]byte, error) {
	return []byte{}, nil
}

// Decode ...
func (g *gateCodec) Decode(b []byte) (int, game.Message, error) {
	total := int32(len(b))
	if total < HeaderLen {
		return 0, nil, nil
	}
	bodySize := binary.LittleEndian.Uint16(b[:bodySizeEnd])
	// 长度验证
	if total < int32(bodySize+HeaderLen) {
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

	/*
		// decod msgid
		if tm == 0 {
			id = 0
		} else {
			id = id / (tm%10000 + 1)
		}
		// body
		body := make([]byte, bodysz)
		copy(body, bs[12:12+bodysz])
		// 构造msg
		msg := &Msg{
			BodySize:  uint16(bodysz),
			CheckSum:  uint16(sum),
			TimeStamp: tm,
			MsgID:     id,
			Body:      body,
		}*/
	return 0, msg, nil
}
