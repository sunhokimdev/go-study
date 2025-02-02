package proxy

import (
	"io"
	"net"
	"sync"
	"testing"
)

// os.Stdout, *bytes.Buffer, *os.File 외에 io.Writer 인터페이스를 구현한 많은 객체들로부터 데이터를 프락시 할 수 있음
func proxy(from io.Reader, to io.Writer) error {
	fromWriter, fromIsWriter := from.(io.Writer)
	toReader, toIsReader := to.(io.Reader)

	if toIsReader && fromIsWriter {
		// 필요한 인터페이스를 모두 구현했으니, from과 to에 대응하는 프락시 생성
		go func() {
			_, _ = io.Copy(fromWriter, toReader)
		}()
	}

	_, err := io.Copy(to, from)

	return err
}

func TestProxy(t *testing.T) {
	var wg sync.WaitGroup

	// 서버는 "ping" 메시지를 대기하고 "pong" 메시지로 응답합니다.
	// 그 외 메시지는 동일하게 클라이언트로 에코잉됩니다.
	server, err := net.Listen("tcp", "127.0.0.1:")	// 연결 요청을 수신할 수 있는 서버를 초기화
	if err != nil {
		t.Fatal(err)
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			conn, err := server.Accept()
			if err != nil {
				return
			}

			go func(c net.Conn) {
				defer c.Close()

				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
						return
					}			
					
					// "ping"이라는 문자열을 받으면 "pon"으로 응답, 그 외 문자열은 그대로 되돌려줌(= echo server)
					switch msg := string(buf[:n]); msg {
						case "ping":
							_, err = c.Write([]byte("pong"))
						default:
							_, err = c.Write(buf[:n])
					}

					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}

						return
					}
				}
			}(conn)
		}
	}()

	// proxyServer는 메시지를 클라이언트 연결로부터 destinationServer로 프락시합니다.
	// destinationServer 서버에서 온 응답 메시지는 역으로 클라이언트에게 프락시됩니다.
	// 클라이언트와 목적지 서버 간의 메시지 전달을 처리해 주는 프락시 서버를 셋업 -> 클라이언트의 연결 요청을 수신
	proxyServer, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			conn, err := proxyServer.Accept()
			if err != nil {
				return
			}

			go func(from net.Conn) {
				defer from.Close()

				to, err := net.Dial("tcp", server.Addr().String())
				if err != nil {
					t.Error(err)
					return
				}

				defer to.Close()

				// 목적지 서버와의 연결을 수립하고 메시지를 프락싱
				err = proxy(from, to)
				if err != nil && err != io.EOF {
					t.Error(err)
				}
			}(conn)
		}
	}()

	conn, err := net.Dial("tcp", proxyServer.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	msgs := []struct{ Message, Replay string } {
		{"ping", "pong"},
		{"pong", "pong"},
		{"echo", "echo"},
		{"ping", "pong"},
	}

	for i, m := range msgs {
		_, err = conn.Write([]byte(m.Message))
		if err != nil {
			t.Fatal(err)
		}

		buf := make([]byte, 1024)

		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}

		actual := string(buf[:n])
		t.Logf("%q -> proxy -> %q", m.Message, actual)

		if actual != m.Replay {
			t.Errorf("%d: expected reply: %q; actual: %q", i, m.Replay, actual)
		}
	}
	_ = conn.Close()
	_ = proxyServer.Close()
	_ = server.Close()

	wg.Wait()
}
