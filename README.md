# 4장 TCP 데이터 전송하기

#### TCP 연결 확인하기
```shell
ss -tn src :8080
```
* `-t`: TCP 연결만 표시
* `-n`: 호스트명 대신 숫자로 표시(DNS 조회 방지)
* `src :8080`: 8080 포트에서 수신한 연결만 표시


### 데이터를 읽고 쓰는 도중 에러 처리
네트워크 연결로부터 데이터를 읽고 쓸 때 발생하는 모든 에러가 치명적인 것은 아닙니다. 다시 연결이 복구 될 수 있는 에러들도 존재합니다. 예를 들어, 네트워크 상태가 좋지 않아서 수신자의 ACK 패킷이 지연되고 수신 대기 중 연결 시간이 초과된 상태의 네트워크 연결에 데이터 쓰기를 시도하면 일시적인 에러가 발생할 수 있습니다. 이는 송신자와 수신자 사이에 네트워크 케이블이 일시적으로 뽑힌 경우 등의 상황에 발생할 수 있습니다. 그런 경우 네트워크 연결은 아직 활성화 상태이므로 에러를 처리하고 복구 시도를 하거나, 연결을 우아하게 종료할 수 있습니다.
#### 네트워크 연결에 데이터를 쓰는 도중 발생하는 일시적인 에러 확인 방법
```go
var (
    err     error
    n       int
    i       = 7         // 최대 재시도수 
)

// 쓰기 관련된 코드를 for 루프로 감싸기
for ; i > 0; i-- {
    n, err = conn.Write([]byte("hello world"))
    if err != nil {
        // 타입 어설션을 통해 에러가 net.Error 인터페이스를 구현했는지, 그리고 에러가 일시적인지 확인 후 일시적이면 또 다른 쓰기 시도
        if nErr, ok := err.(net.Error); ok && nErr.Temporary() {
            log.Println("temporary error:", nErr)
            time.Sleep(10 * time.Second)
            continue
        }
        return err
    }
    break
}

if i == 0 {
    return errors.New("temporary write failure threshold exceeded")
}

log.Printf("wrote %d bytes to %s\n", n, conn.RemoteAddr())
```

### ICMP(Internet Control Message Protocol)

IP 네트워크에서 오류 메시지와 제어 정보를 전달하는 데 사용되는 프로토콜이다. 주로 네트워크 상태를 모니터링하거나 네트워크 문제를 진단하는 데 사용된다.
#### ICMP 메시지 형식


```shell
$ go test -run <method name>
go test -race -v <method name>
```
* 테스트 간 race detector 활성화: `-race` 플래그
* race detector는 프로그램 상에서 주의해야 할 교착 상태를 감지

# 모르는 것
* sync.WaitGroup
