# gRPC-measurement-tool

## Definition
- The Management of grpc channels connection states
  - **CONNECTING**: The channel is trying to establish TCP connection and TLS handshakes.
  - **READY**: The channel has successfully established a connection all the way through TLS handshake and protocol level (HTTP/2, etc)
  - **TRANSIENT_FAILURE**: There has been some transient failure (such as a TCP 3-way handshake timing out or a socket error.)
  - **IDLE**: This is the state where the channel is not even trying to create a connection because of a lack of new or pending RPCs.
  - **SHUTDOWN**: This channel has started shutting down. Any new RPCs should fail immediately. Pending RPCs may continue running till the application cancels them. Channels may enter this state either because the application explicitly requested a shutdown or if a non-recoverable error has happened during attempts to connect communicate


## 측정 데이터
- total request count   -rt
- timeout               -timeout
- tls 인증여부           -isTls
- call method           -call
- rps              

> go build ./client.go
> ./client  -tr 10 -timeout 1000 -call helloworld.Greeter.SayHello


## 결과
>> go run .\client.go -rps 99 -rt 1000
![image](https://user-images.githubusercontent.com/46148739/152949137-34b125dc-858a-46cf-b273-cf23dbd36625.png)


## Reference
- https://grpc.github.io/grpc/core/md_doc_connectivity-semantics-and-api.html
