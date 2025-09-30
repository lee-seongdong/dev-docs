# Kafka
대용량의 실시간 데이터 스트림을 처리하고 저장하기 위한 분산 메시징 시스템  
주로 이벤트 스트리밍 플랫폼으로 사용되며, 데이터의 생산자(Producer)와 소비자(Consumer) 사이에서 메시지를 안전하고 빠르게 전달하는 역할을 함.

## 1. 특징 및 용도
### 특징
- **분산 시스템**: 여러 서버(브로커)로 구성되어 높은 확장성과 내결함성을 가짐
- **pub/sub**: 생산자가 토픽에 메시지를 발행하면, 소비자가 해당 토픽을 구독하여 메시지를 받아감
- **내구성**: 메시지는 디스크에 저장되어 장애 발생 시에도 데이터 유실을 최소화
- **고성능**: 대용량의 데이터를 빠르게 처리할 수 있음

### 용도
- 실시간 로그 수집 및 분석
- 이벤트 소싱(Event Sourcing)
- 데이터 파이프라인 구축
- 마이크로서비스 간 비동기 통신


## 2. 핵심 구성요소
- [**Broker**](./Broker.md): 메시지를 저장하고 관리하는 서버
- [**Producer**](./Producer.md): 메시지를 생성하고 토픽에 전송
- [**Consumer**](./Consumer.md): 토픽에서 메시지를 읽어오는 클라이언트
- [**Zookeeper**](./Zookeeper.md) / [**KRaft**](./KRaft.md): 클러스터 메타데이터 관리 및 코디네이션


## 3. 핵심 개념
- [**Topic**](./Topic.md): 메시지가 저장되는 논리적 카테고리
- [**Partition**](./Partition.md): 토픽의 분산 저장 단위
- [**Consumer Group**](./Consumer%20Group.md): 메시지 처리를 분산하는 컨슈머 집합
- [**Offset**](./Offset.md): 파티션 내 메시지의 위치 정보
- [**Replication**](./Replication.md): 데이터 복제를 통한 내결함성


## 4. 메시지 처리
- [**Message Production**](./Message%20Production.md): 메시지 생성 및 전송 과정
- [**Message Consumption**](./Message%20Consumption.md): 메시지 읽기 및 처리 과정
- [**Stream Processing**](./Stream%20Processing.md): 실시간 스트림 데이터 처리


## 5. 운영 및 관리
- [**Cluster Management**](./Cluster%20Management.md): 클러스터 구성 및 관리
- [**Performance Tuning**](./Performance%20Tuning.md): 성능 최적화 방법
- [**Monitoring**](./Monitoring.md): 모니터링 및 로그 관리
- [**Security**](./Security.md): 보안 설정 및 인증
