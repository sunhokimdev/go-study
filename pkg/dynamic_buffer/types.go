package dynamic_buffer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	BinaryType uint8 = iota + 1
	StringType
)

var MaxPayloadSize uint32 = 10 << 20		// 10MB

var ErrMaxloadSize = errors.New("maximum payload size exceeded")

// 각 타입별 메시지들이 구현해야 할 인터페이스
type Payload interface {
	fmt.Stringer
	io.ReaderFrom
	io.WriterTo
	Bytes()	[]byte
}

func decode(r io.Reader) (Payload, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return nil, err
	}

	var payload Payload

	switch typ {
		case BinaryType:
			payload = new(Binary)
		case StringType:
			payload = new(String)
		default:
			return nil, errors.New("unknown type")
	}

	_, err = payload.ReadFrom(io.MultiReader(bytes.NewReader([]byte{typ}), r))
	if err != nil {
		return nil, err
	}

	return payload, nil
}
