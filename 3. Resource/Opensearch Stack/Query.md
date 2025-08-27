# Query
Opensearch에서 데이터를 검색하고 조회하는 방법

## 기본 쿼리 구조

### Query DSL 기본 형태
```json
{
  "query": {
    "query_type": {
      "field": "value"
    }
  },
  "size": 10,                    // 결과 개수 (기본: 10)
  "from": 0,                     // 시작 위치 (페이징)
  "sort": [                      // 정렬
    { "timestamp": "desc" }
  ],
  "_source": ["field1", "field2"]  // 반환할 필드 지정
}
```

### 간단한 검색 예시
```bash
# 전체 문서 검색
GET /user-logs/_search

# 특정 필드 값 검색
GET /user-logs/_search
{
  "query": {
    "match": {
      "user": "john"
    }
  }
}
```

## 기본 쿼리 타입

### 1. **Match Query** (전문 검색)
```bash
# 기본 match 쿼리
GET /user-logs/_search
{
  "query": {
    "match": {
      "message": "login error"     // "login" 또는 "error" 포함 문서
    }
  }
}

# match_phrase 쿼리 (구문 검색)
GET /user-logs/_search
{
  "query": {
    "match_phrase": {
      "message": "login error"     // "login error" 정확한 구문
    }
  }
}

# match_phrase_prefix (자동완성)
GET /user-logs/_search
{
  "query": {
    "match_phrase_prefix": {
      "message": "user log"        // "user log*" 형태로 검색
    }
  }
}
```

### 2. **Term Query** (정확한 매치)
```bash
# 정확한 값 검색 (keyword 필드용)
GET /user-logs/_search
{
  "query": {
    "term": {
      "status.keyword": "active"   // 정확히 "active"인 문서
    }
  }
}

# 여러 값 중 하나 (terms 쿼리)
GET /user-logs/_search
{
  "query": {
    "terms": {
      "status.keyword": ["active", "pending", "completed"]
    }
  }
}
```

### 3. **Range Query** (범위 검색)
```bash
# 숫자 범위
GET /user-logs/_search
{
  "query": {
    "range": {
      "age": {
        "gte": 18,               // 18 이상
        "lte": 65                // 65 이하
      }
    }
  }
}

# 날짜 범위
GET /user-logs/_search
{
  "query": {
    "range": {
      "timestamp": {
        "gte": "2024-01-01",
        "lte": "2024-01-31",
        "format": "yyyy-MM-dd"
      }
    }
  }
}

# 상대적 날짜
GET /user-logs/_search
{
  "query": {
    "range": {
      "timestamp": {
        "gte": "now-7d",         // 7일 전부터
        "lte": "now"             // 현재까지
      }
    }
  }
}
```

### 4. **Exists & Missing Query**
```bash
# 필드가 존재하는 문서
GET /user-logs/_search
{
  "query": {
    "exists": {
      "field": "email"
    }
  }
}

# 필드가 없거나 null인 문서
GET /user-logs/_search
{
  "query": {
    "bool": {
      "must_not": {
        "exists": {
          "field": "email"
        }
      }
    }
  }
}
```

## Bool Query (복합 쿼리)

### Bool Query 구조
```json
{
  "query": {
    "bool": {
      "must": [/* 반드시 만족 (AND 조건) */],
      "should": [/* 하나라도 만족 (OR 조건) */],
      "must_not": [/* 만족하지 않아야 함 (NOT 조건) */],
      "filter": [/* 점수 계산 없이 필터링 */]
    }
  }
}
```

### 실제 예시
```bash
# 복합 조건 검색
GET /user-logs/_search
{
  "query": {
    "bool": {
      "must": [                          // 반드시 만족
        {
          "match": {
            "status": "active"
          }
        }
      ],
      "should": [                        // 하나라도 만족 (점수 가산)
        {
          "match": {
            "level": "ERROR"
          }
        },
        {
          "match": {
            "level": "WARN"
          }
        }
      ],
      "must_not": [                      // 제외 조건
        {
          "term": {
            "user.keyword": "system"
          }
        }
      ],
      "filter": [                        // 점수 계산 없이 필터
        {
          "range": {
            "timestamp": {
              "gte": "now-1d"
            }
          }
        }
      ],
      "minimum_should_match": 1          // should 절 중 최소 몇 개 만족해야 하는지
    }
  }
}
```

## 필터 vs 쿼리

### 차이점
```bash
# Query Context: 점수 계산 O, 캐싱 X
{
  "query": {
    "match": {
      "title": "search text"            // 관련도 점수 계산
    }
  }
}

# Filter Context: 점수 계산 X, 캐싱 O
{
  "query": {
    "bool": {
      "filter": {
        "term": {
          "status": "published"          // 단순 true/false, 빠른 처리
        }
      }
    }
  }
}
```

### 성능 최적화
```bash
# 권장: filter로 데이터 축소 후 query로 정확도 계산
GET /articles/_search
{
  "query": {
    "bool": {
      "filter": [                       // 먼저 필터로 범위 축소
        {
          "term": { "status": "published" }
        },
        {
          "range": {
            "publish_date": {
              "gte": "2024-01-01"
            }
          }
        }
      ],
      "must": [                         // 축소된 데이터에서 관련도 검색
        {
          "match": {
            "title": "opensearch tutorial"
          }
        }
      ]
    }
  }
}
```

## 정렬 및 페이징

### 정렬
```bash
# 단일 필드 정렬
GET /user-logs/_search
{
  "query": { "match_all": {} },
  "sort": [
    { "timestamp": "desc" }
  ]
}

# 복합 정렬
GET /user-logs/_search
{
  "query": { "match_all": {} },
  "sort": [
    { "priority": "desc" },             // 1순위: priority 내림차순
    { "timestamp": "desc" },            // 2순위: timestamp 내림차순
    { "_score": "desc" }                // 3순위: 관련도 점수
  ]
}

# 조건부 정렬
GET /user-logs/_search
{
  "query": { "match_all": {} },
  "sort": [
    {
      "age": {
        "order": "desc",
        "missing": "_last"              // 값이 없는 문서는 마지막에
      }
    }
  ]
}
```

### 페이징
```bash
# 기본 페이징 (from/size)
GET /user-logs/_search
{
  "query": { "match_all": {} },
  "from": 20,                          // 20번째부터
  "size": 10                           // 10개 반환 (21~30번째)
}

# Search After (대용량 데이터 효율적 페이징)
GET /user-logs/_search
{
  "query": { "match_all": {} },
  "size": 10,
  "sort": [
    { "timestamp": "desc" },
    { "_id": "desc" }                   // 동일한 timestamp 처리용
  ],
  "search_after": ["2024-01-15T10:30:00Z", "doc_id_123"]
}
```

## 고급 검색 기법

### 1. **Wildcard & Regex**
```bash
# 와일드카드 검색
GET /user-logs/_search
{
  "query": {
    "wildcard": {
      "user.keyword": "john*"          // john으로 시작하는 사용자
    }
  }
}

# 정규표현식 검색
GET /user-logs/_search
{
  "query": {
    "regexp": {
      "email.keyword": ".*@company\\.com"
    }
  }
}
```

### 2. **Fuzzy Search** (오타 허용)
```bash
GET /user-logs/_search
{
  "query": {
    "fuzzy": {
      "user": {
        "value": "jhon",               // "john"의 오타
        "fuzziness": "AUTO"            // 자동 오타 보정
      }
    }
  }
}
```

### 3. **Multi-Match** (다중 필드 검색)
```bash
GET /articles/_search
{
  "query": {
    "multi_match": {
      "query": "opensearch tutorial",
      "fields": ["title^2", "content", "tags"],  // title 필드에 가중치 2배
      "type": "best_fields"            // 가장 잘 매치되는 필드 우선
    }
  }
}
```

### 4. **Function Score** (커스텀 점수)
```bash
GET /products/_search
{
  "query": {
    "function_score": {
      "query": {
        "match": { "name": "laptop" }
      },
      "functions": [
        {
          "filter": { "term": { "featured": true }},
          "weight": 2                  // 추천 상품에 가중치
        },
        {
          "field_value_factor": {
            "field": "rating",
            "factor": 0.1              // 평점에 따른 점수 조정
          }
        }
      ],
      "boost_mode": "multiply"
    }
  }
}
```

## 검색 성능 튜닝

### 1. **결과 필드 제한**
```bash
GET /user-logs/_search
{
  "query": { "match_all": {} },
  "_source": {
    "includes": ["user", "action", "timestamp"],
    "excludes": ["large_field", "binary_data"]
  }
}
```

### 2. **타임아웃 설정**
```bash
GET /user-logs/_search?timeout=10s
{
  "query": { "match_all": {} }
}
```

### 3. **샤드 요청 캐시 활용**
```bash
GET /user-logs/_search?request_cache=true
{
  "query": {
    "bool": {
      "filter": {                      // filter는 캐시됨
        "term": { "status": "active" }
      }
    }
  }
}
```

### 4. **Profile API로 성능 분석**
```bash
GET /user-logs/_search
{
  "profile": true,
  "query": {
    "match": { "message": "error" }
  }
}
```

## 검색 결과 분석

### 결과 구조
```json
{
  "took": 5,                           // 검색 소요 시간 (ms)
  "timed_out": false,                  // 타임아웃 여부
  "hits": {
    "total": {
      "value": 1234,                   // 총 매치 문서 수
      "relation": "eq"                 // eq(정확) 또는 gte(이상)
    },
    "max_score": 1.5,                  // 최고 점수
    "hits": [                          // 실제 검색 결과
      {
        "_index": "user-logs",
        "_id": "1",
        "_score": 1.5,
        "_source": {
          "user": "john",
          "action": "login"
        }
      }
    ]
  }
}
```

### 검색 통계 확인
```bash
# 검색 성능 통계
GET /_cat/indices?v&h=index,search.query_total,search.query_time_in_millis

# 느린 쿼리 로그 설정
PUT /_cluster/settings
{
  "transient": {
    "logger.index.search.slowlog.threshold.query.warn": "10s",
    "logger.index.search.slowlog.threshold.fetch.warn": "1s"
  }
}
```
