package dynamic_buffer

import (
	"encoding/binary"
	"io"
	"errors"
)

type Binary []byte

func (m Binary) Bytes() []byte { return m}
func (m Binary) String() string { return string(m) }

func (m Binary) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, BinaryType)	// 1바이트 타입
	if err != nil {
		return 0, nil
	}
	var n int64 = 1

	err = binary.Write(w, binary.BigEndian, uint32(len(m))) // 4바이트 크기

	n += 4
	o, err := w.Write(m)		// 페이로드

	return n + int64(o), err
}

func (m *Binary) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)	// 1바이트 타입
	if err != nil {
		return 0, err
	}

	var n int64 = 1
	if typ != BinaryType {
		return n, errors.New("invalid Binary")
	}

	var size uint32
	err = binary.Read(r, binary.BigEndian, &size)		// 4바이트 크기
	if err != nil {
		return n, err
	}

	n += 4
	if size > MaxPayloadSize {
		return n, ErrMaxloadSize
	}

	*m = make([]byte, size)
	o, err := r.Read(*m)		// 페이로드

	return n + int64(o), err
}


