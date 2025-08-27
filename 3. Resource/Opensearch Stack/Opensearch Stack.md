# Opensearch Stack (구 ELK Stack)
Opensearch 검색엔진을 중심으로 한 로그 수집/분석용 기술스택

## 1. 핵심 구성요소
- [**Opensearch**](./Opensearch.md): 데이터 저장 및 검색 (Elasticsearch 대체)
- [**Filebeat**](./Filebeat.md): 파일 데이터 수집 및 전송
- [**Logstash**](./Logstash.md): 다양한 소스로부터 데이터 수집, 집계 파싱
- **Opensearch Dashboard**: 데이터 시각화 및 관리 (Kibana 대체)

## 2. 주요 개념
- [**Cluster**](./Cluster.md) - 노드들의 집합
- [**Node**](./Node.md) - 클러스터 구성 단위
- [**Index**](./Index.md) - 데이터 저장 공간
- [**Shard**](./Shard.md) - 분산 저장 단위

## 3. 주요 데이터 작업
- [**Indexing**](./Indexing.md) - 데이터 색인
- [**Query**](./Query.md) - 데이터 검색
- [**Aggregation**](./Aggregation.md) - 데이터 집계