# Query
Opensearch에서 데이터를 검색하고 조회하는 방법  
다양한 검색 조건과 필터를 사용하여 원하는 데이터를 찾는 과정

## 1. 기본 검색
```bash
# Match Query - 전문 검색 (분석기 적용)
GET /my-index/_search
{
  "query": {
    "match": {
      "message": "hello world"
    }
  }
}

# Term Query - 정확한 매칭 (분석기 미적용)
GET /my-index/_search
{
  "query": {
    "term": {
      "status": "active"
    }
  }
}

# Range Query - 범위 검색
GET /my-index/_search
{
  "query": {
    "range": {
      "timestamp": {
        "gte": "2024-01-01",
        "lte": "2024-01-31"
      }
    }
  }
}

# Wildcard Query - 패턴 검색
GET /my-index/_search
{
  "query": {
    "wildcard": {
      "user": "john*"
    }
  }
}
```

## 2. 복합 검색
```bash
GET /my-index/_search
{
  "query": {
    "bool": {
      "must": [                    # 반드시 포함 (AND 조건)
        { "match": { "title": "elasticsearch" } }
      ],
      "should": [                  # 포함되면 좋음 (OR 조건) 
        { "match": { "content": "search" } }
      ],
      "must_not": [               # 제외 조건
        { "term": { "status": "deleted" } }
      ],
      "filter": [                 # 스코어에 영향 없는 필터
        { "range": { "created_at": { "gte": "2024-01-01" } } }
      ]
    }
  }
}
```

| 조건 | 용도 | 스코어 영향 |
|------|------|------------|
| `must` | 필수 조건 | 영향 있음 |
| `should` | 선택 조건 | 영향 있음 |
| `must_not` | 제외 조건 | 영향 없음 |
| `filter` | 필터 조건 | 영향 없음 |

## 3. 정렬 및 페이징
```bash
GET /my-index/_search
{
  "query": { "match_all": {} },
  "sort": [
    { "timestamp": { "order": "desc" } },
    { "_score": { "order": "desc" } }
  ],
  "from": 0,                     # 페이징 시작점
  "size": 10,                    # 결과 개수
  "_source": ["title", "timestamp"] # 반환 필드 선택
}

# 대용량 페이징 (Search After)
GET /my-index/_search
{
  "query": { "match_all": {} },
  "sort": [{ "timestamp": "desc" }],
  "search_after": ["2024-01-15T10:30:00Z"],
  "size": 100
}
```

## 4. 고급 검색
```bash
# 다중 인덱스 검색
GET /logs-*,users-*/_search
{
  "query": { "match": { "message": "error" } }
}

# 중첩 객체 검색
GET /my-index/_search
{
  "query": {
    "nested": {
      "path": "tags",
      "query": {
        "bool": {
          "must": [
            { "match": { "tags.name": "important" } },
            { "range": { "tags.score": { "gte": 5 } } }
          ]
        }
      }
    }
  }
}
```

## 5. 최적화
- **Filter Context 활용**: 스코어 계산이 불필요한 조건은 `filter` 사용
- **필드 선택**: `_source`로 필요한 필드만 반환
- **페이징 최적화**: 대용량 데이터는 `search_after` 사용


## 관련 문서
- [Index](./Index.md)
- [Index Mappings](./Index%20Mappings.md)
- [Aggregation](./Aggregation.md)
