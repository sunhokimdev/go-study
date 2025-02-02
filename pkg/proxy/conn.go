package proxy

import (
	"io"
	"net"
)

/*
* io.Copy 함수
* 출발지 노드로부터 전송된 데이터를 목적지 노드로, 목적지 노드에서 전송된 데이터를 출발지 노드로 복제
* io.Copy 함수 동작 도중 reader로부터 모든 데이터를 읽었다는 의미의 io.EOF 에러 이외의 에러가 발생했을 때에만 error를 반환
* Go 1.11 이후 이 함수(+io.CopyN) 사용 시 두 매개변수 모두 *net.TCPConn 객체로 사용 시 데이터는 리눅스상의 유저 스페이스를 경유하지 않음
* -> 마치 리눅스 커널 상에서 애플리케이션을 거치지 않고 소켓상으로 한쪽에서 데이터를 읽어서 그대로 쓰는 것과 같다.
*/

// 두 네트워크 연결 간의 데이터 프록시
func proxyConn(source, destination string) error {
	connSource, err := net.Dial("tcp", source)
	if err != nil {
		return err
	}
	defer connSource.Close()

	connDestination, err := net.Dial("tcp", destination)
	if err != nil {
		return err
	}
	defer connDestination.Close()

	// connSource에 대응하는 connDestination
	go func() {
		// connDestination으로 데이터를 읽고, connSource으로 데이터를 쓴다.
		_, _ = io.Copy(connSource, connDestination)
	}()

	// connDestination으로 메시지를 보내는 connSource
	// 두 노드 중 하나의 연결이 끊어지면 io.Copy는 자동으로 종료되어 고루틴의 메모리 누수는 걱정 X
	_, err = io.Copy(connDestination, connSource)

	return err
}
