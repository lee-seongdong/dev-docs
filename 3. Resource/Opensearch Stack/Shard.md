# Shard
데이터를 분산 저장하기 위한 물리적 단위로, Lucene 인덱스의 기본 단위이다.  
각 인덱스는 여러 샤드로 나뉘어 저장되어 병렬처리와 확장성 제공


## 1. 분류
### 1.1. Primary Shard
원본 데이터를 저장하는 샤드. 여러 노드에 배치되면 병렬 처리 가능
- 역할: 모든 인덱싱(쓰기) 작업의 시작점
- 설정: 인덱스 생성 시 개수 결정 (이후 변경 불가)
- 권장값: 노드 수 * 1~3 (노드 수 보다 많아야 분산처리 가능)

### 1.2. Replica Shard  
Primary Shard의 복제본
- 목적:
  - 데이터 가용성과 장애 복구
  - 읽기 작업의 부하 분산
  - 검색 성능 향상
  - 고가용성 제공
- 설정: 개수는 언제든 동적 변경 가능
- 권장값: 일반적으로 `노드 수 - 1`로 설정 (모든 노드에 복제)


## 2. 샤드의 크기 및 개수
### 2.1. 샤드의 크기
| 크기 범위 | 상태 | 설명 |
|----------|------|------|
| **10-30GB** | 최적 | 성능과 관리의 최적 균형점 |
| **30-50GB** | 양호 | 허용 가능하지만 모니터링 필요 |
| **50-65GB** | 위험 | 성능 저하 시작, 분할 고려 |
| **65GB+** | 심각 | 즉시 샤드 분할 또는 재설계 필요 |

**영향**:
- 너무 클때: 성능저하, 복구시간 증가, 재분배 시간 증가, 메모리 사용량 증가
- 너무 작을때: 오버헤드 증가, 관리 복잡성 증가

### 2.2. 샤드의 수
`예상 데이터 크기 ÷ 원하는 샤드 크기 = 필요한 샤드 수`  
**예시:**
> 노드 수: 5  
> 권장 샤드 수: 5~15  
> 총 데이터: 300GB  
>  
> 목표 샤드 크기: 30GB  
> 필요한 Primary Shard: 300GB ÷ 30GB = 10개(적절함)


## 3. 샤드 관리
### 3.1. 샤드 배치
#### 기본 배치 규칙
- **다른 노드에 배치**: Primary와 Replica는 서로 다른 노드에 배치
- **균등 분산**: 각 노드에 샤드가 고르게 분산
- **랙 인식**: 가능하면 다른 랙(Rack)에 복제본 배치

#### 속성 기반 배치
```yaml
# Hot-Warm 아키텍처에서의 샤드 배치
cluster.routing.allocation.awareness.attributes: tier
cluster.routing.allocation.awareness.force.tier.values: hot, warm

# 특정 노드에만 배치
PUT /my-index/_settings
{
  "index.routing.allocation.require.tier": "hot"
}

# 특정 노드 제외
PUT /my-index/_settings
{
  "index.routing.allocation.exclude._name": "node-1"
}
```

### 3.2. 샤드 상태
#### 샤드 상태 확인
```bash
# 모든 인덱스의 샤드 분포
GET /_cat/shards?v

# 특정 인덱스의 샤드 상태
GET /_cat/shards/my-index?v

# 상세한 샤드 정보
GET /_cluster/health?level=shards

# 할당되지 않은 샤드 확인
GET /_cat/shards?h=index,shard,prirep,state,unassigned.reason&v
```

#### 샤드 재할당
```bash
# 수동 샤드 재할당
POST /_cluster/reroute
{
  "commands": [
    {
      "move": {
        "index": "my-index",
        "shard": 0,
        "from_node": "node-1",
        "to_node": "node-2"
      }
    }
  ]
}

# 할당되지 않은 샤드 강제 할당
POST /_cluster/reroute
{
  "commands": [
    {
      "allocate_empty_primary": {
        "index": "my-index",
        "shard": 1,
        "node": "node-1",
        "accept_data_loss": true
      }
    }
  ]
}
```

### 3.3. 샤드 최적화
```bash
# 1. 적절한 샤드 수 설정
PUT /my-index
{
  "settings": {
    "number_of_shards": 5,        # 노드 수에 맞게 조정
    "number_of_replicas": 1
  }
}

# 2. 샤드 크기 모니터링
GET /_cat/shards?v&h=index,shard,prirep,store&s=store:desc

# 3. 과도한 샤드 통합 (Shrink)
PUT /my-index/_settings
{
  "index.number_of_replicas": 0,
  "index.blocks.write": true
}

POST /my-index/_shrink/smaller-index
{
  "settings": {
    "index.number_of_shards": 2   # 5개에서 2개로 축소
  }
}

# Hot-Warm 아키텍처에서의 샤드 관리
# 4. Hot 단계: 더 많은 샤드로 분산
PUT /logs-hot/_settings
{
  "number_of_shards": 10,
  "number_of_replicas": 1,
  "index.routing.allocation.require.tier": "hot"
}

# 5. Warm 단계: 샤드 통합 및 이동
POST /logs-hot/_shrink/logs-warm
{
  "settings": {
    "index.number_of_shards": 2,  # 샤드 수 축소
    "index.routing.allocation.require.tier": "warm",
    "index.codec": "best_compression"
  }
}
```

## 4. 모니터링
### 4.1. 모니터링 API
```bash
# 샤드 레벨 상세 통계
GET /_stats/store,docs,indexing,search

# 노드별 샤드 분포 확인
GET /_cat/nodeattrs?v&h=node,attr,value

# 클러스터 전체 헬스 체크
GET /_cluster/health?level=indices&pretty

# 샤드 재할당 진행 상황
GET /_cat/recovery?v&active_only=true
```

### 4.2. 주요 모니터링 지표
| 지표 | 정상 범위 | 주의사항 |
|------|-----------|----------|
| **샤드 크기** | 10-50GB | 50GB 초과 시 성능 저하, 65GB+ 시 심각한 문제 |
| **노드당 샤드 수** | < 20개 | 20개 초과 시 메모리 오버헤드 증가 |
| **할당되지 않은 샤드** | 0개 | 1개 이상 시 즉시 원인 분석 필요 |
| **노드별 샤드 분포** | ±20% 이내 | 30% 이상 차이 시 재할당 고려 |
| **Replica 샤드 비율** | Primary:Replica = 1:1 | 노드 장애 시 데이터 손실 방지 |
| **클러스터 상태** | Green | Yellow: Replica 문제, Red: Primary 문제 |
| **샤드 이동 시간** | < 30분 | 장시간 지속 시 네트워크/디스크 확인 |

### 4.3. 권장 알림 설정
- **Critical**:
  - Red 상태: Primary 샤드 할당 실패
  - Unassigned Primary 샤드 발생
  - 클러스터 전체 쓰기 불가

- **Warning**:
  - Yellow 상태: Replica 샤드 할당 실패
  - 샤드 크기 50GB 초과
  - 노드별 샤드 불균형 30% 이상

- **Info**:
  - 샤드 재할당 시작/완료
  - 새로운 인덱스 생성
  - 샤드 통합(Shrink) 작업

## 관련 문서
- [Cluster](./Cluster.md)
- [Node](./Node.md)
- [Index](./Index.md)

