# Index Settings
인덱스의 동작과 성능을 제어하는 설정들

## 설정 유형

### Static Settings (정적 설정)
인덱스 생성 시에만 설정 가능하며, 이후 변경 불가

### Dynamic Settings (동적 설정)  
인덱스 생성 후에도 변경 가능한 설정

## Primary & Replica Shard 설정

### 기본 개념
```json
{
  "settings": {
    "number_of_shards": 3,      // Primary Shard 개수 (Static)
    "number_of_replicas": 1     // Replica Shard 개수 (Dynamic)
  }
}
```

### Primary Shard (Primary Shard)
- **정의**: 원본 데이터를 저장하는 샤드
- **설정**: 인덱스 생성 시 결정, **이후 변경 불가**
- **권장값**: `노드 수 * 1~3` (데이터 분산을 위해)
- **고려사항**: 
  - 너무 많으면: 오버헤드 증가, 메모리 낭비
  - 너무 적으면: 확장성 제한, 성능 병목

### Replica Shard (복제본)
- **정의**: Primary Shard의 복제본
- **설정**: **언제든 동적 변경 가능**
- **권장값**: `노드 수 - 1` (모든 노드에 분산)
- **기능**:
  - 고가용성 보장 (노드 장애 시 서비스 지속)
  - 읽기 성능 향상 (부하 분산)
  - 검색 처리량 증가

### 샤드 크기 고려사항
```bash
# 적절한 샤드 크기: 10~50GB
# 전체 데이터 크기를 기준으로 계산

예시: 300GB 데이터
- 적절한 샤드 개수: 300GB ÷ 30GB = 10개 샤드
- 3노드 클러스터라면: 10개 Primary Shard, 1개 Replica
```

## 용량 계획

### 저장 용량 계산
```
전체 필요 용량 = 원본 데이터 × (1 + 복제본 수) × 오버헤드(1.1)

예시:
- 원본 데이터: 100GB
- Replica 개수: 1개  
- 실제 필요 용량: 100GB × (1 + 1) × 1.1 = 220GB
```

### 성능 기반 용량 계획
```bash
# 1. 색인 성능 우선
- 더 많은 Primary Shard (병렬 색인)
- SSD 스토리지 사용
- 충분한 힙 메모리

# 2. 검색 성능 우선  
- 적절한 Replica 개수 (읽기 부하 분산)
- 캐시 최적화
- 필드별 인덱싱 전략

# 3. 저장 비용 우선
- 적은 Replica 개수
- 압축 설정 활용
- Cold 스토리지 활용
```

## 주요 설정 옵션

### 기본 설정
```json
{
  "settings": {
    // 샤드 설정
    "number_of_shards": 3,
    "number_of_replicas": 1,
    
    // 새로고침 설정
    "refresh_interval": "1s",           // 검색 가능한 새로고침 주기
    
    // 압축 설정  
    "codec": "best_compression",        // 압축 레벨 (default, best_compression)
    
    // 트랜스로그 설정
    "translog.durability": "request",   // request(동기), async(비동기)
    "translog.sync_interval": "5s",     // 비동기 모드의 동기화 주기
    
    // 병합 설정
    "merge.policy.max_merge_at_once": 10,
    "merge.policy.segments_per_tier": 10
  }
}
```

### 성능 최적화 설정
```json
{
  "settings": {
    // 색인 성능 최적화
    "refresh_interval": "30s",          // 배치 색인 시 길게 설정
    "number_of_replicas": 0,            // 색인 중에는 복제본 비활성화
    "translog.durability": "async",     // 비동기 모드로 성능 향상
    
    // 검색 성능 최적화
    "index.queries.cache.enabled": true,
    "index.requests.cache.enable": true,
    
    // 메모리 최적화
    "index.store.preload": ["nvd", "dvd"], // 자주 사용하는 파일 미리 로드
    
    // 분석 설정
    "index.max_result_window": 10000,   // 기본 검색 결과 크기 제한
    "index.max_terms_count": 65536      // 집계 시 최대 terms 개수
  }
}
```

### Hot-Warm 아키텍처 설정
```json
{
  "settings": {
    // Hot 노드 설정
    "index.routing.allocation.require.tier": "hot",
    "number_of_shards": 5,              // 더 많은 샤드로 분산
    "refresh_interval": "1s",           // 빠른 검색 가능
    
    // Warm 노드 설정  
    "index.routing.allocation.require.tier": "warm",
    "number_of_shards": 1,              // 샤드 통합
    "refresh_interval": "30s",          // 느린 새로고침
    "codec": "best_compression",        // 압축 활용
    "number_of_replicas": 0             // 복제본 최소화
  }
}
```

## 동적 설정 변경

### Replica 개수 변경
```bash
# 복제본 개수 증가 (읽기 성능 향상)
PUT /my-index/_settings
{
  "number_of_replicas": 2
}

# 복제본 비활성화 (색인 성능 향상)
PUT /my-index/_settings  
{
  "number_of_replicas": 0
}
```

### 새로고침 주기 변경
```bash
# 실시간 검색을 위한 빠른 새로고침
PUT /my-index/_settings
{
  "refresh_interval": "1s"
}

# 배치 색인을 위한 느린 새로고침
PUT /my-index/_settings
{
  "refresh_interval": "30s"
}

# 수동 새로고침 (성능 최적화)
PUT /my-index/_settings
{
  "refresh_interval": -1
}
```

### 할당 설정 변경
```bash
# 특정 노드에 할당
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

## 인덱스 템플릿 활용

### 기본 템플릿
```json
{
  "index_patterns": ["logs-*"],
  "template": {
    "settings": {
      "number_of_shards": 3,
      "number_of_replicas": 1,
      "refresh_interval": "5s",
      "codec": "best_compression"
    }
  },
  "priority": 100
}
```

### 환경별 템플릿
```json
{
  // 개발 환경: 성능 우선
  "index_patterns": ["dev-logs-*"],
  "template": {
    "settings": {
      "number_of_shards": 1,
      "number_of_replicas": 0,
      "refresh_interval": "1s"
    }
  }
}

{
  // 운영 환경: 안정성 우선
  "index_patterns": ["prod-logs-*"],
  "template": {
    "settings": {
      "number_of_shards": 5,
      "number_of_replicas": 2,
      "refresh_interval": "30s",
      "codec": "best_compression"
    }
  }
}
```

## 설정 확인 및 모니터링

### 설정 확인
```bash
# 인덱스 설정 확인
GET /my-index/_settings

# 모든 인덱스 설정 확인
GET /_all/_settings

# 특정 설정만 확인
GET /my-index/_settings/index.number_of_*
```

### 상태 모니터링
```bash
# 샤드 분포 확인
GET /_cat/shards/my-index?v

# 인덱스 통계 확인
GET /my-index/_stats

# 클러스터 할당 상태
GET /_cluster/allocation/explain
```
