# Index Settings
인덱스의 동작과 성능을 제어하는 설정들  
인덱스 생성 시 개별적으로 설정할 수도 있고, [Index Template](./Index%20Template.md)의 설정을 따를 수도 있다.

## 1. 설정 유형
### 1.1. Static Settings (정적 설정)
인덱스 생성 시에만 설정 가능하며, 이후 변경 불가

| 설정 | 설명 | 기본값 |
|------|------|--------|
| `number_of_shards` | Primary 샤드 개수 | 1 |
| `codec` | 데이터 압축 방식 | `default` |
| `index.sort.field` | 인덱스 정렬 필드 | - |
| `index.sort.order` | 정렬 순서 (asc/desc) | - |
| `shard.check_on_startup` | 시작 시 샤드 검증 | `false` |

### 1.2. Dynamic Settings (동적 설정)
인덱스 생성 후에도 변경 가능한 설정

| 설정 | 설명 | 기본값 |
|------|------|--------|
| `number_of_replicas` | Replica 샤드 개수 | 1 |
| `refresh_interval` | 검색 새로고침 주기 | `1s` |
| `index.max_result_window` | 검색 결과 최대 수 | `10000` |
| `index.routing.allocation.require.tier` | 노드 tier 할당 | - |
| `index.blocks.read_only` | 읽기 전용 모드 | `false` |
| `index.priority` | 복구 우선순위 | 1 |
| `index.auto_expand_replicas` | 자동 복제본 조정 | `false` |

## 2. 주요 설정 예시
### 2.1. 기본 인덱스
```json
PUT /my-index
{
  "settings": {
    "number_of_shards": 3,
    "number_of_replicas": 1,
    "refresh_interval": "5s",
    "codec": "best_compression"
  }
}
```

### 2.2. Hot-Warm 인덱스
```json
// Hot 단계 - 활발한 인덱싱
PUT /hot-index
{
  "settings": {
    "number_of_shards": 5,
    "number_of_replicas": 1,
    "refresh_interval": "1s",
    "index.routing.allocation.require.tier": "hot"
  }
}

// Warm 단계 - 검색 최적화
PUT /warm-index
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0,
    "refresh_interval": "30s",
    "codec": "best_compression",
    "index.routing.allocation.require.tier": "warm"
  }
}
```

## 3. 설정 확인 및 관리
### 3.1. 설정 조회
```bash
# 전체 인덱스 설정 확인
GET /my-index/_settings

# 특정 설정만 확인
GET /my-index/_settings/index.number_of_replicas

# 여러 인덱스 설정 확인
GET /logs-*/_settings

# 모든 인덱스 설정 확인
GET /_all/_settings
```

### 3.2. 동적 설정 변경
```bash
# Replica 수 변경
PUT /my-index/_settings
{
  "number_of_replicas": 2
}

# 새로고침 주기 변경  
PUT /my-index/_settings
{
  "refresh_interval": "10s"
}

# 노드 tier 변경 (Hot → Warm)
PUT /my-index/_settings
{
  "index.routing.allocation.require.tier": "warm"
}

# 읽기 전용 모드 설정
PUT /my-index/_settings
{
  "index.blocks.read_only": true
}
```

### 3.3. 설정 초기화
```bash
# 기본값으로 복원
PUT /my-index/_settings
{
  "refresh_interval": null,
  "number_of_replicas": null
}
```


## 관련 문서
- [Index](./Index.md)
- [Index Template](./Index%20Template.md)
- [Shard](./Shard.md)