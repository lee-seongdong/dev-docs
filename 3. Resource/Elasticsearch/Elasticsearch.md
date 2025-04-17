## Elasticsearch
Apache Lucene 기반의 java 오픈소스 분산 검색엔진  
검색을 위해 단독으로 사용되기도 하지만, 주로 Elastic Stack으로 사용
- Elasticsearch: 데이터 저장 및 검색
- Logstash: 다양한 소스로부터 데이터 수집, 집계, 파싱하여 Elasticsearch로 전달
- Beats: 경량 데이터 수집기로 다양한 소스에서 데이터를 수집하여 Elasticsearch로 전달
- Kibana: Elasticsearch 시각화하여 관리 및 사용


### 주요 기능
- **강력한 검색 및 분석**: 구조화된 데이터와 비구조화된 데이터를 모두 처리 가능
- **분산 처리**: 대규모 데이터 처리 및 확장성 제공
- **실시간 데이터 처리**: 실시간으로 데이터 검색 및 분석 가능

### 구성요소
- **Cluster**: 여러 노드가 모여 하나의 시스템처럼 동작하는 집합
- **Node**: Elasticsearch 클러스터를 구성하는 기본 단위
- **Shard**: 데이터를 분산 저장하기 위한 물리적 단위
- **Replica**: 데이터 가용성과 장애 복구를 위한 샤드 복제본
- **Index**: 데이터를 저장하는 논리적 공간
- **Document**: Elasticsearch에서 저장되는 기본 데이터 단위