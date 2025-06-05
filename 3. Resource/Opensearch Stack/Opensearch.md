## Opensearch
Apache Lucene 기반의 분산 검색엔진  
Elasticsearch의 오픈소스 라이선스 변경 이후 7.10.2를 포크해서 만든 오픈소스 검색엔진이다.  


### 주요 기능
- **검색 및 분석**: 구조화된 데이터와 비구조화된 데이터를 모두 처리 가능
- **분산 처리**: 대규모 데이터 처리 및 확장성 제공
- **실시간 데이터 처리**: 실시간으로 데이터 검색 및 분석 가능


### 구성요소
- **Cluster**: 여러 노드가 모여 하나의 시스템처럼 동작하는 집합
- **Node**: Opensearch 클러스터를 구성하는 기본 단위
- **Index**: 데이터를 저장하는 논리적 공간
- **Shard**: 데이터를 분산 저장하기 위한 물리적 단위
- **Replica**: 데이터 가용성과 장애 복구를 위한 샤드 복제본
- **Document**: Opensearch에 저장되는 기본 데이터 단위


### 주요 플러그인
- **Security**: 사용자 인증, 권한 관리, TLS/SSL 암호화 (기본 포함)
- **Alerting**: 데이터 모니터링 및 알림 기능 (임계값 기반 알림)
- **Index State Management (ISM)**: 인덱스 생명주기 자동 관리 (hot/warm/cold/delete)
- **Cross Cluster Replication (CCR)**: 클러스터 간 데이터 복제로 재해 복구 및 이중화 지원
- **Cross Cluster Search**: 여러 클러스터에서 통합 검색 수행
- **Snapshot and Restore**: 인덱스 및 클러스터 백업/복원 기능
- **ML Commons**: 머신러닝 기능 (이상 탐지, 분류, 예측)
- **Anomaly Detection**: 시계열 데이터의 이상 패턴 탐지
- **Performance Analyzer**: 클러스터 성능 분석 및 모니터링
- **Reports**: 대시보드와 시각화 리포트 생성 및 스케줄링
- **Notebooks**: SQL과 PPL을 사용한 대화형 데이터 분석


### 주요 설정 파일들
- **opensearch.yml**: 메인 설정 파일 (클러스터, 노드, 네트워크 등)
- **jvm.options**: JVM 힙 메모리 및 GC 설정
- **log4j2.properties**: 로깅 레벨 및 출력 형식 설정
- **opensearch-security/**: Security 플러그인 설정 파일들

### Opensearch 설정 (Hot-Warm Architecture 기반)
#### 1. Hot 노드 기본 설정
```yaml
cluster.name: opensearch-cluster
node.name: opensearch-hot-1

# 노드 역할 설정 (Hot 노드)
node.roles: [cluster_manager, data, ingest]
node.attr.tier: hot  # Hot-Warm architecture를 위한 커스텀 속성 (속성 명칭은 자유)

# 경로 설정
path.data: /var/lib/opensearch
path.logs: /var/log/opensearch

# 네트워크 설정
network.host: 0.0.0.0
http.port: 9200
transport.port: 9300

# 클러스터 구성
discovery.seed_hosts: ["opensearch-hot-1:9300", "opensearch-hot-2:9300", "opensearch-warm-1:9301"]
cluster.initial_cluster_manager_nodes: ["opensearch-hot-1", "opensearch-hot-2"]
```

#### 2. Warm 노드 기본 설정
```yaml
cluster.name: opensearch-cluster
node.name: opensearch-warm-1

# 노드 역할 설정 (Warm 노드 - 데이터만 담당)
node.roles: [data]
node.attr.tier: warm

path.data: /var/lib/opensearch
path.logs: /var/log/opensearch

network.host: 0.0.0.0
http.port: 9201
transport.port: 9301

# 클러스터 구성
discovery.seed_hosts: ["opensearch-hot-1:9300", "opensearch-hot-2:9300", "opensearch-warm-1:9301"]
cluster.initial_cluster_manager_nodes: ["opensearch-hot-1", "opensearch-hot-2"]
```

#### 3. Hot-Warm 할당 설정
```yaml
# 샤드를 할당할 때 tier 속성 기반으로 각 노드에 균등하게 할당
cluster.routing.allocation.awareness.attributes: tier
# 샤드 할당 시 강제로 다른 속성에 할당.
# 할당할 곳이 없을 때 클러스터는 yellow 상태로 전환되며 할당 대기하며, 다시 할당 가능해 지는 경우 복제본 할당
cluster.routing.allocation.awareness.force.tier.values: hot, warm
```

#### 4. ISM 플러그인 설정
```yaml
plugins.index_state_management.enabled: true
```

#### 5. Security 플러그인 설정
```yaml
plugins.security.ssl.transport.pemcert_filepath: certs/node.pem
plugins.security.ssl.transport.pemkey_filepath: certs/node-key.pem
plugins.security.ssl.transport.pemtrustedcas_filepath: certs/root-ca.pem
plugins.security.ssl.http.enabled: true
plugins.security.ssl.http.pemcert_filepath: certs/node.pem
plugins.security.ssl.http.pemkey_filepath: certs/node-key.pem
plugins.security.ssl.http.pemtrustedcas_filepath: certs/root-ca.pem
```



### 클러스터 및 노드 상태 확인
```bash
# 1. 클러스터 전체 상태 확인
# green(정상), yellow(경고), red(위험)
GET /_cluster/health

# 2. 클러스터 설정 확인
# 현재 클러스터 동적 설정 확인
GET /_cluster/settings

# 3. 노드 목록 및 역할 확인
# 노드별 CPU, 메모리, 디스크 사용량 및 역할 표시
GET /_cat/nodes?v

# 4. 노드 상세 정보
# 노드별 상세 성능 지표 및 통계
GET /_nodes/stats

# 5. 인덱스 목록 및 상태
# 인덱스별 문서 수, 크기, 상태 확인
GET /_cat/indices?v

# 6. 샤드 할당 상태
# 샤드가 어느 노드에 할당되었는지 확인
GET /_cat/shards?v
```

#### 주로 확인해야 할 상태 항목
| 항목 | 정상 상태 | 주의사항 |
|------|-----------|----------|
| **Cluster Health** | Green | Yellow: 레플리카 누락, Red: 프라이머리 샤드 누락 |
| **노드 상태** | 모든 노드 활성 | 노드 다운 시 즉시 복구 필요 |
| **디스크 사용량** | < 85% | 85% 초과 시 샤드 할당 중단, 95% 초과 시 읽기 전용 |
| **메모리 사용량** | Heap < 75% | 75% 초과 시 GC 빈발, 성능 저하 |
| **샤드 할당** | 모든 샤드 할당됨 | Unassigned 샤드 발생 시 원인 분석 필요 |
| **인덱스 상태** | 정상 open | Close 상태나 오류 발생 시 확인 |
| **Hot-Warm 밸런스** | tier별 균등 분산 | 한쪽 tier에 과부하 집중 시 재할당 |

