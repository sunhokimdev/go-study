package dynamic_buffer

import (
	"bytes"
	"encoding/binary"
	"net"
	"reflect"
	"testing"
)

func TestPayloads(t *testing.T) {
	b1 := Binary("Clear is better than clever.")
	b2 := Binary("Don't panic.")

	s1 := String("Errors are values.")
	payloads := []Payload{&b1, &s1, &b2}

	listener, err := net.Listen("tcp", "127.0.0.1:8082")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		for _, p := range payloads {
			_, err = p.WriteTo(conn)			// 각 타입을 페이로드 슬라이스 형태로 전송
			if err != nil {
				t.Error(err)
				break
			}
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	for i := 0; i < len(payloads); i++ {
		actual, err := decode(conn)
		if err != nil {
			t.Fatal(err)
		}

		if expected := payloads[i]; !reflect.DeepEqual(expected, actual) {
			t.Errorf("value mismatch: %v != %v", expected, actual)
			continue
		}

		t.Logf("[%T] %[1]q", actual)
	}
}

func TestMaxPayloadSize(t *testing.T) {
	buf := new(bytes.Buffer)
	err := buf.WriteByte(BinaryType)

	if err != nil {
		t.Fatal(err)
	}

	// BinaryType 타입을 저장하는 바이트와 페이로드 크기가 최대 1GB임을 나타내는 4바이트의 부호 없는 정수를 포함한 bytes.Buffer를 생성
	err = binary.Write(buf, binary.BigEndian, uint32(1 << 30))
	if err != nil {
		t.Fatal(err)
	}

	var b Binary
	_, err = b.ReadFrom(buf)
	if err != ErrMaxloadSize {
		t.Fatalf("expected ErrMaxPayloadSize; actual: %v", err)
	}
}
