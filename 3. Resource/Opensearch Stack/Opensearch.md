# Opensearch
Apache Lucene 기반의 분산 검색엔진  
Elasticsearch의 오픈소스 라이선스 변경 이후 7.10.2를 포크해서 만든 오픈소스 검색엔진이다.  

## 1. 주요 기능
- **검색 및 분석**: 구조화된 데이터와 비구조화된 데이터를 모두 처리 가능
- **분산 처리**: 대규모 데이터 처리 및 확장성 제공
- **실시간 데이터 처리**: 실시간으로 데이터 검색 및 분석 가능

## 2. 구성요소
- [**Cluster**](./Cluster.md): 여러 노드가 모여 하나의 시스템처럼 동작하는 집합
- [**Node**](./Node.md): Opensearch 클러스터를 구성하는 기본 단위
- [**Index**](./Index.md): 데이터를 저장하는 논리적 공간
- [**Shard**](./Shard.md): 데이터를 분산 저장하기 위한 물리적 단위
- [**Replica**](./Shard.md): 데이터 가용성과 장애 복구를 위한 샤드 복제본
- [**Document**](./Document.md): Opensearch에 저장되는 기본 데이터 단위

## 3. 주요 플러그인
- **Index State Management (ISM)**: 인덱스 생명주기 자동 관리 (hot/warm/cold/delete)
- [**Cross Cluster Replication (CCR)**](./Cross%20Cluster%20Replication.md): 클러스터 간 데이터 복제로 재해 복구 및 이중화 지원
- **Cross Cluster Search**: 여러 클러스터에서 통합 검색 수행
- **Security**: 사용자 인증, 권한 관리, TLS/SSL 암호화 (기본 포함)
- **Alerting**: 데이터 모니터링 및 알림 기능 (임계값 기반 알림)
- **Snapshot and Restore**: 인덱스 및 클러스터 백업/복원 기능
- **ML Commons**: 머신러닝 기능 (이상 탐지, 분류, 예측)
- **Anomaly Detection**: 시계열 데이터의 이상 패턴 탐지
- **Performance Analyzer**: 클러스터 성능 분석 및 모니터링
- **Reports**: 대시보드와 시각화 리포트 생성 및 스케줄링
- **Notebooks**: SQL과 PPL을 사용한 대화형 데이터 분석
