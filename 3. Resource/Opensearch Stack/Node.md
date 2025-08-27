# Node
클러스터를 구성하는 기본 단위. 각 노드는 데이터 저장, 검색, 분석, 클러스터 관리 등 다양한 역할을 함.  

Opensearch의 실행 단위는 **Node**이며, 각 Node는 독립적으로 실행되는 프로세스  
하나의 서버에서 여러 Node를 실행할 수도 있지만, 일반적으로 서버 한 대에 하나의 Node를 실행하는 것이 권장


## 1. Node의 종류 및 역할
설정파일(`opensearch.yml`)에 명시하지 않으면, 노드는 기본적으로 마스터, 데이터, 코디네이터 역할을 모두 수행  
소규모 클러스터에서는 유용하지만, 대규모 클러스터에서는 성능과 안정성을 위해 역할을 분리하는 것이 권장

### 1.1. 마스터 노드 (Master Node)
클러스터를 관리하는 역할  
- 노드 추가/삭제, 인덱스 생성/삭제, 샤드 할당, 클러스터 상태 변경 등의 작업을 수행  
- 클러스터의 안정성을 위해 홀수 개의 마스터 노드를 구성하는 것이 권장됨. (split brain 이슈)

### 1.2. 데이터 노드 (Data Node)
실제 데이터의 저장, 검색, 집계, 색인 작업을 처리하는 역할
- 클러스터의 성능과 용량에 직접적인 영향을 미치며, 데이터 처리량이 많을수록 더 많은 데이터 노드가 필요
- 디스크 I/O, CPU, 메모리 사용량이 높은 작업을 주로 처리

### 1.3. 코디네이터 노드 (Coordinator Node)
클라이언트 요청을 받아 적절한 데이터 노드로 요청을 분산시키는 게이트웨이 역할
- 검색 및 집계 결과를 통합하여 클라이언트에 반환
- 데이터 저장 및 색인 작업을 수행하지 않으므로, 클러스터의 부하를 줄이는 데 유용

### 1.4. 인제스트 노드 (Ingest Node)
pipeline을 통해, 인덱싱 전에 문서 변환 및 데이터 전처리를 담당하는 역할  
- 로그 데이터나 비정형 데이터를 정형화된 형태로 변환하는 데 유용
- `node.ingest: true`으로 명시적 활성화 필요 (default: `node.ingest: false`)

### 1.5. 머신 러닝 노드 (Machine Learning Node)
머신 러닝 기능 사용하는 경우, 전용으로 학습/예측을 처리하는 역할
- 머신 러닝 작업이 클러스터의 다른 작업에 영향을 주지 않도록 전용 노드로 분리하는 것이 권장됨
- `node.ml: true`으로 명시적 활성화 필요 (default: `node.ml: false`)

## 2. 설정
`opensearch.yml`에 클러스터, 노드, 네트워크, 플러그인 등을 설정할 수 있음

### 2.1. Hot 노드 기본 설정
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

### 2.2. Warm 노드 기본 설정
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

### 2.3. Hot-Warm 할당 설정
```yaml
# 샤드를 할당할 때 tier 속성 기반으로 각 노드에 균등하게 할당
cluster.routing.allocation.awareness.attributes: tier
# 샤드 할당 시 강제로 다른 속성에 할당.
# 할당할 곳이 없을 때 클러스터는 yellow 상태로 전환되며 할당 대기하며, 다시 할당 가능해 지는 경우 복제본 할당
cluster.routing.allocation.awareness.force.tier.values: hot, warm
```

### 2.4. 플러그인 설정
```yaml
# ISM 플러그인
plugins.index_state_management.enabled: true

# Security 플러그인
plugins.security.ssl.transport.pemcert_filepath: certs/node.pem
plugins.security.ssl.transport.pemkey_filepath: certs/node-key.pem
plugins.security.ssl.transport.pemtrustedcas_filepath: certs/root-ca.pem
plugins.security.ssl.http.enabled: true
plugins.security.ssl.http.pemcert_filepath: certs/node.pem
plugins.security.ssl.http.pemkey_filepath: certs/node-key.pem
plugins.security.ssl.http.pemtrustedcas_filepath: certs/root-ca.pem
```

### 2.5. jvm 설정
`jvm.options`에 JVM 힙 메모리 및 GC 등을 설정할 수 있음
```bash
# 힙 크기 설정 (물리 메모리의 50% 이하, 최대 31GB)
-Xms16g
-Xmx16g

# GC 설정 (G1GC 권장)
-XX:+UseG1GC
```

## 3. 노드 관리
노드수가 변경되면 클러스터가 자동으로 데이터를 재분배.  
이 과정에서 클러스터 상태가 `yellow`나 `red`로 바뀔 수 있으므로, 재분배 완료될때까지 잠시 기다려야함

### 3.1. 노드 추가
1. 새 노드에 동일한 `cluster.name` 설정
2. `discovery.seed_hosts`에 기존 노드 주소 추가
3. 노드 시작 후 클러스터 자동 조인 확인

### 3.2. 노드 제거
노드 제거 시 데이터 손실을 방지하기 위해 Replica를 줄여야함
1. 안전한 노드 제거를 위한 샤드 배제 
    ```bash
    PUT /_cluster/settings
    {
      "transient": {
        "cluster.routing.allocation.exclude._name": "node_name"
      }
    }
    ```
2. 샤드 재할당 완료 후 노드 종료


### 3.3. 노드 역할 변경
1. `opensearch.yml` 파일 수정  
`node.roles: [data] # 데이터 노드로 변경`
2. 노드 재시작

### 3.4. 장애 복구 시나리오
- **메모리 부족**: 힙 크기 조정 또는 하드웨어 증설
- **디스크 풀**: 오래된 인덱스 삭제 또는 디스크 증설  
- **네트워크 분리**: 방화벽, DNS 설정 확인
- **Split-brain**: 마스터 노드 개수를 홀수로 설정

## 4. 모니터링
**Opensearch Dashboard**나 **REST API**로 모니터링 가능

### 4.1. 모니터링 API
```bash
# 클러스터 전체 상태 확인
# green(정상), yellow(경고), red(위험)
GET /_cluster/health

# 클러스터 설정 확인
# 현재 클러스터 동적 설정 확인
GET /_cluster/settings

# 할당되지 않은 샤드 확인
GET /_cluster/allocation/explain

# 노드 목록 및 역할 확인
# 노드별 CPU, 메모리, 디스크 사용량 및 역할 표시
GET /_cat/nodes?v

# 노드 상세 정보
# 노드별 상세 성능 지표 및 통계
GET /_nodes/stats

# 샤드 할당 상태 확인
GET /_cat/shards?v&h=index,shard,prirep,state,node
```

### 4.2. 주요 모니터링 지표
| 지표 | 정상 범위 | 주의사항 |
|------|-----------|----------|
| **힙 메모리** | < 75% | 75% 초과 시 GC 빈발, 성능 저하 |
| **디스크 사용량** | < 85% | 85% 초과 시 샤드 할당 중단 |
| **CPU 사용률** | < 80% | 지속적 고사용률 시 스케일아웃 고려 |
| **로드 평균** | < CPU 코어 수 | 초과 시 성능 병목 발생 |


## 관련 문서
- [Cluster](./Cluster.md)
- [Shard](./Shard.md)
- [Index](./Index.md)