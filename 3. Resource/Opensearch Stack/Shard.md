# Shard
데이터를 분산 저장하기 위한 물리적 단위

## 샤드 개념
### 기본 정의
- **Primary Shard**: 원본 데이터를 저장하는 샤드
- **Replica Shard**: Primary Shard의 복제본
- 각 인덱스는 여러 샤드로 나뉘어 저장됨
- 샤드는 Lucene 인덱스의 기본 단위

## 샤드 종류
### Primary Shard
- **특징**: 원본 데이터를 저장하는 샤드
- **설정**: 인덱스 생성 시 개수 결정 (이후 변경 불가)
- **역할**: 모든 인덱싱(쓰기) 작업의 시작점
- **분산**: 여러 노드에 분산 배치되어 병렬 처리 가능

### Replica Shard  
- **특징**: Primary Shard의 복제본
- **목적**: 
  - 데이터 가용성과 장애 복구
  - 읽기 작업의 부하 분산
  - 검색 성능 향상
  - 고가용성 제공
- **설정**: 개수는 언제든 동적 변경 가능
- **권장값**: 일반적으로 `노드 수 - 1`로 설정

## 샤드 크기 고려사항
### 크기 기준
```bash
# 적절한 샤드 크기
- 권장 크기: 10~50GB
- 최적 크기: 20~30GB
- 최대 크기: 65GB (성능 저하 시작점)
```

### 크기 영향 요소
- **너무 작으면**: 오버헤드 증가, 관리 복잡성 증가
- **너무 크면**: 
  - 성능 저하 
  - 복구 시간 증가
  - 재분배 시간 증가
  - 메모리 사용량 증가

### 샤드 수 계산
```bash
# 기본 공식
예상 데이터 크기 ÷ 원하는 샤드 크기 = 필요한 샤드 수

# 예시
- 총 데이터: 300GB
- 목표 샤드 크기: 30GB
- 필요한 Primary Shard: 300GB ÷ 30GB = 10개

# 노드 수 고려
- 노드 수보다 샤드 수가 많아야 분산 처리 가능
- 권장: 노드 수 * 1~3배
```

## 샤드 배치 원칙
### 기본 배치 규칙
- **다른 노드에 배치**: Primary와 Replica는 서로 다른 노드에 배치
- **균등 분산**: 각 노드에 샤드가 고르게 분산
- **랙 인식**: 가능하면 다른 랙(Rack)에 복제본 배치

### 속성 기반 배치
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

## 샤드 상태 관리
### 샤드 상태 확인
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

### 샤드 재할당
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

## 샤드 최적화
### 성능 최적화
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
```

### Hot-Warm 아키텍처에서의 샤드 관리
```bash
# Hot 단계: 더 많은 샤드로 분산
PUT /logs-hot/_settings
{
  "number_of_shards": 10,
  "number_of_replicas": 1,
  "index.routing.allocation.require.tier": "hot"
}

# Warm 단계: 샤드 통합 및 이동
POST /logs-hot/_shrink/logs-warm
{
  "settings": {
    "index.number_of_shards": 2,  # 샤드 수 축소
    "index.routing.allocation.require.tier": "warm",
    "index.codec": "best_compression"
  }
}
```

## 모니터링 및 알림

### 주요 모니터링 지표
```bash
# 샤드 레벨 통계
GET /_stats/store,docs

# 노드별 샤드 분포
GET /_cat/nodeattrs?v

# 클러스터 건강상태  
GET /_cluster/health?level=indices
```

### 권장 알림 설정
- **Red 상태**: Primary 샤드 할당 실패
- **Yellow 상태**: Replica 샤드 할당 실패  
- **Unassigned 샤드 발생**: 즉시 알림
- **노드별 샤드 불균형**: 30% 이상 차이
- **샤드 크기**: 50GB 초과 시 경고
