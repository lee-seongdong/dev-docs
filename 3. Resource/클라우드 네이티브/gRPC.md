- 구글에서 개발한 RPC 프레임워크
- 데이터 직렬화에 효율적인 protocol buffers를 사용함
- 사용법
	- 통신에 사용할 데이터와 서비스를 정의하는 프로토콜 버퍼 파일 작성(.proto)
		```proto
		syntax = "proto3";
	
		option go_package = "ex05/keyvalue"
	
		// 통신에 사용할 데이터
		message GetRequest {
			string key = 1;
			// repeated string value = 2; // 배열
		}
	
		message GetResponse {
			string value = 1;
		}
	
		// 통신하기위한 서비스
		service KeyValue {
			rpc Get(GetRequest) returns (GetResponse);
		}
		```
	- 프로토콜 버퍼 파일 컴파일. (각 언어별로 사용가능한 코드 생성)
		```bash
		$ protoc -I=. --go_out=. --go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		keyvalue.proto
		```
	- 서버 코드 및 클라이언트 코드 구현