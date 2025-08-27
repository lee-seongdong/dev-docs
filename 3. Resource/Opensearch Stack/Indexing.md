# Indexing
Opensearch에 Document(문서)를 색인하는 작업  
데이터를 검색 가능한 형태로 저장하는 과정

## 1. 단일 색인
```bash
# 문서 ID 자동 생성
# 장점: 중복 방지, 빠른 색인
# 단점: ID 예측 불가
# 활용: 로그, 이벤트 데이터
POST /my-index/_doc
{
  "user": "john",
  "message": "Hello World",
  "timestamp": "2024-01-15T10:30:00Z"
}

# 문서 ID 지정
# 장점: ID 제어 가능, 업데이트 용이
# 단점: 중복 관리 필요
# 활용: 사용자 정보, 제품 데이터
PUT /my-index/_doc/user-123
{
  "user": "jane",
  "message": "Hello OpenSearch",
  "timestamp": "2024-01-15T10:31:00Z"
}
```

## 2. 벌크 색인
### 2.1. Bulk API
```bash
POST /_bulk
{ "index": { "_index": "my-index", "_id": "1" } }
{ "user": "user1", "message": "First message" }
{ "index": { "_index": "my-index", "_id": "2" } }
{ "user": "user2", "message": "Second message" }
{ "delete": { "_index": "my-index", "_id": "3" } }
{ "update": { "_index": "my-index", "_id": "4" } }
{ "doc": { "message": "Updated message" } }
```

### 2.2. 배치 크기 최적화
| 배치 크기 | 메모리 사용량 | 처리 속도 | 권장 사용 |
|----------|--------------|----------|----------|
| **1-100개** | 낮음 | 느림 | 실시간 소량 데이터 |
| **100-1000개** | 적정 | 최적 | **일반적 권장** |
| **1000개+** | 높음 | 빠르지만 위험 | 대량 초기 적재 |

**권장 설정**:
- 배치 크기: 5-15MB
- 문서 수: 100-1000개
- 동시 요청: 2-8개


## 3. 업데이트
### 3.1. 업데이트
```bash
# 부분 업데이트
POST /my-index/_update/user-123
{
  "doc": {
    "message": "Updated message",
    "last_modified": "2024-01-15T11:00:00Z"
  }
}

# 스크립트 업데이트
POST /my-index/_update/user-123
{
  "script": {
    "source": "ctx._source.view_count += params.increment",
    "params": { "increment": 1 }
  }
}
```

### 3.2. Upsert
```bash
# 문서가 없으면 생성, 있으면 업데이트
POST /my-index/_update/user-456
{
  "doc": {
    "user": "alice",
    "message": "Hello"
  },
  "upsert": {
    "user": "alice",
    "message": "Hello",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

## 4. 삭제
```bash
# 단일 문서 삭제
DELETE /my-index/_doc/user-123

# 조건부 삭제 (Delete by Query)
POST /my-index/_delete_by_query
{
  "query": {
    "range": {
      "timestamp": {
        "lt": "2024-01-01"
      }
    }
  }
}
```


## 4. 모니터링
### 4.1. 모니터링 API
```bash
# 색인 통계 확인
GET /_stats/indexing

# 인덱스별 색인 성능
GET /my-index/_stats/indexing

# 인덱스 상태 확인
GET /_cat/indices?v&s=store.size:desc

# 색인 진행 상황 모니터링
GET /_tasks?actions=*index*&detailed
```

### 4.2. 주요 모니터링 지표
| 지표 | 정상 범위 | 주의사항 |
|------|-----------|----------|
| **색인 속도** | 1000-10000 docs/sec | 급격한 감소 시 원인 분석 |
| **색인 지연** | < 100ms | 지속적 증가 시 최적화 필요 |
| **거부율** | < 1% | 큐 포화 상태 점검 |
| **메모리 사용량** | < 75% | 힙 메모리 부족 주의 |


## 관련 문서
- [Index](./Index.md)
- [Index Settings](./Index%20Settings.md)
- [Index Mappings](./Index%20Mappings.md)
