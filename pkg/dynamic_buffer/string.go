package dynamic_buffer

import (
	"encoding/binary"
	"errors"
	"io"
)

type String string

func (m String) Bytes() []byte { return []byte(m) }
func (m String) String() string { return string(m) }
func (m String) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, StringType) 	// 1바이트 타입
	if err != nil {
		return 0, err
	}

	var n int64 = 1
	err = binary.Write(w, binary.BigEndian, uint32(len(m)))	// 4바이트 크기
	if err != nil {
		return n, err
	}
	n += 4

	o, err := w.Write([]byte(m))								// 페이로드
	
	return n + int64(o), err
}

func (m *String) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)			// 1 바이트 타입
	if err != nil {
		return 0, err
	}

	var n int64 = 1
	if typ != StringType {
		return n, errors.New("invalid String")
	}

	var size uint32
	err = binary.Read(r, binary.BigEndian, &size)			// 4바이트 크기
	if err != nil {
		return n, err
	}

	n += 4

	buf := make([]byte, size)
	o, err := r.Read(buf)									// 페이로드
	if err != nil {
		return n, err
	}

	*m = String(buf)
	
	return n + int64(o), nil
}
