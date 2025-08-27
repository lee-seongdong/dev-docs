# Indexing
Opensearch에 Document(문서)를 색인하는 작업

## Document 구조
### 기본 문서 형태
```json
{
  "_index": "user-logs",           // 소속 인덱스
  "_id": "unique-id-123",          // 문서의 고유 식별자
  "_source": {                     // 원본 JSON 데이터
    "user": "john",
    "action": "login",
    "timestamp": "2024-01-15T10:30:00Z",
    "ip": "192.168.1.100"
  },
  "_score": 1.5                    // 검색 시 관련도 점수
}
```

### 메타 필드
- **_index**: 문서가 저장된 인덱스명
- **_id**: 문서의 고유 식별자 (자동 생성 또는 직접 지정)
- **_source**: 원본 JSON 데이터
- **_score**: 검색 시 관련도 점수 (색인 시에는 없음)
- **_version**: 문서 버전 (업데이트마다 증가)
- **_routing**: 샤드 라우팅 값

## 단일 문서 색인
### 기본 색인 (자동 ID)
```bash
# POST로 자동 ID 생성
POST /user-logs/_doc/
{
  "user": "john",
  "action": "login",
  "timestamp": "2024-01-15T10:30:00Z",
  "ip": "192.168.1.100"
}

# 응답
{
  "_index": "user-logs",
  "_id": "auto-generated-id",
  "_version": 1,
  "result": "created"
}
```

### 수동 ID 지정
```bash
# PUT으로 ID 직접 지정
PUT /user-logs/_doc/user-john-001
{
  "user": "john",
  "action": "login",
  "timestamp": "2024-01-15T10:30:00Z",
  "ip": "192.168.1.100"
}
```

### 조건부 색인 (Create)
```bash
# 문서가 없을 때만 생성 (409 에러 방지)
PUT /user-logs/_create/user-john-001
{
  "user": "john",
  "action": "login"
}

# 또는
POST /user-logs/_doc/user-john-001?op_type=create
{
  "user": "john", 
  "action": "login"
}
```

## 벌크 색인
### 기본 벌크 색인
```bash
POST /_bulk
{"index":{"_index":"user-logs"}}
{"user":"alice","action":"search","timestamp":"2024-01-15T10:31:00Z"}
{"index":{"_index":"user-logs","_id":"user-bob-001"}}
{"user":"bob","action":"logout","timestamp":"2024-01-15T10:32:00Z"}
{"create":{"_index":"user-logs","_id":"user-charlie-001"}}
{"user":"charlie","action":"login","timestamp":"2024-01-15T10:33:00Z"}
```

### 인덱스별 벌크 색인
```bash
POST /user-logs/_bulk
{"index":{}}
{"user":"alice","action":"search"}
{"index":{"_id":"user-bob-001"}}
{"user":"bob","action":"logout"}
{"create":{"_id":"user-charlie-001"}}
{"user":"charlie","action":"login"}
```

### 벌크 응답 처리
```json
{
  "took": 30,
  "errors": false,
  "items": [
    {
      "index": {
        "_index": "user-logs",
        "_id": "auto-id",
        "_version": 1,
        "result": "created",
        "status": 201
      }
    },
    {
      "index": {
        "_index": "user-logs", 
        "_id": "user-bob-001",
        "_version": 1,
        "result": "created",
        "status": 201
      }
    }
  ]
}
```

## 문서 업데이트 전략
### 1. 전체 문서 교체
```bash
# 기존 문서를 완전히 교체
PUT /user-logs/_doc/user-john-001
{
  "user": "john",
  "action": "logout",                    // 변경된 값
  "timestamp": "2024-01-15T11:00:00Z",   // 변경된 값
  "ip": "192.168.1.100"
}
```

### 2. 부분 업데이트 (Partial Update)
```bash
# 특정 필드만 업데이트
POST /user-logs/_update/user-john-001
{
  "doc": {
    "action": "logout",
    "timestamp": "2024-01-15T11:00:00Z"
  }
}

# 스크립트를 이용한 업데이트
POST /user-logs/_update/user-john-001
{
  "script": {
    "source": "ctx._source.login_count += params.increment",
    "params": {
      "increment": 1
    }
  }
}
```

### 3. Upsert (Update or Insert)
```bash
# 문서가 있으면 업데이트, 없으면 생성
POST /user-logs/_update/user-john-001
{
  "doc": {
    "action": "logout",
    "last_seen": "2024-01-15T11:00:00Z"
  },
  "upsert": {
    "user": "john",
    "action": "logout",
    "created": "2024-01-15T11:00:00Z",
    "login_count": 1
  }
}
```

### 4. 벌크 업데이트
```bash
POST /_bulk
{"update":{"_index":"user-logs","_id":"user-john-001"}}
{"doc":{"action":"logout"}}
{"update":{"_index":"user-logs","_id":"user-bob-001"}}
{"doc":{"last_seen":"2024-01-15T11:00:00Z"},"upsert":{"user":"bob","created":"2024-01-15"}}
```

## 색인 성능 최적화
### 1. 색인 설정 최적화
```bash
# 색인 중 성능 최적화 설정
PUT /my-index/_settings
{
  "refresh_interval": "30s",        // 새로고침 간격 증가
  "number_of_replicas": 0,          // 복제본 임시 비활성화
  "translog.durability": "async",   // 비동기 트랜스로그
  "translog.sync_interval": "30s"   // 동기화 간격 증가
}

# 색인 완료 후 복원
PUT /my-index/_settings
{
  "refresh_interval": "1s",
  "number_of_replicas": 1,
  "translog.durability": "request"
}
```

### 2. 벌크 크기 최적화
```bash
# 권장 벌크 크기: 5-15MB 또는 1000-5000 문서
# 너무 크면: 메모리 부족, 타임아웃
# 너무 작으면: 오버헤드 증가

# 최적 크기 찾기 예시
for size in 1000 2000 5000 10000; do
  echo "Testing bulk size: $size"
  # 벌크 색인 수행 및 성능 측정
done
```

### 3. 병렬 처리
```python
# Python 예시: 멀티스레드 벌크 색인
import threading
from opensearchpy import OpenSearch
from opensearchpy.helpers import bulk

def index_worker(documents, thread_id):
    client = OpenSearch([{'host': 'localhost', 'port': 9200}])
    
    def doc_generator():
        for doc in documents:
            yield {
                "_index": "my-index",
                "_source": doc
            }
    
    bulk(client, doc_generator(), chunk_size=1000)

# 데이터를 스레드별로 분할하여 병렬 처리
threads = []
chunk_size = len(documents) // 4  # 4개 스레드

for i in range(4):
    start = i * chunk_size
    end = start + chunk_size if i < 3 else len(documents)
    thread = threading.Thread(target=index_worker, args=(documents[start:end], i))
    threads.append(thread)
    thread.start()

for thread in threads:
    thread.join()
```

## 실시간 vs 배치 색인
### 실시간 색인
```bash
# 특징
- 개별 문서 즉시 색인
- 빠른 검색 가능 (refresh_interval 내)
- 상대적으로 낮은 처리량

# 설정
PUT /realtime-logs/_settings
{
  "refresh_interval": "1s",          // 빠른 검색 가능
  "number_of_replicas": 1            // 고가용성 유지
}

# 사용 예시: 실시간 모니터링, 알림 시스템
POST /realtime-logs/_doc/
{
  "level": "ERROR",
  "message": "Critical system failure",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### 배치 색인
```bash
# 특징  
- 대량 데이터 효율적 처리
- 높은 처리량
- 검색까지 지연 시간 존재

# 최적화 설정
PUT /batch-logs/_settings
{
  "refresh_interval": "-1",          // 수동 새로고침
  "number_of_replicas": 0,           // 복제본 비활성화
  "translog.durability": "async"     // 비동기 모드
}

# 배치 처리 후 복원
POST /batch-logs/_refresh            // 수동 새로고침
PUT /batch-logs/_settings
{
  "refresh_interval": "30s",
  "number_of_replicas": 1,
  "translog.durability": "request"
}
```

## 모니터링 및 성능 측정
### 색인 통계 확인
```bash
# 인덱스별 색인 통계
GET /my-index/_stats/indexing

# 노드별 색인 통계  
GET /_nodes/stats/indices/indexing

# 색인 성능 지표
GET /_cat/indices/my-index?v&h=index,docs.count,store.size,indexing.index_total,indexing.index_time_in_millis
```

### 성능 튜닝 체크포인트
1. **색인 속도**: docs/sec 측정
2. **메모리 사용량**: 힙 메모리 모니터링  
3. **디스크 I/O**: 디스크 사용률 확인
4. **네트워크**: 클러스터 간 통신 지연
5. **에러율**: 실패한 색인 작업 비율
