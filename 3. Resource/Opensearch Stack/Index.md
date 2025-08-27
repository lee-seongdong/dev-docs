# Index
데이터를 저장하는 논리적 공간으로, RDBMS의 table과 유사한 역할을 함

## 1. Index 구성요소
- [**Policy**](./Index%20Policy.md) - 인덱스의 생명주기 관리
- [**Shard / Replica**](./Shard.md) - 인덱스를 구성하는 물리적 단위
- [**Document**](./Document.md) - 인덱스에 저장되는 실제 데이터


## 2. Index 설정
- [**Index Setting**](./Index%20Settings.md) - 샤드, 레플리카 등 설정
- [**Index Mapping**](./Index%20Mappings.md) - 필드 타입 및 분석기 설정

**예시:**
```json
{
  "settings": {
    "number_of_shards": 3,           // Primary Shard 개수 (정적)
    "number_of_replicas": 1,         // Replica Shard 개수 (동적)
    "refresh_interval": "1s",        // 색인 후 검색 가능해지기까지 대기 시간
    "index.routing.allocation.require.tier": "hot", // 샤드 배치 정책 (예: hot 노드에만)
    "analysis": {                    // 분석기 설정 예시
      "analyzer": {
        "my_korean_analyzer": {
          "type": "custom",
          "tokenizer": "seunjeon_tokenizer"
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "user":    { "type": "keyword" },
      "action":  { "type": "keyword" },
      "ip":      { "type": "ip" },
      "timestamp": { "type": "date" }
    }
  },
  "aliases": {
    "user-logs-current": {}
  }
}

```

## 3. Index 운영
- [**Indexing**](./Indexing.md) - 데이터 색인 및 최적화
- [**Query**](./Query.md) - 데이터 검색


## 관련 문서
- [Policy](./Index%20Policy.md)
- [Shard](./Shard.md)
- [Document](./Document.md)